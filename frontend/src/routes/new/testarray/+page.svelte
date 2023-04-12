<script lang="ts">
	import MultiInput from "$lib/components/MultiInput.svelte";
    import KVMultiInput from "$lib/components/KVMultiInput.svelte";

    let testData = ["foo", "bar", "baz"]

    let stringstring: [string, string][] = [
        ["foo", "bar"],
        ["baz", "qux"],
        ["quux", "quuz"],
    ]
</script>

<MultiInput
    id="test"
    bind:values={testData}
    title="Test parameter"
    placeholder="Enter whatever you want"
    minlength={3}
/>
<KVMultiInput
    id="test"
    bind:values={stringstring}
    title="Test parameter"
    placeholder="Enter whatever you want"
    minlength={3}
/>

<p>{JSON.stringify(testData)}</p>
<p>{JSON.stringify(stringstring)}</p>