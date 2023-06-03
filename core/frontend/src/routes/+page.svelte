<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
    import { getLinks } from "$lib/corelib/links";
</script>

<svelte:head>
	<title>Home</title>
</svelte:head>

<section>
	<h2 class="text-xl font-semibold">What do you want to do?</h2>

    {#await getLinks()}
        <h2 class="text-xl">Loading...</h2>
    {:then links}
        {#each links as link}
            <h3 class="text-lg font-semibold mt-5 mb-1">{link.Title}</h3>
            <Button link={link.Href}>{link.LinkText}</Button>
            <GreyText>{link.Description}</GreyText>
        {/each}
    {:catch error}
        <h2 class="text-xl">Error: {error?.message ? error?.message : error}</h2>
    {/await}
</section>