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
		// If already logged in, redirect to appropriate page
		if ($auth.user && $auth.token) {
			if ($auth.user.username && $auth.user.username.startsWith('temp_')) {
				goto('/onboarding/setup-url');
			} else {
				goto('/dashboard');
			}
		}
	});

	async function handleRegister() {
		loading = true;
		error = '';

		try {
			const response = await api.post<{ user: any; token: string }>('/auth/register', {
				email,
				password
			});

			auth.login(response.user, response.token);
			goto('/onboarding/setup-url');
		} catch (err: any) {
			console.error('Registration error:', err);
			const message = err.message || 'Registration failed';
			
			// Handle specific errors
			if (message.includes('email already exists') || message.includes('duplicate')) {
				error = 'This email is already registered. Please login instead.';
			} else {
				error = message;
			}
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Sign Up - LinkBio</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
	<div class="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow">
		<div>
			<h2 class="text-3xl font-bold text-center">Create your account</h2>
			<p class="mt-2 text-center text-gray-600">Start building your link in bio</p>
		</div>

		<form on:submit|preventDefault={handleRegister} class="mt-8 space-y-6">
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
						minlength="8"
						class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
			</div>

			<button
				type="submit"
				disabled={loading}
				class="w-full bg-black text-white py-2 px-4 rounded-md hover:bg-gray-800 disabled:opacity-50"
			>
				{loading ? 'Creating account...' : 'Create account'}
			</button>

			<p class="text-center text-sm text-gray-600">
				Already have an account?
				<a href="/auth/login" class="text-black font-medium hover:underline">Sign in</a>
			</p>
		</form>
	</div>
</div>
