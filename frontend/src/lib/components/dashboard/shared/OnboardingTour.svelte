<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { onMount } from 'svelte';

	let showOnboarding = false;
	let currentStep = 0;

	const steps = [
		{
			title: 'Welcome to LinkBio! 🎉',
			description: 'Let\'s take a quick tour to help you get started with your new profile.',
			image: '🚀'
		},
		{
			title: 'Add Your First Link',
			description: 'Click the "Add Link" button to create your first link. You can add links to your social media, website, or anything you want to share!',
			image: '🔗'
		},
		{
			title: 'Customize Your Profile',
			description: 'Go to the Appearance tab to upload your avatar, write a bio, and choose a theme that matches your style.',
			image: '🎨'
		},
		{
			title: 'Live Preview',
			description: 'See how your profile looks in real-time on the right side. Any changes you make will appear instantly!',
			image: '📱'
		},
		{
			title: 'Track Your Performance',
			description: 'Check your analytics to see how many people are clicking your links and visiting your profile.',
			image: '📊'
		},
		{
			title: 'Share Your Profile',
			description: 'Copy your profile URL and share it everywhere! Your unique link is: linkbio.com/username',
			image: '🌟'
		}
	];

	onMount(() => {
		// Check if user has seen onboarding
		const hasSeenOnboarding = localStorage.getItem('hasSeenOnboarding');
		if (!hasSeenOnboarding) {
			setTimeout(() => {
				showOnboarding = true;
			}, 1000);
		}
	});

	function nextStep() {
		if (currentStep < steps.length - 1) {
			currentStep++;
		} else {
			completeOnboarding();
		}
	}

	function prevStep() {
		if (currentStep > 0) {
			currentStep--;
		}
	}

	function skipOnboarding() {
		completeOnboarding();
	}

	function completeOnboarding() {
		localStorage.setItem('hasSeenOnboarding', 'true');
		showOnboarding = false;
		currentStep = 0;
	}
</script>

<Dialog.Root bind:open={showOnboarding}>
	<Dialog.Content class="sm:max-w-lg">
		<div class="text-center py-6">
			<!-- Progress -->
			<div class="flex justify-center gap-2 mb-6">
				{#each steps as _, i}
					<div class="h-2 w-8 rounded-full transition-all duration-300 {i === currentStep ? 'bg-gradient-to-r from-purple-600 to-blue-600' : i < currentStep ? 'bg-purple-300' : 'bg-gray-200'}"></div>
				{/each}
			</div>

			<!-- Content -->
			<div class="mb-8">
				<div class="text-6xl mb-4">{steps[currentStep].image}</div>
				<h2 class="text-2xl font-bold text-gray-900 mb-3">{steps[currentStep].title}</h2>
				<p class="text-gray-600 max-w-md mx-auto">{steps[currentStep].description}</p>
			</div>

			<!-- Navigation -->
			<div class="flex items-center justify-between gap-4">
				<Button variant="ghost" on:click={skipOnboarding}>
					Skip Tour
				</Button>
				
				<div class="flex gap-2">
					{#if currentStep > 0}
						<Button variant="outline" on:click={prevStep}>
							Back
						</Button>
					{/if}
					<Button on:click={nextStep} class="bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700">
						{currentStep === steps.length - 1 ? 'Get Started' : 'Next'}
					</Button>
				</div>
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
