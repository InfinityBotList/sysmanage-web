<script lang="ts">
	import { goto } from "$app/navigation";
	import ButtonReact from "$lib/components/ButtonReact.svelte";
    import GreyText from "$lib/components/GreyText.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
	import Select from "$lib/components/Select.svelte";
	import { error, success } from "$lib/strings"

    let name: string;
    let command: string = "/usr/bin/";
    let directory: string = "/root/";
    let target: string;
    let description: string;
    let user: string;
    let group: string;
    let after: string = "ibl-maint"; // Usually what you want
    let brokenValue: string = "0";

    interface Meta {
        Targets?: MetaTarget[]
    }

    interface MetaTarget {
        Name: string
        Description: string
    }

    let meta: Meta = {};

    const getMeta = async () => {
        let metaRes = await fetch(`/api/getMeta`, {
            method: "POST",
        });

        if(!metaRes.ok) {
            let error = await metaRes.text()

            throw new Error(error)
        }

        meta = await metaRes.json();

        return null;
    }

    const createService = async () => {
        let createService = await fetch(`/api/createService`, {
            method: "POST",
            body: JSON.stringify({
                name,
                service: {
                    cmd: command,
                    dir: directory,
                    target,
                    description,
                    after,
                    broken: brokenValue === "0" ? true : false,
                    user,
                    group,
                }
            }),
        });

        if(!createService.ok) {
            let errorText = await createService.text()
            error(errorText)
            return
        }

        success("Service created successfully!")

        setTimeout(() => {
            goto(`/plugins/systemd/service?id=${name}`)
        }, 1000)
    }
</script>

<h1 class="text-2xl font-semibold">Create New Service</h1>

<GreyText>If you want to add a build integration or a git deploy hook, you can do so later after creating the service!</GreyText>

<div>
    {#await getMeta()}
        <GreyText>Loading metadata...</GreyText>
    {:then fl}
        <div id={JSON.stringify(fl)}></div>
        <InputSm 
            id="name"
            label="Service Name"
            placeholder="arcadia, ibl-backup etc."
            bind:value={name}
            minlength={1}
        />
        <InputSm 
            id="command"
            label="Command (must start with /usr/bin/)"
            placeholder="E.g. /usr/bin/arcadia"
            bind:value={command}
            minlength={3}
        />
        <InputSm 
            id="directory"
            label="Directory"
            placeholder="E.g. /root/arcadia"
            bind:value={directory}
            minlength={3}
        />
        <Select
            name="target"
            placeholder="Choose Target"
            bind:value={target}
            options={
                new Map(meta?.Targets?.map(target => [
                    target?.Name + " - " + target?.Description, 
                    target?.Name
                ]))
            }
        />
        <InputSm
            id="description"
            label="Description"
            placeholder="E.g. Arcadia"
            bind:value={description}
            minlength={5}
        />
        <InputSm
            id="after"
            label="After"
            placeholder="E.g. ibl-maint"
            bind:value={after}
            minlength={1}
        />
        <Select
            name="broken"
            placeholder="Is the service broken/disabled?"
            bind:value={brokenValue}
            options={new Map([
                ["Yes, it is", "0"],
                ["No, its not", "1"],
            ])}
        />
        <h2 class="text-xl font-semibold mt-4">Service User</h2>
        <GreyText>Defaults to root if unset. Note that this could be a possible security risk to use the wrong user/group!</GreyText>
        <InputSm
            id="user"
            label="User"
            placeholder="E.g. root"
            bind:value={user}
            minlength={1}
        />
        <InputSm
            id="group"
            label="Group"
            placeholder="E.g. root"
            bind:value={group}
            minlength={1}
        />
        <div class="mb-2"></div>
        <ButtonReact
                onclick={() => createService()}
        >
            Create Service
        </ButtonReact>
    {/await}
</div>