<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';
	
	export let group: Link;
	
	const dispatch = createEventDispatcher();
	
	let groupTitle = group.group_title || '';
	let description = group.description || '';
	let isActive = group.is_active !== undefined ? group.is_active : true;
	
	// Sync with prop changes
	$: groupTitle = group.group_title || '';
	$: description = group.description || '';
	$: isActive = group.is_active !== undefined ? group.is_active : true;
	
	function updateGroupTitle() {
		if (groupTitle.trim()) {
			dispatch('update', {
				groupId: group.id,
				group_title: groupTitle.trim()
			});
		}
	}
	
	function updateDescription() {
		dispatch('update', {
			groupId: group.id,
			description: description.trim() || null
		});
	}
	
	function toggleActive() {
		dispatch('update', {
			groupId: group.id,
			is_active: !isActive
		});
	}
</script>

<div class="p-6 space-y-6">
	<!-- Basic Settings -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-4">Basic</h3>
		<div class="space-y-4">
			<!-- Group Title -->
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-2">
					Group Title
				</label>
				<input
					type="text"
					bind:value={groupTitle}
					onblur={updateGroupTitle}
					onkeydown={(e) => e.key === 'Enter' && updateGroupTitle()}
					placeholder="Enter group title"
					class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all"
				/>
			</div>

			<!-- Description -->
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-2">
					Description
				</label>
				<textarea
					bind:value={description}
					onblur={updateDescription}
					placeholder="Add a description for this group (optional)"
					rows="3"
					class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all resize-none"
				></textarea>
			</div>

			<!-- Show Group Toggle -->
			<div class="flex items-center justify-between py-2">
				<div>
					<label class="text-sm font-medium text-gray-700">Show Group</label>
					<p class="text-xs text-gray-500 mt-0.5">Make this group visible on your profile</p>
				</div>
				<button
					type="button"
					onclick={toggleActive}
					class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2"
					class:bg-emerald-600={isActive}
					class:bg-gray-200={!isActive}
					role="switch"
					aria-checked={isActive}
				>
					<span
						class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
						class:translate-x-5={isActive}
						class:translate-x-0={!isActive}
					></span>
				</button>
			</div>
		</div>
	</div>
</div>
