<script>
	import Service from './Service.svelte';

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

	let getServiceListPromise = getServiceList();
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	{#await getServiceListPromise}
		<h2 class="text-red">Loading service list</h2>
	{:then data}
		<div class="flex flex-wrap justify-center items-center justify-evenly">
			{#each data as service}
				<Service 
					service={service} 
				/>
				<!--<h2>{JSON.stringify(service)}</h2>-->
			{/each}
		</div>
	{:catch err}
		<h2 class="text-red-500">{err}</h2>
	{/await}
</section>