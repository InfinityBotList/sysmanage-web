<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
import InputSm from "$lib/components/InputSm.svelte";
	import { error, success } from "$lib/strings";
	import UpdateTarget from "./UpdateTarget.svelte";

    const updateMeta = async (action: string, target: MetaTarget) => {
        let response = await fetch(`/api/updateMeta?action=${action}`, {
            method: "POST",
            body: JSON.stringify(target)
        });

        if(!response.ok) {
            let resp = await response.text()
            error(resp)
        }

        success("Target updated successfully")
    }

    interface Meta {
        Targets?: MetaTarget[]
    }

    interface MetaTarget {
        Name: string
        Description: string
    }

    const getMeta = async () => {
        let metaRes = await fetch(`/api/getMeta`, {
            method: "POST",
        });

        if(!metaRes.ok) {
            let error = await metaRes.text()

            throw new Error(error)
        }

        let meta: Meta = await metaRes.json();

        return meta;
    }

    let addName: string;
    let addDescription: string;
</script>

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
                Name: addName,
                Description: addDescription
            })
        }}
    >
        Create Target
    </ButtonReact>
</div>

<h2 class="text-2xl font-semibold">Meta Editor</h2>

<div>
    <div>
        {#await getMeta()}
            <GreyText>Loading metadata...</GreyText>
        {:then meta}
            {#each (meta?.Targets || []) as target}
                <div>
                    <div class="inline-block w-32">
                        <p class="text-lg font-semibold">{target?.Name}</p>
                        <p class="text-sm">{target?.Description}</p>
                    </div>
                    <div>
                        <ButtonReact
                            onclick={() => {
                                updateMeta("delete", target)
                            }}
                        >
                            Delete
                        </ButtonReact>
                        <span class="ml-2"></span>
                        <UpdateTarget target={target} />
                    </div>
                </div>
            {/each}
        {/await}
    </div>    
</div>