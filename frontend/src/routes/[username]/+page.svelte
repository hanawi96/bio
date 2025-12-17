<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { Link } from '$lib/api/links';
	import { globalTheme, themeStyles } from '$lib/stores/theme';
	
	$: styles = $themeStyles;
	$: currentTheme = $globalTheme;

	let profile: any = null;
	let links: Link[] = [];
	let blocks: any[] = [];
	let loading = true;
	let error = '';
	
	// Helper to get card properties from theme (NEVER use link values for card background)
	// This ensures consistent styling across all link groups based on theme
	function getCardProps() {
		const t = currentTheme;
		return {
			hasCustomBg: t.enableCardBackground ?? true,
			bgColor: t.cardBackground || '#ffffff',
			bgOpacity: t.cardBackgroundOpacity ?? 100,
			borderRadius: t.cardBorderRadius ?? 12,
			textColor: t.cardTextColor || '#000000',
			showShadow: t.cardShadow ?? false,
			shadowX: t.cardShadowX ?? 0,
			shadowY: t.cardShadowY ?? 4,
			shadowBlur: t.cardShadowBlur ?? 10,
			hasBorder: t.cardBorder ?? false,
			borderColor: t.cardBorderColor || '#e5e7eb',
			borderWidth: t.cardBorderWidth ?? 1
		};
	}

	// Helper function to parse padding from style
	function getPaddingStyle(style: string | null | undefined): string {
		try {
			if (!style) return '16px';
			const parsed = JSON.parse(style);
			const p = parsed.padding;
			if (!p) return '16px';
			if (typeof p === 'number') return `${p}px`;
			return `${p.top || 16}px ${p.right || 16}px ${p.bottom || 16}px ${p.left || 16}px`;
		} catch {
			return '16px';
		}
	}

	// Helper function to parse margin from style (only bottom to avoid CSS collapse)
	// Returns margin value in pixels (number) for use in gap or margin-bottom
	function getMarginValue(style: string | null | undefined): number {
		try {
			if (!style) return 0;
			const parsed = JSON.parse(style);
			const m = parsed.margin;
			if (!m) return 0;
			return m.bottom ?? 0;
		} catch {
			return 0;
		}
	}

	onMount(async () => {
		try {
			const data: any = await api.get(`/p/${$page.params.username}`);
			profile = data.profile;
			links = data.links || [];
			blocks = data.blocks || [];
			
			// Load theme if exists
			if (profile?.theme_config) {
				try {
					const themeStr = typeof profile.theme_config === 'string' 
						? profile.theme_config 
						: JSON.stringify(profile.theme_config);
					if (themeStr && themeStr !== '{}' && themeStr !== 'null') {
						globalTheme.loadFromJSON(themeStr);
					}
				} catch (e) {
					console.warn('Failed to load theme:', e);
				}
			}

			// Wait for DOM to render then check carousels
			setTimeout(() => {
				const carousels = document.querySelectorAll('.overflow-x-scroll');
				console.log('🎯 [Public Page] Found carousels:', carousels.length);
				
				carousels.forEach((carousel, idx) => {
					console.log(`📊 [Public Carousel ${idx}] Properties:`, {
						scrollWidth: carousel.scrollWidth,
						clientWidth: carousel.clientWidth,
						overflowX: getComputedStyle(carousel).overflowX,
						canScroll: carousel.scrollWidth > carousel.clientWidth
					});

					carousel.addEventListener('scroll', (e) => {
						console.log(`📜 [Public Carousel ${idx}] Scrolling:`, {
							scrollLeft: carousel.scrollLeft,
							scrollWidth: carousel.scrollWidth,
							clientWidth: carousel.clientWidth
						});
					});

					carousel.addEventListener('touchstart', () => {
						console.log(`👆 [Public Carousel ${idx}] Touch started`);
					});

					carousel.addEventListener('touchmove', () => {
						console.log(`👆 [Public Carousel ${idx}] Touch moving`);
					});
				});
			}, 500);
		} catch (err: any) {
			error = err.message || 'Profile not found';
		} finally {
			loading = false;
		}
	});

	$: activeBlocks = blocks.filter(b => b.is_active);
	
	// Get link_ids from inactive blocks - these links should be hidden
	$: hiddenLinkIds = new Set(
		blocks
			.filter(block => !block.is_active && block.link_id)
			.map(block => block.link_id)
	);
	
	// Filter active links, excluding those hidden by inactive blocks
	// Apply theme styling to links
	$: activeLinks = links
		.filter(l => l.is_active && !hiddenLinkIds.has(l.id))
		.map(link => {
			if (!link.is_group) return link;
			// Override with theme styling
			const theme = $globalTheme;
			return {
				...link,
				text_alignment: theme.textAlignment || link.text_alignment || 'center',
				text_size: theme.textSize || link.text_size || 'M',
				image_shape: theme.imageShape || link.image_shape || 'square'
			};
		});
	
	// Combine and sort: pinned links first, then by position
	$: items = [
		...activeLinks.map(link => ({ type: 'link' as const, data: link, position: link.position, isPinned: link.is_pinned || false })),
		...activeBlocks.map(block => ({ type: 'block' as const, data: block, position: block.position, isPinned: false }))
	].sort((a, b) => {
		// Pinned items first
		if (a.isPinned && !b.isPinned) return -1;
		if (!a.isPinned && b.isPinned) return 1;
		// Then by position
		return a.position - b.position;
	});
