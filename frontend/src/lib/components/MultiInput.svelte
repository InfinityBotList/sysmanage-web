<script lang="ts">
	import ButtonReact from "./ButtonReact.svelte";
	import DangerButton from "./DangerButton.svelte";
	import Input from "./Input.svelte";
    import InputSm from "./InputSm.svelte";

    export let id: string;
    export let values: string[];
    export let title: string;
    export let label: string = title;
    export let placeholder: string;
    export let minlength: number;
    export let small: boolean = true;
    export let showErrors: boolean = false;
    export let showLabel: boolean = true;
    export let allowBulkAdd: boolean = true;

    const deleteValue = (i: number) => {
        values = values.filter((_, index) => index !== i);
    }

    const addValue = (i: number) => {
        values = [...values.slice(0, i + 1), "", ...values.slice(i + 1)];
    }

    let showBulkAdd = false;
    let bulkAddValues = "";

    $: if (bulkAddValues.length > 0) {
        values = bulkAddValues.split("\n");
    }
</script>

{#if showLabel || values.length == 0}
    <label for={id} class="block mb-1 font-medium text-gray-900 dark:text-gray-300">{label}</label>
{:else}
    <label for={id} class="sr-only">{label}</label>
{/if}
<div id={id} class="mt-2 mb-2 ml-4">
    {#if allowBulkAdd}
        <ButtonReact onclick={() => showBulkAdd = !showBulkAdd}>{showBulkAdd ? "Close" : "Import"}</ButtonReact>
        {#if showBulkAdd}
            <Input
                id={`${id}-bulk`}
                label="Bulk Add"
                placeholder="Enter one statement per line"
                bind:value={bulkAddValues} 
                minlength={0}
                showErrors={true}
            />
        {/if}
    {/if}

    {#if values.length > 0}
        <DangerButton onclick={() => values = []}>Clear {title}</DangerButton>
    {/if}

    {#each values as value, i}
        {#if small}
            <InputSm
                id={i.toString()}
                inpClass="mb-1"
                label={title + " " + (i + 1)}
                placeholder={placeholder}
                bind:value={value} 
                minlength={minlength}
                showErrors={showErrors}
            >
                <DangerButton onclick={() => deleteValue(i)}>Delete</DangerButton>
                <ButtonReact onclick={() => addValue(i)}>Add</ButtonReact>    
            </InputSm>
        {:else}
            <Input 
                id={i.toString()}
                inpClass="mb-1"
                label={title + " " + (i + 1)}
                placeholder={placeholder}
                bind:value={value} 
                minlength={minlength}
            >
                <DangerButton onclick={() => deleteValue(i)}>Delete</DangerButton>
                <ButtonReact onclick={() => addValue(i)}>Add</ButtonReact>    
            </Input>
        {/if}
    {/each}

    {#if values.length == 0}
        <ButtonReact onclick={() => addValue(-1)}>New {title}</ButtonReact>
    {/if}

    <div class="mb-3"></div>
</div>