<script lang="ts">
	import { cn } from '$lib/utils/cn';

	export let src: string = '';
	export let alt: string = '';
	export let fallback: string = '';
	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';

	const sizes = {
		sm: 'w-8 h-8 text-xs',
		md: 'w-10 h-10 text-sm',
		lg: 'w-12 h-12 text-base',
		xl: 'w-16 h-16 text-lg'
	};

	let imageLoaded = false;
	let imageError = false;

	function handleLoad() {
		imageLoaded = true;
	}

	function handleError() {
		imageError = true;
	}
</script>

<div
	class={cn(
		'relative flex shrink-0 overflow-hidden rounded-full bg-gray-100',
		sizes[size],
		$$props.class
	)}
>
	{#if src && !imageError}
		<img
			{src}
			{alt}
			class="aspect-square h-full w-full object-cover"
			class:opacity-0={!imageLoaded}
			class:opacity-100={imageLoaded}
			on:load={handleLoad}
			on:error={handleError}
		/>
	{/if}
	{#if (!src || imageError) && fallback}
		<div class="flex h-full w-full items-center justify-center bg-gradient-to-br from-violet-500 to-indigo-500 text-white font-semibold">
			{fallback}
		</div>
	{/if}
</div>
