<script lang="ts">
	import InputSm from '$lib/components/InputSm.svelte';
	import ButtonReact from '$lib/components/ButtonReact.svelte';
	import { error, success } from '$lib/strings';
	import TaskWindow from '../../lib/components/TaskWindow.svelte';
	import { newTask } from '$lib/tasks';

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

    let domain: string;
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

	<div class="mb-3"></div>
	
	<h2 class="text-xl font-semibold">Services</h2>

	{#if buildNginxTaskId != ""}
		<TaskWindow 
			output={buildNginxTaskOutput}
		/>
	{/if}

	<InputSm
		id="domain"
		label="Filter by domain"
		bind:value={domain}
		placeholder="E.g. arcadia"
		showErrors={false}
		minlength={0}
	/>
</section>