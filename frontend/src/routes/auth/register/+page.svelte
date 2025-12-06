<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api } from '$lib/api/client';

	let email = '';
	let username = '';
	let password = '';
	let error = '';
	let loading = false;

	async function handleRegister() {
		loading = true;
		error = '';

		try {
			const response = await api.post<{ user: any; token: string }>('/auth/register', {
				email,
				username,
				password
			});

			auth.login(response.user, response.token);
			goto('/dashboard');
		} catch (err: any) {
			console.error('Registration error:', err);
			error = err.message || 'Registration failed';
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
					<label for="username" class="block text-sm font-medium text-gray-700">Username</label>
					<input
						id="username"
						type="text"
						bind:value={username}
						required
						pattern="[a-zA-Z0-9_]+"
						class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
					<p class="mt-1 text-sm text-gray-500">Your profile URL: linkbio.com/{username || 'username'}</p>
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
