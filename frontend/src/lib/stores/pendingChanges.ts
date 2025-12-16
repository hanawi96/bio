import { writable } from 'svelte/store';

interface PendingChanges {
	theme?: Record<string, any>;
	header?: Record<string, any>;
	linkStyles?: Record<string, any>;
	hasChanges: boolean;
}

interface OriginalValues {
	theme?: Record<string, any>;
	header?: Record<string, any>;
	linkStyles?: Record<string, any>;
}

function deepEqual(obj1: any, obj2: any): boolean {
	if (obj1 === obj2) return true;
	if (!obj1 || !obj2) return false;
	if (typeof obj1 !== 'object' || typeof obj2 !== 'object') return false;
	
	const keys1 = Object.keys(obj1);
	const keys2 = Object.keys(obj2);
	
	if (keys1.length !== keys2.length) return false;
	
	for (const key of keys1) {
		if (!keys2.includes(key)) return false;
		if (!deepEqual(obj1[key], obj2[key])) return false;
	}
	
	return true;
}

function createPendingChangesStore() {
	const { subscribe, set, update } = writable<PendingChanges>({
		hasChanges: false
	});
	
	let originalValues: OriginalValues = {};

	function checkHasChanges(current: PendingChanges): boolean {
		if (current.theme && !deepEqual(current.theme, originalValues.theme)) return true;
		if (current.header && !deepEqual(current.header, originalValues.header)) return true;
		if (current.linkStyles && !deepEqual(current.linkStyles, originalValues.linkStyles)) return true;
		return false;
	}

	return {
		subscribe,
		setOriginal(values: OriginalValues) {
			originalValues = JSON.parse(JSON.stringify(values));
		},
		updateTheme(changes: Record<string, any>) {
			update(state => {
				const newState = {
					...state,
					theme: changes
				};
				newState.hasChanges = checkHasChanges(newState);
				return newState;
			});
		},
		updateHeader(changes: Record<string, any>) {
			update(state => {
				const newState = {
					...state,
					header: changes
				};
				newState.hasChanges = checkHasChanges(newState);
				return newState;
			});
		},
		updateLinkStyles(changes: Record<string, any>) {
			update(state => {
				const newState = {
					...state,
					linkStyles: { ...state.linkStyles, ...changes }
				};
				newState.hasChanges = checkHasChanges(newState);
				return newState;
			});
		},
		reset() {
			set({ hasChanges: false });
			originalValues = {};
		},
		getAll(): PendingChanges {
			let current: PendingChanges = { hasChanges: false };
			subscribe(v => current = v)();
			return current;
		}
	};
}

export const pendingChanges = createPendingChangesStore();
