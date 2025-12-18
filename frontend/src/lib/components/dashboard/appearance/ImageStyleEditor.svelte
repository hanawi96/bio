<script lang="ts">
	import { previewStyles } from '$lib/stores/previewStyles';
	import { globalTheme } from '$lib/stores/theme';

	let { links = [] }: { links?: any[] } = $props();

	// Read from globalTheme (single source of truth)
	const imageShape = $derived<'sharp' | 'square' | 'circle'>($globalTheme.imageShape || 'square');

	function updateImageShape(shape: 'sharp' | 'square' | 'circle') {
		// Update theme config
		globalTheme.update({ imageShape: shape });
		
		// Update preview
		previewStyles.update({ image_shape: shape });
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h3 class="text-lg font-bold text-gray-900">Image Shape</h3>
		{#if links.filter(l => l.is_group && l.image_shape !== null && l.image_shape !== undefined).length > 0}
			<span class="text-xs text-amber-600 font-medium flex items-center gap-1">
				<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
					<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
				</svg>
				{links.filter(l => l.is_group && l.image_shape !== null && l.image_shape !== undefined).length} group(s) have custom values
			</span>
		{/if}
	</div>
	
	<div class="flex gap-2">
		<button
			onclick={() => updateImageShape('sharp')}
			class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md"
			class:border-indigo-600={imageShape === 'sharp'}
			class:bg-indigo-50={imageShape === 'sharp'}
			class:border-gray-200={imageShape !== 'sharp'}
		>
			<div class="flex items-center gap-3">
				<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 flex-shrink-0"></div>
				<span class="text-sm font-medium text-gray-700">Sharp</span>
			</div>
		</button>
		
		<button
			onclick={() => updateImageShape('square')}
			class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md"
			class:border-indigo-600={imageShape === 'square'}
			class:bg-indigo-50={imageShape === 'square'}
			class:border-gray-200={imageShape !== 'square'}
		>
			<div class="flex items-center gap-3">
				<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-lg flex-shrink-0"></div>
				<span class="text-sm font-medium text-gray-700">Rounded</span>
			</div>
		</button>
		
		<button
			onclick={() => updateImageShape('circle')}
			class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md"
			class:border-indigo-600={imageShape === 'circle'}
			class:bg-indigo-50={imageShape === 'circle'}
			class:border-gray-200={imageShape !== 'circle'}
		>
			<div class="flex items-center gap-3">
				<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-full flex-shrink-0"></div>
				<span class="text-sm font-medium text-gray-700">Circle</span>
			</div>
		</button>
	</div>
</div>
