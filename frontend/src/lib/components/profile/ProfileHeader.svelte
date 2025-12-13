<script lang="ts">
	import type { HeaderStyles } from '$lib/stores/theme';
	
	let { profile, headerStyle, textColor, textSecondaryColor }: {
		profile: any;
		headerStyle: HeaderStyles;
		textColor: string;
		textSecondaryColor: string;
	} = $props();
	
	const avatarSrc = $derived(profile?.avatar_url || `https://api.dicebear.com/7.x/avataaars/svg?seed=${profile?.username || 'user'}`);
	const username = $derived(profile?.username || 'username');
	const bio = $derived(profile?.bio || '');
	const bioFontSize = $derived(headerStyle.bioSize === 'sm' ? '0.875rem' : headerStyle.bioSize === 'lg' ? '1.125rem' : '1rem');
</script>

{#if headerStyle.layout === 'centered'}
	<div class="relative">
		{#if headerStyle.showCover}
			<div style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};"></div>
		{/if}
		<div class="px-6 pb-6" style="padding-top: {headerStyle.showCover ? '20px' : '80px'};">
			<div class="flex justify-center" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2}px` : '0'};">
				<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover shadow-xl" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
			</div>
			<div class="mt-3" style="text-align: {headerStyle.bioAlign};">
				<h2 class="text-xl font-bold" style="color: {textColor};">@{username}</h2>
				{#if bio}<p class="mt-1" style="color: {textSecondaryColor}; font-size: {bioFontSize};">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'overlap'}
	<div class="relative pb-6">
		{#if headerStyle.showCover}
			<div style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};"></div>
		{/if}
		<div class="px-6" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2 + 10}px` : '80px'};">
			<div class="flex justify-center">
				<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover shadow-2xl" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
			</div>
			<div class="mt-4 text-center">
				<h2 class="text-xl font-bold" style="color: {textColor};">@{username}</h2>
				{#if bio}<p class="text-sm mt-1" style="color: {textSecondaryColor};">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'card'}
	<div class="px-6 pt-20 pb-6">
		<div class="bg-white/90 backdrop-blur-sm rounded-2xl p-4 shadow-xl border border-gray-200">
			{#if headerStyle.showCover}
				<div class="rounded-xl mb-4" style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};"></div>
			{/if}
			<div class="flex justify-center" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2 + 16}px` : '0'};">
				<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover shadow-lg" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
			</div>
			<div class="mt-3 text-center">
				<h2 class="text-lg font-bold text-gray-900">@{username}</h2>
				{#if bio}<p class="text-sm mt-1 text-gray-600">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'glass'}
	<div class="relative pb-6">
		{#if headerStyle.showCover}
			<div style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};"></div>
		{/if}
		<div class="px-6" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2 + 20}px` : '80px'};">
			<div class="bg-white/20 backdrop-blur-xl rounded-2xl p-5 border border-white/30 shadow-2xl">
				<div class="flex justify-center">
					<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover shadow-xl ring-2 ring-white/50" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
				</div>
				<div class="mt-4 text-center">
					<h2 class="text-xl font-bold text-white drop-shadow-lg">@{username}</h2>
					{#if bio}<p class="text-sm mt-1 text-white/90 drop-shadow">{bio}</p>{/if}
				</div>
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'gradient'}
	<div class="relative pb-6">
		{#if headerStyle.showCover}
			<div class="relative overflow-hidden" style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};">
				<div class="absolute inset-0 opacity-30" style="background: radial-gradient(circle at 20% 50%, rgba(255,255,255,0.3) 0%, transparent 50%), radial-gradient(circle at 80% 80%, rgba(255,255,255,0.2) 0%, transparent 50%);"></div>
			</div>
		{/if}
		<div class="px-6" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2}px` : '80px'};">
			<div class="flex justify-center">
				<div class="relative">
					<div class="absolute inset-0 rounded-full" style="background: linear-gradient(135deg, {headerStyle.coverGradientFrom}, {headerStyle.coverGradientTo}); filter: blur(8px); opacity: 0.6;"></div>
					<img src={avatarSrc} alt="Avatar" class="relative rounded-full object-cover shadow-2xl" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
				</div>
			</div>
			<div class="mt-4 text-center">
				<h2 class="text-2xl font-bold" style="color: {textColor};">@{username}</h2>
				{#if bio}<p class="text-base mt-2" style="color: {textSecondaryColor};">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'minimal'}
	<div class="px-6 pt-20 pb-6">
		<div class="flex items-center gap-3">
			<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
			<div>
				<h2 class="text-lg font-bold" style="color: {textColor};">@{username}</h2>
				{#if bio}<p class="text-xs mt-0.5" style="color: {textSecondaryColor};">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'full'}
	<div class="relative">
		{#if headerStyle.showCover}
			<div class="flex items-center justify-center" style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};">
				<div class="text-center">
					<img src={avatarSrc} alt="Avatar" class="mx-auto rounded-full object-cover shadow-2xl" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
					<h2 class="text-2xl font-bold text-white mt-4 drop-shadow-lg">@{username}</h2>
					{#if bio}<p class="text-base text-white/90 mt-2 drop-shadow">{bio}</p>{/if}
				</div>
			</div>
		{/if}
		<div class="px-6 pt-6"></div>
	</div>
{:else if headerStyle.layout === 'side'}
	<div class="px-6 pt-20 pb-6">
		<div class="flex items-start gap-4">
			<img src={avatarSrc} alt="Avatar" class="flex-shrink-0 rounded-full object-cover" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
			<div class="flex-1 pt-1">
				<h2 class="text-lg font-bold" style="color: {textColor};">@{username}</h2>
				{#if bio}<p class="text-sm mt-1" style="color: {textSecondaryColor};">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'split'}
	<div class="relative pb-6">
		{#if headerStyle.showCover}
			<div class="flex" style="height: {headerStyle.coverHeight}px;">
				<div class="w-1/2" style="background: {headerStyle.coverGradientFrom};"></div>
				<div class="w-1/2" style="background: {headerStyle.coverGradientTo};"></div>
			</div>
		{/if}
		<div class="px-6" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2}px` : '80px'};">
			<div class="flex justify-center">
				<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover shadow-2xl" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
			</div>
			<div class="mt-4 text-center">
				<h2 class="text-xl font-bold" style="color: {textColor};">@{username}</h2>
				{#if bio}<p class="text-sm mt-1" style="color: {textSecondaryColor};">{bio}</p>{/if}
			</div>
		</div>
	</div>
{:else if headerStyle.layout === 'asymmetric'}
	<div class="relative pb-6">
		{#if headerStyle.showCover}
			<div class="relative" style="height: {headerStyle.coverHeight}px; background: {headerStyle.coverType === 'gradient' ? `linear-gradient(135deg, ${headerStyle.coverGradientFrom}, ${headerStyle.coverGradientTo})` : headerStyle.coverColor};">
				<div class="absolute bottom-0 right-0 w-2/3 h-2/3" style="background: linear-gradient(135deg, transparent, rgba(0,0,0,0.1));"></div>
			</div>
		{/if}
		<div class="px-6" style="margin-top: {headerStyle.showCover ? `-${headerStyle.avatarSize/2 + 20}px` : '80px'};">
			<div class="flex items-end gap-4">
				<img src={avatarSrc} alt="Avatar" class="rounded-full object-cover shadow-2xl" style="width: {headerStyle.avatarSize}px !important; height: {headerStyle.avatarSize}px !important; border: {headerStyle.avatarBorder}px solid {headerStyle.avatarBorderColor};" />
				<div class="flex-1 pb-2">
					<h2 class="text-2xl font-bold" style="color: {textColor};">@{username}</h2>
					{#if bio}<p class="text-base mt-1" style="color: {textSecondaryColor};">{bio}</p>{/if}
				</div>
			</div>
		</div>
	</div>
{/if}
