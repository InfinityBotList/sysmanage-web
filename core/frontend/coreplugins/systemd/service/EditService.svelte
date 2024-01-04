<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import DangerButton from "$lib/components/DangerButton.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
    import Input from "$lib/components/Input.svelte";
	import TaskWindow from "$lib/components/TaskWindow.svelte";
	import { error, success } from "$lib/corelib/strings";
	import { newTask } from "$lib/corelib/tasks";
	import Icon from "@iconify/svelte";
    import Select from "$lib/components/Select.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
    import Service from './Service.svelte';

    export let service: any;

    let deleteServiceTaskId: string = "";
	let deleteServiceTaskOutput: string[] = [];
	const deleteService = async () => {
		let confirm = window.prompt("Are you sure you want to delete this service? (YES to confirm))")

		if(confirm != "YES") {
			return
		}

		let res = await fetch(`/api/systemd/deleteService`, {
			method: "POST",
			body: JSON.stringify({
				name: service?.RawService?.FileName || service?.ID,
			})
		});

		if(!res.ok) {
			let errorStr = await res.text()
			error(errorStr)
			return
		}

		deleteServiceTaskId = await res.text()

		newTask(deleteServiceTaskId, (output: string[]) => {
			deleteServiceTaskOutput = output
		})
	}

	let deployTaskId: string = "";
	let deployTaskOutput: string[] = [];
	const initDeploy = async () => {
		let res = await fetch(`/api/systemd/initDeploy?id=${service?.ID}`, {
			method: "POST",
		});

		if(!res.ok) {
			let errorStr = await res.text()
			error(errorStr)
			return
		}

		deployTaskId = await res.text()

		newTask(deployTaskId, (output: string[]) => {
			deployTaskOutput = output
		})
	}

	let getServiceLogTaskId: string = "";
	let getServiceLogTaskOutput: string[] = [];
	const getServiceLogs = async () => {
		let res = await fetch(`/api/systemd/getServiceLogs?id=${service?.ID}`, {
			method: "POST",
		});

		if(!res.ok) {
			let errorText = await res.text()

			error(errorText)
		}

		getServiceLogTaskId = await res.text();
		newTask(getServiceLogTaskId, (output: string[]) => {
			getServiceLogTaskOutput = output
		})
    }

    const parseMap = (map: Record<string, string>): [string, string][] => {
        let arr: [string, string][] = [];

        for(let key in map) {
            arr.push([key, map[key]])
        }

        return arr;
    }

    const parseMapReverse = (map: [string, string][]): Record<string, string> => {
        let obj: Record<string, string> = {};

        for(let [key, value] of map) {
            obj[key] = value;
        }

        return obj;
    }

    let serviceDataYaml = {
        name: service?.ID || "",
        service: {
            cmd: service?.Service?.Command || "",
            dir: service?.Service?.Directory || "",
            target: service?.Service?.Target || "ibl-maint",
            description: service?.Service?.Description || "",
            after: service?.Service?.After,
            broken: service?.Service?.Broken ? true : false,
            user: service?.Service?.User || "",
            group: service?.Service?.Group || "",
        }
    }
    let brokenValue = service?.Service?.Broken ? "0" : "1";

    interface Meta {
        Targets?: MetaTarget[]
    }

    interface MetaTarget {
        Name: string
        Description: string
    }

    const getMeta = async () => {
        let metaRes = await fetch(`/api/systemd/getMeta`, {
            method: "POST",
        });

        if(!metaRes.ok) {
            let error = await metaRes.text()

            throw new Error(error)
        }

        let meta: Meta = await metaRes.json();

        return meta;
    }

    const editServiceRaw = async () => {
        let editService = await fetch(`/api/systemd/createService?update=true`, {
            method: "POST",
            body: JSON.stringify({
                raw_service: {
                    filename: service?.RawService?.FileName,
                    body: service?.RawService?.Body,
                }
            }),
        });

        if(!editService.ok) {
            let errorText = await editService.text()
            error(errorText)
            return
        }

        success("Service editted successfully!")
    }

    const editServiceYaml = async () => {
        let editService = await fetch(`/api/systemd/createService?update=true`, {
            method: "POST",
            body: JSON.stringify({
                name: serviceDataYaml.name,
                service: {
                    ...serviceDataYaml.service,
                    broken: brokenValue === "0" ? true : false,
                }
            }),
        });

        if(!editService.ok) {
            let errorText = await editService.text()
            error(errorText)
            return
        }

        success("Service editted successfully!")
    }
