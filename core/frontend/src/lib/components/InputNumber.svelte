<script lang="ts">
    export let id: string;
    export let label: string;
    export let placeholder: string;
    export let value: number = 0;
    export let description: string = "";
    export let inpClass: string = "mb-6";
</script>

<div class={inpClass}>
    <label for={id} class="block mb-1 font-medium text-gray-900 dark:text-gray-300">{label}</label>
    {#if description}
        <span class="text-md text-gray-500 dark:text-gray-400 mb-2">{{description}}</span>
    {/if}
    <input type="number" id={id} class="bg-gray-50 border border-gray-300 text-gray-900 text-md rounded-md focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder={placeholder} required bind:value />

    <slot />
</div> 
