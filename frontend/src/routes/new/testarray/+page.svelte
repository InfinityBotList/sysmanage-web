<script lang="ts">
	import MultiInput from "$lib/components/MultiInput.svelte";

    let testData = ["foo", "bar", "baz"]
</script>

<MultiInput
    bind:values={testData}
    title="Test parameter"
    placeholder="Enter whatever you want"
    minlength={3}
/>

<p>{JSON.stringify(testData)}</p>