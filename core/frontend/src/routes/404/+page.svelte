<script lang="ts">
	import Button from "$lib/components/Button.svelte";

	const findUrl = async () => {
		let searchParams = new URLSearchParams(window.location.search);

		return searchParams.get("from") || "";
	}
</script>

<svelte:head>
	<title>404 - Not Found</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<h1 class="font-semibold">404 - Not Found</h1>
<p>Sorry, the page you were looking for was not found.</p>

{#await findUrl() then url}
	<p>URL: {url}</p>
	<Button link={url}>Retry</Button>
{/await}