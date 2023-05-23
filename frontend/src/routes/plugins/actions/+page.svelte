<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import TaskWindow from "$lib/components/TaskWindow.svelte";
	import { error, success } from "$lib/strings";
	import { newTask } from "$lib/tasks";

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

    let execActionTaskIds: string[] = []
    let execActionTaskOutputs: string[][] = []
    const executeAction = async (action: Action) => {
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
            if(res.headers.get("X-Task-ID")) {
                // Execute task
                execActionTaskIds.push(res.headers.get("X-Task-ID") || "")
                execActionTaskOutputs[execActionTaskIds.length - 1] = [action.Name + "\n"]

                newTask(execActionTaskIds[execActionTaskIds.length - 1], (output: string[]) => {
			        execActionTaskOutputs[execActionTaskIds.length - 1] = [action.Name + "\n", ...output]
		        })
            }

            success("Action executed successfully")
        } else {
            let errorStr = await res.text()

            error(errorStr)
        }
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
                            onclick={() => executeAction(action)}
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

    {#each execActionTaskIds as _, i}
        <TaskWindow 
            output={execActionTaskOutputs[i]}
        />
    {/each}
</section>