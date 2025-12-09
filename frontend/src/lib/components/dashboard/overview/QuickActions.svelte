<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';

	const actions = [
		{
			title: 'Share Profile',
			description: 'Share your LinkBio with QR code',
			icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/></svg>',
			color: 'from-blue-500 to-cyan-500'
		},
		{
			title: 'Analytics',
			description: 'View detailed insights',
			icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>',
			color: 'from-purple-500 to-pink-500'
		},
		{
			title: 'Integrations',
			description: 'Connect your tools',
			icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z"/></svg>',
			color: 'from-green-500 to-teal-500'
		}
	];
</script>

<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
	{#each actions as action}
		<button class="group relative bg-white rounded-xl p-6 text-left hover:shadow-lg transition-all duration-200 border border-gray-100 overflow-hidden">
			<div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br {action.color} opacity-10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-300"></div>
			<div class="relative">
				<div class="w-12 h-12 bg-gradient-to-br {action.color} rounded-xl flex items-center justify-center text-white mb-4 group-hover:scale-110 transition-transform">
					{@html action.icon}
				</div>
				<h3 class="font-semibold text-gray-900 mb-1">{action.title}</h3>
				<p class="text-sm text-gray-600">{action.description}</p>
			</div>
		</button>
	{/each}
</div>
