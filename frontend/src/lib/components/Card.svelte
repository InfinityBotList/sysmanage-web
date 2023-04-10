<!--Card made using tailwind-->

<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    export let title: string;
    export let link: string;
    export let linkText: string;

    export let showArrow: boolean = true;

    export let showButton: boolean = true;
</script>

<div class="w-96 h-[26rem] p-6 m-3 bg-white shadow-md rounded-lg border border-gray-200 dark:bg-gray-800 dark:border-gray-700">
    <h5 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">{title}</h5>
    <p class="mb-3 font-normal text-gray-700 dark:text-gray-400 h-72"><slot /></p>
    
    {#if showButton}
        <Button link={link} showArrow={showArrow}>{linkText}</Button>
    {/if}
</div>
