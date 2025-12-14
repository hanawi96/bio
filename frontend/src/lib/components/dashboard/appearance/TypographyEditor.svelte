<script lang="ts">
	import { linksApi } from '$lib/api/links';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { previewStyles } from '$lib/stores/previewStyles';

	let { onUpdate, links = [] }: { onUpdate?: () => Promise<void>; links?: any[] } = $props();

	// Sync from first group in links
	const firstGroup = $derived(links.find(l => l.is_group));
	const textAlignment = $derived<'left' | 'center' | 'right'>(firstGroup?.text_alignment || 'center');
	const textSize = $derived<'S' | 'M' | 'L' | 'XL'>(firstGroup?.text_size || 'M');

	async function updateTypography(field: string, value: any) {
		previewStyles.update({ [field]: value });
		
		try {
			await linksApi.updateAllGroupStyles({
				[field]: value
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}
</script>

<div class="space-y-6">
	<h3 class="text-lg font-bold text-gray-900">Typography</h3>
	
	<!-- Text Alignment -->
	<div>
		<label class="text-sm font-semibold text-gray-900 mb-3 block">Text alignment</label>
		<div class="flex gap-2">
			<button
				onclick={() => updateTypography('text_alignment', 'left')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center"
				class:border-gray-900={textAlignment === 'left'}
				class:bg-gray-900={textAlignment === 'left'}
				class:text-white={textAlignment === 'left'}
				class:border-gray-200={textAlignment !== 'left'}
				class:text-gray-700={textAlignment !== 'left'}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h14"/>
				</svg>
			</button>
			
			<button
				onclick={() => updateTypography('text_alignment', 'center')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center"
				class:border-gray-900={textAlignment === 'center'}
				class:bg-gray-900={textAlignment === 'center'}
				class:text-white={textAlignment === 'center'}
				class:border-gray-200={textAlignment !== 'center'}
				class:text-gray-700={textAlignment !== 'center'}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M6 18h12"/>
				</svg>
			</button>
			
			<button
				onclick={() => updateTypography('text_alignment', 'right')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center"
				class:border-gray-900={textAlignment === 'right'}
				class:bg-gray-900={textAlignment === 'right'}
				class:text-white={textAlignment === 'right'}
				class:border-gray-200={textAlignment !== 'right'}
				class:text-gray-700={textAlignment !== 'right'}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M8 18h12"/>
				</svg>
			</button>
		</div>
	</div>

	<!-- Text Size -->
	<div>
		<label class="text-sm font-semibold text-gray-900 mb-3 block">Text size</label>
		<div class="flex gap-2">
			{#each ['S', 'M', 'L', 'XL'] as size}
				<button
					onclick={() => updateTypography('text_size', size)}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all font-semibold"
					class:border-gray-900={textSize === size}
					class:bg-gray-900={textSize === size}
					class:text-white={textSize === size}
					class:border-gray-200={textSize !== size}
					class:text-gray-700={textSize !== size}
					class:text-xs={size === 'S'}
					class:text-sm={size === 'M'}
					class:text-base={size === 'L'}
					class:text-lg={size === 'XL'}
				>
					{size}
				</button>
			{/each}
		</div>
	</div>
</div>
