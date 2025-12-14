<script lang="ts">
	import { profileApi } from '$lib/api/profile';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	let { 
		settings
	}: {
		settings: {
			showShareButton: boolean;
			showSubscribeButton: boolean;
			hideBranding: boolean;
		};
	} = $props();

	let saving = $state(false);

	// Map camelCase to snake_case for API
	const fieldMap: Record<string, string> = {
		showShareButton: 'show_share_button',
		showSubscribeButton: 'show_subscribe_button',
		hideBranding: 'hide_branding'
	};

	async function updateSetting(field: string, value: boolean) {
		// Save to API in background (no await, no blocking)
		const apiField = fieldMap[field];
		saving = true;
		profileApi.updateProfile({ [apiField]: value }, get(auth).token!)
			.then(() => {
				saving = false;
			})
			.catch((e: any) => {
				// Rollback on error
				(settings as any)[field] = !value;
				toast.error(e.message || 'Failed to update setting');
				saving = false;
			});
	}
</script>

<div class="space-y-4">
	<h3 class="text-lg font-bold text-gray-900">Page Settings</h3>
	
	<!-- Show Share Button -->
	<div class="flex items-start justify-between p-4 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors">
		<div class="flex-1">
			<h4 class="text-sm font-semibold text-gray-900">Show Share Button</h4>
			<p class="text-xs text-gray-500 mt-0.5">Let visitors easily share your profile</p>
		</div>
		<label class="relative inline-flex items-center cursor-pointer">
			<input
				type="checkbox"
				bind:checked={settings.showShareButton}
				onchange={(e) => updateSetting('showShareButton', e.currentTarget.checked)}
				class="sr-only peer"
			/>
			<div class="w-11 h-6 bg-gray-300 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
		</label>
	</div>

	<!-- Show Subscribe Button -->
	<div class="flex items-start justify-between p-4 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors">
		<div class="flex-1">
			<h4 class="text-sm font-semibold text-gray-900">Show Subscribe Button</h4>
			<p class="text-xs text-gray-500 mt-0.5">Let visitors subscribe to your mailing list</p>
		</div>
		<label class="relative inline-flex items-center cursor-pointer">
			<input
				type="checkbox"
				bind:checked={settings.showSubscribeButton}
				onchange={(e) => updateSetting('showSubscribeButton', e.currentTarget.checked)}
				class="sr-only peer"
			/>
			<div class="w-11 h-6 bg-gray-300 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
		</label>
	</div>

	<!-- Hide Branding (Pro) -->
	<div class="flex items-start justify-between p-4 bg-gray-50 rounded-xl opacity-60 cursor-not-allowed">
		<div class="flex-1">
			<div class="flex items-center gap-2">
				<h4 class="text-sm font-semibold text-gray-900">Hide Beacons logo and footer</h4>
				<span class="inline-flex items-center gap-1 px-2 py-0.5 bg-gradient-to-r from-pink-500 to-red-500 text-white text-xs font-bold rounded-full">
					<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
					</svg>
					Pro
				</span>
			</div>
		</div>
		<label class="relative inline-flex items-center cursor-not-allowed">
			<input
				type="checkbox"
				bind:checked={settings.hideBranding}
				disabled
				class="sr-only peer"
			/>
			<div class="w-11 h-6 bg-gray-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
		</label>
	</div>

	{#if saving}
		<div class="flex items-center gap-2 text-xs text-indigo-600">
			<svg class="w-3 h-3 animate-spin" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			Saving...
		</div>
	{/if}
</div>
