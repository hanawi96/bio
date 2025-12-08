<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let searchQuery = '';
	export let showFilters = false;
	export let hasActiveFilters = false;
	export let linksCount = 0;

	const dispatch = createEventDispatcher();

	function toggleFilters() {
		dispatch('toggleFilters');
	}

	function selectAll() {
		dispatch('selectAll');
	}

	function addCollection() {
		dispatch('addCollection');
	}

	function viewArchive() {
		dispatch('viewArchive');
	}
</script>

<!-- Compact Professional Toolbar -->
<div class="bg-white border border-gray-200 rounded-xl shadow-sm">
	<div class="px-6 py-3">
		<!-- Main Toolbar Row -->
		<div class="flex items-center gap-3">
			<!-- Search Bar -->
			<div class="flex-1 max-w-md relative">
				<svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
				</svg>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search by title or URL..."
					class="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-gray-50 hover:bg-white transition-colors"
				/>
			</div>

			<!-- Filters Button -->
			<button
				onclick={toggleFilters}
				class="flex items-center gap-2 px-3 py-2 text-sm font-medium border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors {hasActiveFilters ? 'bg-indigo-50 border-indigo-300 text-indigo-700' : 'bg-white text-gray-700'}"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/>
				</svg>
				<span>Filters</span>
				{#if hasActiveFilters}
					<span class="px-1.5 py-0.5 bg-indigo-600 text-white text-xs rounded-full font-semibold">
						{hasActiveFilters}
					</span>
				{/if}
			</button>

			<!-- Divider -->
			<div class="w-px h-6 bg-gray-300"></div>

			<!-- Add Collection -->
			<button
				onclick={addCollection}
				class="flex items-center gap-2 px-3 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
				</svg>
				<span class="hidden sm:inline">Add collection</span>
			</button>

			<!-- Select All -->
			<button
				onclick={selectAll}
				class="flex items-center gap-2 px-3 py-2 text-sm font-medium text-indigo-600 border border-indigo-300 rounded-lg hover:bg-indigo-50 transition-colors"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
				</svg>
				<span class="hidden sm:inline">Select all</span>
			</button>

			<!-- View Archive -->
			<button
				onclick={viewArchive}
				class="ml-auto flex items-center gap-2 px-3 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/>
				</svg>
				<span class="hidden sm:inline">View archive</span>
				<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
				</svg>
			</button>

			<!-- Links Count -->
			<div class="text-sm text-gray-500 font-medium whitespace-nowrap">
				{linksCount} {linksCount === 1 ? 'link' : 'links'}
			</div>
		</div>
	</div>
</div>
