<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let editingBlock: any = null;
	
	const dispatch = createEventDispatcher();

	let content = '';
	let fontSize: 'text-small' | 'text-medium' | 'text-large' | 'headline-small' | 'headline-medium' | 'headline-large' = 'text-medium';
	let textAlign: 'left' | 'center' | 'right' | 'justify' = 'left';
	let isBold = false;
	let isItalic = false;
	let isUnderline = false;
	let isStrikethrough = false;
	let textColor = '#3b82f6';
	let backgroundColor = '';
	let showFontSizeDropdown = false;

	const fontSizeOptions = {
		'Text': [
			{ value: 'text-small', label: 'Small text' },
			{ value: 'text-medium', label: 'Medium text' },
			{ value: 'text-large', label: 'Large text' }
		],
		'Headline': [
			{ value: 'headline-small', label: 'Small headline' },
			{ value: 'headline-medium', label: 'Medium headline' },
			{ value: 'headline-large', label: 'Large headline' }
		]
	};

	$: selectedLabel = Object.values(fontSizeOptions).flat().find(opt => opt.value === fontSize)?.label || 'Medium text';

	$: if (editingBlock) {
		loadData(editingBlock);
	}

	function loadData(block: any) {
		content = block.content || '';
		let styleConfig: any = {};
		try {
			if (block.style) {
				styleConfig = JSON.parse(block.style);
			}
		} catch (e) {
			console.error('Failed to parse style:', e);
		}
		fontSize = styleConfig.fontSize || 'text-medium';
		textAlign = styleConfig.textAlign || 'left';
		isBold = styleConfig.isBold || false;
		isItalic = styleConfig.isItalic || false;
		isUnderline = styleConfig.isUnderline || false;
		isStrikethrough = styleConfig.isStrikethrough || false;
		textColor = styleConfig.textColor || '#3b82f6';
		backgroundColor = styleConfig.backgroundColor || '';
	}

	function handleSave() {
		if (!content.trim()) {
			alert('Please enter text content');
			return;
		}

		dispatch('save', {
			content,
			fontSize,
			textAlign,
			isBold,
			isItalic,
			isUnderline,
			isStrikethrough,
			textColor,
			backgroundColor
		});
		
		// Reset form
		content = '';
		fontSize = 'text-medium';
		textAlign = 'left';
		isBold = false;
		isItalic = false;
		isUnderline = false;
		isStrikethrough = false;
		textColor = '#3b82f6';
		backgroundColor = '';
	}
</script>

<style>
	input[type="color"]::-webkit-color-swatch-wrapper {
		padding: 0;
		border-radius: 9999px;
	}
	input[type="color"]::-webkit-color-swatch {
		border: none;
		border-radius: 9999px;
	}
	input[type="color"]::-moz-color-swatch {
		border: none;
		border-radius: 9999px;
	}
</style>

