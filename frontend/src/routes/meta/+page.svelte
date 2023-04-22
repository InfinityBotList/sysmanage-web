<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
import InputSm from "$lib/components/InputSm.svelte";
	import { error } from "$lib/strings";

    const updateMeta = async (action: string, target: any) => {
        let response = await fetch(`/api/updateMeta?action=${action}`, {
            method: "POST",
            body: JSON.stringify(target)
        });

        if(!response.ok) {
            let resp = await response.text()
            error(resp)
        }

        return await response.json();
    }

    let addName: string;
    let addDescription: string;
</script>

<h2 class="text-2xl font-semibold">Meta Editor</h2>

<h3 class="text-xl font-semibold">Add Target</h3>

<div>
    <InputSm 
        id="addName"
        label="Target Name"
        placeholder="ibl, artie etc."
        bind:value={addName}
        minlength={1}
    />
    <InputSm 
        id="addDescription"
        label="Target Description"
        placeholder="Whoa"
        bind:value={addDescription}
        minlength={1}
    />
    <ButtonReact
        onclick={() => {
            updateMeta("create", {
                name: addName,
                description: addDescription
            })
        }}
    >
        Create Target
    </ButtonReact>
</div>

<h3 class="text-xl font-semibold">Update Target</h3>

<div>
    <InputSm 
        id="addName"
        label="Target Name"
        placeholder="ibl, artie etc."
        bind:value={addName}
        minlength={1}
    />
    <InputSm 
        id="addDescription"
        label="Target Description"
        placeholder="Whoa"
        bind:value={addDescription}
        minlength={1}
    />
    <ButtonReact
        onclick={() => {
            updateMeta("create", {
                name: addName,
                description: addDescription
            })
        }}
    >
        Create Target
    </ButtonReact>
</div>