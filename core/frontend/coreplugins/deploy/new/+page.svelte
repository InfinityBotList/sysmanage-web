<script lang="ts">
	import DeploySettings from "../DeploySettings.svelte";
	import type { Data } from "../dpsettings";

    export let id: string = ""
    export let dpSettings: Data = {
        allowedIps: [],
        type: "git",
        url: "",
        token: "",
        ref: "",
        brokenValue: "1",
        outputPath: "",
        commands: [],
        webhooks: [],
        timeout: 60,
        env: [],
        configFiles: [],
    }
</script>

<svelte:head>
	<title>New Deploy</title>
	<meta name="description" content="Create a new deploy here" />
</svelte:head>

<section>
	<h2 class="text-xl font-semibold">New Deploy</h2>

    <DeploySettings bind:data={dpSettings} bind:id={id} />
</section>