import { writable } from 'svelte/store';

interface PendingChanges {
	theme?: Record<string, any>;
	header?: Record<string, any>;
	linkStyles?: Record<string, any>;
	hasChanges: boolean;
}

function createPendingChangesStore() {
	const { subscribe, set, update } = writable<PendingChanges>({
		hasChanges: false
	});

	return {
		subscribe,
		updateTheme(changes: Record<string, any>) {
			update(state => ({
				...state,
				theme: { ...state.theme, ...changes },
				hasChanges: true
			}));
		},
		updateHeader(changes: Record<string, any>) {
			update(state => ({
				...state,
				header: { ...state.header, ...changes },
				hasChanges: true
			}));
		},
		updateLinkStyles(changes: Record<string, any>) {
			update(state => ({
				...state,
				linkStyles: { ...state.linkStyles, ...changes },
				hasChanges: true
			}));
		},
		reset() {
			set({ hasChanges: false });
		},
		getAll(): PendingChanges {
			let current: PendingChanges = { hasChanges: false };
			subscribe(v => current = v)();
			return current;
		}
	};
}

export const pendingChanges = createPendingChangesStore();
