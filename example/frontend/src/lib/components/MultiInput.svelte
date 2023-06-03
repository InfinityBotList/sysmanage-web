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

    const deleteValue = (i: number) => {
        values = values.filter((_, index) => index !== i);
    }

    const addValue = (i: number) => {
        values = [...values.slice(0, i + 1), "", ...values.slice(i + 1)];
    }
</script>

{#if showLabel || values.length == 0}
    <label for={id} class="block mb-1 font-medium text-gray-900 dark:text-gray-300">{label}</label>
{:else}
    <label for={id} class="sr-only">{label}</label>
{/if}
<div id={id} class="mt-2 mb-2">
    <div class="ml-4">
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
                    <div class="mt-1">
                        <DangerButton onclick={() => deleteValue(i)}>Delete</DangerButton>
                        <ButtonReact onclick={() => addValue(i)}>Add</ButtonReact>        
                    </div>
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
                    <div class="mt-1">
                        <DangerButton onclick={() => deleteValue(i)}>Delete</DangerButton>
                        <ButtonReact onclick={() => addValue(i)}>Add</ButtonReact>        
                    </div>
                </Input>
            {/if}
        {/each}
    </div>

    {#if values.length == 0}
        <ButtonReact onclick={() => addValue(-1)}>New {title}</ButtonReact>
    {/if}
</div>