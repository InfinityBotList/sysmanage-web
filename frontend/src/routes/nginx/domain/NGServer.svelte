<script lang="ts">
	import InputSm from "$lib/components/InputSm.svelte";
import MultiInput from "$lib/components/MultiInput.svelte";

    interface NGServer {
        ID: string,
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
    <p class="font-semibold">{server.ID}</p>

    <MultiInput 
        id={`s-names-${i}`} 
        title="Subdomains" 
        bind:values={server.Names} 
        placeholder="example.com, www.example.com etc."
        minlength={3}
    />

    <div class="mb-2"></div>

    <InputSm
        id={`s-comment-${i}`}
        label="Comment"
        placeholder="E.g. Popplio Web API"
        minlength={1}
        bind:value={server.Comment}
    />
</div>