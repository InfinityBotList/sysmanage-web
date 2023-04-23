<script lang="ts">
    import LinkCard from '$lib/components/LinkCard.svelte';
    import GreyText from '$lib/components/GreyText.svelte';
	import ObjectRender from '$lib/components/ObjectRender.svelte';

    export let domain: any;

    export let showDomainInfo: boolean = false;
</script>

<LinkCard 
    title={domain?.Domain}
    link={`/nginx/domain?domain=${domain?.Domain}`}
    linkText="Edit"
    onClickTitle={() => showDomainInfo = !showDomainInfo}
>
    <GreyText>Click the domain title to view more information. Click "Edit" to edit this domain</GreyText>

    {#if showDomainInfo}
        {#each domain?.Server?.Servers as server}
            <ObjectRender object={server} />
            <div class="mt-7"></div>
        {/each}
    {/if}
</LinkCard>
