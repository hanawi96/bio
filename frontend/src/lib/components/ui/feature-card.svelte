<script lang="ts">
	import { cn } from '$lib/utils/cn';

	export let icon: string = '';
	export let title: string;
	export let description: string;
	export let gradient: string = 'from-violet-500 to-indigo-500';
</script>

<div
	class={cn(
		'group bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300',
		$$props.class
	)}
>
	<div
		class={cn(
			'w-14 h-14 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform bg-gradient-to-br',
			gradient
		)}
	>
		<slot name="icon">
			{#if icon}
				<div class="text-white text-2xl">{icon}</div>
			{/if}
		</slot>
	</div>
	<h3 class="text-xl font-bold text-gray-900 mb-3">{title}</h3>
	<p class="text-gray-600 leading-relaxed">
		<slot>{description}</slot>
	</p>
</div>
