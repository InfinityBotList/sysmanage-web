<script lang="ts">
	import ButtonReact from "./ButtonReact.svelte";
	import DangerButton from "./DangerButton.svelte";
	import KvMultiInputElement from "./KVMultiInputElement.svelte";
    
    export let id: string;
    export let values: [string, string][] = [];
    export let title: string;
    export let label: string = title;
    export let placeholder: string;
    export let minlength: number;
    export let showErrors: boolean = false;

    const deleteValue = (i: number) => {
        values = values.filter((_, index) => index !== i);
    }

    const addValue = (i: number) => {
        values = [...values.slice(0, i + 1), "", ...values.slice(i + 1)] as [string, string][];
    }
</script>

<label for={id} class="block mb-1 font-medium text-gray-900 dark:text-gray-300">{label}</label>
<div id={id}>
    {#each values as value, i}
        <KvMultiInputElement 
            title={title}
            placeholder={placeholder}
            minlength={minlength}
            showErrors={showErrors}
            i={i}
            bind:value={value}
        />
        <DangerButton onclick={() => deleteValue(i)}>Delete</DangerButton>
        <ButtonReact onclick={() => addValue(i)}>Add</ButtonReact>
    {/each}

    {#if values.length == 0}
        <ButtonReact onclick={() => addValue(-1)}>New {title}</ButtonReact>
    {/if}
</div>