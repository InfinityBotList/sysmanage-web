<script lang="ts">
    export let id: string;
    export let label: string;
    export let placeholder: string;
    export let minlength: number;
    export let value: string = "";
    export let showErrors: boolean = true;
    export let description: string = "";
    export let inpClass: string = "mb-6";

    let success: boolean | null = null;

    let errorMsg = ""

    function checkLength() {
        if(!showErrors) return

        console.log(value)
        if(!value) {
            success = null;
            return
        }
        
        if (value.length < minlength) {
            success = false;
            errorMsg = `Must be at least ${minlength} characters long`
        } else {
            success = true;
        }
    }
</script>

<div class={inpClass}>
    <label for={id} class="block mb-1 font-medium text-gray-900 dark:text-gray-300">{label}</label>
    {#if description}
        <span class="text-md text-gray-500 dark:text-gray-400 mb-2">{{description}}</span>
    {/if}
    <textarea on:change={checkLength} minlength={minlength} id={id} class="h-36 bg-gray-50 border border-gray-300 text-gray-900 text-md rounded-md focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder={placeholder} required bind:value></textarea>

    {#if success == true}
        <p class="mt-2 text-sm text-green-600 dark:text-green-500"><span class="font-medium">Looks good!</span></p>
    {:else if success == false}
        <p class="mt-2 text-sm text-red-600 dark:text-red-500"><span class="font-medium">{errorMsg}</span></p>
    {/if}

    <slot />
</div> 
