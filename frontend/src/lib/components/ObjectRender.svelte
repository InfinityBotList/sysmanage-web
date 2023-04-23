<script lang="ts">
    export let object: any;
</script>

{#each Object.entries(object) as [key, value]}
    {#if typeof value === "object" && value !== null}
        <div class="ml-4">
            <svelte:self object={value} />
        </div>
    {:else}
        <p>
            <span class="font-semibold">{key}:</span> {value}
        </p>
    {/if}
{/each}