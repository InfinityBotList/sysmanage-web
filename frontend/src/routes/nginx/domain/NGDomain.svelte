<script lang="ts">
	import NgServer from "./NGServer.svelte";
    import ObjectRender from "$lib/components/ObjectRender.svelte";
    import Section from "$lib/components/Section.svelte";

    export let domain: any;
</script>

<h1 class="text-2xl font-semibold">Viewing {domain?.Domain}</h1>

<h2 class="text-xl font-semibold">Server List</h2>

<div class="flex flex-col space-y-2">
    {#each domain?.Server?.Servers as server, i}
        <Section title={server.ID}>
            <NgServer bind:server={server} i={i} />
        </Section>
    {/each}
</div>
<Section title="Tree View">
    <ObjectRender object={domain} />
</Section>