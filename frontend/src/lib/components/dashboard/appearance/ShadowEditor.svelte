<script lang="ts">
	import { linksApi } from '$lib/api/links';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { previewStyles } from '$lib/stores/previewStyles';

	let { onUpdate, links = [] }: { onUpdate?: () => Promise<void>; links?: any[] } = $props();

	const shadowPresets = {
		subtle: { x: 0, y: 2, blur: 4 },
		medium: { x: 0, y: 4, blur: 10 },
		strong: { x: 0, y: 8, blur: 20 }
	};

	// Sync from first group in links
	const firstGroup = $derived(links.find(l => l.is_group));
	const showShadow = $derived(firstGroup?.show_shadow ?? false);
	const shadowX = $derived(firstGroup?.shadow_x ?? 0);
	const shadowY = $derived(firstGroup?.shadow_y ?? 4);
	const shadowBlur = $derived(firstGroup?.shadow_blur ?? 10);
	
	// Detect which preset matches current shadow values
	const currentPreset = $derived<'subtle' | 'medium' | 'strong' | 'custom'>(
		shadowX === 0 && shadowY === 2 && shadowBlur === 4 ? 'subtle' :
		shadowX === 0 && shadowY === 4 && shadowBlur === 10 ? 'medium' :
		shadowX === 0 && shadowY === 8 && shadowBlur === 20 ? 'strong' :
		'custom'
	);
	
	let showCustomShadow = $state(false);

	// Log for debugging
	$effect(() => {
		console.log('🔍 [Shadow] Values:', { shadowX, shadowY, shadowBlur, currentPreset, showShadow });
	});

	async function updateShadow(enabled: boolean, preset?: 'subtle' | 'medium' | 'strong') {
		let x = shadowX, y = shadowY, blur = shadowBlur;
		if (preset) {
			const shadow = shadowPresets[preset];
			x = shadow.x;
			y = shadow.y;
			blur = shadow.blur;
		}
		
		console.log('🔧 [Shadow] Updating:', { enabled, preset, x, y, blur });
		previewStyles.update({ show_shadow: enabled, shadow_x: x, shadow_y: y, shadow_blur: blur });
		
		try {
			await linksApi.updateAllGroupStyles({
				show_shadow: enabled,
				shadow_x: x,
				shadow_y: y,
				shadow_blur: blur
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}

	async function updateCustomShadow(x: number, y: number, blur: number) {
		previewStyles.update({ show_shadow: true, shadow_x: x, shadow_y: y, shadow_blur: blur });
		
		try {
			await linksApi.updateAllGroupStyles({
				show_shadow: true,
				shadow_x: x,
				shadow_y: y,
				shadow_blur: blur
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
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
					showCustomShadow = false;
					updateShadow(true, 'subtle');
				}}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={!showCustomShadow && currentPreset === 'subtle'}
				class:bg-indigo-50={!showCustomShadow && currentPreset === 'subtle'}
				class:border-gray-200={showCustomShadow || currentPreset !== 'subtle'}
			>
				Subtle
			</button>
			<button 
				onclick={() => {
					showCustomShadow = false;
					updateShadow(true, 'medium');
				}}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={!showCustomShadow && currentPreset === 'medium'}
				class:bg-indigo-50={!showCustomShadow && currentPreset === 'medium'}
				class:border-gray-200={showCustomShadow || currentPreset !== 'medium'}
			>
				Medium
			</button>
			<button 
				onclick={() => {
					showCustomShadow = false;
					updateShadow(true, 'strong');
				}}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={!showCustomShadow && currentPreset === 'strong'}
				class:bg-indigo-50={!showCustomShadow && currentPreset === 'strong'}
				class:border-gray-200={showCustomShadow || currentPreset !== 'strong'}
			>
				Strong
			</button>
			<button 
				onclick={() => showCustomShadow = !showCustomShadow}
				class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
				class:border-indigo-600={currentPreset === 'custom' || showCustomShadow}
				class:bg-indigo-50={currentPreset === 'custom' || showCustomShadow}
				class:border-gray-200={currentPreset !== 'custom' && !showCustomShadow}
			>
				Custom
			</button>
		</div>

		{#if showCustomShadow}
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
