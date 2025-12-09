<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';

	const activities = [
		{ type: 'click', link: 'My Website', time: '2 minutes ago', icon: '🖱️' },
		{ type: 'view', time: '5 minutes ago', icon: '👁️' },
		{ type: 'click', link: 'YouTube Channel', time: '12 minutes ago', icon: '🖱️' },
		{ type: 'edit', action: 'Updated bio', time: '1 hour ago', icon: '✏️' },
		{ type: 'click', link: 'Twitter', time: '2 hours ago', icon: '🖱️' }
	];
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Recent Activity</Card.Title>
		<Card.Description>Latest interactions with your profile</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class="space-y-4">
			{#each activities as activity}
				<div class="flex items-start gap-3 pb-4 border-b border-gray-100 last:border-0 last:pb-0">
					<div class="text-2xl">{activity.icon}</div>
					<div class="flex-1 min-w-0">
						{#if activity.type === 'click'}
							<p class="text-sm text-gray-900">
								Someone clicked <span class="font-semibold">{activity.link}</span>
							</p>
						{:else if activity.type === 'view'}
							<p class="text-sm text-gray-900">Your profile was viewed</p>
						{:else if activity.type === 'edit'}
							<p class="text-sm text-gray-900">{activity.action}</p>
						{/if}
						<p class="text-xs text-gray-500 mt-1">{activity.time}</p>
					</div>
				</div>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
