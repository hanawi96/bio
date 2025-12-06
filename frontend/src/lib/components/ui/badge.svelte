<script lang="ts">
	import { cn } from '$lib/utils/cn';
	import type { HTMLAttributes } from 'svelte/elements';

	type $$Props = HTMLAttributes<HTMLSpanElement> & {
		variant?: 'default' | 'secondary' | 'outline' | 'success';
	};

	export let variant: $$Props['variant'] = 'default';

	const variants = {
		default: 'bg-violet-100 text-violet-700 border-violet-200',
		secondary: 'bg-gray-100 text-gray-700 border-gray-200',
		outline: 'border-gray-300 text-gray-700',
		success: 'bg-green-100 text-green-700 border-green-200'
	};
</script>

<span
	class={cn(
		'inline-flex items-center rounded-full border px-3 py-1 text-sm font-medium transition-colors',
		variants[variant],
		$$props.class
	)}
	{...$$restProps}
>
	<slot />
</span>
