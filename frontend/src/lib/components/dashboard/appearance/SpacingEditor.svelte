<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';
	import { previewStyles } from '$lib/stores/previewStyles';
	import { pendingChanges } from '$lib/stores/pendingChanges';

	let { links = [] }: { links?: any[] } = $props();

	// Read from globalTheme (now spacing is theme-level)
	const paddingValue = $derived($globalTheme.cardPadding || 16);
	const spacingValue = $derived($globalTheme.cardSpacing || 8);
	
	// Auto-detect if current value is custom (not matching any preset)
	const paddingPresets = [8, 16, 24, 32];
	const spacingPresets = [4, 8, 16];
	const isCustomPadding = $derived(!paddingPresets.includes(paddingValue));
	const isCustomSpacing = $derived(!spacingPresets.includes(spacingValue));
	
	let showCustomPadding = $state(false);
	let showCustomSpacing = $state(false);

	function updatePadding(value: number) {
		console.log('🔧 updatePadding called with value:', value);
		
		// Update theme
		globalTheme.update({ cardPadding: value });
		console.log('✅ Updated globalTheme.cardPadding to:', value);
		
		// Get current style and merge with new padding
		let currentStyle: any = {};
		if ($previewStyles.style) {
			try {
				currentStyle = JSON.parse($previewStyles.style);
				console.log('📦 Current previewStyles.style:', currentStyle);
			} catch (e) {
				console.error('❌ Failed to parse previewStyles.style:', e);
			}
		} else {
			console.log('⚠️ No previewStyles.style found');
		}
		
		const newStyle = {
			...currentStyle,
			padding: { top: value, right: value, bottom: value, left: value }
		};
		console.log('🆕 New style after merge:', newStyle);
		
		const styleStr = JSON.stringify(newStyle);
		previewStyles.update({ style: styleStr });
		console.log('✅ Updated previewStyles.style to:', styleStr);
		
		pendingChanges.updateTheme({ cardPadding: value });
		pendingChanges.updateLinkStyles({ style: styleStr });
	}

	function updateSpacing(value: number) {
		console.log('🔧 updateSpacing called with value:', value);
		
		// Update theme
		globalTheme.update({ cardSpacing: value });
		console.log('✅ Updated globalTheme.cardSpacing to:', value);
		
		// Get current style and merge with new spacing
		let currentStyle: any = {};
		if ($previewStyles.style) {
			try {
				currentStyle = JSON.parse($previewStyles.style);
				console.log('📦 Current previewStyles.style:', currentStyle);
			} catch (e) {
				console.error('❌ Failed to parse previewStyles.style:', e);
			}
		} else {
			console.log('⚠️ No previewStyles.style found');
		}
		
		const newStyle = {
			...currentStyle,
			margin: { top: 0, bottom: value }
		};
		console.log('🆕 New style after merge:', newStyle);
		
		const styleStr = JSON.stringify(newStyle);
		previewStyles.update({ style: styleStr });
		console.log('✅ Updated previewStyles.style to:', styleStr);
		
		pendingChanges.updateTheme({ cardSpacing: value });
		pendingChanges.updateLinkStyles({ style: styleStr });
	}
</script>

<div class="space-y-6">
	<h3 class="text-lg font-bold text-gray-900">Spacing</h3>
	
	<!-- Padding -->
	<div>
		<label class="text-sm font-semibold text-gray-900 mb-3 block">Card padding</label>
		<div class="flex gap-2">
			{#each [
				{ value: 8, label: 'Small' },
				{ value: 16, label: 'Medium' },
				{ value: 24, label: 'Large' },
				{ value: 32, label: 'XL' }
			] as preset}
				<button 
					onclick={() => {
						showCustomPadding = false;
						updatePadding(preset.value);
					}}
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-indigo-600={!showCustomPadding && paddingValue === preset.value}
					class:bg-indigo-50={!showCustomPadding && paddingValue === preset.value}
					class:border-gray-200={showCustomPadding || paddingValue !== preset.value}
				>
					{preset.label}
				</button>
			{/each}
			<button 
				onclick={() => showCustomPadding = !showCustomPadding}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={showCustomPadding || isCustomPadding}
				class:bg-indigo-50={showCustomPadding || isCustomPadding}
				class:border-gray-200={!showCustomPadding && !isCustomPadding}
			>
				Custom
			</button>
		</div>
		
		{#if showCustomPadding}
			<div class="mt-3 p-3 bg-indigo-50 rounded-lg border border-indigo-200">
				<div class="flex items-center justify-between mb-2">
					<span class="text-xs font-medium text-gray-700">Padding</span>
					<span class="text-xs font-mono text-indigo-600">{paddingValue}px</span>
				</div>
				<input 
					type="range" 
					min="0" 
					max="64" 
					value={paddingValue}
					oninput={(e) => updatePadding(parseInt(e.currentTarget.value))}
					class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
				/>
			</div>
		{/if}
	</div>

	<!-- Spacing -->
	<div>
		<label class="text-sm font-semibold text-gray-900 mb-3 block">Card spacing</label>
		<p class="text-xs text-gray-500 mb-3">Space between cards</p>
		<div class="flex gap-2">
			{#each [
				{ value: 4, label: 'Small' },
				{ value: 8, label: 'Medium' },
				{ value: 16, label: 'Large' }
			] as preset}
				<button 
					onclick={() => {
						showCustomSpacing = false;
						updateSpacing(preset.value);
					}}
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-indigo-600={!showCustomSpacing && spacingValue === preset.value}
					class:bg-indigo-50={!showCustomSpacing && spacingValue === preset.value}
					class:border-gray-200={showCustomSpacing || spacingValue !== preset.value}
				>
					{preset.label}
				</button>
			{/each}
			<button 
				onclick={() => showCustomSpacing = !showCustomSpacing}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={showCustomSpacing || isCustomSpacing}
				class:bg-indigo-50={showCustomSpacing || isCustomSpacing}
				class:border-gray-200={!showCustomSpacing && !isCustomSpacing}
			>
				Custom
			</button>
		</div>
		
		{#if showCustomSpacing}
			<div class="mt-3 p-3 bg-indigo-50 rounded-lg border border-indigo-200">
				<div class="flex items-center justify-between mb-2">
					<span class="text-xs font-medium text-gray-700">Spacing</span>
					<span class="text-xs font-mono text-indigo-600">{spacingValue}px</span>
				</div>
				<input 
					type="range" 
					min="0" 
					max="32" 
					value={spacingValue}
					oninput={(e) => updateSpacing(parseInt(e.currentTarget.value))}
					class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
				/>
			</div>
		{/if}
	</div>
</div>
