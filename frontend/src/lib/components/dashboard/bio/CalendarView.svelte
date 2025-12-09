<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';

	export let links: Link[] = [];

	const dispatch = createEventDispatcher();

	let currentDate = new Date();
	let currentMonth = currentDate.getMonth();
	let currentYear = currentDate.getFullYear();

	$: scheduledLinks = links.filter(l => l.scheduled_at || l.expires_at);
	$: calendarDays = generateCalendar(currentYear, currentMonth);

	function generateCalendar(year: number, month: number) {
		const firstDay = new Date(year, month, 1);
		const lastDay = new Date(year, month + 1, 0);
		const daysInMonth = lastDay.getDate();
		const startingDayOfWeek = firstDay.getDay();

		const days = [];
		
		// Previous month padding
		for (let i = 0; i < startingDayOfWeek; i++) {
			days.push({ date: null, isCurrentMonth: false });
		}

		// Current month days
		for (let i = 1; i <= daysInMonth; i++) {
			days.push({ date: i, isCurrentMonth: true });
		}

		return days;
	}

	function getLinksForDate(day: number | null) {
		if (!day) return [];
		
		const dateStr = new Date(currentYear, currentMonth, day).toISOString().split('T')[0];
		
		return scheduledLinks.filter(link => {
			const scheduled = link.scheduled_at ? link.scheduled_at.split('T')[0] : null;
			const expires = link.expires_at ? link.expires_at.split('T')[0] : null;
			return scheduled === dateStr || expires === dateStr;
		});
	}

	function previousMonth() {
		if (currentMonth === 0) {
			currentMonth = 11;
			currentYear--;
		} else {
			currentMonth--;
		}
	}

	function nextMonth() {
		if (currentMonth === 11) {
			currentMonth = 0;
			currentYear++;
		} else {
			currentMonth++;
		}
	}

	function goToToday() {
		const today = new Date();
		currentMonth = today.getMonth();
		currentYear = today.getFullYear();
	}

	const monthNames = [
		'January', 'February', 'March', 'April', 'May', 'June',
		'July', 'August', 'September', 'October', 'November', 'December'
	];

	const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
</script>

<div class="bg-white rounded-2xl border border-gray-200 overflow-hidden">
	<!-- Header -->
	<div class="bg-gradient-to-r from-indigo-600 to-blue-600 px-6 py-4">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-bold text-white">
				{monthNames[currentMonth]} {currentYear}
			</h3>
			<div class="flex items-center gap-2">
				<button
					onclick={goToToday}
					class="px-3 py-1.5 bg-white/20 hover:bg-white/30 text-white text-sm font-medium rounded-lg transition-colors"
				>
					Today
				</button>
				<button
					onclick={previousMonth}
					class="p-2 bg-white/20 hover:bg-white/30 text-white rounded-lg transition-colors"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
					</svg>
				</button>
				<button
					onclick={nextMonth}
					class="p-2 bg-white/20 hover:bg-white/30 text-white rounded-lg transition-colors"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
					</svg>
				</button>
			</div>
		</div>
	</div>

	<!-- Calendar Grid -->
	<div class="p-4">
		<!-- Day Names -->
		<div class="grid grid-cols-7 gap-2 mb-2">
			{#each dayNames as day}
				<div class="text-center text-xs font-semibold text-gray-500 py-2">
					{day}
				</div>
			{/each}
		</div>

		<!-- Calendar Days -->
		<div class="grid grid-cols-7 gap-2">
			{#each calendarDays as { date, isCurrentMonth }}
				{@const dayLinks = getLinksForDate(date)}
				{@const isToday = date === new Date().getDate() && 
					currentMonth === new Date().getMonth() && 
					currentYear === new Date().getFullYear()}
				
				<div
					class="min-h-[80px] p-2 rounded-lg border transition-all {
						isCurrentMonth 
							? dayLinks.length > 0 
								? 'bg-indigo-50 border-indigo-200 hover:bg-indigo-100 cursor-pointer' 
								: isToday
									? 'bg-blue-50 border-blue-300'
									: 'bg-white border-gray-200 hover:bg-gray-50'
							: 'bg-gray-50 border-gray-100'
					}"
					onclick={() => dayLinks.length > 0 && dispatch('selectDate', { date, links: dayLinks })}
				>
					{#if date}
						<div class="text-sm font-semibold mb-1 {
							isCurrentMonth 
								? isToday 
									? 'text-blue-600' 
									: 'text-gray-900'
								: 'text-gray-400'
						}">
							{date}
						</div>
						
						{#if dayLinks.length > 0}
							<div class="space-y-1">
								{#each dayLinks.slice(0, 2) as link}
									{@const hasScheduled = link.scheduled_at && link.scheduled_at.split('T')[0] === new Date(currentYear, currentMonth, date).toISOString().split('T')[0]}
									{@const hasExpiry = link.expires_at && link.expires_at.split('T')[0] === new Date(currentYear, currentMonth, date).toISOString().split('T')[0]}
									
									<div class="text-xs px-2 py-1 rounded {
										hasScheduled ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
									} truncate">
										<span class="font-medium">{hasScheduled ? '▲' : '▼'}</span>
										{link.title}
									</div>
								{/each}
								{#if dayLinks.length > 2}
									<div class="text-xs text-gray-500 font-medium px-2">
										+{dayLinks.length - 2} more
									</div>
								{/if}
							</div>
						{/if}
					{/if}
				</div>
			{/each}
		</div>
	</div>

	<!-- Legend -->
	<div class="border-t border-gray-200 px-6 py-3 bg-gray-50">
		<div class="flex items-center gap-6 text-xs">
			<div class="flex items-center gap-2">
				<div class="w-3 h-3 bg-green-100 border border-green-200 rounded"></div>
				<span class="text-gray-600">Scheduled Publish</span>
			</div>
			<div class="flex items-center gap-2">
				<div class="w-3 h-3 bg-red-100 border border-red-200 rounded"></div>
				<span class="text-gray-600">Expires</span>
			</div>
			<div class="ml-auto text-gray-500">
				{scheduledLinks.length} scheduled link{scheduledLinks.length !== 1 ? 's' : ''}
			</div>
		</div>
	</div>
</div>
