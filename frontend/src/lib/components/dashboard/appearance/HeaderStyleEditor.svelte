<script lang="ts">
	import { currentHeaderStyle, themePresets, type HeaderStyles } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	
	// Use derived to always sync with store
	const headerStyle = $derived($currentHeaderStyle);
	let selectedCategory = $state<string>('classic');
	let applying = $state(false);
	let manuallySelected = $state<string | null>(null);

	const categories = [
		{ id: 'classic', label: 'Classic' },
		{ id: 'modern', label: 'Modern' },
		{ id: 'creative', label: 'Creative' },
		{ id: 'minimal', label: 'Minimal' }
	];

	const headerPresets: Record<string, Array<{ id: string; name: string; layout: HeaderStyles['layout']; preview: { cover: string; avatar: string } }>> = {
		classic: [
			{ id: 'centered', name: 'Centered', layout: 'centered', preview: { cover: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', avatar: '#ffffff' } },
			{ id: 'overlap', name: 'Overlap', layout: 'overlap', preview: { cover: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', avatar: '#ffffff' } },
			{ id: 'side', name: 'Side Profile', layout: 'side', preview: { cover: '#ffffff', avatar: '#6366f1' } }
		],
		modern: [
			{ id: 'card', name: 'Card Float', layout: 'card', preview: { cover: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', avatar: '#ffffff' } },
			{ id: 'glass', name: 'Glass', layout: 'glass', preview: { cover: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)', avatar: 'rgba(255,255,255,0.2)' } },
			{ id: 'gradient', name: 'Gradient', layout: 'gradient', preview: { cover: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)', avatar: '#ffffff' } }
		],
		creative: [
			{ id: 'split', name: 'Split Screen', layout: 'split', preview: { cover: 'linear-gradient(90deg, #ff9a9e 0%, #fecfef 100%)', avatar: '#ffffff' } },
			{ id: 'asymmetric', name: 'Asymmetric', layout: 'asymmetric', preview: { cover: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)', avatar: '#ffffff' } },
			{ id: 'full', name: 'Full Bleed', layout: 'full', preview: { cover: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)', avatar: '#ffffff' } }
		],
		minimal: [
			{ id: 'minimal', name: 'Minimal', layout: 'minimal', preview: { cover: '#f5f5f5', avatar: '#000000' } }
		]
	};

	async function selectLayout(layout: HeaderStyles['layout']) {
		if (applying) return;
		
		const presetMap: Record<HeaderStyles['layout'], Partial<HeaderStyles>> = {
			centered: { layout: 'centered', coverType: 'gradient', coverGradientFrom: '#667eea', coverGradientTo: '#764ba2', coverHeight: 160, avatarSize: 120, avatarBorder: 4, showCover: true, bioAlign: 'center', bioSize: 'md' },
			overlap: { layout: 'overlap', coverType: 'gradient', coverGradientFrom: '#f093fb', coverGradientTo: '#f5576c', coverHeight: 180, avatarSize: 130, avatarBorder: 4, showCover: true, bioAlign: 'center', bioSize: 'md' },
			card: { layout: 'card', coverType: 'gradient', coverGradientFrom: '#4facfe', coverGradientTo: '#00f2fe', coverHeight: 140, avatarSize: 110, avatarBorder: 4, showCover: true, bioAlign: 'center', bioSize: 'md' },
			glass: { layout: 'glass', coverType: 'gradient', coverGradientFrom: '#fa709a', coverGradientTo: '#fee140', coverHeight: 170, avatarSize: 120, avatarBorder: 4, avatarBorderColor: 'rgba(255,255,255,0.2)', showCover: true, bioAlign: 'center', bioSize: 'md' },
			gradient: { layout: 'gradient', coverType: 'gradient', coverGradientFrom: '#30cfd0', coverGradientTo: '#330867', coverHeight: 190, avatarSize: 130, avatarBorder: 5, showCover: true, bioAlign: 'center', bioSize: 'lg' },
			minimal: { layout: 'minimal', coverType: 'color', coverColor: '#f5f5f5', coverHeight: 0, avatarSize: 90, avatarBorder: 2, avatarBorderColor: '#000000', showCover: false, bioAlign: 'left', bioSize: 'sm' },
			full: { layout: 'full', coverType: 'gradient', coverGradientFrom: '#a8edea', coverGradientTo: '#fed6e3', coverHeight: 240, avatarSize: 150, avatarBorder: 6, showCover: true, bioAlign: 'center', bioSize: 'lg' },
			side: { layout: 'side', coverType: 'color', coverColor: '#ffffff', coverHeight: 0, avatarSize: 90, avatarBorder: 2, avatarBorderColor: '#6366f1', showCover: false, bioAlign: 'left', bioSize: 'sm' },
			split: { layout: 'split', coverType: 'gradient', coverGradientFrom: '#ff9a9e', coverGradientTo: '#fecfef', coverHeight: 200, avatarSize: 120, avatarBorder: 4, showCover: true, bioAlign: 'center', bioSize: 'md' },
			asymmetric: { layout: 'asymmetric', coverType: 'gradient', coverGradientFrom: '#ffecd2', coverGradientTo: '#fcb69f', coverHeight: 220, avatarSize: 140, avatarBorder: 5, showCover: true, bioAlign: 'left', bioSize: 'lg' }
		};
		
		const oldStyle = { ...headerStyle };
		const newStyle = { ...headerStyle, ...presetMap[layout], avatarBorderColor: headerStyle.avatarBorderColor || '#ffffff' };
		
		// Mark as manually selected
		manuallySelected = layout;
		
		// Optimistic update
		currentHeaderStyle.set(newStyle);
		
		// Auto-save to backend
		applying = true;
		try {
			await profileApi.updateProfile({ header_config: newStyle }, get(auth).token!);
			toast.success('Header style updated!');
		} catch (e: any) {
			// Rollback on error
			currentHeaderStyle.set(oldStyle);
			manuallySelected = null;
			toast.error(e.message || 'Failed to update header');
		} finally {
			applying = false;
		}
	}
	
	// Reset manual selection when theme changes externally
	$effect(() => {
		if (headerStyle.layout && manuallySelected && headerStyle.layout !== manuallySelected) {
			manuallySelected = null;
		}
	});
</script>

<div class="space-y-6">
	<div>
		<div class="flex items-center justify-between mb-4">
			<h3 class="text-lg font-semibold text-gray-900">Header Style</h3>
			{#if applying}
				<span class="text-sm text-indigo-600 flex items-center gap-2">
					<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Saving...
				</span>
			{/if}
		</div>
		
		<!-- Category Tabs -->
		<div class="flex gap-2 mb-6">
			{#each categories as cat}
				<button
					onclick={() => selectedCategory = cat.id}
					class="px-4 py-2 rounded-full text-sm font-medium transition-all {selectedCategory === cat.id
						? 'bg-gray-900 text-white'
						: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
				>
					{cat.label}
				</button>
			{/each}
		</div>

		<!-- Header Presets Grid -->
		<div class="grid grid-cols-2 md:grid-cols-3 gap-3">
			{#each headerPresets[selectedCategory] || [] as preset}
				<button
					onclick={() => selectLayout(preset.layout)}
					class="group relative"
				>
					<div class="aspect-[3/4] rounded-xl overflow-hidden border-2 transition-all {manuallySelected === preset.layout ? 'border-indigo-600 shadow-lg' : 'border-gray-200 hover:border-gray-300'}"
						style="background: {preset.preview.cover};"
					>
						<div class="h-full flex flex-col items-center justify-center p-3">
							<div class="w-10 h-10 rounded-full" style="border: 2px solid {preset.preview.avatar}; background: rgba(255,255,255,0.3);"></div>
							<div class="mt-2 w-12 h-1.5 rounded" style="background: {preset.preview.avatar}; opacity: 0.7;"></div>
							<div class="mt-1 w-8 h-1 rounded" style="background: {preset.preview.avatar}; opacity: 0.5;"></div>
						</div>
					</div>
					{#if manuallySelected === preset.layout}
						<div class="absolute top-2 right-2 w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center">
							<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
							</svg>
						</div>
					{/if}
					<p class="mt-2 text-xs font-medium text-gray-700 text-center">{preset.name}</p>
				</button>
			{/each}
		</div>
	</div>
</div>
