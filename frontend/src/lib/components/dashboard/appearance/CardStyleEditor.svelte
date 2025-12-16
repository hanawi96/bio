<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';
	import { previewStyles } from '$lib/stores/previewStyles';
	import { pendingChanges } from '$lib/stores/pendingChanges';

	const enableBg = $derived($globalTheme.enableCardBackground ?? true);
	const currentBgColor = $derived($globalTheme.cardBackground);
	const currentOpacity = $derived($globalTheme.cardBackgroundOpacity);
	
	let showCustomOpacity = $state(false);

	function updateCardStyles(updates: Record<string, any>) {
		// Separate theme updates (camelCase) from preview/link updates (snake_case)
		const themeUpdates: Record<string, any> = {};
		const previewUpdates: Record<string, any> = {};
		
		for (const [key, value] of Object.entries(updates)) {
			if (key.includes('_')) {
				// snake_case = for preview/links
				previewUpdates[key] = value;
			} else {
				// camelCase = for theme
				themeUpdates[key] = value;
			}
		}
		
		if (Object.keys(themeUpdates).length > 0) {
			globalTheme.update(themeUpdates);
		}
		if (Object.keys(previewUpdates).length > 0) {
			previewStyles.update(previewUpdates);
			pendingChanges.updateLinkStyles(previewUpdates);
		}
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
	<!-- Title with Toggle Switch -->
	<div class="flex items-center justify-between">
		<h3 class="text-lg font-bold text-gray-900">Card Background</h3>
		<label class="relative inline-flex items-center cursor-pointer">
			<input
				type="checkbox"
				checked={enableBg}
				onchange={(e) => updateCardStyles({ enableCardBackground: e.currentTarget.checked })}
				class="sr-only peer"
			/>
			<div class="w-11 h-6 bg-gray-300 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
		</label>
	</div>

	{#if enableBg}
		<!-- Background Color -->
		<div>
			<label class="text-sm font-semibold text-gray-900 mb-3 block">Background color</label>
			<div class="flex gap-3 items-center">
				<input 
					type="color" 
					value={currentBgColor}
					onchange={(e) => updateCardStyles({ cardBackground: e.currentTarget.value, card_background_color: e.currentTarget.value })}
					class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer" 
				/>
				<span class="text-xs text-gray-600 font-mono">{currentBgColor}</span>
				
				<div class="flex gap-1.5 flex-wrap flex-1">
					{#each colorPresets as preset}
						<button 
							onclick={() => updateCardStyles({ cardBackground: preset.color, card_background_color: preset.color })}
							class="w-6 h-6 rounded-full border-2 hover:scale-110 transition-transform" 
							style="background-color: {preset.color};"
							title={preset.label}
						></button>
					{/each}
				</div>
			</div>
		</div>

		<!-- Background Opacity -->
		<div>
			<label class="text-sm font-semibold text-gray-900 mb-3 block">Background opacity</label>
			<div class="flex gap-2">
				{#each [25, 50, 75, 100] as opacity}
					<button 
						onclick={() => {
							showCustomOpacity = false;
							updateCardStyles({ cardBackgroundOpacity: opacity, card_background_opacity: opacity });
						}}
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-indigo-600={!showCustomOpacity && currentOpacity === opacity}
						class:bg-indigo-50={!showCustomOpacity && currentOpacity === opacity}
						class:border-gray-200={showCustomOpacity || currentOpacity !== opacity}
					>
						{opacity}%
					</button>
				{/each}
				<button 
					onclick={() => showCustomOpacity = !showCustomOpacity}
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-indigo-600={showCustomOpacity}
					class:bg-indigo-50={showCustomOpacity}
					class:border-gray-200={!showCustomOpacity}
				>
					Custom
				</button>
			</div>
			
			{#if showCustomOpacity}
				<div class="mt-3 p-3 bg-indigo-50 rounded-lg border border-indigo-200">
					<div class="flex items-center justify-between mb-2">
						<span class="text-xs font-medium text-gray-700">Opacity</span>
						<span class="text-xs font-mono text-indigo-600">{currentOpacity}%</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="100" 
						value={currentOpacity}
						oninput={(e) => {
							const val = parseInt(e.currentTarget.value);
							updateCardStyles({ cardBackgroundOpacity: val, card_background_opacity: val });
						}}
						class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
					/>
				</div>
			{/if}
		</div>
	{:else}
		<p class="text-sm text-gray-500 italic">Background is disabled. Only border will be visible.</p>
	{/if}
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
