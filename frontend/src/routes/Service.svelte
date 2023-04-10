<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import Card from "$lib/components/DefaultCard.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import { title, error, warning } from "$lib/strings";
	import Icon from "@iconify/svelte"

	export let service: any;

	let showServiceInfo = false;

	const restartService = async () => {
		let restartService = await fetch(`/api/restartService?id=${service?.ID}`, {
			method: "POST",
		});

		if(!restartService.ok) {
			let errorText = await restartService.text()

			error(errorText)
		}

		let out = await restartService.text();

		if(out) {
			warning(out)
		}
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
			<span class="mr-2">Restart</span>
		</ButtonReact>
	{/if}
</Card>
