<script lang="ts">
	import { onMount } from 'svelte';
	import { dndzone } from 'svelte-dnd-action';
	import { auth } from '$lib/stores/auth';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Dialog from '$lib/components/ui/dialog';
	import { toast } from 'svelte-sonner';
	
	import { LinkGroupCard, TextGroupCard, TextBlock, ImageBlock, VideoBlock, SocialBlock, DividerBlock, EmailCollectorBlock, EmbedBlock } from '$lib/components/dashboard/bio/blocks';
	import { AddBlockDialog } from '$lib/components/dashboard/bio/dialogs';
	import CreateGroupDialog from '$lib/components/dashboard/bio/dialogs/CreateGroupDialog.svelte';
	import EditGroupLinkDialog from '$lib/components/dashboard/bio/dialogs/EditGroupLinkDialog.svelte';
	import DuplicateLinkDialog from '$lib/components/dashboard/bio/dialogs/DuplicateLinkDialog.svelte';
	import ProfilePreview from '$lib/components/dashboard/preview/ProfilePreview.svelte';
	import CalendarView from '$lib/components/dashboard/bio/CalendarView.svelte';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import { blocksApi, type Block } from '$lib/api/blocks';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';

	let profile: Profile | null = null;
	let links: Link[] = [];
	let allLinks: Link[] = [];
	let blocks: Block[] = [];
	let initialLoading = true;
	let selectedIds = new Set<string>();

	// Combined items for drag & drop
	type CombinedItem = 
		| { type: 'link'; data: Link; id: string; position: number }
		| { type: 'block'; data: Block; id: string; position: number };
	let items: CombinedItem[] = [];
	
	$: {
		const linkItems = Array.isArray(links) ? links.map(link => ({ type: 'link' as const, data: link, id: link.id, position: link.position })) : [];
		const blockItems = Array.isArray(blocks) ? blocks.map(block => ({ type: 'block' as const, data: block, id: block.id, position: block.position })) : [];
		items = [...linkItems, ...blockItems].sort((a, b) => a.position - b.position);
	}

	// Get link_ids from inactive blocks - these links should be hidden
	$: hiddenLinkIds = new Set(
		Array.isArray(blocks) 
			? blocks
				.filter(block => !block.is_active && block.link_id)
				.map(block => block.link_id!)
			: []
	);

	$: selectedCount = selectedIds.size;

	function handleCopyUrl() {
		if (!profile?.username) {
			toast.error('Profile not loaded yet');
			return;
		}
		const url = `${window.location.origin}/${profile.username}`;
		navigator.clipboard.writeText(url).then(() => {
			toast.success('Profile URL copied to clipboard!');
		}).catch(err => {
			console.error('Copy failed:', err);
			toast.error('Failed to copy URL');
		});
	}

	let showAddBlockDialog = false;
	let showCalendarView = false;
	
	let showCreateGroupDialog = false;
	let showGroupLinkDialog = false;
	let currentGroupId: string | null = null;
	let currentGroupTitle: string = '';
	let editingGroupLink: Link | null = null;
	let isUploadingThumbnail = false;
	
	// Duplicate link dialog
	let showDuplicateLinkDialog = false;
	let duplicatingLink: Link | null = null;
	let duplicatingFromGroupId: string = '';
	
	// Delete confirmation modal
	let showDeleteConfirmModal = false;
	let deleteTarget: { 
		type: 'link-group' | 'text-group' | 'block'; 
		id: string; 
		title: string;
	} | null = null;
	
	// Track which groups are expanded (for both link and text groups)
	let expandedGroupIds = new Set<string>();
	let expandedTextGroupIds = new Set<string>();

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			initialLoading = true;
			
			let profileData = null;
			let linksData: Link[] = [];
			let blocksData: Block[] = [];
			
			try {
				profileData = await profileApi.getMyProfile($auth.token!);
			} catch (profileError: any) {
				console.warn('Profile not found, will be created automatically:', profileError);
			}
			
			try {
				linksData = await linksApi.getLinks($auth.token!);
			} catch (linksError: any) {
				console.warn('No links found:', linksError);
				linksData = [];
			}
			
			try {
				blocksData = await blocksApi.getBlocks($auth.token!);
			} catch (blocksError: any) {
				console.error('Failed to load blocks:', blocksError);
				blocksData = [];
			}
			
			profile = profileData;
			links = linksData;
			allLinks = linksData;
			blocks = blocksData;
		} catch (error: any) {
			console.error('Failed to load dashboard data:', error);
			toast.error(error.message || 'Failed to load data');
		} finally {
			initialLoading = false;
		}
	}

	function handleSelect(event: CustomEvent) {
		const id = event.detail;
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = selectedIds;
	}

	function selectAll() {
		selectedIds = new Set(links.map(l => l.id));
	}

	function clearSelection() {
		selectedIds.clear();
		selectedIds = selectedIds;
	}

	async function bulkDelete() {
		if (selectedIds.size === 0) return;
		
		try {
			await linksApi.bulkAction(Array.from(selectedIds), 'delete', $auth.token!);
			links = [...links.filter(l => !selectedIds.has(l.id))];
			allLinks = [...allLinks.filter(l => !selectedIds.has(l.id))];
			toast.success(`Deleted ${selectedIds.size} links`);
			clearSelection();
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete links');
		}
	}

	async function bulkActivate() {
		if (selectedIds.size === 0) return;
		
		try {
			await linksApi.bulkAction(Array.from(selectedIds), 'activate', $auth.token!);
			links = [...links.map(l => selectedIds.has(l.id) ? {...l, is_active: true} : l)];
			allLinks = [...allLinks.map(l => selectedIds.has(l.id) ? {...l, is_active: true} : l)];
			toast.success(`Activated ${selectedIds.size} links`);
			clearSelection();
		} catch (error: any) {
			toast.error(error.message || 'Failed to activate links');
		}
	}

	async function bulkDeactivate() {
		if (selectedIds.size === 0) return;
		
		try {
			await linksApi.bulkAction(Array.from(selectedIds), 'deactivate', $auth.token!);
			links = [...links.map(l => selectedIds.has(l.id) ? {...l, is_active: false} : l)];
			allLinks = [...allLinks.map(l => selectedIds.has(l.id) ? {...l, is_active: false} : l)];
			toast.success(`Deactivated ${selectedIds.size} links`);
			clearSelection();
		} catch (error: any) {
			toast.error(error.message || 'Failed to deactivate links');
		}
	}

	async function bulkChangeLayout(layoutType: 'classic' | 'featured') {
		if (selectedIds.size === 0) return;
		
		try {
			const updatePromises = Array.from(selectedIds).map(id => 
				linksApi.updateLink(id, { layout_type: layoutType }, $auth.token!)
			);
			await Promise.all(updatePromises);
			
			links = [...links.map(l => selectedIds.has(l.id) ? {...l, layout_type: layoutType} : l)];
			allLinks = allLinks.map(l => selectedIds.has(l.id) ? {...l, layout_type: layoutType} : l);
			clearSelection();
		} catch (error: any) {
			// Silent error handling
		}
	}



	// Group handlers
	async function handleCreateGroup(event: CustomEvent) {
		const { title, layout } = event.detail;
		try {
			const newGroup = await linksApi.createGroup(title, layout, $auth.token!);
			await loadData();
			showCreateGroupDialog = false;
			
			// Auto-expand the newly created group (close others)
			expandedGroupIds.clear();
			expandedGroupIds.add(newGroup.id);
			expandedGroupIds = expandedGroupIds; // Trigger reactivity
			
			toast.success('Group created!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to create group');
		}
	}

	function handleAddLinkToGroup(event: CustomEvent) {
		const groupId = event.detail;
		const group = links.find(l => l.id === groupId);
		if (group) {
			currentGroupId = groupId;
			currentGroupTitle = group.group_title || 'Group';
			editingGroupLink = null; // null = add mode
			isUploadingThumbnail = false;
			showGroupLinkDialog = true;
		}
	}

	async function handleRemoveFromGroup(event: CustomEvent) {
		const linkId = event.detail;
		try {
			await linksApi.removeFromGroup(linkId, $auth.token!);
			// Update local state - force reactivity with spread
			links = [...links.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.filter(child => child.id !== linkId)
					};
				}
				return link;
			})];
			allLinks = [...allLinks.filter(l => l.id !== linkId)];
			toast.success('Link removed!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to remove link');
		}
	}

	async function handleToggleNewTab(event: CustomEvent) {
		const linkId = event.detail;
		console.log('🎯 handleToggleNewTab called with linkId:', linkId);
		try {
			// Find current state
			let currentState = false;
			for (const link of links) {
				if (link.is_group && link.children) {
					const child = link.children.find(c => c.id === linkId);
					if (child) {
						currentState = child.open_in_new_tab || false;
						console.log('📍 Found link in group, current state:', currentState);
						break;
					}
				}
			}
			
			// Toggle
			const newState = !currentState;
			console.log('🔄 Toggling from', currentState, 'to', newState);
			await linksApi.updateLink(linkId, { open_in_new_tab: newState }, $auth.token!);
			
			// Update local state without reloading - keep UI expanded
			links = links.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.map(child => 
							child.id === linkId ? { ...child, open_in_new_tab: newState } : child
						)
					};
				}
				return link;
			});
			
			toast.success(newState ? 'Will open in new tab' : 'Will open in same tab');
		} catch (error: any) {
			toast.error(error.message || 'Failed to update setting');
		}
	}

	async function handleToggleLinkVisibility(event: CustomEvent) {
		const linkId = event.detail;
		console.log('🔄 Toggle visibility for linkId:', linkId);
		
		try {
			// Find current state
			let currentState = false;
			for (const link of links) {
				if (link.is_group && link.children) {
					const child = link.children.find(c => c.id === linkId);
					if (child) {
						currentState = child.is_active || false;
						console.log('📌 Found link in group:', { 
							linkId, 
							currentState, 
							groupId: link.id,
							groupTitle: link.group_title 
						});
						break;
					}
				}
			}
			
			// Toggle
			const newState = !currentState;
			console.log('🔀 Toggling from', currentState, 'to', newState);
			
			await linksApi.updateLink(linkId, { is_active: newState }, $auth.token!);
			console.log('✅ Backend updated successfully');
			
			// Update local state without reloading - keep UI expanded
			const updateLinks = (linksList: Link[]) => linksList.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.map(child => 
							child.id === linkId ? { ...child, is_active: newState } : child
						)
					};
				}
				return link;
			});
			
			const oldLinksLength = links.length;
			const oldAllLinksLength = allLinks.length;
			
			links = updateLinks(links);
			allLinks = updateLinks(allLinks);
			
			console.log('🔄 State updated:', { 
				linksLength: links.length,
				allLinksLength: allLinks.length,
				linksChanged: links !== links, // Will always be true due to reassignment
				newState 
			});
			
			// Verify the update
			for (const link of links) {
				if (link.is_group && link.children) {
					const child = link.children.find(c => c.id === linkId);
					if (child) {
						console.log('✅ Verified in links array:', { 
							linkId, 
							is_active: child.is_active,
							expected: newState,
							match: child.is_active === newState 
						});
					}
				}
			}
			
			// No toast for visibility toggle - visual feedback (icon change) is sufficient
		} catch (error: any) {
			console.error('❌ Toggle visibility error:', error);
			toast.error(error.message || 'Failed to toggle visibility');
		}
	}

	function handleEditGroupLink(event: CustomEvent) {
		editingGroupLink = event.detail;
		isUploadingThumbnail = false;
		showGroupLinkDialog = true;
	}

	async function handleSaveGroupLink(event: CustomEvent) {
		const { id, title, url, description, file } = event.detail;
		console.log('🔄 handleSaveGroupLink called', { id, hasFile: !!file });
		
		try {
			await linksApi.updateLink(id, { title, url, description }, $auth.token!);
			
			// Update local state immediately
			links = links.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.map(child => 
							child.id === id 
								? { ...child, title, url, description } 
								: child
						)
					};
				}
				return link;
			});
			allLinks = allLinks.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.map(child => 
							child.id === id 
								? { ...child, title, url, description } 
								: child
						)
					};
				}
				return link;
			});
			
			toast.success('Link updated!');
			
			// Upload thumbnail after save if file exists
			if (file) {
				console.log('📤 Starting thumbnail upload...');
				isUploadingThumbnail = true;
				console.log('✅ isUploadingThumbnail set to TRUE');
				
				await linksApi.uploadThumbnail(id, file, $auth.token!)
					.then(uploadedLink => {
						console.log('✅ Thumbnail uploaded successfully');
						// Update with actual uploaded thumbnail
						links = links.map(link => {
							if (link.is_group && link.children) {
								return {
									...link,
									children: link.children.map(child => 
										child.id === id 
											? { ...child, thumbnail_url: uploadedLink.thumbnail_url }
											: child
									)
								};
							}
							return link;
						});
						allLinks = allLinks.map(link => 
							link.id === id 
								? { ...link, thumbnail_url: uploadedLink.thumbnail_url }
								: link
						);
						toast.success('Thumbnail uploaded!');
					})
					.catch(error => {
						console.error('❌ Thumbnail upload failed:', error);
						toast.error('Failed to upload thumbnail');
					})
					.finally(() => {
						console.log('🏁 Upload complete, setting isUploadingThumbnail to FALSE');
						isUploadingThumbnail = false;
						showGroupLinkDialog = false; // Close modal after upload
						console.log('🚪 Modal closed');
					});
			}
		} catch (error: any) {
			console.error('❌ Update link error:', error);
			toast.error(error.message || 'Failed to update link');
		}
	}

	async function handleAddGroupLink(event: CustomEvent) {
		const { title, url, description, file } = event.detail;
		if (!currentGroupId) return;
		
		try {
			// Create link first
			const newLink = await linksApi.addToGroup(currentGroupId, { title, url, description }, $auth.token!);
			
			// If file selected, create preview URL for instant display
			let displayLink = newLink;
			if (file) {
				const previewUrl = URL.createObjectURL(file);
				displayLink = { ...newLink, thumbnail_url: previewUrl };
			}
			
			// Update UI immediately with preview (instant feedback)
			links = links.map(link => {
				if (link.id === currentGroupId && link.is_group) {
					return {
						...link,
						children: [...(link.children || []), displayLink]
					};
				}
				return link;
			});
			allLinks = [...allLinks, displayLink];
			
			toast.success('Link added to group!');
			
			// Upload thumbnail in background if file exists
			if (file) {
				console.log('📤 Starting thumbnail upload for new link...');
				isUploadingThumbnail = true;
				console.log('✅ isUploadingThumbnail set to TRUE');
				
				linksApi.uploadThumbnail(newLink.id, file, $auth.token!)
					.then(uploadedLink => {
						console.log('✅ Thumbnail uploaded successfully for new link');
						// Replace temporary preview with actual uploaded link (backend returns full object)
						links = links.map(link => {
							if (link.is_group && link.children) {
								return {
									...link,
									children: link.children.map(child => 
										child.id === newLink.id 
											? uploadedLink
											: child
									)
								};
							}
							return link;
						});
						allLinks = allLinks.map(link => 
							link.id === newLink.id 
								? uploadedLink
								: link
						);
						toast.success('Thumbnail uploaded!');
					})
					.catch(error => {
						console.error('❌ Thumbnail upload failed for new link:', error);
						toast.error('Failed to upload thumbnail');
						// Remove preview URL on error, keep link without thumbnail
						links = links.map(link => {
							if (link.is_group && link.children) {
								return {
									...link,
									children: link.children.map(child => 
										child.id === newLink.id ? { ...child, thumbnail_url: null } : child
									)
								};
							}
							return link;
						});
					})
					.finally(() => {
						console.log('🏁 Upload complete for new link, closing modal');
						isUploadingThumbnail = false;
						showGroupLinkDialog = false; // Close modal after upload
					});
			}
		} catch (error: any) {
			toast.error(error.message || 'Failed to add link');
		}
	}

	async function handleRemoveGroupLinkThumbnail(event: CustomEvent) {
		const linkId = event.detail;
		try {
			await linksApi.deleteThumbnail(linkId, $auth.token!);
			
			// Update editingGroupLink to trigger dialog update
			if (editingGroupLink && editingGroupLink.id === linkId) {
				editingGroupLink = { ...editingGroupLink, thumbnail_url: null };
			}
			
			// Update local state immediately
			links = links.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.map(child => 
							child.id === linkId 
								? { ...child, thumbnail_url: null } 
								: child
						)
					};
				}
				return link;
			});
			allLinks = allLinks.map(link => {
				if (link.is_group && link.children) {
					return {
						...link,
						children: link.children.map(child => 
							child.id === linkId 
								? { ...child, thumbnail_url: null } 
								: child
						)
					};
				}
				return link;
			});
			
			toast.success('Thumbnail removed!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to remove thumbnail');
		}
	}

	async function handlePinGroupLink(event: CustomEvent) {
		const linkId = event.detail;
		try {
			const updatedLink = await linksApi.togglePin(linkId, $auth.token!);
			
			// Update local state - only update pin status, keep position unchanged
			const updateLinks = (linksList: Link[]) => linksList.map(link => {
				if (link.is_group && link.children) {
					// Update children based on backend response
					const updatedChildren = link.children.map(child => {
						if (child.id === linkId) {
							// Update this link with new pin status
							return { ...child, is_pinned: updatedLink.is_pinned };
						} else if (updatedLink.is_pinned && child.is_pinned) {
							// Backend unpins others in the same group, sync frontend
							return { ...child, is_pinned: false };
						}
						return child;
					});
					
					// Keep original order - ProfilePreview will sort by is_pinned when rendering
					return { ...link, children: updatedChildren };
				}
				return link;
			});
			
			links = updateLinks(links);
			allLinks = updateLinks(allLinks);
			
			// No toast for pin/unpin - visual feedback (icon color change) is sufficient
		} catch (error: any) {
			toast.error(error.message || 'Failed to toggle pin');
		}
	}

	function handleDuplicateGroupLink(event: CustomEvent) {
		const { linkId, groupId } = event.detail;
		
		// Find the link - it's inside group's children, not in allLinks
		let foundLink: Link | null = null;
		for (const link of links) {
			if (link.is_group && link.children) {
				const childLink = link.children.find(c => c.id === linkId);
				if (childLink) {
					foundLink = childLink;
					break;
				}
			}
		}
		
		if (!foundLink) return;
		
		duplicatingLink = foundLink;
		duplicatingFromGroupId = groupId;
		showDuplicateLinkDialog = true;
	}
	
	async function handleDuplicateLinkConfirm(event: CustomEvent) {
		const { linkId, targetGroupId } = event.detail;
		
		try {
			// Step 1: Duplicate the link (will be in same group)
			const newLink = await linksApi.duplicateLink(linkId, $auth.token!);
			
			// Step 2: If target is different group, move it
			if (targetGroupId !== duplicatingFromGroupId) {
				const movedLink = await linksApi.moveToGroup(newLink.id, targetGroupId, $auth.token!);
				
				// Update UI: Add to target group
				links = links.map(link => {
					if (link.id === targetGroupId && link.is_group) {
						return {
							...link,
							children: [...(link.children || []), movedLink]
						};
					}
					return link;
				});
				allLinks = [...allLinks, movedLink];
				
				toast.success(`Link duplicated to another group!`);
			} else {
				// Same group: Add next to original
				links = links.map(link => {
					if (link.is_group && link.children) {
						const originalIndex = link.children.findIndex(c => c.id === linkId);
						if (originalIndex !== -1) {
							const newChildren = [...link.children];
							newChildren.splice(originalIndex + 1, 0, newLink);
							return { ...link, children: newChildren };
						}
					}
					return link;
				});
				allLinks = [...allLinks, newLink];
				
				toast.success('Link duplicated in same group!');
			}
		} catch (error: any) {
			toast.error(error.message || 'Failed to duplicate link');
		}
	}

	async function handleDuplicateGroup(event: CustomEvent) {
		const groupId = event.detail;
		try {
			const newGroup = await linksApi.duplicateGroup(groupId, $auth.token!);
			
			// Optimistic update - insert new group after original
			const originalIndex = links.findIndex(l => l.id === groupId);
			if (originalIndex !== -1) {
				const newLinks = [...links];
				newLinks.splice(originalIndex + 1, 0, newGroup);
				links = newLinks;
				allLinks = [...allLinks, newGroup];
			} else {
				links = [...links, newGroup];
				allLinks = [...allLinks, newGroup];
			}
			
			toast.success('Group duplicated!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to duplicate group');
		}
	}
	
	async function handleDuplicateTextGroup(event: CustomEvent) {
		const groupId = event.detail;
		try {
			const newGroup = await blocksApi.duplicateGroup(groupId, $auth.token!);
			
			// Optimistic update - insert new group after original
			const originalIndex = Array.isArray(blocks) ? blocks.findIndex(b => b.id === groupId) : -1;
			if (originalIndex !== -1) {
				const newBlocks = [...(Array.isArray(blocks) ? blocks : [])];
				newBlocks.splice(originalIndex + 1, 0, newGroup);
				blocks = newBlocks;
			} else {
				blocks = Array.isArray(blocks) ? [...blocks, newGroup] : [newGroup];
			}
			
			toast.success('Text group duplicated!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to duplicate text group');
		}
	}

	async function handleReorderGroupLinks(event: CustomEvent) {
		const { groupId, linkIds } = event.detail;
		try {
			await linksApi.reorderGroupLinks(groupId, linkIds, $auth.token!);
			
			// Update local state without reloading - keep UI expanded
			links = links.map(link => {
				if (link.id === groupId && link.children) {
					// Reorder children based on linkIds AND update position
					const orderedChildren = linkIds.map((id: string, index: number) => {
						const child = link.children!.find(c => c.id === id);
						if (child) {
							return { ...child, position: index };
						}
						return null;
					}).filter(Boolean) as Link[];
					return { ...link, children: orderedChildren };
				}
				return link;
			});
		} catch (error: any) {
			toast.error(error.message || 'Failed to reorder links');
		}
	}

	async function handleToggleGroup(event: CustomEvent) {
		const id = event.detail;
		console.log('👁️ handleToggleGroup called!', { id });
		const group = links.find(l => l.id === id);
		if (!group) {
			console.log('❌ Group not found!', { id });
			return;
		}
		
		// Find the block that contains this link group
		const block = Array.isArray(blocks) ? blocks.find(b => b.link_id === id) : undefined;
		
		console.log('🔄 Toggling group visibility:', { 
			groupId: id, 
			groupTitle: group.group_title,
			currentState: group.is_active, 
			newState: !group.is_active,
			hasBlock: !!block,
			blockId: block?.id
		});
		
		try {
			// Toggle the link group
			await linksApi.updateLink(id, { is_active: !group.is_active }, $auth.token!);
			
			links = links.map(l => l.id === id ? {...l, is_active: !l.is_active} : l);
			allLinks = allLinks.map(l => l.id === id ? {...l, is_active: !l.is_active} : l);
			
			// Also toggle the block if it exists
			if (block) {
				await blocksApi.updateBlock(block.id, { is_active: !group.is_active }, $auth.token!);
				blocks = Array.isArray(blocks) ? blocks.map(b => b.id === block.id ? {...b, is_active: !group.is_active} : b) : [];
			}
			
			console.log('✅ Group visibility toggled!', { 
				updatedLinks: links.length, 
				updatedAllLinks: allLinks.length,
				blockToggled: !!block
			});
		} catch (error: any) {
			console.error('❌ Failed to toggle group:', error);
			toast.error(error.message || 'Failed to toggle group');
		}
	}
	
	async function handleUpdateGroupTitle(event: CustomEvent) {
		const { groupId, title } = event.detail;
		
		// Optimistic update - update UI immediately
		const oldLinks = links;
		const oldAllLinks = allLinks;
		links = links.map(l => l.id === groupId ? {...l, group_title: title} : l);
		allLinks = allLinks.map(l => l.id === groupId ? {...l, group_title: title} : l);
		
		try {
			await linksApi.updateLink(groupId, { group_title: title }, $auth.token!);
			toast.success('Group title updated!');
		} catch (error: any) {
			// Revert on error
			links = oldLinks;
			allLinks = oldAllLinks;
			toast.error(error.message || 'Failed to update title');
		}
	}
	
	async function handleUpdateTextGroupTitle(event: CustomEvent) {
		const { groupId, title } = event.detail;
		
		// Optimistic update - update UI immediately
		const oldBlocks = blocks;
		blocks = Array.isArray(blocks) ? blocks.map(b => b.id === groupId ? {...b, group_title: title} : b) : [];
		
		try {
			await blocksApi.updateBlock(groupId, { group_title: title }, $auth.token!);
			toast.success('Text group title updated!');
		} catch (error: any) {
			// Revert on error
			blocks = oldBlocks;
			toast.error(error.message || 'Failed to update title');
		}
	}

	async function handleUpdateGroupLayout(event: CustomEvent) {
		const { groupId, ...settings } = event.detail;
		
		// Optimistic update - update local state immediately
		links = links.map(l => l.id === groupId ? { ...l, ...settings } : l);
		allLinks = allLinks.map(l => l.id === groupId ? { ...l, ...settings } : l);
		
		try {
			await linksApi.updateLink(groupId, settings, $auth.token!);
			// No need to update again - optimistic update is sufficient
		} catch (error: any) {
			// Reload to revert optimistic update on error
			await loadData();
		}
	}

	function handleExpandGroup(event: CustomEvent) {
		const groupId = event.detail;
		// Close ALL groups (both link and text groups)
		expandedGroupIds.clear();
		expandedTextGroupIds.clear();
		expandedGroupIds.add(groupId);
		expandedGroupIds = expandedGroupIds;
		expandedTextGroupIds = expandedTextGroupIds;
	}

	function handleCollapseGroup(event: CustomEvent) {
		const groupId = event.detail;
		expandedGroupIds.delete(groupId);
		expandedGroupIds = expandedGroupIds;
	}

	async function handleDeleteGroup(event: CustomEvent) {
		const id = event.detail;
		const group = links.find(l => l.id === id);
		if (!group) return;
		
		// Show confirmation modal
		deleteTarget = { type: 'link-group', id, title: group.group_title || 'Group' };
		showDeleteConfirmModal = true;
	}
	
	async function confirmDeleteGroup(id: string) {
		const group = links.find(l => l.id === id);
		if (!group) return;
		
		try {
			await linksApi.deleteLink(id, $auth.token!);
			
			// Remove group and update state - force reactivity
			const childIds = new Set(group.children?.map(c => c.id) || []);
			links = [...links.filter(l => l.id !== id)];
			allLinks = [...allLinks.filter(l => l.id !== id && !childIds.has(l.id))];
			
			toast.success('Group deleted!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete group');
		}
	}

	// Block handlers
	async function handleAddBlock(event: CustomEvent) {
		const { type } = event.detail;
		
		// If group, show create group dialog
		if (type === 'group') {
			showCreateGroupDialog = true;
			return;
		}
		
		// If text block, create text group and expand it
		if (type === 'text') {
			try {
				const newTextGroup = await blocksApi.createBlock({
					block_type: 'text',
					is_group: true,
					group_title: 'New Text Group',
					is_active: true
				}, $auth.token!);
				
				blocks = Array.isArray(blocks) ? [...blocks, newTextGroup] : [newTextGroup];
				
				// Auto-expand the new group (close others)
				expandedTextGroupIds.clear();
				expandedTextGroupIds.add(newTextGroup.id);
				expandedTextGroupIds = expandedTextGroupIds;
				
				toast.success('Text group created!');
			} catch (error: any) {
				toast.error(error.message || 'Failed to create text group');
			}
			return;
		}
		
		try {
			let blockData: any = {
				block_type: type,
				is_active: true
			};

			// Set default values based on type - only include relevant fields
			switch (type) {
				case 'image':
					blockData.image_url = '';
					blockData.alt_text = '';
					break;
				case 'video':
					blockData.video_url = '';
					break;
				case 'social':
					blockData.social_links = [{ platform: 'facebook', url: '' }];
					break;
				case 'divider':
					blockData.divider_style = 'line';
					break;
				case 'email':
					blockData.content = 'Subscribe to my newsletter';
					blockData.placeholder = 'Enter your email';
					break;
				case 'embed':
					blockData.embed_url = '';
					blockData.embed_type = 'spotify';
					break;
			}

			const newBlock = await blocksApi.createBlock(blockData, $auth.token!);
			blocks = Array.isArray(blocks) ? [...blocks, newBlock] : [newBlock];
			toast.success('Block added!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to add block');
		}
	}

	// Text group handlers
	function handleExpandTextGroup(event: CustomEvent) {
		const groupId = event.detail;
		// Close ALL groups (both link and text groups)
		expandedGroupIds.clear();
		expandedTextGroupIds.clear();
		expandedTextGroupIds.add(groupId);
		expandedGroupIds = expandedGroupIds;
		expandedTextGroupIds = expandedTextGroupIds;
	}

	function handleCollapseTextGroup(event: CustomEvent) {
		const groupId = event.detail;
		expandedTextGroupIds.delete(groupId);
		expandedTextGroupIds = expandedTextGroupIds;
	}

	async function handleAddTextToGroup(event: CustomEvent) {
		const { groupId, content } = event.detail;
		
		try {
			const newTextBlock = await blocksApi.createBlock({
				block_type: 'text',
				parent_id: groupId,
				is_group: false,
				is_active: true,
				content
			}, $auth.token!);
			
			// Update the group's children
			blocks = Array.isArray(blocks) ? blocks.map(b => {
				if (b.id === groupId) {
					return {
						...b,
						children: [...(b.children || []), newTextBlock]
					};
				}
				return b;
			}) : [];
			
			toast.success('Text added to group!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to add text');
		}
	}

	async function handleUpdateGroupStyle(event: CustomEvent) {
		const { groupId, style } = event.detail;
		
		try {
			const styleString = JSON.stringify(style);
			await blocksApi.updateBlock(groupId, { style: styleString }, $auth.token!);
			
			// Update local state - force reactivity by creating new array
			blocks = blocks.map(b => 
				b.id === groupId ? { ...b, style: styleString } : b
			);
		} catch (error: any) {
			console.error('Failed to update style:', error);
			toast.error(error.message || 'Failed to update style');
		}
	}

	async function handleDeleteTextGroup(event: CustomEvent) {
		const id = event.detail;
		const group = Array.isArray(blocks) ? blocks.find(b => b.id === id) : null;
		if (!group) return;
		
		// Show confirmation modal
		deleteTarget = { type: 'text-group', id, title: group.group_title || 'Text Group' };
		showDeleteConfirmModal = true;
	}
	
	async function confirmDeleteTextGroup(id: string) {
		try {
			await blocksApi.deleteBlock(id, $auth.token!);
			blocks = Array.isArray(blocks) ? blocks.filter(b => b.id !== id) : [];
			expandedTextGroupIds.delete(id);
			expandedTextGroupIds = expandedTextGroupIds;
			toast.success('Text group deleted!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete text group');
		}
	}

	async function handleDeleteTextFromGroup(event: CustomEvent) {
		const { groupId, textId } = event.detail;
		try {
			await blocksApi.deleteBlock(textId, $auth.token!);
			
			// Update group's children
			blocks = Array.isArray(blocks) ? blocks.map(b => {
				if (b.id === groupId && b.children) {
					return {
						...b,
						children: b.children.filter(c => c.id !== textId)
					};
				}
				return b;
			}) : [];
			
			toast.success('Text deleted!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete text');
		}
	}

	async function handleReorderTextsInGroup(event: CustomEvent) {
		const { groupId, blockIds } = event.detail;
		try {
			await blocksApi.reorderGroupBlocks(groupId, blockIds, $auth.token!);
			
			// Update local state
			blocks = blocks.map(b => {
				if (b.id === groupId && b.children) {
					const orderedChildren = blockIds.map((id: string, index: number) => {
						const child = b.children!.find(c => c.id === id);
						if (child) {
							return { ...child, position: index };
						}
						return null;
					}).filter(Boolean) as Block[];
					return { ...b, children: orderedChildren };
				}
				return b;
			});
		} catch (error: any) {
			toast.error(error.message || 'Failed to reorder texts');
		}
	}

	async function handleEditTextInGroup(event: CustomEvent) {
		const { groupId, textId, content } = event.detail;
		try {
			// Update the text block content
			await blocksApi.updateBlock(textId, { content }, $auth.token!);
			
			// Update local state - find and update the text block in the group's children
			blocks = blocks.map(block => {
				if (block.id === groupId && block.children) {
					return {
						...block,
						children: block.children.map(child =>
							child.id === textId ? { ...child, content } : child
						)
					};
				}
				return block;
			});
			
			toast.success('Text updated!');
		} catch (error: any) {
			console.error('handleEditTextInGroup error:', error);
			toast.error(error.message || 'Failed to update text');
		}
	}
	
	async function handleToggleTextGroup(event: CustomEvent) {
		const groupId = event.detail;
		try {
			const group = blocks.find(b => b.id === groupId);
			if (!group) return;
			
			const newState = !group.is_active;
			await blocksApi.updateBlock(groupId, { is_active: newState }, $auth.token!);
			
			blocks = blocks.map(b => 
				b.id === groupId ? { ...b, is_active: newState } : b
			);
			
			toast.success(newState ? 'Text group shown!' : 'Text group hidden!');
		} catch (error: any) {
			console.error('handleToggleTextGroup error:', error);
			toast.error(error.message || 'Failed to toggle text group');
		}
	}
	
	async function handleToggleTextInGroup(event: CustomEvent) {
		const { groupId, textId } = event.detail;
		try {
			// Find the text block
			const group = blocks.find(b => b.id === groupId);
			if (!group || !group.children) return;
			
			const textBlock = group.children.find(c => c.id === textId);
			if (!textBlock) return;
			
			const newState = !textBlock.is_active;
			await blocksApi.updateBlock(textId, { is_active: newState }, $auth.token!);
			
			// Update local state
			blocks = blocks.map(block => {
				if (block.id === groupId && block.children) {
					return {
						...block,
						children: block.children.map(child =>
							child.id === textId ? { ...child, is_active: newState } : child
						)
					};
				}
				return block;
			});
			
			toast.success(newState ? 'Text shown!' : 'Text hidden!');
		} catch (error: any) {
			console.error('handleToggleTextInGroup error:', error);
			toast.error(error.message || 'Failed to toggle text');
		}
	}

	async function handleUpdateBlock(event: CustomEvent) {
		const { id, ...data } = event.detail;
		try {
			await blocksApi.updateBlock(id, data, $auth.token!);
			blocks = Array.isArray(blocks) ? blocks.map(b => b.id === id ? { ...b, ...data } : b) : [];
			toast.success('Block updated!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to update block');
		}
	}

	async function handleDeleteBlock(event: CustomEvent) {
		const id = event.detail;
		const block = Array.isArray(blocks) ? blocks.find(b => b.id === id) : null;
		if (!block) return;
		
		// Show confirmation modal
		deleteTarget = { 
			type: 'block', 
			id, 
			title: block.content || block.block_type || 'Block'
		};
		showDeleteConfirmModal = true;
	}
	
	async function confirmDeleteBlock(id: string) {
		try {
			await blocksApi.deleteBlock(id, $auth.token!);
			blocks = Array.isArray(blocks) ? blocks.filter(b => b.id !== id) : [];
			toast.success('Block deleted!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete block');
		}
	}

	function handleSelectBlock(event: CustomEvent) {
		const id = event.detail;
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = selectedIds;
	}

	async function handleToggleBlock(event: CustomEvent) {
		const id = event.detail;
		const block = Array.isArray(blocks) ? blocks.find(b => b.id === id) : undefined;
		if (!block) return;
		
		try {
			await blocksApi.updateBlock(id, { is_active: !block.is_active }, $auth.token!);
			blocks = Array.isArray(blocks) ? blocks.map(b => b.id === id ? { ...b, is_active: !b.is_active } : b) : [];
		} catch (error: any) {
			toast.error('Failed to toggle block');
		}
	}

	async function handleDuplicateBlock(event: CustomEvent) {
		const id = event.detail;
		const block = Array.isArray(blocks) ? blocks.find(b => b.id === id) : undefined;
		if (!block) return;
		
		try {
			const { id: _, created_at, updated_at, position, ...blockData } = block;
			const newBlock = await blocksApi.createBlock({ ...blockData, is_active: true }, $auth.token!);
			
			// Update local state: insert new block after original
			const targetPos = block.position + 1;
			
			// Update positions for items after insertion point
			links = links.map(l => 
				l.position >= targetPos ? { ...l, position: l.position + 1 } : l
			);
			blocks = Array.isArray(blocks) ? blocks.map(b => 
				b.position >= targetPos ? { ...b, position: b.position + 1 } : b
			) : [];
			
			// Add new block with correct position
			blocks = Array.isArray(blocks) ? [...blocks, { ...newBlock, position: targetPos }] : [{ ...newBlock, position: targetPos }];
			
			// Reorder on backend
			const newItems = [
				...items.filter(item => item.position < targetPos).map(item => ({ id: item.id, type: item.type })),
				{ id: newBlock.id, type: 'block' },
				...items.filter(item => item.position >= targetPos).map(item => ({ id: item.id, type: item.type }))
			];
			
			await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:3000/api'}/items/reorder`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${$auth.token}`
				},
				body: JSON.stringify({ items: newItems })
			});
			
			toast.success('Block duplicated!');
		} catch (error: any) {
			toast.error('Failed to duplicate block');
			await loadData();
		}
	}

	// Drag & Drop handlers
	function handleDndConsider(e: CustomEvent) {
		items = e.detail.items;
	}

	async function handleDndFinalize(e: CustomEvent) {
		items = e.detail.items;
		
		try {
			// Send unified reorder request with all items
			const itemsData = items.map(item => ({
				id: item.id,
				type: item.type
			}));
			
			const res = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:3000/api'}/items/reorder`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${$auth.token}`
				},
				body: JSON.stringify({ items: itemsData })
			});
			
			if (!res.ok) throw new Error('Failed to reorder');
			
			// Update position in original arrays to match new order
			items.forEach((item, idx) => {
				item.data.position = idx;
				if (item.type === 'link') {
					const linkIndex = links.findIndex(l => l.id === item.id);
					if (linkIndex !== -1) links[linkIndex].position = idx;
				} else if (Array.isArray(blocks)) {
					const blockIndex = blocks.findIndex(b => b.id === item.id);
					if (blockIndex !== -1) blocks[blockIndex].position = idx;
				}
			});
			
			// Trigger reactivity
			links = links;
			blocks = blocks || [];
		} catch (error: any) {
			toast.error('Failed to reorder items');
			console.error(error);
			// Only reload on error to restore correct state
			await loadData();
		}
	}

	$: displayProfile = profile || { 
		username: 'loading...', 
		bio: '', 
		avatar_url: undefined,
		id: '',
		user_id: '',
		theme_config: {},
		custom_css: undefined,
		created_at: new Date(),
		updated_at: new Date()
	};
