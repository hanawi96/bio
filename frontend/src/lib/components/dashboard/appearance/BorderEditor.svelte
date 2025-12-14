<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { previewStyles } from '$lib/stores/previewStyles';

	let { onUpdate, links = [] }: { onUpdate?: () => Promise<void>; links?: any[] } = $props();

	const currentBorderRadius = $derived($globalTheme.cardBorderRadius);
	
	// Sync from first group
	const firstGroup = $derived(links.find(l => l.is_group));
	const hasBorder = $derived(firstGroup?.has_card_border ?? false);
	const borderColor = $derived(firstGroup?.card_border_color || '#e5e7eb');
	const borderWidth = $derived(firstGroup?.card_border_width || 1);
	
	let showCustomRadius = $state(false);

	async function updateBorderRadius(radius: number) {
		globalTheme.update({ cardBorderRadius: radius });
		previewStyles.update({ card_border_radius: radius });
		
		try {
			const updatedTheme = globalTheme.getCurrent();
			await profileApi.updateProfile({ 
				theme_config: JSON.stringify(updatedTheme) 
			}, get(auth).token!);
			
			await linksApi.updateAllGroupStyles({
				card_border_radius: radius
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}

	async function updateBorder(enabled: boolean) {
		previewStyles.update({ has_card_border: enabled, card_border_color: borderColor, card_border_width: borderWidth });
		
		try {
			await linksApi.updateAllGroupStyles({
				has_card_border: enabled,
				card_border_color: borderColor,
				card_border_width: borderWidth
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}

	async function updateBorderColor(color: string) {
		previewStyles.update({ card_border_color: color });
		
		try {
			await linksApi.updateAllGroupStyles({
				has_card_border: true,
				card_border_color: color
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}

	async function updateBorderWidth(width: number) {
		previewStyles.update({ card_border_width: width });
		
		try {
			await linksApi.updateAllGroupStyles({
				has_card_border: true,
				card_border_width: width
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}

	const colorPresets = [
		{ color: '#e5e7eb', label: 'Light Gray' },
		{ color: '#9ca3af', label: 'Gray' },
		{ color: '#000000', label: 'Black' },
		{ color: '#3b82f6', label: 'Blue' },
		{ color: '#8b5cf6', label: 'Purple' },
		{ color: '#10b981', label: 'Green' },
		{ color: '#ef4444', label: 'Red' },
		{ color: '#f59e0b', label: 'Amber' }
	];
</script>

<div class="space-y-6">
	<h3 class="text-lg font-bold text-gray-900">Border & Radius</h3>
	
	<!-- Border Radius -->
	<div>
		<label class="text-sm font-semibold text-gray-900 mb-3 block">Border radius</label>
		<div class="flex gap-2">
			{#each [
				{ value: 0, label: 'None' },
				{ value: 8, label: 'Small' },
				{ value: 12, label: 'Medium' },
				{ value: 24, label: 'Large' }
			] as preset}
				<button 
					onclick={() => {
						showCustomRadius = false;
						updateBorderRadius(preset.value);
					}}
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-indigo-600={!showCustomRadius && currentBorderRadius === preset.value}
					class:bg-indigo-50={!showCustomRadius && currentBorderRadius === preset.value}
					class:border-gray-200={showCustomRadius || currentBorderRadius !== preset.value}
				>
					{preset.label}
				</button>
			{/each}
			<button 
				onclick={() => showCustomRadius = !showCustomRadius}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={showCustomRadius}
				class:bg-indigo-50={showCustomRadius}
				class:border-gray-200={!showCustomRadius}
			>
				Custom
			</button>
		</div>
		
		{#if showCustomRadius}
			<div class="mt-3 p-3 bg-indigo-50 rounded-lg border border-indigo-200">
				<div class="flex items-center justify-between mb-2">
					<span class="text-xs font-medium text-gray-700">Radius</span>
					<span class="text-xs font-mono text-indigo-600">{currentBorderRadius}px</span>
				</div>
				<input 
					type="range" 
					min="0" 
					max="32" 
					value={currentBorderRadius}
					oninput={(e) => updateBorderRadius(parseInt(e.currentTarget.value))}
					class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
				/>
			</div>
		{/if}
	</div>

	<!-- Card Border -->
	<div>
		<div class="flex items-center gap-3 mb-3">
			<label class="text-sm font-semibold text-gray-900">Card border</label>
			<button
				onclick={() => updateBorder(!hasBorder)}
				class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors"
				class:bg-gray-900={hasBorder}
				class:bg-gray-300={!hasBorder}
			>
				<span
					class="inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform"
					class:translate-x-5={hasBorder}
					class:translate-x-0.5={!hasBorder}
				></span>
			</button>
		</div>

		{#if hasBorder}
			<div class="space-y-4">
				<!-- Border Color -->
				<div>
					<label class="text-xs font-medium text-gray-700 mb-2 block">Border color</label>
					<div class="flex gap-3 items-center">
						<input 
							type="color" 
							value={borderColor}
							onchange={(e) => updateBorderColor(e.currentTarget.value)}
							class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer" 
						/>
						<span class="text-xs text-gray-600 font-mono">{borderColor}</span>
						
						<div class="flex gap-1.5 flex-wrap flex-1">
							{#each colorPresets as preset}
								<button 
									onclick={() => updateBorderColor(preset.color)}
									class="w-6 h-6 rounded-full border-2 hover:scale-110 transition-transform" 
									style="background-color: {preset.color};"
									title={preset.label}
								></button>
							{/each}
						</div>
					</div>
				</div>

				<!-- Border Width -->
				<div>
					<label class="text-xs font-medium text-gray-700 mb-2 block">Border width</label>
					<div class="flex gap-2">
						{#each [1, 2, 3, 4] as width}
							<button 
								onclick={() => updateBorderWidth(width)}
								class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
								class:border-indigo-600={borderWidth === width}
								class:bg-indigo-50={borderWidth === width}
								class:border-gray-200={borderWidth !== width}
							>
								{width}px
							</button>
						{/each}
					</div>
				</div>
			</div>
		{/if}
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
