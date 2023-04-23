<script lang="ts">
    import LinkCard from '$lib/components/LinkCard.svelte';
    import GreyText from '$lib/components/GreyText.svelte';

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
            {#each Object.entries(server) as [key, value]}
                <p>
                    <span class="font-semibold">{key}:</span> {value}
                </p>
            {/each}
            <div class="mt-3"></div>
        {/each}
    {/if}
</LinkCard>
