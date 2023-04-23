<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    export let title: string;
    export let link: string;
    export let linkText: string;

    export let showArrow: boolean = true;

    export let showButton: boolean = true;
</script>

<div class="w-full block p-6 m-3 md:w-1/2 sm:w-96 bg-white rounded-lg shadow-md border border-gray-200 dark:bg-gray-800 dark:border-gray-700">
    <h5 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">{title}</h5>
    <p class="mb-3 font-normal text-gray-700 dark:text-gray-400 h-72"><slot /></p>
    
    {#if showButton}
        <Button link={link} showArrow={showArrow}>{linkText}</Button>
    {/if}
</div>