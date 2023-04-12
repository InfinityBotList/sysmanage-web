<script lang="ts">
	import ButtonReact from "./ButtonReact.svelte";
	import DangerButton from "./DangerButton.svelte";
	import Input from "./Input.svelte";
    import InputSm from "./InputSm.svelte";

    export let values: string[];
    export let title: string;
    export let placeholder: string;
    export let minlength: number;
    export let small: boolean = true;
    export let showErrors: boolean = false;

    const deleteValue = (i: number) => {
        values = values.filter((_, index) => index !== i);
    }

    const addValue = (i: number) => {
        values = [...values.slice(0, i + 1), "", ...values.slice(i + 1)];
    }
</script>

{#each values as value, i}
    {#if small}
        <InputSm
            id={i.toString()}
            label={title + " " + (i + 1)}
            placeholder={placeholder}
            bind:value={value} 
            minlength={minlength}
            showErrors={showErrors}
        />
    {:else}
        <Input 
            id={i.toString()}
            label={title + " " + (i + 1)}
            placeholder={placeholder}
            bind:value={value} 
            minlength={minlength}
        />
    {/if}

    <DangerButton onclick={() => deleteValue(i)}>Delete</DangerButton>
    <ButtonReact onclick={() => addValue(i)}>Add</ButtonReact>
{/each}

{#if values.length == 0}
    <ButtonReact onclick={() => addValue(-1)}>New {title}</ButtonReact>
{/if}