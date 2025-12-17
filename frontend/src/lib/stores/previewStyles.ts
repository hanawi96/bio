import { writable } from 'svelte/store';

export interface PreviewStyles {
	text_alignment?: 'left' | 'center' | 'right';
	text_size?: 'S' | 'M' | 'L' | 'XL';
	card_text_color?: string;
	image_shape?: 'sharp' | 'square' | 'circle';
	enableCardBackground?: boolean;
	show_shadow?: boolean;
	shadow_x?: number;
	shadow_y?: number;
	shadow_blur?: number;
	card_border_radius?: number;
	has_card_border?: boolean;
	card_border_color?: string;
	card_border_width?: number;
	card_background_color?: string;
	card_background_opacity?: number;
	style?: string;
}

function createPreviewStylesStore() {
	const { subscribe, set, update } = writable<PreviewStyles>({});
	
	return {
		subscribe,
		update: (styles: PreviewStyles) => update(current => ({ ...current, ...styles })),
		reset: () => set({})
	};
}

export const previewStyles = createPreviewStylesStore();
