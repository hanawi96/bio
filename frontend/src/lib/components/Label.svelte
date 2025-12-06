<script lang="ts">
	import { cn } from '$lib/utils/cn';
</script>

<label class={cn('block text-sm font-medium text-gray-700 mb-1', $$props.class)}>
	<slot />
</label>
