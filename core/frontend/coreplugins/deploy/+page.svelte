<script lang="ts">
	import InputSm from '$lib/components/InputSm.svelte';

	const getDeployList = async () => {
		let depList = await fetch(`/api/deploy/getDeployList`, {
			method: "POST",
		});

		if(!depList.ok) {
			let error = await depList.text()

			throw new Error(error)
		} 

		return await depList.json();
	}

	const showDeployMeta = (
		deployObj: any, 
		deploy: string,
	): boolean => {
		let flag = true

		if(deploy != "" && !deployObj?.ID?.includes(deploy.toLowerCase())) {
			flag = false
		}

		return flag
	}

    let depQuery: string;
</script>

<svelte:head>
	<title>Deploy Management</title>
	<meta name="description" content="Deploy view" />
</svelte:head>

<section>	
	<h2 class="text-xl font-semibold">Deploy List</h2>

	<InputSm
		id="id"
		label="Filter by id"
		bind:value={depQuery}
		placeholder="E.g. arcadia"
		showErrors={false}
		minlength={0}
	/>

	{#await getDeployList()}
		<h2 class="text-xl">Loading deploy list</h2>
	{:then data}
		<div class="flex flex-wrap justify-center items-center justify-evenly">
			{#each data as deploy}
				{#if showDeployMeta(deploy, depQuery)}
                    {JSON.stringify(deploy)}
				{/if}
			{/each}
		</div>
	{:catch err}
		<h2 class="text-red-500">{err}</h2>
	{/await}
</section>
