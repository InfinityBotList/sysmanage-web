<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import { error, success } from "$lib/strings";

    interface Action {
        Name: string,
        Description:   string,
        ConfirmDialog: string // If unset, no confirm dialog will be shown
    }

    const getActionList = async () => {
		let serviceList = await fetch(`/api/getActionList`, {
			method: "POST",
		});

		if(!serviceList.ok) {
			let error = await serviceList.text()

			throw new Error(error)
		} 

        let json: Action[] = await serviceList.json();

		return await json;
	}
</script>

<svelte:head>
	<title>Custom Actions</title>
</svelte:head>

<section>
    {#await getActionList()}
        <GreyText>Loading actions...</GreyText>
    {:then actionList}
        <h1 class="text-2xl font-semibold">Custom Actions</h1>
        <div class="mt-4">
            {#each actionList as action}
                <div class="flex flex-row items-center justify-between">
                    <div class="flex flex-col">
                        <h2 class="text-lg font-semibold">{action.Name}</h2>
                        <GreyText>{action.Description}</GreyText>
                    </div>
                    <div class="flex flex-row items-center">
                        <ButtonReact
                            onclick={async () => {
                                if(action.ConfirmDialog) {
                                    let confirm = window.prompt(action.ConfirmDialog + " (yes to confirm)")

                                    if(confirm?.toLowerCase() != "yes") {
                                        return
                                    }
                                }

                                let res = await fetch(`/api/executeAction?actionName=${action.Name}`, {
                                    method: "POST"
                                })

                                if(res.ok) {
                                    success("Action executed successfully")
                                } else {
                                    let errorStr = await res.text()

                                    error(errorStr)
                                }
                            }}
                        >
                            {action.Name}
                        </ButtonReact>
                    </div>
                </div>
                <hr class="my-4" />
            {/each}
        </div>
    {:catch err}
        <p class="text-red-500">{err}</p>
    {/await}
</section>