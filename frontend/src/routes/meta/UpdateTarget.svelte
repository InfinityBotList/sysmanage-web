<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
	import { error, success } from "$lib/strings";

    interface MetaTarget {
        Name: string;
        Description: string;
    }

    export let target: MetaTarget;

    let open: boolean = false;

    let editName = target?.Name;
    let editDescription = target?.Description;

    const update = async (target: MetaTarget) => {
        let response = await fetch(`/api/updateMeta?action=update&name=${target?.Name}`, {
            method: "POST",
            body: JSON.stringify(target)
        });

        if(!response.ok) {
            let resp = await response.text()
            error(resp)
        }

        success("Target updated successfully")
    }
</script>

<ButtonReact onclick={() => open = !open}>{open ? "Close Editor" : "Edit Target"}</ButtonReact>

{#if open}
    <div class="mt-2">
        <InputSm
            id="editName"
            label="New Target Name"
            placeholder="ibl, artie etc."
            bind:value={editName}
            minlength={1}
        />

        <InputSm
            id="editDescription"
            label="Target Description"
            placeholder="A description of the target"
            bind:value={editDescription}
            minlength={1}
        />

        <div class="flex justify-end mt-2">
            <ButtonReact onclick={() => update({
                Name: editName,
                Description: editDescription
            })}>Save</ButtonReact>
        </div>
    </div>
{/if}