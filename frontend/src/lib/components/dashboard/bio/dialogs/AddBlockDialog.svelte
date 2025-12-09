<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';

	export let open = false;
	
	const dispatch = createEventDispatcher();

	type BlockType = 'text' | 'image' | 'video' | 'social' | 'divider' | 'email' | 'embed';

	const blockTypes: Array<{
		type: BlockType;
		icon: string;
		title: string;
		description: string;
		badge?: string;
	}> = [
		{
			type: 'text',
			icon: 'M4 6h16M4 12h16M4 18h7',
			title: 'Text',
			description: 'Add heading or paragraph'
		},
		{
			type: 'image',
			icon: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z',
			title: 'Image',
			description: 'Upload image or banner'
		},
		{
			type: 'video',
			icon: 'M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z',
			title: 'Video',
			description: 'Embed YouTube, Vimeo, TikTok'
		},
		{
			type: 'social',
			icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z',
			title: 'Social Links',
			description: 'Add social media icons'
		},
		{
			type: 'divider',
			icon: 'M20 12H4',
			title: 'Divider',
			description: 'Add spacing or line'
		},
		{
			type: 'email',
			icon: 'M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z',
			title: 'Email Collector',
			description: 'Collect subscriber emails',
			badge: 'NEW'
		},
		{
			type: 'embed',
			icon: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4',
			title: 'Embed',
			description: 'Spotify, Maps, SoundCloud',
			badge: 'NEW'
		}
	];

	function selectBlock(type: BlockType) {
		dispatch('select', { type });
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-2xl">
		<Dialog.Header>
			<Dialog.Title>Add Block</Dialog.Title>
			<Dialog.Description>
				Choose a block type to add to your bio page
			</Dialog.Description>
		</Dialog.Header>
		
		<div class="grid grid-cols-2 gap-3 py-4">
			{#each blockTypes as block}
				<button
					onclick={() => selectBlock(block.type)}
					class="group relative p-4 bg-gray-50 hover:bg-indigo-50 rounded-xl text-left transition-all border-2 border-transparent hover:border-indigo-200"
				>
					{#if block.badge}
						<span class="absolute top-2 right-2 px-2 py-0.5 bg-gradient-to-r from-pink-500 to-rose-500 text-white text-[10px] font-bold rounded-full">
							{block.badge}
						</span>
					{/if}
					<div class="flex items-start gap-3">
						<div class="flex-shrink-0 w-10 h-10 bg-white group-hover:bg-indigo-100 rounded-lg flex items-center justify-center transition-colors">
							<svg class="w-5 h-5 text-gray-600 group-hover:text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={block.icon}/>
							</svg>
						</div>
						<div class="flex-1 min-w-0">
							<h3 class="text-sm font-semibold text-gray-900 mb-0.5">{block.title}</h3>
							<p class="text-xs text-gray-500">{block.description}</p>
						</div>
					</div>
				</button>
			{/each}
		</div>
	</Dialog.Content>
</Dialog.Root>