<div class="bg-white rounded-lg border border-gray-200 p-4 space-y-3">
	<!-- Toolbar -->
	<div class="flex items-center gap-2 flex-wrap">
		<!-- Font Size -->
		<div class="relative min-w-[140px]">
			<button
				onclick={() => showFontSizeDropdown = !showFontSizeDropdown}
				class="w-full bg-gray-50 text-gray-900 px-3 py-1.5 pr-8 rounded-lg text-sm border border-gray-200 hover:border-gray-300 focus:outline-none text-left"
			>
				{selectedLabel}
			</button>
			<svg class="absolute right-2 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
			</svg>

			{#if showFontSizeDropdown}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				<div class="fixed inset-0 z-40" onclick={() => showFontSizeDropdown = false}></div>
				<div class="absolute top-full left-0 mt-1 w-full bg-white rounded-lg shadow-xl z-50 py-1 border border-gray-200">
					{#each Object.entries(fontSizeOptions) as [group, options]}
						<div class="px-2 py-1">
							<div class="text-xs font-medium text-gray-500 px-2 mb-1">{group}</div>
							{#each options as option}
								<button
									onclick={() => { fontSize = option.value as typeof fontSize; showFontSizeDropdown = false; }}
									class="w-full text-left px-2 py-1.5 rounded text-sm text-gray-700 hover:bg-gray-100 transition-colors flex items-center justify-between"
								>
									<span>{option.label}</span>
									{#if fontSize === option.value}
										<svg class="w-3 h-3 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
										</svg>
									{/if}
								</button>
							{/each}
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<div class="h-6 w-px bg-gray-200"></div>

		<!-- Text Align -->
		<div class="flex gap-1">
			<button
				onclick={() => textAlign = 'left'}
				class="p-1.5 rounded {textAlign === 'left' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Align left"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h14"/>
				</svg>
			</button>
			<button
				onclick={() => textAlign = 'center'}
				class="p-1.5 rounded {textAlign === 'center' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Align center"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M5 18h14"/>
				</svg>
			</button>
			<button
				onclick={() => textAlign = 'right'}
				class="p-1.5 rounded {textAlign === 'right' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Align right"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M6 18h14"/>
				</svg>
			</button>
			<button
				onclick={() => textAlign = 'justify'}
				class="p-1.5 rounded {textAlign === 'justify' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Justify"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
				</svg>
			</button>
		</div>

		<div class="h-6 w-px bg-gray-200"></div>

		<!-- Text Style -->
		<div class="flex gap-1">
			<button
				onclick={() => isBold = !isBold}
				class="p-1.5 rounded font-bold text-sm {isBold ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Bold"
			>
				B
			</button>
			<button
				onclick={() => isItalic = !isItalic}
				class="p-1.5 rounded italic text-sm {isItalic ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Italic"
			>
				I
			</button>
			<button
				onclick={() => isUnderline = !isUnderline}
				class="p-1.5 rounded underline text-sm {isUnderline ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Underline"
			>
				U
			</button>
			<button
				onclick={() => isStrikethrough = !isStrikethrough}
				class="p-1.5 rounded line-through text-sm {isStrikethrough ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'}"
				title="Strikethrough"
			>
				S
			</button>
		</div>

		<div class="h-6 w-px bg-gray-200"></div>

		<!-- Colors -->
		<div class="flex items-center gap-1.5">
			<input
				type="color"
				bind:value={textColor}
				class="w-7 h-7 rounded-full cursor-pointer border-2 border-gray-200 hover:border-gray-300"
				title="Text color"
				style="padding: 0; -webkit-appearance: none; -moz-appearance: none; appearance: none;"
			/>
			<div class="relative">
				<input
					type="color"
					bind:value={backgroundColor}
					class="w-7 h-7 rounded-full cursor-pointer border-2 border-gray-200 hover:border-gray-300"
					title="Background color"
					style="padding: 0; -webkit-appearance: none; -moz-appearance: none; appearance: none;"
				/>
				{#if backgroundColor}
					<button
						onclick={() => backgroundColor = ''}
						class="absolute -top-0.5 -right-0.5 w-3.5 h-3.5 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600"
						title="Remove background"
					>
						×
					</button>
				{/if}
			</div>
		</div>
	</div>

	<!-- Text Input -->
	<textarea
		bind:value={content}
		placeholder="Enter your text..."
		rows="4"
		class="w-full px-3 py-2 border border-gray-200 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 resize-none"
		style="
			font-size: {
				fontSize === 'text-small' ? '14px' : 
				fontSize === 'text-medium' ? '16px' : 
				fontSize === 'text-large' ? '18px' :
				fontSize === 'headline-small' ? '20px' :
				fontSize === 'headline-medium' ? '24px' :
				fontSize === 'headline-large' ? '32px' : '16px'
			};
			font-weight: {fontSize.startsWith('headline') ? 'bold' : isBold ? 'bold' : 'normal'};
			text-align: {textAlign};
			font-style: {isItalic ? 'italic' : 'normal'};
			text-decoration: {isUnderline && isStrikethrough ? 'underline line-through' : isUnderline ? 'underline' : isStrikethrough ? 'line-through' : 'none'};
			color: {textColor};
			background-color: {backgroundColor || 'transparent'};
		"
	></textarea>

	<!-- Actions -->
	<div class="flex justify-end">
		<button
			onclick={handleSave}
			class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition-colors"
		>
			{editingBlock ? 'Update Text' : 'Add Text'}
		</button>
	</div>
</div>
