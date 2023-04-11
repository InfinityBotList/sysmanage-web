<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import InputSm from "$lib/components/InputSm.svelte";

    let name: string;
</script>

<h1 class="text-2xl font-semibold">Create New Service</h1>

<GreyText>If you want to add a build integration or a git deploy hook, you can do so later after creating the service!</GreyText>

<div>
    <InputSm 
        id="name"
        label="Service Name"
        placeholder="arcadia, ibl-backup etc."
        bind:value={name}
        minlength={1}
    />
</div>