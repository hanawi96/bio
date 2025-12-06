<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';

	let profile: any = null;
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			profile = await api.get(`/p/${$page.params.username}`);
		} catch (err: any) {
			error = err.message || 'Profile not found';
		} finally {
			loading = false;
		}
	});
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
			<div class="text-center mb-6">
				<div class="w-24 h-24 bg-gray-300 rounded-full mx-auto mb-4"></div>
				<h1 class="text-2xl font-bold">@{$page.params.username}</h1>
				<p class="text-gray-600 mt-2">{profile?.bio || 'No bio yet'}</p>
			</div>

			<!-- Links -->
			<div class="space-y-3">
				<div class="bg-white rounded-lg p-4 text-center shadow hover:shadow-md transition">
					<p class="font-medium">Example Link</p>
				</div>
			</div>
		</div>
	{/if}
</div>
