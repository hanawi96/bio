<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { Link } from '$lib/api/links';

	let profile: any = null;
	let links: Link[] = [];
	let blocks: any[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			const data: any = await api.get(`/p/${$page.params.username}`);
			profile = data.profile;
			links = data.links || [];
			blocks = data.blocks || [];
		} catch (err: any) {
			error = err.message || 'Profile not found';
		} finally {
			loading = false;
		}
	});

	$: activeLinks = links.filter(l => l.is_active);
	$: activeBlocks = blocks.filter(b => b.is_active);
	
	// Combine and sort by position
	$: items = [
		...activeLinks.map(link => ({ type: 'link' as const, data: link, position: link.position })),
		...activeBlocks.map(block => ({ type: 'block' as const, data: block, position: block.position }))
	].sort((a, b) => a.position - b.position);
</script>

<svelte:head>
	<title>{$page.params.username} - LinkBio</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-purple-50 to-blue-50 py-12 px-4">
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
		<div class="max-w-md mx-auto">
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
					{#if item.type === 'block' && item.data.block_type === 'text'}
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
					{#if link.layout_type === 'featured'}
						<!-- Featured Layout -->
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
					{:else}
						<!-- Classic Layout -->
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
								<svg class="w-6 h-6 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
								</svg>
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