</script>

<svelte:head>
	<title>{$page.params.username} - LinkBio</title>
</svelte:head>

<div class="min-h-screen py-12 px-4 relative" style="{
	styles.pageBackgroundType === 'image' && styles.pageBackgroundImage 
		? `background: url(${styles.pageBackgroundImage}) center/cover no-repeat fixed;` 
		: styles.pageBackgroundType === 'video' && styles.pageBackgroundVideo
		? `background: ${styles.pageBackground};`
		: `background: ${styles.pageBackground};`
}">
	{#if styles.pageBackgroundType === 'video' && styles.pageBackgroundVideo}
		<video autoplay muted loop playsinline class="fixed inset-0 w-full h-full object-cover -z-10">
			<source src={styles.pageBackgroundVideo} type="video/mp4" />
		</video>
		<div class="fixed inset-0 bg-black/20 -z-10"></div>
	{/if}
	<!-- Fixed Action Buttons -->
	<div class="fixed top-6 left-6 right-6 flex justify-between items-center z-50 max-w-md mx-auto">
		<button 
			class="flex items-center gap-2 px-5 py-2.5 bg-white/90 backdrop-blur-sm hover:bg-white rounded-full shadow-lg transition-all hover:shadow-xl border border-gray-200"
			onclick={() => alert('Subscribe functionality')}
		>
			<svg class="w-5 h-5 text-gray-700" fill="currentColor" viewBox="0 0 20 20">
				<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
			</svg>
			<span class="text-sm font-semibold text-gray-800">Subscribe</span>
		</button>
		
		<button 
			class="flex items-center justify-center w-11 h-11 bg-white/90 backdrop-blur-sm hover:bg-white rounded-full shadow-lg transition-all hover:shadow-xl border border-gray-200 mr-1"
			onclick={() => {
				if (navigator.share) {
					navigator.share({
						title: `@${$page.params.username}`,
						url: window.location.href
					}).catch(() => {});
				} else {
					navigator.clipboard.writeText(window.location.href);
					alert('Link copied to clipboard!');
				}
			}}
			title="Share"
		>
			<svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
			</svg>
		</button>
	</div>
	
	{#if loading}
		<div class="flex items-center justify-center min-h-[60vh]">
			<div class="relative">
				<div class="w-16 h-16 border-4 border-indigo-200 rounded-full"></div>
				<div class="w-16 h-16 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin absolute top-0 left-0"></div>
			</div>
		</div>
	{:else if error}
		<div class="max-w-md mx-auto text-center">
			<h1 class="text-4xl font-bold mb-4">404</h1>
			<p class="text-gray-600">{error}</p>
			<a href="/" class="mt-4 inline-block text-black hover:underline">Go home</a>
		</div>
	{:else}
		<div class="max-w-md mx-auto pt-16">
			<!-- Avatar -->
			<div class="text-center mb-8">
				{#if profile?.avatar_url}
					<img 
						src={profile.avatar_url} 
						alt={profile.username}
						class="w-24 h-24 rounded-full mx-auto mb-4 border-4 border-white shadow-lg object-cover"
					/>
				{:else}
					<div class="w-24 h-24 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-full mx-auto mb-4 border-4 border-white shadow-lg"></div>
				{/if}
				<h1 class="text-2xl font-bold text-gray-900">@{$page.params.username}</h1>
				{#if profile?.bio}
					<p class="text-gray-600 mt-2 max-w-sm mx-auto">{profile.bio}</p>
				{/if}
			</div>

			<!-- Links & Blocks -->
			<div class="space-y-4">
				{#each items as item}
					{#if item.type === 'block' && item.data.block_type === 'text' && !item.data.is_group}
						{@const styleConfig = item.data.style ? JSON.parse(item.data.style) : {}}
						{@const content = item.data.content || ''}
						{@const parts = content.split(/(\[([^\]]+)\]\(([^)]+)\))/g)}
						<div 
							class="w-full rounded-xl"
							class:p-4={styleConfig.backgroundColor}
							style="
								font-size: {
									styleConfig.fontSize === 'text-small' ? '14px' : 
									styleConfig.fontSize === 'text-medium' ? '16px' : 
									styleConfig.fontSize === 'text-large' ? '18px' :
									styleConfig.fontSize === 'headline-small' ? '20px' :
									styleConfig.fontSize === 'headline-medium' ? '24px' :
									styleConfig.fontSize === 'headline-large' ? '32px' : '16px'
								};
								text-align: {styleConfig.textAlign || 'left'};
								font-weight: {styleConfig.fontSize?.startsWith('headline') || styleConfig.isBold ? 'bold' : 'normal'};
								font-style: {styleConfig.isItalic ? 'italic' : 'normal'};
								text-decoration: {styleConfig.isUnderline && styleConfig.isStrikethrough ? 'underline line-through' : styleConfig.isUnderline ? 'underline' : styleConfig.isStrikethrough ? 'line-through' : 'none'};
								color: {styleConfig.textColor || '#000000'};
								background-color: {styleConfig.backgroundColor || 'transparent'};
							"
						>
							{#each parts as part, i}
								{#if i % 4 === 0}
									{part}
								{:else if i % 4 === 2}
									<a href={parts[i + 1]} target="_blank" rel="noopener noreferrer" class="underline hover:opacity-70 transition-opacity">{part}</a>
								{/if}
							{/each}
						</div>
					{:else if item.type === 'block' && item.data.block_type === 'text' && item.data.is_group}
						<!-- Text Group: Render all children with group style -->
						{#if item.data.children && item.data.children.length > 0}
							{@const sortedChildren = item.data.children
								.filter(c => c.is_active)
								.sort((a, b) => a.position - b.position)}
							{@const groupStyle = (() => {
								try {
									if (!item.data.style) {
										return { textAlign: 'left', fontSize: 'text-medium', textColor: '#000000', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none', hasBackground: false };
									}
									const parsed = typeof item.data.style === 'string' ? JSON.parse(item.data.style) : item.data.style;
									return {
										textAlign: parsed.textAlign || 'left',
										fontSize: parsed.fontSize || 'text-medium',
										textColor: parsed.textColor || '#000000',
										isBold: parsed.isBold || false,
										isItalic: parsed.isItalic || false,
										isUnderline: parsed.isUnderline || false,
										isStrikethrough: parsed.isStrikethrough || false,
										textTransform: parsed.textTransform || 'none',
										hasBackground: parsed.hasBackground || false,
										backgroundColor: parsed.backgroundColor || '#ffffff',
										backgroundOpacity: parsed.backgroundOpacity !== undefined ? parsed.backgroundOpacity : 90,
										borderRadius: parsed.borderRadius !== undefined ? parsed.borderRadius : 12,
										padding: parsed.padding || 16,
										shadow: parsed.shadow || 'none',
										hasBorder: parsed.hasBorder || false,
										borderColor: parsed.borderColor || '#e5e7eb',
										borderWidth: parsed.borderWidth || 1,
										borderStyle: parsed.borderStyle || 'solid'
									};
								} catch (e) {
									console.error('Failed to parse groupStyle:', item.data.style, e);
									return { textAlign: 'left', fontSize: 'text-medium', textColor: '#000000', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none', hasBackground: false };
								}
							})()}

							<div class="space-y-2">
								{#each sortedChildren as child}
									{#if groupStyle.hasBackground}
										{@const r = parseInt(groupStyle.backgroundColor?.slice(1,3) || 'ff', 16)}
										{@const g = parseInt(groupStyle.backgroundColor?.slice(3,5) || 'ff', 16)}
										{@const b = parseInt(groupStyle.backgroundColor?.slice(5,7) || 'ff', 16)}
										{@const opacity = (groupStyle.backgroundOpacity ?? 90) / 100}
										{@const bgColor = `rgba(${r}, ${g}, ${b}, ${opacity})`}
										{@const borderStyle = groupStyle.hasBorder ? `${groupStyle.borderWidth ?? 1}px ${groupStyle.borderStyle || 'solid'} ${groupStyle.borderColor || '#e5e7eb'}` : 'none'}
										<div 
											class="rounded-xl"
											class:shadow-sm={groupStyle.shadow === 'sm'}
											class:shadow-md={groupStyle.shadow === 'md'}
											class:shadow-lg={groupStyle.shadow === 'lg'}
											style:background-color={bgColor}
											style:border-radius="{groupStyle.borderRadius ?? 12}px"
											style:padding="{groupStyle.padding ?? 16}px"
											style:border={borderStyle}
										>
											<p 
												class="break-words"
												style="
													font-size: {groupStyle.fontSize === 'headline-large' ? '1.5rem' : groupStyle.fontSize === 'headline-medium' ? '1.25rem' : groupStyle.fontSize === 'headline-small' ? '1.125rem' : groupStyle.fontSize === 'text-large' ? '1.125rem' : groupStyle.fontSize === 'text-small' ? '0.875rem' : '1rem'};
													text-align: {groupStyle.textAlign || 'left'};
													font-weight: {groupStyle.isBold ? 'bold' : 'normal'};
													font-style: {groupStyle.isItalic ? 'italic' : 'normal'};
													text-decoration: {groupStyle.isUnderline && groupStyle.isStrikethrough ? 'underline line-through' : groupStyle.isUnderline ? 'underline' : groupStyle.isStrikethrough ? 'line-through' : 'none'};
													text-transform: {groupStyle.textTransform || 'none'};
													color: {groupStyle.textColor || '#000000'};
												"
											>
												{child.content || 'Empty text'}
											</p>
										</div>
									{:else}
										<p 
											class="break-words"
											style="
												font-size: {groupStyle.fontSize === 'headline-large' ? '1.5rem' : groupStyle.fontSize === 'headline-medium' ? '1.25rem' : groupStyle.fontSize === 'headline-small' ? '1.125rem' : groupStyle.fontSize === 'text-large' ? '1.125rem' : groupStyle.fontSize === 'text-small' ? '0.875rem' : '1rem'};
												text-align: {groupStyle.textAlign || 'left'};
												font-weight: {groupStyle.isBold ? 'bold' : 'normal'};
												font-style: {groupStyle.isItalic ? 'italic' : 'normal'};
												text-decoration: {groupStyle.isUnderline && groupStyle.isStrikethrough ? 'underline line-through' : groupStyle.isUnderline ? 'underline' : groupStyle.isStrikethrough ? 'line-through' : 'none'};
												text-transform: {groupStyle.textTransform || 'none'};
												color: {groupStyle.textColor || '#000000'};
											"
										>
											{child.content || 'Empty text'}
										</p>
									{/if}
								{/each}
							</div>
						{/if}
					{:else if item.type === 'block' && item.data.block_type === 'divider'}
						{#if item.data.divider_style === 'line'}
							<div class="border-t border-gray-200"></div>
						{:else if item.data.divider_style === 'dots'}
							<div class="flex justify-center gap-1">
								<div class="w-1.5 h-1.5 rounded-full bg-gray-300"></div>
								<div class="w-1.5 h-1.5 rounded-full bg-gray-300"></div>
								<div class="w-1.5 h-1.5 rounded-full bg-gray-300"></div>
							</div>
						{:else}
							<div class="h-4"></div>
						{/if}
					{:else if item.type === 'link'}
						{@const link = item.data}
					{#if link.is_group}
						<!-- Group Link -->
						{#if link.children && link.children.length > 0}
							{@const sortedChildren = link.children
								.filter(c => c.is_active)
								.sort((a, b) => {
									if (a.is_pinned && !b.is_pinned) return -1;
									if (!a.is_pinned && b.is_pinned) return 1;
									return a.position - b.position;
								})}
							{console.log('🌐 [Public Page] Group Link:', {
								id: link.id,
								title: link.group_title,
								layout: link.group_layout,
								children: sortedChildren.length,
								childrenData: sortedChildren
							})}
							<div class="space-y-3">
								{#if link.group_layout === 'grid'}
									{@const gridCols = link.grid_columns || 2}
									{@const aspectRatio = link.grid_aspect_ratio || '3:2'}
									{@const textSize = link.text_size || 'M'}
									{@const cardProps = getCardProps()}
									{@const r = parseInt(cardProps.bgColor.slice(1,3), 16)}
									{@const g = parseInt(cardProps.bgColor.slice(3,5), 16)}
									{@const b = parseInt(cardProps.bgColor.slice(5,7), 16)}
									{@const shadowStyle = cardProps.showShadow ? `box-shadow: ${cardProps.shadowX}px ${cardProps.shadowY}px ${cardProps.shadowBlur}px rgba(0,0,0,0.2);` : ''}
									{@const borderStyle = cardProps.hasBorder ? `border: ${cardProps.borderWidth}px solid ${cardProps.borderColor};` : ''}
									{@const cardSpacing = getMarginValue(link.style)}
									{@const bgStyle = cardProps.hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${cardProps.bgOpacity / 100}); border-radius: ${cardProps.borderRadius}px; padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}` : `padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}`}
									{@const getAspectStyle = (ratio) => {
										const map = { '1:1': '1/1', '3:2': '3/2', '16:9': '16/9', '3:1': '3/1', '2:3': '2/3' };
										return `aspect-ratio: ${map[ratio] || '3/2'}`;
									}}
									<div class="grid" class:grid-cols-1={gridCols === 1} class:grid-cols-2={gridCols === 2} class:grid-cols-3={gridCols === 3} class:grid-cols-4={gridCols === 4} style="gap: {cardSpacing ?? 12}px;">
										{#each sortedChildren as child}
											<a
												href={child.url}
												target="{child.open_in_new_tab ? '_blank' : '_self'}"
												rel="noopener noreferrer"
												class="block transition-all"
												class:hover:bg-gray-50={cardProps.hasCustomBg}
												class:border-2={link.show_outline && !cardProps.hasBorder}
												class:border-gray-200={link.show_outline && !cardProps.hasBorder}
												style={bgStyle}
											>
												{#if child.thumbnail_url}
													<img src={child.thumbnail_url} alt={child.title} class="w-full object-cover rounded-lg mb-2" style={getAspectStyle(aspectRatio)} />
												{/if}
												<p class="font-medium truncate"
													class:text-xs={textSize === 'S'}
													class:text-sm={textSize === 'M'}
													class:text-base={textSize === 'L'}
													class:text-lg={textSize === 'XL'}
													style="text-align: {link.text_alignment || 'center'}; color: {cardProps.textColor};"
												>{child.title}</p>
												{#if link.show_description !== false && child.description}
													<p class="text-xs mt-1 line-clamp-2" style="text-align: {link.text_alignment || 'center'}; color: {cardProps.textColor}; opacity: 0.7;">{child.description}</p>
												{/if}
											</a>
										{/each}
									</div>
								{:else if link.group_layout === 'carousel'}
									{@const textSize = link.text_size || 'M'}
									{@const aspectRatio = link.grid_aspect_ratio || '3:2'}
									{@const cardSpacing = getMarginValue(link.style)}
									{@const carouselId = `public-carousel-${link.id}`}
									{@const cardProps = getCardProps()}
									{@const r = parseInt(cardProps.bgColor.slice(1,3), 16)}
									{@const g = parseInt(cardProps.bgColor.slice(3,5), 16)}
									{@const b = parseInt(cardProps.bgColor.slice(5,7), 16)}
									{@const shadowStyle = cardProps.showShadow ? `box-shadow: ${cardProps.shadowX}px ${cardProps.shadowY}px ${cardProps.shadowBlur}px rgba(0,0,0,0.2);` : ''}
									{@const borderStyle = cardProps.hasBorder ? `border: ${cardProps.borderWidth}px solid ${cardProps.borderColor};` : ''}
									{@const bgStyle = cardProps.hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${cardProps.bgOpacity / 100}); border-radius: ${cardProps.borderRadius}px; padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}` : `padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}`}
									{@const getAspectStyle = (ratio) => {
										const map = { '1:1': '1/1', '3:2': '3/2', '16:9': '16/9', '3:1': '3/1', '2:3': '2/3' };
										return `aspect-ratio: ${map[ratio] || '3/2'}`;
									}}
									<div class="relative group">
										<!-- Carousel Container -->
										<div id={carouselId} class="overflow-x-scroll snap-x snap-mandatory scrollbar-hide" style="scroll-behavior: smooth;">
											<div class="flex px-4" style="gap: {cardSpacing ?? 12}px;">
												{#each sortedChildren as child, idx}
													<a
														href={child.url}
														target="{child.open_in_new_tab ? '_blank' : '_self'}"
														rel="noopener noreferrer"
														class="block transition-all flex-shrink-0 snap-center w-[85%]"
														class:hover:bg-gray-50={cardProps.hasCustomBg}
														style={bgStyle}
													>
														{#if child.thumbnail_url}
															<img src={child.thumbnail_url} alt={child.title} class="w-full object-cover rounded-lg mb-3" style={getAspectStyle(aspectRatio)}/>
														{/if}
														{#if link.show_text !== false}
															<p class="font-medium mb-1"
																class:text-xs={textSize === 'S'}
																class:text-sm={textSize === 'M'}
																class:text-base={textSize === 'L'}
																class:text-lg={textSize === 'XL'}
																style="text-align: {link.text_alignment || 'center'}; color: {cardProps.textColor};"
															>{child.title}</p>
															{#if child.description}
																<p class="text-xs line-clamp-2" style="text-align: {link.text_alignment || 'center'}; color: {cardProps.textColor}; opacity: 0.7;">{child.description}</p>
															{/if}
														{/if}
													</a>
												{/each}
											</div>
										</div>

										<!-- Navigation Arrows -->
										{#if sortedChildren.length > 1}
											<button
												onclick={() => {
													const el = document.getElementById(carouselId);
													if (el) el.scrollBy({ left: -el.clientWidth * 0.85, behavior: 'smooth' });
												}}
												class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-white/90 hover:bg-white rounded-full shadow-lg flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
												aria-label="Previous"
											>
												<svg class="w-5 h-5 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
												</svg>
											</button>
											<button
												onclick={() => {
													const el = document.getElementById(carouselId);
													if (el) el.scrollBy({ left: el.clientWidth * 0.85, behavior: 'smooth' });
												}}
												class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-white/90 hover:bg-white rounded-full shadow-lg flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
												aria-label="Next"
											>
												<svg class="w-5 h-5 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
												</svg>
											</button>

											<!-- Dots Indicator -->
											<div class="flex justify-center gap-1.5 mt-3">
												{#each sortedChildren as _, idx}
													<button
														onclick={() => {
															const el = document.getElementById(carouselId);
															if (el) {
																const cardWidth = el.clientWidth * 0.85;
																el.scrollTo({ left: cardWidth * idx, behavior: 'smooth' });
															}
														}}
														class="w-1.5 h-1.5 rounded-full bg-gray-300 hover:bg-gray-500 transition-colors cursor-pointer"
														aria-label={`Go to slide ${idx + 1}`}
													></button>
												{/each}
											</div>
										{/if}
									</div>
								{:else if link.group_layout === 'card'}
									<!-- Card layout -->
									{@const textSize = link.text_size || 'M'}
									{@const imagePlacement = link.image_placement || 'alternating'}
									{@const cardSpacing = getMarginValue(link.style)}
									{@const cardProps = getCardProps()}
									{@const r = parseInt(cardProps.bgColor.slice(1,3), 16)}
									{@const g = parseInt(cardProps.bgColor.slice(3,5), 16)}
									{@const b = parseInt(cardProps.bgColor.slice(5,7), 16)}
									{@const shadowStyle = cardProps.showShadow ? `box-shadow: ${cardProps.shadowX}px ${cardProps.shadowY}px ${cardProps.shadowBlur}px rgba(0,0,0,0.1);` : ''}
									{@const borderStyle = cardProps.hasBorder ? `border: ${cardProps.borderWidth}px solid ${cardProps.borderColor};` : ''}
									{@const bgStyle = cardProps.hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${cardProps.bgOpacity / 100}); border-radius: ${cardProps.borderRadius}px; ${shadowStyle} ${borderStyle}` : `${shadowStyle} ${borderStyle}`}
									<div style="display: flex; flex-direction: column; gap: {cardSpacing ?? 12}px;">
										{#each sortedChildren as child, index}
											{@const shouldReverse = imagePlacement === 'right' || (imagePlacement === 'alternating' && index % 2 === 0)}
											<a
												href={child.url}
												target="{child.open_in_new_tab ? '_blank' : '_self'}"
												rel="noopener noreferrer"
												class="block overflow-hidden transition-all"
												class:hover:bg-gray-50={cardProps.hasCustomBg}
												class:border-2={link.show_outline && !cardProps.hasBorder}
												class:border-gray-200={link.show_outline && !cardProps.hasBorder}
												style={bgStyle}
											>
												<div class="flex items-stretch" class:flex-row-reverse={shouldReverse}>
													{#if child.thumbnail_url}
														<div class="w-1/2 bg-gray-100">
															<img src={child.thumbnail_url} alt={child.title} class="w-full h-full object-cover"/>
														</div>
													{/if}
													<div class="flex-1 flex flex-col justify-center" style="padding: {getPaddingStyle(link.style)}">
														<p class="font-bold mb-1"
															class:text-sm={textSize === 'S'}
															class:text-base={textSize === 'M'}
															class:text-lg={textSize === 'L'}
															class:text-xl={textSize === 'XL'}
															style="text-align: {link.text_alignment || 'left'}; color: {cardProps.textColor};"
														>{child.title}</p>
														{#if link.show_description !== false && child.description}
															<p class="text-xs" style="text-align: {link.text_alignment || 'left'}; color: {cardProps.textColor}; opacity: 0.7;">{child.description}</p>
														{/if}
													</div>
												</div>
											</a>
										{/each}
									</div>
								{:else}
									<!-- List layout (default) -->
									{@const textSize = link.text_size || 'M'}
									{@const imageShape = link.image_shape || 'square'}
									{@const cardSpacing = getMarginValue(link.style)}
									{@const cardProps = getCardProps()}
									{@const r = parseInt(cardProps.bgColor.slice(1,3), 16)}
									{@const g = parseInt(cardProps.bgColor.slice(3,5), 16)}
									{@const b = parseInt(cardProps.bgColor.slice(5,7), 16)}
									{@const shadowStyle = cardProps.showShadow ? `box-shadow: ${cardProps.shadowX}px ${cardProps.shadowY}px ${cardProps.shadowBlur}px rgba(0,0,0,0.2);` : ''}
									{@const borderStyle = cardProps.hasBorder ? `border: ${cardProps.borderWidth}px solid ${cardProps.borderColor};` : ''}
									{@const bgStyle = cardProps.hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${cardProps.bgOpacity / 100}); border-radius: ${cardProps.borderRadius}px; padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}` : `padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}`}
									<div style="display: flex; flex-direction: column; gap: {cardSpacing ?? 12}px;">
										{#each sortedChildren as child}
											<a
												href={child.url}
												target="{child.open_in_new_tab ? '_blank' : '_self'}"
												rel="noopener noreferrer"
												class="block transition-all"
												class:hover:bg-gray-50={cardProps.hasCustomBg}
												class:border-2={link.show_outline && !cardProps.hasBorder}
												class:border-gray-200={link.show_outline && !cardProps.hasBorder}
												style={bgStyle}
											>
												<div class="flex items-center gap-3">
													{#if child.thumbnail_url}
														<img src={child.thumbnail_url} alt={child.title} class="w-10 h-10 object-cover flex-shrink-0" class:rounded-lg={imageShape === 'square'} class:rounded-full={imageShape === 'circle'} class:rounded-none={imageShape === 'sharp'}/>
													{/if}
													<div class="flex-1">
														<p class="font-medium"
															class:text-xs={textSize === 'S'}
															class:text-sm={textSize === 'M'}
															class:text-base={textSize === 'L'}
															class:text-lg={textSize === 'XL'}
															style="text-align: {link.text_alignment || 'left'}; color: {cardProps.textColor};"
														>{child.title}</p>
														{#if link.show_description !== false && child.description}
															<p class="text-xs mt-0.5" style="text-align: {link.text_alignment || 'left'}; color: {cardProps.textColor}; opacity: 0.7;">{child.description}</p>
														{/if}
													</div>
												</div>
											</a>
										{/each}
									</div>
								{/if}
							</div>
						{/if}
					{:else if !link.is_group && link.layout_type === 'featured'}
						<!-- Featured Layout (Non-group links only) -->
						<a
							href={link.url}
							target="_blank"
							rel="noopener noreferrer"
							class="block w-full bg-white hover:scale-[1.02] rounded-2xl overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300 border border-gray-100"
						>
							{#if link.thumbnail_url}
								<div class="relative h-48 bg-gradient-to-br from-indigo-100 to-blue-100">
									<img 
										src={link.thumbnail_url} 
										alt={link.title}
										class="w-full h-full object-cover"
									/>
									<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent"></div>
									<div class="absolute bottom-0 left-0 right-0 p-6">
										<h3 class="font-bold text-white text-xl drop-shadow-lg">{link.title}</h3>
									</div>
								</div>
							{:else}
								<div class="relative h-48 bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 flex items-center justify-center">
									<div class="absolute bottom-0 left-0 right-0 p-6">
										<h3 class="font-bold text-white text-xl drop-shadow-lg">{link.title}</h3>
									</div>
								</div>
							{/if}
						</a>
					{:else if !link.is_group}
						<!-- Classic Layout (Non-group links only) -->
						<a
							href={link.url}
							target="_blank"
							rel="noopener noreferrer"
							class="block w-full bg-white hover:scale-[1.02] rounded-2xl p-5 shadow-md hover:shadow-lg transition-all duration-300 border border-gray-100"
						>
							<div class="flex items-center gap-4">
								{#if link.thumbnail_url}
									<div class="w-12 h-12 rounded-xl overflow-hidden flex-shrink-0 bg-gray-100 border-2 border-gray-200">
										<img 
											src={link.thumbnail_url} 
											alt={link.title}
											class="w-full h-full object-cover"
										/>
									</div>
								{/if}
								<span class="font-semibold text-gray-900 flex-1 text-lg">{link.title}</span>
							</div>
						</a>
					{/if}
					{/if}
				{/each}
				
				{#if items.length === 0}
					<div class="text-center py-12 text-gray-400">
						<svg class="w-16 h-16 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
						</svg>
						<p class="text-sm">No links yet</p>
					</div>
				{/if}
			</div>

			<!-- Footer -->
			<div class="text-center mt-12 pt-8 border-t border-gray-200">
				<p class="text-sm text-gray-500">Made with <span class="text-red-500">♥</span> using LinkBio</p>
			</div>
		</div>
	{/if}
</div>

<style>
	.scrollbar-hide {
		-ms-overflow-style: none;
		scrollbar-width: none;
		-webkit-overflow-scrolling: touch;
	}
	.scrollbar-hide::-webkit-scrollbar {
		display: none;
	}
</style>
