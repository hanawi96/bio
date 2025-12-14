<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	
	let { onUpdate }: { onUpdate?: () => Promise<void> } = $props();

	let activeTab = $state<'popular' | 'custom'>('popular');
	let saving = $state(false);
	let customMode = $state(false);
	let customValue = $state(12);

	// Border radius presets (removed Full Width)
	const borderRadiusPresets = [
		{ id: 'rounded', label: 'Rounded', value: 12 },
		{ id: 'oval', label: 'Oval', value: 24 },
		{ id: 'square', label: 'Square', value: 0 }
	];
	
	// Get current values directly from store - reactive
	const currentBorderRadius = $derived($globalTheme.cardBorderRadius);
	const showShadow = $derived($globalTheme.cardShadow);
	const showOutline = $derived($globalTheme.cardBorder);
	
	// Check if current radius is custom (doesn't match any preset)
	const isCustomRadius = $derived(() => {
		return !borderRadiusPresets.some(p => p.value === currentBorderRadius);
	});

	async function updateBorderRadius(value: number, isCustom = false) {
		// Update theme immediately for instant UI
		globalTheme.update({ cardBorderRadius: value });
		
		// If custom, update custom value
		if (isCustom) {
			customValue = value;
		}
		
		// Save in background
		saving = true;
		try {
			const updatedTheme = globalTheme.getCurrent();
			const themeConfig = JSON.stringify(updatedTheme);
			
			await profileApi.updateProfile({ theme_config: themeConfig }, get(auth).token!);
			await linksApi.updateAllGroupStyles({
				card_border_radius: value
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		} finally {
			saving = false;
		}
	}
	
	function selectPreset(value: number) {
		customMode = false;
		updateBorderRadius(value);
	}
	
	function enableCustomMode() {
		if (!customMode) {
			customMode = true;
			customValue = currentBorderRadius;
		}
	}
	
	function handleCustomSlider(e: Event) {
		const value = parseInt((e.target as HTMLInputElement).value);
		customValue = value;
		updateBorderRadius(value, true);
	}

	async function toggleShadow() {
		const newValue = !showShadow;
		saving = true;
		try {
			globalTheme.update({ cardShadow: newValue });
			
			const updatedTheme = globalTheme.getCurrent();
			const themeConfig = JSON.stringify(updatedTheme);
			await profileApi.updateProfile({ theme_config: themeConfig }, get(auth).token!);
			
			await linksApi.updateAllGroupStyles({
				show_shadow: newValue
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		} finally {
			saving = false;
		}
	}

	async function toggleOutline() {
		const newValue = !showOutline;
		saving = true;
		try {
			globalTheme.update({ cardBorder: newValue });
			
			const updatedTheme = globalTheme.getCurrent();
			const themeConfig = JSON.stringify(updatedTheme);
			await profileApi.updateProfile({ theme_config: themeConfig }, get(auth).token!);
			
			await linksApi.updateAllGroupStyles({
				has_card_border: newValue
			}, get(auth).token!);
			
			if (onUpdate) await onUpdate();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update');
		} finally {
			saving = false;
		}
	}
</script>

<div class="space-y-6">
	<h3 class="text-lg font-bold text-gray-900">Block style</h3>

	<!-- Tabs -->
	<div class="flex gap-2">
		<button
			onclick={() => activeTab = 'popular'}
			class="px-6 py-2 rounded-full text-sm font-medium transition-all {activeTab === 'popular'
				? 'bg-gray-900 text-white'
				: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
		>
			Popular
		</button>
		<button
			onclick={() => activeTab = 'custom'}
			class="px-6 py-2 rounded-full text-sm font-medium transition-all {activeTab === 'custom'
				? 'bg-gray-900 text-white'
				: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
		>
			Custom
		</button>
	</div>

	{#if activeTab === 'popular'}
		<!-- Border Radius Presets -->
		<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
			{#each borderRadiusPresets as preset}
				{@const isActive = currentBorderRadius === preset.value && !customMode}
				<button
					onclick={() => selectPreset(preset.value)}
					class="relative group"
				>
					<div
						class="aspect-[4/3] rounded-xl border-2 transition-all flex items-center justify-center {isActive
							? 'border-indigo-600 bg-indigo-50'
							: 'border-gray-200 hover:border-gray-300 bg-white'}"
					>
						<div
							class="w-24 h-12 bg-gray-400 transition-all"
							style="border-radius: {preset.value}px;"
						></div>
					</div>
					<p class="mt-2 text-sm font-medium text-gray-700 text-center">{preset.label}</p>
					{#if isActive}
						<div class="absolute top-2 right-2 w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center">
							<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
							</svg>
						</div>
					{/if}
				</button>
			{/each}
			
			<!-- Custom Preset -->
			<button
				onclick={enableCustomMode}
				class="relative group"
			>
				<div
					class="aspect-[4/3] rounded-xl border-2 transition-all flex items-center justify-center {customMode
						? 'border-indigo-600 bg-indigo-50'
						: 'border-gray-200 hover:border-gray-300 bg-white'}"
				>
					<div class="text-center">
						<svg class="w-8 h-8 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
						</svg>
					</div>
				</div>
				<p class="mt-2 text-sm font-medium text-gray-700 text-center">Custom</p>
				{#if customMode}
					<div class="absolute top-2 right-2 w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center">
						<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
					</div>
				{/if}
			</button>
		</div>
		
		<!-- Custom Slider (shown when custom mode is active) -->
		{#if customMode}
			<div class="mt-6 p-4 bg-indigo-50 rounded-xl border-2 border-indigo-200">
				<div class="flex items-center justify-between mb-3">
					<label class="text-sm font-semibold text-gray-900">Border Radius</label>
					<span class="text-sm font-mono font-bold text-indigo-600">{customValue}px</span>
				</div>
				<input
					type="range"
					min="0"
					max="50"
					value={customValue}
					oninput={handleCustomSlider}
					class="w-full h-2 bg-indigo-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
				/>
				<div class="flex justify-between mt-2 text-xs text-gray-500">
					<span>0px</span>
					<span>25px</span>
					<span>50px</span>
				</div>
			</div>
		{/if}

		<!-- Toggle Options -->
		<div class="space-y-3 pt-4">
			<!-- Block Shadow -->
			<div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors">
				<div class="flex items-center gap-3">
					<label class="relative inline-flex items-center cursor-pointer">
						<input
							type="checkbox"
							checked={showShadow}
							onchange={toggleShadow}
							class="sr-only peer"
						/>
						<div class="w-11 h-6 bg-gray-300 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
					</label>
					<span class="text-sm font-medium text-gray-900">Block shadow</span>
				</div>
			</div>

			<!-- Block Outline -->
			<div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors">
				<div class="flex items-center gap-3">
					<label class="relative inline-flex items-center cursor-pointer">
						<input
							type="checkbox"
							checked={showOutline}
							onchange={toggleOutline}
							class="sr-only peer"
						/>
						<div class="w-11 h-6 bg-gray-300 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
					</label>
					<span class="text-sm font-medium text-gray-900">Block outline</span>
				</div>
			</div>
		</div>
	{:else}
		<div class="text-center py-12 text-gray-500">
			<p class="text-sm">Custom styles coming soon...</p>
		</div>
	{/if}
</div>
