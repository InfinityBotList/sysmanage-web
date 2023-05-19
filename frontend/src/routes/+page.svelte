<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import GreyText from "$lib/components/GreyText.svelte";

    interface Link {
        name: string;
        description: string;
        link: string;
        plugin?: string;
    }

    const links: Link[] = [
        {
            name: "Service Management",
            description: "Systemd service management",
            link: "/plugins/systemd",
            plugin: "systemd"
        },
        {
            name: "Nginx Management",
            description: "Add, update, remove and manage nginx-proxied domains",
            link: "/plugins/nginx",
            plugin: "nginx"
        },
    ]
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