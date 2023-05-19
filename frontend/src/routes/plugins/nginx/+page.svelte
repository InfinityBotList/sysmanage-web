<script lang="ts">
	import InputSm from '$lib/components/InputSm.svelte';
	import ButtonReact from '$lib/components/ButtonReact.svelte';
	import { error, success } from '$lib/strings';
	import TaskWindow from '../../../lib/components/TaskWindow.svelte';
	import { newTask } from '$lib/tasks';
    import Domain from './Domain.svelte';
	import Button from '$lib/components/Button.svelte';

	let buildNginxTaskId: string = "";
	let buildNginxTaskOutput: string[] = [];
	const buildNginx = async () => {
		let taskId = await fetch(`/api/nginx/buildNginx`, {
			method: "POST"
		});

		if(!taskId.ok) {
			let errorStr = await taskId.text()

			error(errorStr)

			return
		}

		buildNginxTaskId = await taskId.text()

		newTask(buildNginxTaskId, (output: string[]) => {
			buildNginxTaskOutput = output
		})
	}

	let updateDnsRecordCfId: string = "";
	let updateDnsRecordCfOutput: string[] = [];
	const updateDnsRecordCf = async () => {
		let taskId = await fetch(`/api/nginx/updateDnsRecordCf`, {
			method: "POST"
		});

		if(!taskId.ok) {
			let errorStr = await taskId.text()

			error(errorStr)

			return
		}

		updateDnsRecordCfId = await taskId.text()

		newTask(updateDnsRecordCfId, (output: string[]) => {
			updateDnsRecordCfOutput = output
		})
	}

	const getNginxDomainList = async () => {
		let domList = await fetch(`/api/nginx/getDomainList`, {
			method: "POST",
		});

		if(!domList.ok) {
			let error = await domList.text()

			throw new Error(error)
		} 

		return await domList.json();
	}

	const showDomain = (
		domainObj: any, 
		domain: string,
	): boolean => {
		let flag = true

		if(domain != "" && !domainObj?.Domain?.toLowerCase().includes(domain.toLowerCase())) {
			flag = false
		}

		return flag
	}

    let domainQuery: string;
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<h2 class="text-xl font-semibold">Actions</h2>
	<ButtonReact 
		onclick={() => buildNginx()}
	>
		Build Nginx
	</ButtonReact>
	<ButtonReact 
		onclick={() => updateDnsRecordCf()}
	>
		Update Cloudflare
	</ButtonReact>
	<Button 
		link="/plugins/nginx/new"
	>
		Add Domain
	</Button>

	<div class="mb-3"></div>
	
	<h2 class="text-xl font-semibold">Services</h2>

	{#if buildNginxTaskId != ""}
		<TaskWindow 
			output={buildNginxTaskOutput}
		/>
	{/if}

	{#if updateDnsRecordCfId != ""}
		<TaskWindow 
			output={updateDnsRecordCfOutput}
		/>
	{/if}

	<InputSm
		id="domain"
		label="Filter by domain"
		bind:value={domainQuery}
		placeholder="E.g. arcadia"
		showErrors={false}
		minlength={0}
	/>

	{#await getNginxDomainList()}
		<h2 class="text-xl">Loading domain list</h2>
	{:then data}
		<div class="flex flex-wrap justify-center items-center justify-evenly">
			{#each data as domain}
				{#if showDomain(domain, domainQuery)}
                    <Domain domain={domain} />
				{/if}
			{/each}
		</div>
	{:catch err}
		<h2 class="text-red-500">{err}</h2>
	{/await}
</section>
