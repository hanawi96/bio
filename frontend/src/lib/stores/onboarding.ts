import { writable } from 'svelte/store';

interface OnboardingData {
	avatarPreview: string;
	bio: string;
}

function createOnboardingStore() {
	const { subscribe, set, update } = writable<OnboardingData>({
		avatarPreview: '',
		bio: ''
	});

	return {
		subscribe,
		setAvatar: (preview: string) => update(state => ({ ...state, avatarPreview: preview })),
		setBio: (bio: string) => update(state => ({ ...state, bio })),
		setData: (data: Partial<OnboardingData>) => update(state => ({ ...state, ...data })),
		clear: () => set({ avatarPreview: '', bio: '' })
	};
}

export const onboardingData = createOnboardingStore();
