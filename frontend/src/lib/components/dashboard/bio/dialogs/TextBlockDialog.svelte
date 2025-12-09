<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';

	export let open = false;
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
	let activeTab: 'content' | 'settings' | 'section' = 'content';
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

	$: if (open && editingBlock) {
		loadInitialData();
	} else if (open && !editingBlock) {
		resetForm();
	}

	function loadInitialData() {
		content = editingBlock.content || '';
		
		let styleConfig: any = {};
		try {
			if (editingBlock.style) {
				styleConfig = JSON.parse(editingBlock.style);
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
		activeTab = 'content';
	}

	function resetForm() {
		content = '';
		fontSize = 'text-medium';
		textAlign = 'left';
		isBold = false;
		isItalic = false;
		isUnderline = false;
		isStrikethrough = false;
		textColor = '#3b82f6';
		activeTab = 'content';
	}

	function handleSave() {
		if (!content.trim()) {
			alert('Vui lòng nhập nội dung');
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
			textColor
		});
		
		resetForm();
		open = false;
	}

	function handleCancel() {
		resetForm();
		open = false;
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

<Dialog.Root bind:open onOpenChange={(isOpen) => { if (!isOpen) handleCancel(); }}>
	<Dialog.Content class="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title class="text-xl font-semibold">Text</Dialog.Title>
		</Dialog.Header>

		<!-- Tabs -->
		<div class="flex gap-6 mb-6 border-b border-gray-200">
			<button
				onclick={() => activeTab = 'content'}
				class="pb-3 px-1 text-sm font-medium transition-colors relative {activeTab === 'content' ? 'text-indigo-600' : 'text-gray-500 hover:text-gray-700'}"
			>
				CONTENT
				{#if activeTab === 'content'}
					<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600"></div>
				{/if}
			</button>
			<button
				onclick={() => activeTab = 'settings'}
				class="pb-3 px-1 text-sm font-medium transition-colors relative {activeTab === 'settings' ? 'text-indigo-600' : 'text-gray-500 hover:text-gray-700'}"
			>
				SETTINGS
				{#if activeTab === 'settings'}
					<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600"></div>
				{/if}
			</button>
			<button
				onclick={() => activeTab = 'section'}
				class="pb-3 px-1 text-sm font-medium transition-colors relative {activeTab === 'section' ? 'text-indigo-600' : 'text-gray-500 hover:text-gray-700'}"
			>
				SECTION
				{#if activeTab === 'section'}
					<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600"></div>
				{/if}
			</button>
		</div>

		<!-- Content Tab -->
		{#if activeTab === 'content'}
			<div class="space-y-4">
				<!-- Font Size & Formatting Toolbar -->
				<div class="flex items-center gap-3 flex-wrap">
					<!-- Font Size Dropdown -->
					<div class="relative min-w-[160px]">
						<button
							onclick={() => showFontSizeDropdown = !showFontSizeDropdown}
							class="w-full bg-gray-50 text-gray-900 px-4 py-2 pr-10 rounded-lg text-sm border border-gray-200 hover:border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer text-left"
						>
							{selectedLabel}
						</button>
						<svg class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
						</svg>

						{#if showFontSizeDropdown}
							<div 
								class="fixed inset-0 z-40" 
								onclick={() => showFontSizeDropdown = false}
							></div>
							<div class="absolute top-full left-0 mt-1 w-64 bg-gray-800 rounded-lg shadow-xl z-50 py-2 border border-gray-700">
								{#each Object.entries(fontSizeOptions) as [group, options]}
									<div class="px-3 py-2">
										<div class="text-xs font-medium text-gray-400 mb-2">{group}</div>
										{#each options as option}
											<button
												onclick={() => { fontSize = option.value; showFontSizeDropdown = false; }}
												class="w-full text-left px-3 py-2 rounded text-sm text-gray-200 hover:bg-gray-700 transition-colors flex items-center justify-between"
											>
												<span>{option.label}</span>
												{#if fontSize === option.value}
													<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

					<!-- Text Align Buttons -->
					<div class="flex gap-1">
						<button
							onclick={() => textAlign = 'left'}
							class="p-2 rounded {textAlign === 'left' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Align left"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h14"/>
							</svg>
						</button>
						<button
							onclick={() => textAlign = 'center'}
							class="p-2 rounded {textAlign === 'center' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Align center"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M5 18h14"/>
							</svg>
						</button>
						<button
							onclick={() => textAlign = 'right'}
							class="p-2 rounded {textAlign === 'right' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Align right"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M6 18h14"/>
							</svg>
						</button>
						<button
							onclick={() => textAlign = 'justify'}
							class="p-2 rounded {textAlign === 'justify' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Justify"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
							</svg>
						</button>
					</div>

					<div class="h-6 w-px bg-gray-200"></div>

					<!-- Text Style Buttons -->
					<div class="flex gap-1">
						<button
							onclick={() => isBold = !isBold}
							class="p-2 rounded font-bold {isBold ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Bold"
						>
							B
						</button>
						<button
							onclick={() => isItalic = !isItalic}
							class="p-2 rounded italic {isItalic ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Italic"
						>
							I
						</button>
						<button
							onclick={() => isStrikethrough = !isStrikethrough}
							class="p-2 rounded line-through {isStrikethrough ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Strikethrough"
						>
							S
						</button>
						<button
							onclick={() => isUnderline = !isUnderline}
							class="p-2 rounded underline {isUnderline ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
							title="Underline"
						>
							U
						</button>
					</div>

					<div class="h-6 w-px bg-gray-200"></div>

					<!-- Color Picker -->
					<div class="relative">
						<input
							type="color"
							bind:value={textColor}
							class="w-8 h-8 rounded-full cursor-pointer border-2 border-gray-200 hover:border-gray-300 overflow-hidden"
							style="padding: 0; -webkit-appearance: none; -moz-appearance: none; appearance: none;"
							title="Text color"
						/>
					</div>
				</div>

				<!-- Text Input Area -->
				<div class="relative">
					<textarea
						bind:value={content}
						placeholder="Enter text"
						rows="8"
						class="w-full px-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 resize-none"
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
						"
					></textarea>
					<button class="absolute bottom-3 right-3 p-2 text-gray-500 hover:text-gray-700">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
						</svg>
					</button>
				</div>
			</div>
		{/if}

		<!-- Settings Tab -->
		{#if activeTab === 'settings'}
			<div class="space-y-4">
				<p class="text-gray-500 text-sm">Additional settings coming soon...</p>
			</div>
		{/if}

		<!-- Section Tab -->
		{#if activeTab === 'section'}
			<div class="space-y-4">
				<p class="text-gray-500 text-sm">Section settings coming soon...</p>
			</div>
		{/if}

		<!-- Footer Actions -->
		<div class="flex items-center justify-between mt-8 pt-6 border-t border-gray-200">
			<button
				onclick={handleCancel}
				class="flex items-center gap-2 px-4 py-2 text-gray-500 hover:text-gray-700 transition-colors"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
				</svg>
				Back
			</button>
			<button
				onclick={handleSave}
				class="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg transition-colors"
			>
				Save changes
			</button>
		</div>
	</Dialog.Content>
</Dialog.Root>
