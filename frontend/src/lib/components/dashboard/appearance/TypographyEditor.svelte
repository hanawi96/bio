<script lang="ts">
	import { previewStyles } from '$lib/stores/previewStyles';
	import { globalTheme } from '$lib/stores/theme';
	import { pendingChanges } from '$lib/stores/pendingChanges';

	let { links = [] }: { links?: any[] } = $props();

	// Read from globalTheme (single source of truth)
	const textAlignment = $derived($globalTheme.textAlignment || 'center');
	const textSize = $derived($globalTheme.textSize || 'M');
	const currentTextColor = $derived($globalTheme.cardTextColor);

	function updateTypography(field: string, value: any) {
		console.log('📝 TypographyEditor - updateTypography:', {
			field,
			value,
			'current textAlignment': textAlignment
		});
		
		// Update theme config (single source of truth)
		if (field === 'text_alignment') {
			globalTheme.update({ textAlignment: value });
			console.log('✅ Updated globalTheme.textAlignment to:', value);
		} else if (field === 'text_size') {
			globalTheme.update({ textSize: value });
		}
		
		// Update preview
		previewStyles.update({ [field]: value });
		console.log('✅ Updated previewStyles:', { [field]: value });
	}

	function updateTextColor(color: string) {
		globalTheme.update({ cardTextColor: color });
		previewStyles.update({ card_text_color: color });
		pendingChanges.updateLinkStyles({ card_text_color: color });
	}

	const colorPresets = [
		{ color: '#ffffff', label: 'White' },
		{ color: '#f3f4f6', label: 'Light Gray' },
		{ color: '#000000', label: 'Black' },
		{ color: '#3b82f6', label: 'Blue' },
		{ color: '#8b5cf6', label: 'Purple' },
		{ color: '#10b981', label: 'Green' },
		{ color: '#ef4444', label: 'Red' },
		{ color: '#f59e0b', label: 'Amber' },
		{ color: '#ec4899', label: 'Pink' }
	];
</script>

<div class="space-y-6">
	<h3 class="text-lg font-bold text-gray-900">Typography</h3>
	
	<!-- Text Color -->
	<div>
		<label class="text-sm font-semibold text-gray-900 mb-3 block">Text color</label>
		<div class="flex gap-3 items-center">
			<input 
				type="color" 
				value={currentTextColor}
				onchange={(e) => updateTextColor(e.currentTarget.value)}
				class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer" 
			/>
			<span class="text-xs text-gray-600 font-mono">{currentTextColor}</span>
			
			<div class="flex gap-1.5 flex-wrap flex-1">
				{#each colorPresets as preset}
					<button 
						onclick={() => updateTextColor(preset.color)}
						class="w-6 h-6 rounded-full border-2 hover:scale-110 transition-transform" 
						style="background-color: {preset.color};"
						title={preset.label}
					></button>
				{/each}
			</div>
		</div>
	</div>

	<!-- Text Alignment -->
	<div>
		<div class="flex items-center justify-between mb-3">
			<label class="text-sm font-semibold text-gray-900">Text alignment</label>
			{#if links.filter(l => l.is_group && l.text_alignment !== null && l.text_alignment !== undefined).length > 0}
				<span class="text-xs text-amber-600 font-medium flex items-center gap-1">
					<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
					</svg>
					{links.filter(l => l.is_group && l.text_alignment !== null && l.text_alignment !== undefined).length} group(s) have custom values
				</span>
			{/if}
		</div>
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
		<div class="flex items-center justify-between mb-3">
			<label class="text-sm font-semibold text-gray-900">Text size</label>
			{#if links.filter(l => l.is_group && l.text_size !== null && l.text_size !== undefined).length > 0}
				<span class="text-xs text-amber-600 font-medium flex items-center gap-1">
					<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
					</svg>
					{links.filter(l => l.is_group && l.text_size !== null && l.text_size !== undefined).length} group(s) have custom values
				</span>
			{/if}
		</div>
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

<style>
	input[type="color"] {
		padding: 0;
		background: none;
	}
	
	input[type="color"]::-webkit-color-swatch-wrapper {
		padding: 0;
		border: none;
	}
	
	input[type="color"]::-webkit-color-swatch {
		border: none;
		border-radius: 8px;
	}
</style>
