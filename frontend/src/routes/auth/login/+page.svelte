<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api } from '$lib/api/client';
	import { onMount } from 'svelte';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;

	onMount(() => {
		// If already logged in, redirect to dashboard
		if ($auth.user && $auth.token) {
			goto('/dashboard');
		}
	});

	async function handleLogin() {
		loading = true;
		error = '';

		try {
			const response = await api.post<{ user: any; token: string }>('/auth/login', {
				email,
				password
			});

			auth.login(response.user, response.token);
			
			// Check if user needs to complete onboarding
			if (response.user.username && response.user.username.startsWith('temp_')) {
				goto('/onboarding/setup-url');
			} else {
				goto('/dashboard');
			}
		} catch (err: any) {
			error = err.message || 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Login - LinkBio</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
	<div class="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow">
		<div>
			<h2 class="text-3xl font-bold text-center">Welcome back</h2>
			<p class="mt-2 text-center text-gray-600">Sign in to your account</p>
		</div>

		<form on:submit|preventDefault={handleLogin} class="mt-8 space-y-6">
			{#if error}
				<div class="bg-red-50 text-red-600 p-3 rounded">{error}</div>
			{/if}

			<div class="space-y-4">
				<div>
					<label for="email" class="block text-sm font-medium text-gray-700">Email</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						required
						class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>

				<div>
					<label for="password" class="block text-sm font-medium text-gray-700">Password</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						required
						class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
			</div>

			<button
				type="submit"
				disabled={loading}
				class="w-full bg-black text-white py-2 px-4 rounded-md hover:bg-gray-800 disabled:opacity-50"
			>
				{loading ? 'Signing in...' : 'Sign in'}
			</button>

			<p class="text-center text-sm text-gray-600">
				Don't have an account?
				<a href="/auth/register" class="text-black font-medium hover:underline">Sign up</a>
			</p>
		</form>
	</div>
</div>
