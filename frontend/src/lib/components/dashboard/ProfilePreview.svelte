<script lang="ts">
	import Avatar from '$lib/components/ui/avatar.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';

	export let profile: Partial<Profile> = {};
	export let links: Link[] = [];

	$: activeLinks = links.filter(l => l.is_active);
</script>

<div class="w-full max-w-md mx-auto">
	<!-- Phone Frame -->
	<div class="relative bg-gray-900 rounded-[3rem] p-3 shadow-2xl">
		<!-- Notch -->
		<div class="absolute top-0 left-1/2 -translate-x-1/2 w-40 h-7 bg-gray-900 rounded-b-3xl z-10"></div>
		
		<!-- Screen -->
		<div class="relative bg-gradient-to-br from-purple-50 to-blue-50 rounded-[2.5rem] overflow-hidden h-[600px]">
			<div class="h-full overflow-y-auto p-6 space-y-6">
				<!-- Profile Header -->
				<div class="text-center space-y-3">
					<div class="flex justify-center">
						<div class="relative">
							<Avatar 
								src={profile.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + (profile.username || 'user')}
								alt="Avatar"
								class="w-24 h-24 border-4 border-white shadow-lg"
							/>
							<div class="absolute -bottom-1 -right-1 w-6 h-6 bg-green-500 border-2 border-white rounded-full"></div>
						</div>
					</div>
					<div>
						<h2 class="text-xl font-bold text-gray-900">@{profile.username || 'username'}</h2>
						{#if profile.bio}
							<p class="text-sm text-gray-600 mt-1">{profile.bio}</p>
						{/if}
					</div>
				</div>

				<!-- Links -->
				<div class="space-y-3">
					{#each activeLinks as link}
						<a
							href={link.url}
							target="_blank"
							rel="noopener noreferrer"
							class="block w-full bg-white hover:bg-gray-50 rounded-2xl p-4 shadow-sm hover:shadow-md transition-all duration-200 border border-gray-100"
						>
							<div class="flex items-center justify-between">
								<span class="font-medium text-gray-900">{link.title}</span>
								<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
								</svg>
							</div>
						</a>
					{/each}
					
					{#if activeLinks.length === 0}
						<div class="text-center py-8 text-gray-400">
							<svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
							</svg>
							<p class="text-sm">No links yet</p>
						</div>
					{/if}
				</div>

				<!-- Footer -->
				<div class="text-center pt-4">
					<Badge variant="secondary" class="text-xs">
						Made with LinkBio
					</Badge>
				</div>
			</div>
		</div>
	</div>
</div>
