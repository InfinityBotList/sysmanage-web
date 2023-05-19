<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
    import { links } from "$lib/links";
</script>

<svelte:head>
	<title>Home</title>
</svelte:head>

<section>
	<h2 class="text-xl font-semibold">What do you want to do?</h2>

    {#each links as link}
        <GreyText>{link.description}</GreyText>
        <Button link={link.link}>{link.name}</Button>
    {/each}
</section>