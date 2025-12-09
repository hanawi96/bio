<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { toast } from 'svelte-sonner';
	
	import StatsCard from '$lib/components/dashboard/overview/StatsCard.svelte';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import type { Link } from '$lib/api/links';

	let links: Link[] = [];
	let initialLoading = true;

	$: stats = {
		totalClicks: links.reduce((sum, link) => sum + link.clicks, 0),
		totalLinks: links.length,
		activeLinks: links.filter(l => l.is_active).length,
		profileViews: 0
	};

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			initialLoading = true;
			
			try {
				const linksData = await linksApi.getLinks($auth.token!);
				links = linksData.sort((a, b) => a.position - b.position);
			} catch (linksError: any) {
				console.warn('No links found:', linksError);
				links = [];
			}
		} catch (error: any) {
			console.error('Failed to load dashboard data:', error);
			toast.error(error.message || 'Failed to load data');
		} finally {
			initialLoading = false;
		}
	}

	const icons = {
		clicks: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"/></svg>',
		links: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/></svg>',
		active: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>',
		views: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>'
	};
</script>

<svelte:head>
	<title>Dashboard - LinkBio</title>
</svelte:head>

{#if initialLoading}
	<div class="h-full flex items-center justify-center bg-gray-50">
		<div class="text-center">
			<div class="relative w-20 h-20 mx-auto mb-6">
				<div class="absolute inset-0 border-4 border-purple-100 rounded-full"></div>
				<div class="absolute inset-0 border-4 border-transparent border-t-purple-600 border-r-blue-600 rounded-full animate-spin"></div>
			</div>
			<h3 class="text-lg font-semibold text-gray-900 mb-2">Loading Dashboard</h3>
			<p class="text-sm text-gray-600">Please wait a moment...</p>
		</div>
	</div>
{:else}
<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 py-6 sticky top-0 z-10">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 via-indigo-900 to-gray-900 bg-clip-text text-transparent">
					Dashboard
				</h1>
				<p class="text-sm text-gray-500 mt-1.5">Overview of your LinkBio performance</p>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8">
		<!-- Stats -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<StatsCard title="Total Clicks" value={stats.totalClicks} icon={icons.clicks} />
			<StatsCard title="Total Links" value={stats.totalLinks} icon={icons.links} />
			<StatsCard title="Active Links" value={stats.activeLinks} icon={icons.active} />
			<StatsCard title="Profile Views" value={stats.profileViews} icon={icons.views} />
		</div>

		<!-- Quick Actions -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<div class="bg-white rounded-xl p-6 shadow-sm">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
				<div class="space-y-3">
					<a href="/dashboard/bio" class="flex items-center gap-3 p-4 rounded-lg border border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 transition-colors">
						<div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center">
							<svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
							</svg>
						</div>
						<div>
							<h3 class="font-semibold text-gray-900">Build My Bio</h3>
							<p class="text-sm text-gray-500">Add content to your bio page</p>
						</div>
					</a>
					<a href="/dashboard/analytics" class="flex items-center gap-3 p-4 rounded-lg border border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 transition-colors">
						<div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
							<svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
							</svg>
						</div>
						<div>
							<h3 class="font-semibold text-gray-900">View Analytics</h3>
							<p class="text-sm text-gray-500">Track your performance</p>
						</div>
					</a>
				</div>
			</div>

			<div class="bg-white rounded-xl p-6 shadow-sm">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Recent Activity</h2>
				<div class="text-center py-8 text-gray-500">
					<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
					</svg>
					<p>No recent activity</p>
				</div>
			</div>
		</div>
	</div>
</div>
{/if}

