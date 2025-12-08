<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';

	export let open = false;
	export let link: Link;

	const dispatch = createEventDispatcher();
	
	let activeTab: 'settings' | 'layout' = 'layout';
	let selectedLayout: 'classic' | 'featured' = 'classic';

	function handleLayoutSelect(layoutType: 'classic' | 'featured') {
		selectedLayout = layoutType;
	}

	function handleSave() {
		if (selectedLayout !== link.layout_type) {
			dispatch('updateLayout', { id: link.id, layout_type: selectedLayout });
		}
		open = false;
	}

	// Sync selectedLayout với link.layout_type khi dialog mở hoặc link thay đổi
	$: if (open && link) {
		selectedLayout = link.layout_type || 'classic';
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-3xl p-0 gap-0 overflow-hidden">
		<!-- Header with Tabs -->
		<div class="border-b border-gray-200 bg-gray-50">
			<div class="px-6 pt-6 pb-4">
				<h2 class="text-xl font-bold text-gray-900 mb-4">{link.title}</h2>
				
				<!-- Tabs -->
				<div class="flex gap-8 border-b border-gray-200">
					<button
						onclick={() => activeTab = 'settings'}
						class="pb-3 px-1 text-sm font-medium transition-colors relative"
						class:text-gray-900={activeTab === 'settings'}
						class:text-gray-500={activeTab !== 'settings'}
					>
						Link Settings
						{#if activeTab === 'settings'}
							<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gray-900"></div>
						{/if}
					</button>
					<button
						onclick={() => activeTab = 'layout'}
						class="pb-3 px-1 text-sm font-medium transition-colors relative"
						class:text-gray-900={activeTab === 'layout'}
						class:text-gray-500={activeTab !== 'layout'}
					>
						Layout
						{#if activeTab === 'layout'}
							<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gray-900"></div>
						{/if}
					</button>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="p-6">
			{#if activeTab === 'settings'}
				<!-- Link Settings Tab -->
				<div class="space-y-6">
					<div>
						<h3 class="text-lg font-semibold text-gray-900 mb-4">Link Settings</h3>
						<p class="text-sm text-gray-500">Configure your link properties and behavior</p>
					</div>
					<!-- Add more settings here later -->
				</div>
			{:else}
				<!-- Layout Tab -->
				<div class="space-y-6">
					<div>
						<h3 class="text-lg font-semibold text-gray-900 mb-2">Choose a layout for your link</h3>
						<p class="text-sm text-gray-500">Select how you want this link to appear on your profile</p>
					</div>

					<!-- Layout Options -->
					<div class="space-y-4">
						<!-- Classic Layout Option -->
						<button
							onclick={() => handleLayoutSelect('classic')}
							class="w-full text-left p-6 rounded-2xl border-2 transition-all hover:border-gray-400 hover:shadow-md group"
							class:border-gray-900={selectedLayout === 'classic'}
							class:bg-gray-50={selectedLayout === 'classic'}
							class:border-gray-200={selectedLayout !== 'classic'}
							class:shadow-sm={selectedLayout === 'classic'}
						>
							<div class="flex items-start gap-6">
								<!-- Radio Button -->
								<div class="flex-shrink-0 mt-1">
									<div class="w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all"
										class:border-gray-900={selectedLayout === 'classic'}
										class:bg-gray-900={selectedLayout === 'classic'}
										class:border-gray-300={selectedLayout !== 'classic'}
									>
										{#if selectedLayout === 'classic'}
											<div class="w-2.5 h-2.5 bg-white rounded-full"></div>
										{/if}
									</div>
								</div>

								<!-- Content -->
								<div class="flex-1">
									<h4 class="text-lg font-bold text-gray-900 mb-1">Classic</h4>
									<p class="text-sm text-gray-600">Efficient, direct and compact.</p>
								</div>

								<!-- Preview -->
								<div class="flex-shrink-0">
									<div class="w-64 h-16 bg-gradient-to-br from-teal-600 to-teal-700 rounded-2xl flex items-center px-4 gap-3 shadow-lg transition-transform group-hover:scale-105">
										<div class="w-10 h-10 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
											<svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
												<path d="M10 12a2 2 0 100-4 2 2 0 000 4z"/>
												<path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"/>
											</svg>
										</div>
										<div class="flex-1 text-white">
											<div class="text-sm font-semibold">Link Title</div>
										</div>
										<svg class="w-4 h-4 text-white/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/>
										</svg>
									</div>
								</div>
							</div>
						</button>

						<!-- Featured Layout Option -->
						<button
							onclick={() => handleLayoutSelect('featured')}
							class="w-full text-left p-6 rounded-2xl border-2 transition-all hover:border-gray-400 hover:shadow-md group"
							class:border-gray-900={selectedLayout === 'featured'}
							class:bg-gray-50={selectedLayout === 'featured'}
							class:border-gray-200={selectedLayout !== 'featured'}
							class:shadow-sm={selectedLayout === 'featured'}
						>
							<div class="flex items-start gap-6">
								<!-- Radio Button -->
								<div class="flex-shrink-0 mt-1">
									<div class="w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all"
										class:border-gray-900={selectedLayout === 'featured'}
										class:bg-gray-900={selectedLayout === 'featured'}
										class:border-gray-300={selectedLayout !== 'featured'}
									>
										{#if selectedLayout === 'featured'}
											<div class="w-2.5 h-2.5 bg-white rounded-full"></div>
										{/if}
									</div>
								</div>

								<!-- Content -->
								<div class="flex-1">
									<h4 class="text-lg font-bold text-gray-900 mb-1">Featured</h4>
									<p class="text-sm text-gray-600">Make your link stand out with a larger, more attractive display.</p>
								</div>

								<!-- Preview -->
								<div class="flex-shrink-0">
									<div class="w-64 h-32 bg-gradient-to-br from-orange-400 via-orange-500 to-orange-600 rounded-2xl overflow-hidden shadow-lg relative transition-transform group-hover:scale-105">
										<!-- Decorative pattern -->
										<div class="absolute inset-0 opacity-10">
											<div class="absolute top-4 right-4 w-20 h-20 bg-white rounded-full blur-2xl"></div>
											<div class="absolute bottom-4 left-4 w-16 h-16 bg-white rounded-full blur-xl"></div>
										</div>
										<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent"></div>
										<div class="absolute bottom-0 left-0 right-0 p-4 flex items-end justify-between">
											<div class="text-white">
												<div class="text-sm font-bold mb-0.5">Now touring, get your tickets</div>
												<div class="text-xs opacity-80">Featured Link</div>
											</div>
											<svg class="w-4 h-4 text-white/80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/>
											</svg>
										</div>
									</div>
								</div>
							</div>
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="border-t border-gray-200 px-6 py-4 bg-gray-50 flex justify-end gap-3">
			<button
				onclick={() => open = false}
				class="px-6 py-2.5 text-sm font-medium text-gray-700 hover:text-gray-900 transition-colors"
			>
				Cancel
			</button>
			<button
				onclick={handleSave}
				class="px-6 py-2.5 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-xl transition-all shadow-sm hover:shadow-md"
			>
				Save Changes
			</button>
		</div>
	</Dialog.Content>
</Dialog.Root>

<style>
	button:focus {
		outline: none;
	}
</style>
