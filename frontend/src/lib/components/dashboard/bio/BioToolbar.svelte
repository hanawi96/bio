<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let searchQuery = '';
	export let showFilters = false;
	export let showCalendar = false;
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

	function toggleCalendar() {
		dispatch('toggleCalendar');
	}
</script>

<!-- Modern Flat Toolbar -->
<div class="flex items-center gap-3 py-4">
	<!-- Search Bar - Minimalist -->
	<div class="flex-1 max-w-sm relative group">
		<svg class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 group-focus-within:text-indigo-500 transition-colors pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
		</svg>
		<input
			type="text"
			bind:value={searchQuery}
			placeholder="Search by title or URL..."
			class="w-full pl-10 pr-4 py-2.5 text-sm bg-gray-50 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white transition-all placeholder:text-gray-400"
		/>
	</div>

	<!-- Filters Button - Minimal -->
	<button
		onclick={toggleFilters}
		class="relative flex items-center gap-2 px-4 py-2.5 text-sm font-medium rounded-xl transition-all {hasActiveFilters ? 'bg-indigo-50 text-indigo-600' : 'text-gray-600 hover:bg-gray-50'}"
	>
		<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/>
		</svg>
		<span>Filters</span>
		{#if hasActiveFilters}
			<span class="absolute -top-1 -right-1 w-5 h-5 bg-indigo-600 text-white text-xs rounded-full flex items-center justify-center font-semibold">
				{hasActiveFilters}
			</span>
		{/if}
	</button>

	<!-- Calendar Button -->
	<button
		onclick={toggleCalendar}
		class="flex items-center gap-2 px-4 py-2.5 text-sm font-medium rounded-xl transition-all {showCalendar ? 'bg-purple-50 text-purple-600' : 'text-gray-600 hover:bg-gray-50'}"
	>
		<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
		</svg>
		<span>Calendar</span>
	</button>

	<!-- Action Buttons - Clean Icons -->
	<div class="flex items-center gap-1.5 ml-auto">
		<!-- Select All -->
		<button
			onclick={selectAll}
			class="p-2.5 text-gray-600 hover:bg-gray-50 rounded-xl transition-all"
			title="Select all"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
			</svg>
		</button>

		<!-- Add Collection -->
		<button
			onclick={addCollection}
			class="p-2.5 text-gray-600 hover:bg-gray-50 rounded-xl transition-all"
			title="Add collection"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
			</svg>
		</button>

		<!-- View Calendar -->
		<button
			onclick={viewCalendar}
			class="p-2.5 text-gray-600 hover:bg-gray-50 rounded-xl transition-all"
			title="View calendar"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
			</svg>
		</button>

		<!-- View Archive -->
		<button
			onclick={viewArchive}
			class="p-2.5 text-gray-600 hover:bg-gray-50 rounded-xl transition-all"
			title="View archive"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/>
			</svg>
		</button>

		<!-- Divider -->
		<div class="w-px h-5 bg-gray-200 mx-1"></div>

		<!-- Links Count -->
		<div class="px-3 py-1.5 text-xs font-medium text-gray-500 bg-gray-50 rounded-lg">
			{linksCount} {linksCount === 1 ? 'link' : 'links'}
		</div>
	</div>
</div>
