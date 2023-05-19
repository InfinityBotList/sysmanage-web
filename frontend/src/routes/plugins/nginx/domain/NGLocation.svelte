<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
	import MultiInput from "$lib/components/MultiInput.svelte";

    interface NGLocation {
        Path: string,
        Proxy?: string,
        Opts?: string[],
    }

    export let location: NGLocation;
    export let i: number;

    let opts: string[] = location.Opts || [];

    // As Opts may be null, we use this reactive array to set it to an empty array if it is null
    $: location.Opts = opts;
</script>

<div>
    <InputSm 
        id={`l-path-${i}`}
        label="Path"
        description="You must have a path called '/' in order for stuff to work correctly"
        placeholder="E.g. /api"
        minlength={1}
        bind:value={location.Path}
    />
    <InputSm 
        id={`l-proxy-${i}`}
        label="Proxy (Optional, leave blank to not proxy)"
        placeholder="E.g. http://localhost:8080"
        minlength={0}
        showErrors={false}
        bind:value={location.Proxy}
    />

    <MultiInput 
        id={`l-opts-${i}`}
        title="Options"
        bind:values={opts}
        minlength={0}
        showErrors={false}
        placeholder="E.g. proxy_set_header and X-Forwarded-For etc."
    />
</div>