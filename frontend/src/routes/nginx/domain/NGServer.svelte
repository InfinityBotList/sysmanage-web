<script lang="ts">
	import InputSm from "$lib/components/InputSm.svelte";
    import MultiInput from "$lib/components/MultiInput.svelte";
    import Select from "$lib/components/Select.svelte";

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
    <h3 class="text-xl font-semibold">Editting {server.ID}</h3>

    <MultiInput 
        id={`s-names-${i}`} 
        title="Subdomain"
        bind:values={server.Names} 
        placeholder="example.com, www.example.com etc."
        minlength={3}
    />

    <div class="mb-2"></div>

    <InputSm
        id={`s-id-${i}`}
        label="ID"
        placeholder="E.g. popplio, arcadia-rpc etc."
        minlength={1}
        bind:value={server.ID}
    />

    <InputSm
        id={`s-comment-${i}`}
        label="Comment"
        placeholder="E.g. Popplio Web API"
        minlength={1}
        bind:value={server.Comment}
    />
    <Select
        name="broken"
        placeholder="Is the server broken/disabled?"
        bind:valueBool={server.Broken}
        options={new Map([
            ["Yes, it is", "true"],
            ["No, its not", "false"],
        ])}
    />
</div>