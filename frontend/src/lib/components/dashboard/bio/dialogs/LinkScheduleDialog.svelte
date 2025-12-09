<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import type { Link } from '$lib/api/links';

	export let open = false;
	export let link: Link;

	const dispatch = createEventDispatcher();

	let scheduledDate = '';
	let scheduledTime = '';
	let expiresDate = '';
	let expiresTime = '';
	let enableSchedule = false;
	let enableExpiry = false;

	// Initialize from link data
	$: if (open && link) {
		if (link.scheduled_at) {
			const date = new Date(link.scheduled_at);
			scheduledDate = date.toISOString().split('T')[0];
			scheduledTime = date.toTimeString().slice(0, 5);
			enableSchedule = true;
		} else {
			scheduledDate = '';
			scheduledTime = '';
			enableSchedule = false;
		}

		if (link.expires_at) {
			const date = new Date(link.expires_at);
			expiresDate = date.toISOString().split('T')[0];
			expiresTime = date.toTimeString().slice(0, 5);
			enableExpiry = true;
		} else {
			expiresDate = '';
			expiresTime = '';
			enableExpiry = false;
		}
	}

	function handleSave() {
		let scheduled_at = null;
		let expires_at = null;

		if (enableSchedule && scheduledDate && scheduledTime) {
			scheduled_at = new Date(`${scheduledDate}T${scheduledTime}`).toISOString();
		}

		if (enableExpiry && expiresDate && expiresTime) {
			expires_at = new Date(`${expiresDate}T${expiresTime}`).toISOString();
		}

		dispatch('save', { id: link.id, scheduled_at, expires_at });
		open = false;
	}

	function clearSchedule() {
		enableSchedule = false;
		enableExpiry = false;
		scheduledDate = '';
		scheduledTime = '';
		expiresDate = '';
		expiresTime = '';
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Schedule Link</Dialog.Title>
			<Dialog.Description>
				Set when this link should be published and when it should expire
			</Dialog.Description>
		</Dialog.Header>

		<div class="space-y-6 py-4">
			<!-- Scheduled Publish -->
			<div class="space-y-3">
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						id="enable-schedule"
						bind:checked={enableSchedule}
						class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all
							checked:bg-indigo-600 checked:border-indigo-600 
							hover:border-indigo-400
							focus:ring-2 focus:ring-indigo-500/20"
					/>
					<Label for="enable-schedule" class="cursor-pointer font-semibold">
						Schedule Publish
					</Label>
				</div>

				{#if enableSchedule}
					<div class="pl-6 space-y-3 animate-in">
						<div class="grid grid-cols-2 gap-3">
							<div>
								<Label for="scheduled-date" class="text-xs text-gray-600">Date</Label>
								<input
									type="date"
									id="scheduled-date"
									bind:value={scheduledDate}
									class="w-full mt-1 px-3 py-2 bg-gray-50 border-0 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white"
								/>
							</div>
							<div>
								<Label for="scheduled-time" class="text-xs text-gray-600">Time</Label>
								<input
									type="time"
									id="scheduled-time"
									bind:value={scheduledTime}
									class="w-full mt-1 px-3 py-2 bg-gray-50 border-0 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white"
								/>
							</div>
						</div>
						<p class="text-xs text-gray-500">
							Link will automatically become active at this time
						</p>
					</div>
				{/if}
			</div>

			<!-- Divider -->
			<div class="border-t border-gray-200"></div>

			<!-- Expiry -->
			<div class="space-y-3">
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						id="enable-expiry"
						bind:checked={enableExpiry}
						class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all
							checked:bg-indigo-600 checked:border-indigo-600 
							hover:border-indigo-400
							focus:ring-2 focus:ring-indigo-500/20"
					/>
					<Label for="enable-expiry" class="cursor-pointer font-semibold">
						Set Expiry
					</Label>
				</div>

				{#if enableExpiry}
					<div class="pl-6 space-y-3 animate-in">
						<div class="grid grid-cols-2 gap-3">
							<div>
								<Label for="expires-date" class="text-xs text-gray-600">Date</Label>
								<input
									type="date"
									id="expires-date"
									bind:value={expiresDate}
									class="w-full mt-1 px-3 py-2 bg-gray-50 border-0 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white"
								/>
							</div>
							<div>
								<Label for="expires-time" class="text-xs text-gray-600">Time</Label>
								<input
									type="time"
									id="expires-time"
									bind:value={expiresTime}
									class="w-full mt-1 px-3 py-2 bg-gray-50 border-0 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white"
								/>
							</div>
						</div>
						<p class="text-xs text-gray-500">
							Link will automatically become inactive at this time
						</p>
					</div>
				{/if}
			</div>

			<!-- Preview -->
			{#if enableSchedule || enableExpiry}
				<div class="bg-indigo-50 border border-indigo-200 rounded-lg p-3">
					<div class="flex items-start gap-2">
						<svg class="w-5 h-5 text-indigo-600 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
						</svg>
						<div class="text-sm text-indigo-900 space-y-1">
							{#if enableSchedule && scheduledDate && scheduledTime}
								<p>
									<strong>Publishes:</strong> {new Date(`${scheduledDate}T${scheduledTime}`).toLocaleString()}
								</p>
							{/if}
							{#if enableExpiry && expiresDate && expiresTime}
								<p>
									<strong>Expires:</strong> {new Date(`${expiresDate}T${expiresTime}`).toLocaleString()}
								</p>
							{/if}
						</div>
					</div>
				</div>
			{/if}
		</div>

		<Dialog.Footer>
			<div class="flex items-center justify-between w-full">
				<Button variant="ghost" onclick={clearSchedule} class="text-gray-600">
					Clear All
				</Button>
				<div class="flex gap-2">
					<Button variant="outline" onclick={() => open = false}>Cancel</Button>
					<Button onclick={handleSave} class="bg-indigo-600 hover:bg-indigo-700">
						Save Schedule
					</Button>
				</div>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<style>
	@keyframes fade-in {
		from { opacity: 0; transform: translateY(-10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	
	.animate-in {
		animation: fade-in 0.2s ease-out;
	}
</style>
