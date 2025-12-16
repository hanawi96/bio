import { writable } from 'svelte/store';

interface ThemeSnapshot {
	theme: Record<string, any>;
	header: Record<string, any>;
	customHeaderPresets?: any[];
}

interface PendingChanges {
	hasChanges: boolean;
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
	
	let savedSnapshot: ThemeSnapshot | null = null;
	let currentSnapshot: ThemeSnapshot | null = null;

	function checkHasChanges(): boolean {
		if (!savedSnapshot || !currentSnapshot) return false;
		
		// Compare theme and header
		const themeChanged = !deepEqual(currentSnapshot.theme, savedSnapshot.theme);
		const headerChanged = !deepEqual(currentSnapshot.header, savedSnapshot.header);
		
		return themeChanged || headerChanged;
	}

	return {
		subscribe,
		
		// Set saved snapshot (from DB)
		setSavedSnapshot(snapshot: ThemeSnapshot) {
			savedSnapshot = JSON.parse(JSON.stringify(snapshot));
			currentSnapshot = JSON.parse(JSON.stringify(snapshot));
			set({ hasChanges: false });
		},
		
		// Update current snapshot (when user edits)
		updateCurrentSnapshot(snapshot: ThemeSnapshot) {
			currentSnapshot = JSON.parse(JSON.stringify(snapshot));
			const hasChanges = checkHasChanges();
			update(state => ({ ...state, hasChanges }));
		},
		
		// For link styles (separate tracking)
		updateLinkStyles(changes: Record<string, any>) {
			update(state => ({
				...state,
				linkStyles: { ...state.linkStyles, ...changes },
				hasChanges: true
			}));
		},
		
		reset() {
			savedSnapshot = null;
			currentSnapshot = null;
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
