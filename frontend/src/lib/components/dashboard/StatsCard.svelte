<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	export let title: string;
	export let value: string | number;
	export let icon: string;
	export let trend: string = '';
	export let trendUp: boolean = true;
</script>

<Card.Root class="group relative border-0 hover:shadow-xl hover:shadow-indigo-100/50 transition-all duration-300">
	<div class="absolute inset-0 bg-gradient-to-br from-indigo-50/50 to-blue-50/50 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
	
	<Card.Content class="relative p-6">
		<div class="flex items-start justify-between">
			<div class="space-y-3">
				<p class="text-sm font-medium text-gray-500">{title}</p>
				<div class="flex items-baseline gap-2">
					<h3 class="text-4xl font-bold bg-gradient-to-br from-gray-900 to-gray-600 bg-clip-text text-transparent">
						{value}
					</h3>
					{#if trend}
						<Badge variant={trendUp ? 'default' : 'secondary'} class={trendUp ? 'bg-emerald-50 text-emerald-600 hover:bg-emerald-50' : 'bg-red-50 text-red-600 hover:bg-red-50'}>
							{trendUp ? '↑' : '↓'} {trend}
						</Badge>
					{/if}
				</div>
			</div>
			<div class="p-3 bg-gradient-to-br from-indigo-600 via-indigo-700 to-blue-700 rounded-xl shadow-lg shadow-indigo-500/30 group-hover:scale-110 transition-transform duration-300">
				<div class="w-6 h-6 text-white">
					{@html icon}
				</div>
			</div>
		</div>
	</Card.Content>
</Card.Root>
