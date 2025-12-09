<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Tabs from '$lib/components/ui/tabs';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import StatsCard from '$lib/components/dashboard/overview/StatsCard.svelte';

	const icons = {
		clicks: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"/></svg>',
		views: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>',
		ctr: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/></svg>',
		time: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>'
	};

	const topLinks = [
		{ title: 'My Website', clicks: 142, ctr: '24.5%', trend: '+12%' },
		{ title: 'YouTube Channel', clicks: 89, ctr: '18.2%', trend: '+8%' },
		{ title: 'Twitter', clicks: 56, ctr: '12.1%', trend: '-3%' },
		{ title: 'Instagram', clicks: 23, ctr: '8.4%', trend: '+5%' }
	];

	const countries = [
		{ name: 'United States', flag: '🇺🇸', views: 342, percentage: 45 },
		{ name: 'United Kingdom', flag: '🇬🇧', views: 156, percentage: 20 },
		{ name: 'Canada', flag: '🇨🇦', views: 98, percentage: 13 },
		{ name: 'Germany', flag: '🇩🇪', views: 87, percentage: 11 },
		{ name: 'Others', flag: '🌍', views: 84, percentage: 11 }
	];

	const devices = [
		{ name: 'Mobile', percentage: 68, color: 'bg-purple-500' },
		{ name: 'Desktop', percentage: 28, color: 'bg-blue-500' },
		{ name: 'Tablet', percentage: 4, color: 'bg-green-500' }
	];
</script>

<svelte:head>
	<title>Analytics - LinkBio</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-gray-50 via-purple-50 to-blue-50">
	<div class="container mx-auto px-4 py-8">
		<!-- Header -->
		<div class="flex items-center justify-between mb-8">
			<div>
				<h1 class="text-3xl font-bold text-gray-900">Analytics</h1>
				<p class="text-gray-600 mt-1">Track your profile performance</p>
			</div>
			<div class="flex gap-2">
				<Button variant="outline">Last 7 days</Button>
				<Button variant="outline">Export</Button>
			</div>
		</div>

		<!-- Stats Overview -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
			<StatsCard title="Total Clicks" value="310" icon={icons.clicks} trend="12%" trendUp={true} />
			<StatsCard title="Profile Views" value="1,247" icon={icons.views} trend="8%" trendUp={true} />
			<StatsCard title="Click Rate" value="24.8%" icon={icons.ctr} trend="3%" trendUp={true} />
			<StatsCard title="Avg. Time" value="2m 34s" icon={icons.time} />
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Top Links -->
			<div class="lg:col-span-2">
				<Card.Root>
					<Card.Header>
						<Card.Title>Top Performing Links</Card.Title>
						<Card.Description>Your most clicked links this week</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="space-y-4">
							{#each topLinks as link, i}
								<div class="flex items-center gap-4 p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
									<div class="flex items-center justify-center w-8 h-8 bg-gradient-to-br from-purple-500 to-blue-500 text-white rounded-lg font-bold text-sm">
										{i + 1}
									</div>
									<div class="flex-1 min-w-0">
										<h4 class="font-semibold text-gray-900 truncate">{link.title}</h4>
										<div class="flex items-center gap-4 mt-1">
											<span class="text-sm text-gray-600">{link.clicks} clicks</span>
											<span class="text-sm text-gray-600">CTR: {link.ctr}</span>
										</div>
									</div>
									<Badge variant={link.trend.startsWith('+') ? 'default' : 'secondary'}>
										{link.trend}
									</Badge>
								</div>
							{/each}
						</div>
					</Card.Content>
				</Card.Root>

				<!-- Chart Placeholder -->
				<Card.Root class="mt-6">
					<Card.Header>
						<Card.Title>Clicks Over Time</Card.Title>
						<Card.Description>Daily click trends</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="h-64 flex items-end justify-between gap-2">
							{#each Array(7) as _, i}
								<div class="flex-1 bg-gradient-to-t from-purple-500 to-blue-500 rounded-t-lg hover:opacity-80 transition-opacity cursor-pointer" 
									style="height: {Math.random() * 80 + 20}%">
								</div>
							{/each}
						</div>
						<div class="flex justify-between mt-2 text-xs text-gray-500">
							<span>Mon</span>
							<span>Tue</span>
							<span>Wed</span>
							<span>Thu</span>
							<span>Fri</span>
							<span>Sat</span>
							<span>Sun</span>
						</div>
					</Card.Content>
				</Card.Root>
			</div>

			<!-- Sidebar Stats -->
			<div class="space-y-6">
				<!-- Countries -->
				<Card.Root>
					<Card.Header>
						<Card.Title>Top Countries</Card.Title>
						<Card.Description>Where your visitors are from</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="space-y-4">
							{#each countries as country}
								<div>
									<div class="flex items-center justify-between mb-2">
										<div class="flex items-center gap-2">
											<span class="text-xl">{country.flag}</span>
											<span class="text-sm font-medium text-gray-900">{country.name}</span>
										</div>
										<span class="text-sm text-gray-600">{country.views}</span>
									</div>
									<div class="w-full bg-gray-200 rounded-full h-2">
										<div class="bg-gradient-to-r from-purple-500 to-blue-500 h-2 rounded-full transition-all duration-500" 
											style="width: {country.percentage}%">
										</div>
									</div>
								</div>
							{/each}
						</div>
					</Card.Content>
				</Card.Root>

				<!-- Devices -->
				<Card.Root>
					<Card.Header>
						<Card.Title>Devices</Card.Title>
						<Card.Description>How visitors access your profile</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="space-y-4">
							{#each devices as device}
								<div>
									<div class="flex items-center justify-between mb-2">
										<span class="text-sm font-medium text-gray-900">{device.name}</span>
										<span class="text-sm text-gray-600">{device.percentage}%</span>
									</div>
									<div class="w-full bg-gray-200 rounded-full h-2">
										<div class="{device.color} h-2 rounded-full transition-all duration-500" 
											style="width: {device.percentage}%">
										</div>
									</div>
								</div>
							{/each}
						</div>

						<!-- Pie Chart Visual -->
						<div class="mt-6 flex justify-center">
							<div class="relative w-32 h-32">
								<svg viewBox="0 0 100 100" class="transform -rotate-90">
									<circle cx="50" cy="50" r="40" fill="none" stroke="#e5e7eb" stroke-width="20"/>
									<circle cx="50" cy="50" r="40" fill="none" stroke="#a855f7" stroke-width="20" 
										stroke-dasharray="68 100" stroke-dashoffset="0"/>
									<circle cx="50" cy="50" r="40" fill="none" stroke="#3b82f6" stroke-width="20" 
										stroke-dasharray="28 100" stroke-dashoffset="-68"/>
								</svg>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
	</div>
</div>