</script>

<svelte:head>
	<title>My Bio - LinkBio</title>
</svelte:head>

<style>
	@keyframes fade-in {
		from { opacity: 0; transform: translateY(10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	
	.animate-in {
		animation: fade-in 0.3s ease-out;
	}

	/* Remove all borders and outlines when dragging */
	:global([aria-grabbed="true"]),
	:global([aria-grabbed="true"] *),
	:global(section div:focus),
	:global(section div:focus-visible),
	:global(section div[tabindex]:focus),
	:global(section div[tabindex]:focus-visible) {
		outline: none !important;
		border-color: transparent !important;
		box-shadow: none !important;
	}

	/* Ensure no focus styles on wrapper divs */
	:global(section > div) {
		outline: none !important;
	}

	/* Tắt hoàn toàn mọi transition/animation trong drag zone */
	:global(section.space-y-3),
	:global(section.space-y-3 *) {
		transition: none !important;
		animation: none !important;
	}

	/* Chỉ giữ opacity cho item đang drag */
	:global([aria-grabbed="true"]) {
		opacity: 0.5;
		cursor: grabbing !important;
	}
</style>

<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 h-16 sticky top-0 z-10">
		<div class="flex items-center justify-between h-full">
			<div>
				<h1 class="text-xl font-bold text-gray-900">My Bio</h1>
				<p class="text-xs text-gray-500">Build and customize your bio page</p>
			</div>
			<div class="flex items-center gap-3">
				<div class="flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm min-w-[280px]">
					<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
					</svg>
					<span class="font-medium text-gray-900 flex-1">
						linkbio.com/{profile?.username || 'loading...'}
					</span>
				</div>

				<Button
					variant="outline"
					size="icon"
					onclick={handleCopyUrl}
					class="h-10 w-10 border-gray-200 hover:bg-gray-50"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
					</svg>
				</Button>

				<Button
					variant="outline"
					size="icon"
					onclick={() => {
						if (profile?.username) {
							window.open(`${window.location.origin}/${profile.username}`, '_blank');
						}
					}}
					class="h-10 w-10 border-gray-200 hover:bg-gray-50"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
					</svg>
				</Button>

				<Button
					onclick={() => {
						if (profile?.username) {
							const url = `${window.location.origin}/${profile.username}`;
							if (navigator.share) {
								navigator.share({ title: 'My LinkBio', url }).catch(() => {});
							} else {
								handleCopyUrl();
							}
						}
					}}
					class="h-10 px-6 bg-gray-900 hover:bg-gray-800 text-white font-semibold"
				>
					SHARE
				</Button>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8 relative min-h-[calc(100vh-4rem)]">
		{#if initialLoading}
			<div class="absolute inset-0 flex items-center justify-center bg-gray-50">
				<div class="text-center">
					<div class="relative w-14 h-14 mx-auto mb-4">
						<div class="absolute inset-0 border-4 border-indigo-100 rounded-full"></div>
						<div class="absolute inset-0 border-4 border-transparent border-t-indigo-600 border-r-blue-600 rounded-full animate-spin"></div>
					</div>
					<p class="text-sm text-gray-500">Loading...</p>
				</div>
			</div>
		{:else}
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8 animate-in">
			<!-- Links Section (2/3) -->
			<div class="lg:col-span-2 space-y-6">
				<!-- Action Buttons (Linktree style) -->
				<div class="grid grid-cols-2 gap-3">
					<!-- Add Link Button - Gradient Red to Purple -->
					<button
						onclick={() => showCreateGroupDialog = true}
						class="flex items-center justify-center gap-2 px-6 py-4 rounded-xl text-white text-sm font-semibold transition-all shadow-md hover:shadow-lg transform hover:-translate-y-0.5"
						style="background: linear-gradient(135deg, #FF5757 0%, #FF6B9D 50%, #C471ED 100%);"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
						</svg>
						<span>ADD LINK</span>
					</button>

					<!-- Add Embed Button - Blue -->
					<button
						onclick={() => showAddBlockDialog = true}
						class="flex items-center justify-center gap-2 px-6 py-4 bg-[#0096FF] hover:bg-[#0086E6] rounded-xl text-white text-sm font-semibold transition-all shadow-md hover:shadow-lg transform hover:-translate-y-0.5"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
						</svg>
						<span>ADD EMBED</span>
					</button>
				</div>

				<!-- Add Header Link -->
				<button
					onclick={() => toast.info('Add header feature coming soon!')}
					class="flex items-center gap-2 text-gray-500 hover:text-gray-700 text-sm font-medium transition-colors"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					<span>Add header</span>
				</button>

				<!-- Combined Links & Blocks List -->
				<section 
					class="space-y-3"
					use:dndzone={{items, flipDurationMs: 200, dropTargetStyle: {}}}
					onconsider={handleDndConsider}
					onfinalize={handleDndFinalize}
				>
					{#each items as item (item.id)}
						<div style="outline: none;">
							{#if item.type === 'link' && item.data.is_group && !hiddenLinkIds.has(item.data.id)}
								<LinkGroupCard 
									group={item.data}
									selected={selectedIds.has(item.id)}
									expanded={expandedGroupIds.has(item.data.id)}
									on:select={handleSelect}
									on:expand={handleExpandGroup}
									on:collapse={handleCollapseGroup}
									on:toggle={handleToggleGroup}
									on:delete={handleDeleteGroup}
									on:addlink={handleAddLinkToGroup}
									on:removefromgroup={handleRemoveFromGroup}
									on:editlink={handleEditGroupLink}
									on:pinlink={handlePinGroupLink}
									on:duplicatelink={handleDuplicateGroupLink}
									on:duplicate={handleDuplicateGroup}
									on:toggleNewTab={handleToggleNewTab}
									on:reorderlinks={handleReorderGroupLinks}
									on:togglelinkvisibility={handleToggleLinkVisibility}
									on:updatelayout={handleUpdateGroupLayout}
									on:updatetitle={handleUpdateGroupTitle}
								/>
							{:else if item.type === 'block' && item.data.is_group && item.data.block_type === 'text'}
								<TextGroupCard 
									group={item.data}
									expanded={expandedTextGroupIds.has(item.data.id)}
									on:expand={handleExpandTextGroup}
									on:collapse={handleCollapseTextGroup}
									on:delete={handleDeleteTextGroup}
									on:duplicate={handleDuplicateTextGroup}
									on:addtext={handleAddTextToGroup}
									on:updatestyle={handleUpdateGroupStyle}
									on:edittext={handleEditTextInGroup}
									on:deletetext={handleDeleteTextFromGroup}
									on:reordertexts={handleReorderTextsInGroup}
									on:togglevisibility={handleToggleTextGroup}
									on:toggletextvisibility={handleToggleTextInGroup}
									on:updatetitle={handleUpdateTextGroupTitle}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'text' && !item.data.is_group}
								<TextBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
									on:toggle={handleToggleBlock}
									on:duplicate={handleDuplicateBlock}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'image'}
								<ImageBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:update={handleUpdateBlock}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'video'}
								<VideoBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:update={handleUpdateBlock}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'social'}
								<SocialBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:update={handleUpdateBlock}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'divider'}
								<DividerBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'email'}
								<EmailCollectorBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:update={handleUpdateBlock}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
								/>
							{:else if item.type === 'block' && item.data.block_type === 'embed'}
								<EmbedBlock 
									block={item.data}
									selected={selectedIds.has(item.id)}
									on:update={handleUpdateBlock}
									on:delete={handleDeleteBlock}
									on:select={handleSelectBlock}
								/>
							{/if}
						</div>
					{/each}
				</section>

				{#if items.length === 0}
					<!-- No Content at All -->
					<div class="relative bg-gradient-to-br from-indigo-50 via-white to-blue-50 rounded-2xl p-16 text-center border-2 border-dashed border-indigo-200 overflow-hidden">
						<div class="absolute top-0 right-0 w-32 h-32 bg-indigo-100 rounded-full -mr-16 -mt-16 opacity-50"></div>
						<div class="absolute bottom-0 left-0 w-24 h-24 bg-blue-100 rounded-full -ml-12 -mb-12 opacity-50"></div>
						
						<div class="relative">
							<div class="w-20 h-20 mx-auto mb-6 bg-gradient-to-br from-indigo-500 to-blue-700 rounded-2xl flex items-center justify-center shadow-xl shadow-indigo-500/30">
								<svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
								</svg>
							</div>
							<h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent mb-2">
								No content yet
							</h3>
							<p class="text-gray-600 mb-6 max-w-sm mx-auto">Get started by adding your first block to build your bio page</p>
							<Button 
								onclick={() => showAddBlockDialog = true}
								class="bg-gradient-to-r from-indigo-700 to-blue-700 hover:from-indigo-800 hover:to-blue-800 shadow-lg shadow-indigo-500/30"
							>
								<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
								</svg>
								Add Your First Block
							</Button>
						</div>
					</div>
				{/if}
			</div>

			<!-- Preview Section (1/3) -->
			<div class="lg:col-span-1">
				<div class="sticky top-32 space-y-6">
					<div class="text-center mb-6">
						<h3 class="text-lg font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
							Live Preview
						</h3>
						<p class="text-sm text-gray-500 mt-1">See how your profile looks</p>
					</div>
					<div class="relative">
						<div class="absolute inset-0 bg-gradient-to-r from-indigo-500 to-blue-500 rounded-3xl blur-2xl opacity-20"></div>
						<div class="relative">
							<ProfilePreview profile={displayProfile} {links} {blocks} showInactive={false} />
						</div>
					</div>

					<!-- Branding Toggle -->
					<div class="flex items-center justify-center gap-3">
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" class="sr-only peer" disabled>
							<div class="w-11 h-6 bg-gray-200 rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all"></div>
						</label>
						<span class="text-sm font-medium text-gray-700">Hide Beacons logo and footer</span>
						<button class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-500 hover:from-pink-600 hover:to-rose-600 text-white font-bold rounded-full text-sm transition-all shadow-lg shadow-pink-500/30">
							<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
								<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
							</svg>
							Pro
						</button>
					</div>
				</div>
			</div>
		</div>
		{/if}
	</div>
</div>

<!-- Floating Bulk Action Bar -->
{#if selectedCount > 0}
	<div class="fixed bottom-8 left-1/2 -translate-x-1/2 z-50 animate-in">
		<div class="bg-gray-900 text-white rounded-2xl shadow-2xl px-6 py-4 flex items-center gap-6">
			<div class="flex items-center gap-2">
				<span class="font-semibold">{selectedCount}</span>
				<span class="text-gray-300">selected</span>
			</div>
			
			<div class="w-px h-6 bg-gray-700"></div>
			
			<div class="flex items-center gap-2">
				<button
					onclick={bulkActivate}
					class="px-4 py-2 bg-green-600 hover:bg-green-700 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
					</svg>
					Show
				</button>
				
				<button
					onclick={bulkDeactivate}
					class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
					</svg>
					Hide
				</button>
				
				<button
					onclick={() => bulkChangeLayout('classic')}
					class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
					</svg>
					Classic
				</button>
				
				<button
					onclick={() => bulkChangeLayout('featured')}
					class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
					</svg>
					Featured
				</button>
				
				<button
					onclick={bulkDelete}
					class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
					</svg>
					Delete
				</button>
			</div>
			
			<div class="w-px h-6 bg-gray-700"></div>
			
			<button
				onclick={clearSelection}
				class="p-2 hover:bg-gray-800 rounded-lg transition-colors"
				title="Clear selection"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
	</div>
{/if}

<!-- Add Block Dialog -->
<AddBlockDialog bind:open={showAddBlockDialog} on:select={handleAddBlock} />

<!-- Calendar View -->
{#if showCalendarView}
	<CalendarView links={allLinks} on:close={() => showCalendarView = false} />
{/if}

<!-- Create Group Dialog -->
<CreateGroupDialog 
	bind:open={showCreateGroupDialog}
	on:create={handleCreateGroup}
/>

<!-- Group Link Dialog (Add/Edit) -->
<EditGroupLinkDialog 
	bind:open={showGroupLinkDialog}
	link={editingGroupLink}
	isUploading={isUploadingThumbnail}
	on:save={handleSaveGroupLink}
	on:add={handleAddGroupLink}
	on:removeThumbnail={handleRemoveGroupLinkThumbnail}
/>

<!-- Duplicate Link Dialog -->
<DuplicateLinkDialog 
	bind:open={showDuplicateLinkDialog}
	link={duplicatingLink}
	groups={links.filter(l => l.is_group)}
	currentGroupId={duplicatingFromGroupId}
	on:duplicate={handleDuplicateLinkConfirm}
/>

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirmModal && deleteTarget}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div 
		class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
		onclick={() => {
			showDeleteConfirmModal = false;
			deleteTarget = null;
		}}
	>
		<!-- svelte-ignore a11y-click-events-have-key-events -->
		<!-- svelte-ignore a11y-no-static-element-interactions -->
		<div 
			class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-6 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
		>
			<!-- Icon -->
			<div class="flex items-center justify-center w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/30 mx-auto mb-4">
				<svg class="w-6 h-6 text-red-600 dark:text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
				</svg>
			</div>
			
			<!-- Title -->
			<h3 class="text-xl font-bold text-center text-gray-900 dark:text-white mb-2">
				Delete {
					deleteTarget.type === 'link-group' ? 'Link Group' : 
					deleteTarget.type === 'text-group' ? 'Text Group' :
					'Block'
				}?
			</h3>
			
			<!-- Description -->
			<p class="text-center text-gray-600 dark:text-gray-400 mb-6">
				{#if deleteTarget.type === 'link-group' || deleteTarget.type === 'text-group'}
					Are you sure you want to delete "<span class="font-semibold">{deleteTarget.title}</span>" and all its content? This action cannot be undone.
				{:else}
					Are you sure you want to delete "<span class="font-semibold">{deleteTarget.title}</span>"? This action cannot be undone.
				{/if}
			</p>
			
			<!-- Actions -->
			<div class="flex gap-3">
				<button
					onclick={() => {
						showDeleteConfirmModal = false;
						deleteTarget = null;
					}}
					class="flex-1 px-4 py-2.5 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg font-medium transition-colors"
				>
					Cancel
				</button>
				<button
					onclick={() => {
						if (deleteTarget) {
							if (deleteTarget.type === 'link-group') {
								confirmDeleteGroup(deleteTarget.id);
							} else if (deleteTarget.type === 'text-group') {
								confirmDeleteTextGroup(deleteTarget.id);
							} else if (deleteTarget.type === 'block') {
								confirmDeleteBlock(deleteTarget.id);
							}
						}
						showDeleteConfirmModal = false;
						deleteTarget = null;
					}}
					class="flex-1 px-4 py-2.5 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition-colors"
				>
					Delete
				</button>
			</div>
		</div>
	</div>
{/if}

