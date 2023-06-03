<script lang="ts">
    import LinkCard from '$lib/components/LinkCard.svelte';
    import GreyText from '$lib/components/GreyText.svelte';
	import ObjectRender from '$lib/components/ObjectRender.svelte';

    export let domain: any;

    export let showDomainInfo: boolean = false;
</script>

<LinkCard 
    title={domain?.Domain}
    link={`/plugins/nginx/domain?id=${domain?.Domain}`}
    linkText="Edit"
    onClickTitle={() => showDomainInfo = !showDomainInfo}
>
    <GreyText>Click the domain title to view more information. Click "Edit" to edit this domain</GreyText>

    {#if showDomainInfo}
        <p class="font-semibold text-lg">More information</p>
        <div class="text-sm">
            {#if domain?.Server?.Servers.length == 0}
                <p class="text-red-500">This is a brand new domain with no servers attached to it</p>
            {:else}
                {#each domain?.Server?.Servers as server}
                    <ObjectRender object={server} />
                    <div class="mt-7"></div>
                {/each}
            {/if}
        </div>
    {/if}
</LinkCard>
