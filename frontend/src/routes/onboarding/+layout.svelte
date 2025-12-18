<script lang="ts">
	import { onNavigate } from '$app/navigation';
	import { page } from '$app/stores';
	
	// Track navigation direction
	const pageOrder = [
		'/onboarding/setup-url',
		'/onboarding/setup-profile',
		'/onboarding/choose-theme'
	];
	
	onNavigate((navigation) => {
		if (!document.startViewTransition) return;
		
		const currentPath = $page.url.pathname;
		const targetPath = navigation.to?.url.pathname || '';
		
		// Determine direction
		const currentIndex = pageOrder.indexOf(currentPath);
		const targetIndex = pageOrder.indexOf(targetPath);
		const isBackward = targetIndex < currentIndex;
		
		console.log('🔄 Navigation:', {
			from: currentPath,
			to: targetPath,
			currentIndex,
			targetIndex,
			isBackward
		});
		
		// Remove previous class
		document.documentElement.classList.remove('back-navigation');
		
		// Set direction class on document BEFORE transition
		if (isBackward) {
			document.documentElement.classList.add('back-navigation');
			console.log('⬅️ BACKWARD navigation - added .back-navigation class');
		} else {
			console.log('➡️ FORWARD navigation - no class');
		}
		
		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});
</script>

<div class="onboarding-page">
	<slot />
</div>

<style>
	.onboarding-page {
		width: 100%;
		min-height: 100vh;
		view-transition-name: onboarding-page;
	}
	
	:global(body) {
		overflow-x: hidden;
	}
	
	/* Forward navigation (default) */
	@keyframes slide-from-right {
		from {
			transform: translateX(100%);
		}
	}
	
	@keyframes slide-to-left {
		to {
			transform: translateX(-100%);
		}
	}
	
	/* Backward navigation */
	@keyframes slide-from-left {
		from {
			transform: translateX(-100%);
		}
	}
	
	@keyframes slide-to-right {
		to {
			transform: translateX(100%);
		}
	}
	
	/* Default: Forward animation */
	::view-transition-old(onboarding-page) {
		animation: 350ms cubic-bezier(0.4, 0, 0.2, 1) both slide-to-left;
	}
	
	::view-transition-new(onboarding-page) {
		animation: 350ms cubic-bezier(0.4, 0, 0.2, 1) both slide-from-right;
	}
	
	/* Backward: Reverse animation - use :root instead of :global */
	:root.back-navigation::view-transition-old(onboarding-page) {
		animation: 350ms cubic-bezier(0.4, 0, 0.2, 1) both slide-to-right;
	}
	
	:root.back-navigation::view-transition-new(onboarding-page) {
		animation: 350ms cubic-bezier(0.4, 0, 0.2, 1) both slide-from-left;
	}
</style>