</script>

<DangerButton 
    onclick={() => deleteService()}
>
    <Icon icon="material-symbols:delete-outline-sharp" color="white" />
    <span class="ml-2">Delete</span>
</DangerButton>

{#if deleteServiceTaskId != ""}
    <TaskWindow 
        output={deleteServiceTaskOutput}
    />
{/if}

<ButtonReact 
    onclick={() => initDeploy()}
>
    <Icon icon="material-symbols:deployed-code" color="white" />
    <span class="ml-2">Trigger Deploy</span>
</ButtonReact>

{#if deployTaskId != ""}
    <TaskWindow 
        output={deployTaskOutput}
    />
{/if}

<ButtonReact 
    onclick={() => getServiceLogs()}
>
    <Icon icon="ph:read-cv-logo-bold" color="white" />
    <span class="ml-2">Get Service Logs</span>
</ButtonReact>

{#if getServiceLogTaskId != ""}
    <TaskWindow 
        output={getServiceLogTaskOutput}
    />
{/if}

<h2 class="font-semibold text-xl">Service Info</h2>

{#if service?.RawService}
    <div class="edit-service-raw-container">
        <h2 class="text-xl font-semibold mt-4">Raw Service</h2>
        <GreyText>Warning: Manually managed units should be manually verified for correctness!</GreyText>

        <InputSm 
            id="name"
            label="File Name"
            placeholder="zfsmongo.service"
            value={service?.RawService?.FileName || ""}
            disabled={true}
            minlength={1}
        />

        <Input 
            id="body"
            label="Body"
            placeholder="[Unit]\nDescription=Arcadia\n\n[Service]\nExecStart=/usr/bin/arcadia\nWorkingDirectory=/root/arcadia\n\n[Install]\nWantedBy=ibl-maint"
            bind:value={service.RawService.Body}
            minlength={3}
        />

        <ButtonReact
            onclick={() => editServiceRaw()}
        >
            Edit Service
        </ButtonReact>
    </div>
{:else}
    <div class="edit-service-yaml-container">
        {#await getMeta()}
            <GreyText>Loading metadata...</GreyText>
        {:then meta}
            <Service service={service} />
            
            <InputSm 
                id="name"
                label="Service Name"
                placeholder="arcadia, ibl-backup etc."
                value={serviceDataYaml.name}
                disabled={true}
                minlength={1}
            />
            <InputSm 
                id="command"
                label="Command (must start with /usr/bin/)"
                placeholder="E.g. /usr/bin/arcadia"
                bind:value={serviceDataYaml.service.cmd}
                minlength={3}
            />
            <InputSm 
                id="directory"
                label="Directory"
                placeholder="E.g. /root/arcadia"
                bind:value={serviceDataYaml.service.dir}
                minlength={3}
            />
            <Select
                name="target"
                placeholder="Choose Target"
                bind:value={serviceDataYaml.service.target}
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
                bind:value={serviceDataYaml.service.description}
                minlength={5}
            />
            <InputSm
                id="after"
                label="After"
                placeholder="E.g. ibl-maint"
                bind:value={serviceDataYaml.service.after}
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
                bind:value={serviceDataYaml.service.user}
                minlength={1}
            />
            <InputSm
                id="group"
                label="Group"
                placeholder="E.g. root"
                bind:value={serviceDataYaml.service.group}
                minlength={1}
            />
            <div class="mb-2"></div>
            <ButtonReact
                    onclick={() => editServiceYaml()}
            >
                Edit Service
            </ButtonReact>
        {/await}
    </div>
{/if}