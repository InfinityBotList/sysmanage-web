<script lang="ts">
	import Service from './service/Service.svelte';
	import InputSm from '$lib/components/InputSm.svelte';
	import ButtonReact from '$lib/components/ButtonReact.svelte';
	import { error, success } from '$lib/corelib/strings';
	import TaskWindow from '$lib/components/TaskWindow.svelte';
	import { newTask } from '$lib/corelib/tasks';
	import DangerButton from '$lib/components/DangerButton.svelte';
	import Button from '$lib/components/Button.svelte';

	const getServiceList = async () => {
        const categories: {
            category: string,
            services: any[]
        }[] = [];

		let serviceList = await fetch(`/api/systemd/getServiceList`, {
			method: "POST",
		});

		if(!serviceList.ok) {
			let error = await serviceList.text()
			throw new Error(error)
		} 


        const p = await serviceList.json();
		p.forEach((service: any) => {
            if (!categories.find((i) => i.category === service.Service.Target)) categories.push({
                category: service.Service.Target,
                services: [...service]
            });
            else {
                let category = categories.find((i) => i.category === service.Service.Target);
                category?.services.push(service);
            }
        });

        return categories;
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
		let taskId = await fetch(`/api/systemd/buildServices`, {
			method: "POST"
		});

		if(!taskId.ok) {
			let errorStr = await taskId.text()

			error(errorStr)

			return
		}

		buildServicesTaskId = await taskId.text()

		newTask(buildServicesTaskId, (output: string[]) => {
			buildServicesTaskOutput = output
		})
	}

	const restartServer = async () => {
		let confirm = window.prompt("Are you sure you want to restart the server? (YES to confirm))")

		if(confirm != "YES") {
			return
		}

		let res = await fetch(`/api/systemd/restartServer`, {
			method: "POST"
		});

		if(!res.ok) {
			let errorStr = await res.text()

			error(errorStr)

			return
		}

		success("Server is now restarting...")
	}

	const srvmod = async (action: string) => {
		let confirm = window.prompt(`Are you sure you want to ${action} all services? (YES to confirm)`)

		if(confirm != "YES") {
			return
		}

		let res = await fetch(`/api/systemd/serviceMod?act=${action}`, {
			method: "POST",
		});

		if(!res.ok) {
			let errorStr = await res.text()

			error(errorStr)

			return
		}

		success(`All services are now ${action}...`)
	}

	let showDangerous = false;
</script>

<svelte:head>
	<title>Systemd service list</title>
</svelte:head>

<section>
	<h2 class="text-xl font-semibold">Actions</h2>
	<ButtonReact 
		onclick={() => buildServices()}
	>
		Build Services
	</ButtonReact>
	<Button 
		link="/plugins/systemd/new"
	>
		New Service
	</Button>
	<Button 
		link="/plugins/systemd/meta"
	>
		Meta Editor
	</Button>
	{#if showDangerous}
		<DangerButton 
			onclick={() => showDangerous = false}
		>
			Hide Dangerous Actions
		</DangerButton>
	{:else}
		<DangerButton 
			onclick={() => showDangerous = true}
		>
			Show Dangerous Actions
		</DangerButton>
	{/if}

	{#if showDangerous}
		<h2 class="mt-2 text-xl font-semibold text-red-400">Dangerous Actions</h2>
		<DangerButton 
			onclick={() => restartServer()}
		>
			Restart Server
		</DangerButton>
		<DangerButton 
			onclick={() => srvmod("killall")}
		>
			Kill Services For Maintenance
		</DangerButton>
		<DangerButton 
			onclick={() => srvmod("startall")}
		>
			Start All Services
		</DangerButton>
	{/if}

	<div class="mb-3"></div>
	
	<h2 class="text-xl font-semibold">Services</h2>

	{#if buildServicesTaskId != ""}
		<TaskWindow 
			output={buildServicesTaskOutput}
		/>
	{/if}

	<InputSm
		id="query"
		label="Filter by ID"
		bind:value={query}
		placeholder="E.g. arcadia"
		showErrors={false}
		minlength={0}
	/>
	<InputSm
		id="target-filter"
		label="Filter by systemd target"
		bind:value={targetFilter}
		placeholder="E.g. arcadia"
		showErrors={false}
		minlength={0}
	/>

	{#await getServiceList()}
		<h2 class="text-xl">Loading service list</h2>
	{:then data}
        {#each data as category}
            <h2 class="text-xl font-semibold text-white">{category.category}</h2>

            {#each category.services as service}
                <div class="mt-2 flex flex-wrap justify-center items-center justify-evenly">
				    {#if showService(service, query, targetFilter)}
					    <Service 
						    service={service} 
					    />
				    {/if}
                </div>
            {/each}
		{/each}
	{:catch err}
		<h2 class="text-red-500">{err}</h2>
	{/await}
</section>
