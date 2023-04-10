<script>
	const getServiceList = async () => {
		let serviceList = await fetch(`/api/getServiceList`, {
			method: "POST" 
		});

		if(!serviceList.ok) {
			let error = await serviceList.text()

			throw new Error(error)
		} 

		return await serviceList.json();
	}
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	{#await getServiceList()}
		<h2>Loading service list</h2>
	{:then data}
		{#each data as service}
			<h2>{JSON.stringify(service)}</h2>
		{/each}
	{:catch err}
		<h2 class="text-red-500">{err}</h2>
	{/await}
</section>