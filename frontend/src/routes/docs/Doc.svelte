<script lang="ts">
    import Card from "$lib/components/DefaultCard.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import { title } from "$lib/strings";

    let showContent = false;

    export let docTitle: string;
</script>

<Card 
    title={title(docTitle.replaceAll("-", " ").replaceAll(".md", ""))}
    onclick={() => {
        showContent = !showContent
    }}
>
    {#if showContent}
        <slot />
    {:else}
        <GreyText>Click for more info</GreyText>
    {/if}
</Card>	