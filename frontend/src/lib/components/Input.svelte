<script lang="ts">
	import { cn } from '$lib/utils/cn';

	export let type = 'text';
	export let value = '';
	export let placeholder = '';
	export let disabled = false;
	export let required = false;
</script>

<input
	{type}
	{placeholder}
	{disabled}
	{required}
	bind:value
	class={cn(
		'w-full px-3 py-2 border border-gray-300 rounded-md',
		'focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent',
		'disabled:bg-gray-100 disabled:cursor-not-allowed',
		$$props.class
	)}
	{...$$restProps}
/>
