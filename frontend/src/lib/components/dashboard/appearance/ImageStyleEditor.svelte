<script lang="ts">
	import { previewStyles } from '$lib/stores/previewStyles';
	import { pendingChanges } from '$lib/stores/pendingChanges';

	let { links = [] }: { links?: any[] } = $props();

	const firstGroup = $derived(links.find(l => l.is_group));
	const imageShape = $derived<'sharp' | 'square' | 'circle'>(firstGroup?.image_shape || 'square');

	function updateImageShape(shape: 'sharp' | 'square' | 'circle') {
		previewStyles.update({ image_shape: shape });
		pendingChanges.updateLinkStyles({ image_shape: shape });
	}
</script>

<div class="space-y-4">
	<h3 class="text-lg font-bold text-gray-900">Image Shape</h3>
	
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
