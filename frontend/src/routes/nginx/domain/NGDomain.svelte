<script lang="ts">
	import NgServer from "./NGServer.svelte";
    import ObjectRender from "$lib/components/ObjectRender.svelte";

    export let domain: any;
</script>

<h1 class="text-2xl font-semibold">Viewing {domain?.Domain}</h1>

<h2 class="text-xl font-semibold">Server List</h2>

<div class="flex flex-col space-y-2">
    {#each domain?.Server?.Servers as server, i}
        <NgServer bind:server={server} i={i} />
    {/each}
</div>
<small>
    <p class="font-semibold mt-5">Tree View</p>
    <ObjectRender object={domain} />
</small>