<script lang="ts">
	export let output: string[];
</script>

<pre class="p-4 bg-gray-300 text-black dark:bg-gray-400 break-all whitespace-pre-wrap">{output?.map(f => f.replaceAll("\r", "\n"))?.join("")}</pre>