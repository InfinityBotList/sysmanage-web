<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import DangerButton from "$lib/components/DangerButton.svelte";
	import Card from "$lib/components/DefaultCard.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import TaskWindow from "$lib/components/TaskWindow.svelte";
	import { title, error, warning, success } from "$lib/strings";
	import { newTask } from "$lib/tasks";
	import Icon from "@iconify/svelte"

	export let service: any;

	let showServiceInfo = false;

	const restartService = async () => {
		let res = await fetch(`/api/systemctl?tgt=${service?.ID}&act=restart`, {
			method: "POST",
		});

		if(!res.ok) {
			let errorText = await res.text()

			error(errorText)
		}

		let out = await res.text();

		if(out) {
			warning(out)
		}

		success("Service restarted successfully")
	}

	const stopService = async () => {
		let res = await fetch(`/api/systemctl?tgt=${service?.ID}&act=stop`, {
			method: "POST",
		});

		if(!res.ok) {
			let errorText = await res.text()

			error(errorText)
		}

		let out = await res.text();

		if(out) {
			warning(out)
		}

		success("Service stopped successfully")
	}

	const startService = async () => {
		let res = await fetch(`/api/systemctl?tgt=${service?.ID}&act=start`, {
			method: "POST",
		});

		if(!res.ok) {
			let errorText = await res.text()

			error(errorText)
		}

		let out = await res.text();

		if(out) {
			warning(out)
		}

		success("Service started successfully")
	}

	let deleteServiceTaskId: string = "";
	let deleteServiceTaskOutput: string[] = [];
	const deleteService = async () => {
		let confirm = window.prompt("Are you sure you want to delete this service? (YES to confirm))")

		if(confirm != "YES") {
			return
		}

		let res = await fetch(`/api/deleteService`, {
			method: "POST",
			body: JSON.stringify({
				name: service?.ID,
			})
		});

		if(!res.ok) {
			let errorStr = await res.text()
			error(errorStr)
			return
		}

		deleteServiceTaskId = await res.text()

		newTask(deleteServiceTaskId, (output: string[]) => {
			deleteServiceTaskOutput = output
		})
	}
</script>

<Card 
	title={service?.ID}
	onclick={() => showServiceInfo = !showServiceInfo}
>
	<GreyText>{service?.Service?.Description}</GreyText>

	<!--Activity-->
	{#if service?.Status == "active"}
		<p class="text-green-500 font-semibold">
			<Icon icon="carbon:dot-mark" style="display:inline" color="green" />
			Active (Running)
		</p>
	{:else if service?.Status != "inactive"}
		<p class="text-yellow-500 font-semibold">
			<Icon icon="carbon:dot-mark" style="display:inline" color="yellow" />
			{title(service?.Status)}
		</p>
	{:else}
		<p class="text-red-500 font-semibold">
			<Icon icon="carbon:dot-mark" style="display:inline" color="red" />
			Inactive
		</p>
	{/if}

	{#if showServiceInfo}
		<p class="font-semibold text-lg">More information</p>
		<div class="text-sm">
			{#each Object.entries(service?.Service) as [key, value]}
				{#if key != "Description"}
					<p>
						<span class="font-semibold">{key}:</span> {value}
					</p>
				{/if}
			{/each}
		</div>

		<ButtonReact 
			onclick={() => restartService()}
		>
			<Icon icon="carbon:restart" color="white" />
			<span class="ml-2">Restart</span>
		</ButtonReact>
		<ButtonReact 
			onclick={() => startService()}
		>
			<Icon icon="mdi:auto-start" color="white" />
			<span class="ml-2">Start</span>
		</ButtonReact>
		<ButtonReact 
			onclick={() => stopService()}
		>
			<Icon icon="material-symbols:stop" color="white" />
			<span class="ml-2">Stop</span>
		</ButtonReact>
		<DangerButton 
			onclick={() => deleteService()}
		>
			<Icon icon="material-symbols:delete-outline-sharp" color="white" />
			<span class="ml-2">Delete</span>
		</DangerButton>

		{#if deleteServiceTaskId != ""}
			<h2 class="text-red-500">Delete service log ID: {deleteServiceTaskId}</h2>
			<TaskWindow 
				output={deleteServiceTaskOutput}
			/>
		{/if}
	{/if}
</Card>
