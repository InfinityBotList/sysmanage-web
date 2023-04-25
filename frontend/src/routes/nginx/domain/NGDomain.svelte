<script lang="ts">
	import NgServer from "./NGServer.svelte";
    import ObjectRender from "$lib/components/ObjectRender.svelte";
    import Section from "$lib/components/Section.svelte";
	import DangerButton from "$lib/components/DangerButton.svelte";

    interface NgDomain {
        Domain: string,
        Server: NgServerList
    }

    interface NGServerInterface {
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

    interface NgServerList {
        Servers: NGServerInterface[]
    }

    export let domain: NgDomain;
</script>

<h1 class="text-2xl font-semibold">Viewing {domain?.Domain}</h1>

<h2 class="text-xl font-semibold">Server List</h2>

<div class="flex flex-col space-y-2">
    {#each domain?.Server?.Servers as server, i}
        <Section title={server.ID}>
            <DangerButton 
                onclick={() => {
                    // Delete the server
                    domain.Server.Servers = domain.Server.Servers.filter((_, index) => index !== i);      
                }}
            >
                Delete Server
            </DangerButton>    
            <NgServer bind:server={server} i={i} />
        </Section>
    {/each}
    <hr class="mt-2 mb-2" />
    <Section title="Tree View">
        <ObjectRender object={domain} />
    </Section>
</div>
