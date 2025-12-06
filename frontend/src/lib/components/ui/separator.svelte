<script lang="ts">
	import { cn } from '$lib/utils/cn';
	import type { HTMLAttributes } from 'svelte/elements';

	type $$Props = HTMLAttributes<HTMLDivElement> & {
		orientation?: 'horizontal' | 'vertical';
	};

	export let orientation: $$Props['orientation'] = 'horizontal';
</script>

<div
	role="separator"
	class={cn(
		'shrink-0 bg-gray-200',
		orientation === 'horizontal' ? 'h-[1px] w-full' : 'h-full w-[1px]',
		$$props.class
	)}
	{...$$restProps}
/>
