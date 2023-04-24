<script lang="ts">
	import InputSm from "$lib/components/InputSm.svelte";
import MultiInput from "$lib/components/MultiInput.svelte";

    interface NGServer {
        Names: string[],
        Comment: string,
        Broken: boolean,
        Location: NGLocation,
    }

    interface NGLocation {
        Path: string,
        Proxy?: string,
        Opts?: KV[],
    }

    interface KV {
        Name: string,
        Value: string,
    }

    export let server: NGServer;
    export let i: number;
</script>

<div>
    <MultiInput 
        id={`s-names-${i}`} 
        title="Names" 
        bind:values={server.Names} 
        placeholder="example.com, www.example.com etc."
        minlength={3}
    />

    <InputSm
        id={`s-comment-${i}`}
        label="Comment"
        placeholder="E.g. Popplio Web API"
        minlength={1}
        bind:value={server.Comment}
    />
</div>