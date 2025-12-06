<script lang="ts">
	import { cn } from '$lib/utils/cn';
</script>

<div class={cn('bg-white rounded-lg shadow-sm border border-gray-200 p-6', $$props.class)}>
	<slot />
</div>
