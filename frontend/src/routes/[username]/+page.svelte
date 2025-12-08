<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { Link } from '$lib/api/links';

	let profile: any = null;
	let links: Link[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			const data: any = await api.get(`/p/${$page.params.username}`);
			profile = data.profile;
			links = data.links || [];
		} catch (err: any) {
			error = err.message || 'Profile not found';
		} finally {
			loading = false;
		}
	});

	$: activeLinks = links.filter(l => l.is_active);
</script>

<svelte:head>
	<title>{$page.params.username} - LinkBio</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-purple-50 to-blue-50 py-12 px-4">
	{#if loading}
		<div class="max-w-md mx-auto text-center">
			<p class="text-gray-600">Loading...</p>
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

			<!-- Links -->
			<div class="space-y-4">
				{#each activeLinks as link}
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
				{/each}
				
				{#if activeLinks.length === 0}
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
