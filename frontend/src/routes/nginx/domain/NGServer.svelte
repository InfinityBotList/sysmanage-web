<script lang="ts">
	import InputSm from "$lib/components/InputSm.svelte";
    import MultiInput from "$lib/components/MultiInput.svelte";
	import Section from "$lib/components/Section.svelte";
    import Select from "$lib/components/Select.svelte";
	import NgLocation from "./NGLocation.svelte";

    interface NGServer {
        ID: string,
        Names: string[],
        Comment: string,
        Broken: boolean,
        Locations: NGLocation[],
    }

    interface NGLocation {
        Path: string,
        Proxy?: string,
        Opts?: string[],
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

    <h3 class="text-xl font-semibold">Locations</h3>

    {#each server.Locations as loc, i}
        <Section title={loc.Path}>
            <NgLocation bind:location={loc} i={i} />
        </Section>
    {/each}
</div>