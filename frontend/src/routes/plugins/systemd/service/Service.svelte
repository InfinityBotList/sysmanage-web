<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import Card from "$lib/components/DefaultCard.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import ObjectRender from "$lib/components/ObjectRender.svelte";
	import { title, error, warning, success } from "$lib/strings";
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
			<ObjectRender object={service?.Service} />
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

		<Button link={`/plugins/systemd/service?id=${service?.ID}`}>
			<Icon icon="mdi:file-edit" color="white" />
			Edit
		</Button>	
	{/if}
</Card>
