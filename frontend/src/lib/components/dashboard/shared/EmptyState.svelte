<script lang="ts">
	import { Button } from '$lib/components/ui/button';

	export let title: string = 'No data yet';
	export let description: string = 'Get started by adding your first item';
	export let actionText: string = 'Get Started';
	export let icon: string = '';
	export let onAction: () => void = () => {};
</script>

<div class="flex flex-col items-center justify-center py-12 px-4 text-center">
	{#if icon}
		<div class="w-20 h-20 mb-6 text-gray-300">
			{@html icon}
		</div>
	{:else}
		<svg class="w-20 h-20 mb-6 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
		</svg>
	{/if}
	
	<h3 class="text-xl font-semibold text-gray-900 mb-2">{title}</h3>
	<p class="text-gray-600 mb-6 max-w-sm">{description}</p>
	
	{#if actionText}
		<Button on:click={onAction} class="bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700">
			{actionText}
		</Button>
	{/if}
</div>
