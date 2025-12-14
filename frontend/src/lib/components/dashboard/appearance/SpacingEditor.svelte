<script lang="ts">
	import { linksApi } from '$lib/api/links';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { previewStyles } from '$lib/stores/previewStyles';

	let { onUpdate, links = [] }: { onUpdate?: () => Promise<void>; links?: any[] } = $props();

	// Sync from first group
	const firstGroup = $derived(links.find(l => l.is_group));
	const getPaddingFromStyle = () => {
		if (!firstGroup?.style) return 16;
		try {
			const style = JSON.parse(firstGroup.style);
			return typeof style.padding === 'number' ? style.padding : style.padding?.top || 16;
		} catch { return 16; }
	};
	const getSpacingFromStyle = () => {
		if (!firstGroup?.style) return 8;
		try {
			const style = JSON.parse(firstGroup.style);
			return style.margin?.bottom || 8;
		} catch { return 8; }
	};
	const paddingValue = $derived(getPaddingFromStyle());
	const spacingValue = $derived(getSpacingFromStyle());
	
	let showCustomPadding = $state(false);
	let showCustomSpacing = $state(false);

	async function updatePadding(value: number) {
		const styleStr = JSON.stringify({
			padding: { top: value, right: value, bottom: value, left: value }
		});
		previewStyles.update({ style: styleStr });
		
		try {
			await linksApi.updateAllGroupStyles({
				style: styleStr
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
	}

	async function updateSpacing(value: number) {
		const styleStr = JSON.stringify({
			margin: { top: 0, bottom: value }
		});
		previewStyles.update({ style: styleStr });
		
		try {
			await linksApi.updateAllGroupStyles({
				style: styleStr
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		}
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
				class:border-indigo-600={showCustomPadding}
				class:bg-indigo-50={showCustomPadding}
				class:border-gray-200={!showCustomPadding}
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
				class:border-indigo-600={showCustomSpacing}
				class:bg-indigo-50={showCustomSpacing}
				class:border-gray-200={!showCustomSpacing}
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
