<script lang="ts">
	import { cn } from '$lib/utils/cn';

	export let variant: 'default' | 'outline' | 'ghost' = 'default';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let type: 'button' | 'submit' = 'button';

	const variants = {
		default: 'bg-black text-white hover:bg-gray-800',
		outline: 'border border-gray-300 hover:bg-gray-50',
		ghost: 'hover:bg-gray-100'
	};

	const sizes = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-4 py-2',
		lg: 'px-6 py-3 text-lg'
	};
</script>

<button
	{type}
	{disabled}
	class={cn(
		'rounded-lg font-medium transition disabled:opacity-50 disabled:cursor-not-allowed',
		variants[variant],
		sizes[size],
		$$props.class
	)}
	on:click
>
	<slot />
</button>
