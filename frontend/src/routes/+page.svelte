<script lang="ts">
	import Service from './Service.svelte';
	import InputSm from '$lib/components/InputSm.svelte';
	import ButtonReact from '$lib/components/ButtonReact.svelte';
	import { error } from '$lib/strings';
	import TaskWindow from '../lib/components/TaskWindow.svelte';
	import { newTask } from '$lib/tasks';

	const getServiceList = async () => {
		let serviceList = await fetch(`/api/getServiceList`, {
			method: "POST",
		});

		if(!serviceList.ok) {
			let error = await serviceList.text()

			throw new Error(error)
		} 

		return await serviceList.json();
	}

	let query: string = "";
	let targetFilter: string = "";

	const showService = (
		service: any, 
		query: string,
		targetFilter: string,
	): boolean => {
		let flag = true

		if(query != "" && !service?.ID?.toLowerCase().includes(query.toLowerCase())) {
			flag = false
		}

		if (targetFilter != "") {
			if(targetFilter.startsWith("!")) {
				let target = targetFilter.substring(1)

				if(service?.Service?.Target?.toLowerCase() == target) {
					flag = false
				}
			} else {
				if(service?.Service?.Target?.toLowerCase() != targetFilter.toLowerCase()) {
					flag = false
				}
			}
		}

		return flag
	}

	let buildServicesTaskId: string = "";
	let buildServicesTaskOutput: string[] = [];
	const buildServices = async () => {
		let taskId = await fetch(`/api/buildServices`, {
			method: "POST"
		});

		if(!taskId.ok) {
			let errorStr = await taskId.text()

			error(errorStr)
		}

		buildServicesTaskId = await taskId.text()

		newTask(buildServicesTaskId, (output: string[]) => {
			buildServicesTaskOutput = output
		})
	}
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<h2 class="text-xl font-semibold">Actions</h2>
	<ButtonReact 
		onclick={() => buildServices()}
	>
		Build Services
	</ButtonReact>

	<div class="mb-3"></div>
	
	<h2 class="text-xl font-semibold">Services</h2>

	{#if buildServicesTaskId != ""}
		<h2 class="text-red-500">Build services ID: {buildServicesTaskId}</h2>
		<TaskWindow 
			output={buildServicesTaskOutput}
		/>
	{/if}

	<InputSm
		id="query"
		label="Filter by ID"
		bind:value={query}
		placeholder="E.g. arcadia"
		minlength={0}
	/>
	<InputSm
		id="target-filter"
		label="Filter by systemd target"
		bind:value={targetFilter}
		placeholder="E.g. arcadia"
		minlength={0}
	/>

	{#await getServiceList()}
		<h2 class="text-xl">Loading service list</h2>
	{:then data}
		<div class="flex flex-wrap justify-center items-center justify-evenly">
			{#each data as service}
				{#if showService(service, query, targetFilter)}
					<Service 
						service={service} 
					/>
				{/if}
			{/each}
		</div>
	{:catch err}
		<h2 class="text-red-500">{err}</h2>
	{/await}
</section>