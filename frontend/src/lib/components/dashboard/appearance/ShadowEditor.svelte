<script lang="ts">
	import { previewStyles } from '$lib/stores/previewStyles';
	import { globalTheme } from '$lib/stores/theme';
	import { pendingChanges } from '$lib/stores/pendingChanges';

	const shadowPresets = {
		subtle: { x: 0, y: 2, blur: 4 },
		medium: { x: 0, y: 4, blur: 10 },
		strong: { x: 0, y: 8, blur: 20 }
	};

	// Read from previewStyles first, fallback to globalTheme (NOT from links!)
	const showShadow = $derived($previewStyles.show_shadow ?? $globalTheme.cardShadow ?? false);
	const shadowX = $derived(Number($previewStyles.shadow_x ?? $globalTheme.cardShadowX ?? 0));
	const shadowY = $derived(Number($previewStyles.shadow_y ?? $globalTheme.cardShadowY ?? 4));
	const shadowBlur = $derived(Number($previewStyles.shadow_blur ?? $globalTheme.cardShadowBlur ?? 10));
	
	// Detect which preset matches current shadow values (flexible comparison)
	const currentPreset = $derived<'subtle' | 'medium' | 'strong' | 'custom'>(
		Math.round(shadowX) === 0 && Math.round(shadowY) === 2 && Math.round(shadowBlur) === 4 ? 'subtle' :
		Math.round(shadowX) === 0 && Math.round(shadowY) === 4 && Math.round(shadowBlur) === 10 ? 'medium' :
		Math.round(shadowX) === 0 && Math.round(shadowY) === 8 && Math.round(shadowBlur) === 20 ? 'strong' :
		'custom'
	);
	
	let showCustomPanel = $state(false);

	function updateShadow(enabled: boolean, preset?: 'subtle' | 'medium' | 'strong') {
		let x = shadowX, y = shadowY, blur = shadowBlur;
		if (preset) {
			const shadow = shadowPresets[preset];
			x = shadow.x;
			y = shadow.y;
			blur = shadow.blur;
		}
		
		// Update theme config (single source of truth)
		globalTheme.update({ 
			cardShadow: enabled, 
			cardShadowX: x, 
			cardShadowY: y, 
			cardShadowBlur: blur 
		});
		
		// Update preview
		previewStyles.update({ show_shadow: enabled, shadow_x: x, shadow_y: y, shadow_blur: blur });
		
		// Mark as pending
		pendingChanges.updateTheme({ 
			cardShadow: enabled, 
			cardShadowX: x, 
			cardShadowY: y, 
			cardShadowBlur: blur 
		});
	}

	function updateCustomShadow(x: number, y: number, blur: number) {
		// Update theme config
		globalTheme.update({ 
			cardShadow: true, 
			cardShadowX: x, 
			cardShadowY: y, 
			cardShadowBlur: blur 
		});
		
		// Update preview
		previewStyles.update({ show_shadow: true, shadow_x: x, shadow_y: y, shadow_blur: blur });
		
		// Mark as pending
		pendingChanges.updateTheme({ 
			cardShadow: true, 
			cardShadowX: x, 
			cardShadowY: y, 
			cardShadowBlur: blur 
		});
	}
</script>

<div class="space-y-4">
	<div class="flex items-center gap-3">
		<h3 class="text-lg font-bold text-gray-900">Link Shadow</h3>
		<button
			onclick={() => {
				const newValue = !showShadow;
				updateShadow(newValue, newValue ? 'medium' : undefined);
			}}
			class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors"
			class:bg-gray-900={showShadow}
			class:bg-gray-300={!showShadow}
		>
			<span
				class="inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform"
				class:translate-x-5={showShadow}
				class:translate-x-0.5={!showShadow}
			></span>
		</button>
	</div>
	
	{#if showShadow}
		<div class="flex gap-2">
			<button 
				onclick={() => {
					showCustomPanel = false;
					updateShadow(true, 'subtle');
				}}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={currentPreset === 'subtle'}
				class:bg-indigo-50={currentPreset === 'subtle'}
				class:border-gray-200={currentPreset !== 'subtle'}
			>
				Subtle
			</button>
			<button 
				onclick={() => {
					showCustomPanel = false;
					updateShadow(true, 'medium');
				}}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={currentPreset === 'medium'}
				class:bg-indigo-50={currentPreset === 'medium'}
				class:border-gray-200={currentPreset !== 'medium'}
			>
				Medium
			</button>
			<button 
				onclick={() => {
					showCustomPanel = false;
					updateShadow(true, 'strong');
				}}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={currentPreset === 'strong'}
				class:bg-indigo-50={currentPreset === 'strong'}
				class:border-gray-200={currentPreset !== 'strong'}
			>
				Strong
			</button>
			<button 
				onclick={() => showCustomPanel = !showCustomPanel}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={showCustomPanel || currentPreset === 'custom'}
				class:bg-indigo-50={showCustomPanel || currentPreset === 'custom'}
				class:border-gray-200={!showCustomPanel && currentPreset !== 'custom'}
			>
				Custom
			</button>
		</div>

		{#if showCustomPanel}
			<div class="p-3 bg-indigo-50 rounded-lg border border-indigo-200 space-y-3">
				<!-- Shadow X -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-xs font-medium text-gray-700">Horizontal</span>
						<span class="text-xs font-mono text-indigo-600">{shadowX}px</span>
					</div>
					<input 
						type="range" 
						min="-20" 
						max="20" 
						value={shadowX}
						oninput={(e) => updateCustomShadow(parseInt(e.currentTarget.value), shadowY, shadowBlur)}
						class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
					/>
				</div>
				
				<!-- Shadow Y -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-xs font-medium text-gray-700">Vertical</span>
						<span class="text-xs font-mono text-indigo-600">{shadowY}px</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="20" 
						value={shadowY}
						oninput={(e) => updateCustomShadow(shadowX, parseInt(e.currentTarget.value), shadowBlur)}
						class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
					/>
				</div>
				
				<!-- Shadow Blur -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-xs font-medium text-gray-700">Blur</span>
						<span class="text-xs font-mono text-indigo-600">{shadowBlur}px</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="40" 
						value={shadowBlur}
						oninput={(e) => updateCustomShadow(shadowX, shadowY, parseInt(e.currentTarget.value))}
						class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
					/>
				</div>
			</div>
		{/if}
	{/if}
</div>
