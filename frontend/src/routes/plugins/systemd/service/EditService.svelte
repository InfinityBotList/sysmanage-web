<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import DangerButton from "$lib/components/DangerButton.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
	import KvMultiInput from "$lib/components/KVMultiInput.svelte";
	import MultiInput from "$lib/components/MultiInput.svelte";
	import TaskWindow from "$lib/components/TaskWindow.svelte";
	import { error, success } from "$lib/strings";
	import { newTask } from "$lib/tasks";
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

		let res = await fetch(`/api/deleteService`, {
			method: "POST",
			body: JSON.stringify({
				name: service?.ID,
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
		let res = await fetch(`/api/initDeploy?id=${service?.ID}`, {
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
		let res = await fetch(`/api/getServiceLogs?id=${service?.ID}`, {
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

    interface Preset {
        [key: string]: {
            buildCmds: string[],
            env: [string, string][],
            allowDirty: boolean,
            configFiles: string[],
        }
    }

    const gitPresets: Preset = {
        "NPM": {
            buildCmds: [
                "npm install",
                "npm run build",
            ],
            env: [],
            allowDirty: true,
            configFiles: []
        },
        "Yarn": {
            buildCmds: [
                "yarn install",
                "yarn install --dev",
                "yarn run build"
            ],
            env: [],
            allowDirty: true,
            configFiles: []
        },
        "Go": {
            buildCmds: [
                "go build -v"
            ],
            env: [
                ["CGO_ENABLED", "0"],
            ],
            allowDirty: false,
            configFiles: [
                "config.yaml",
                "secrets.yaml"
            ]
        },
        "Rust": {
            buildCmds: [
                "/root/.cargo/bin/cargo build --release",
                "systemctl stop $NAME",
                "rm -vf $NAME",
                "mv -vf target/release/$NAME .",
                "systemctl start $NAME",
            ],
            env: [
                ["RUSTFLAGS", "-C target-cpu=native -C link-arg=-fuse-ld=lld"]
            ],
            allowDirty: true,
            configFiles: [
                "config.yaml"
            ]
        }
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

    let gitRepo: string = service?.Service?.Git?.Repo || "";
    let gitRef: string = service?.Service?.Git?.Ref || "refs/heads/";
    let gitBuildCommands: string[] = service?.Service?.Git?.BuildCommands || [];
    let configFiles: string[] = service?.Service?.Git?.ConfigFiles || [];
    let gitEnv: [string, string][] = parseMap(service?.Service?.Git?.Env) || [];
    let allowDirty: string = service?.Service?.Git?.AllowDirty?.toString() || false;


    const createGit = async () => {
        let res = await fetch(`/api/createGit?id=${service?.ID}`, {
            method: "POST",
            body: JSON.stringify({
                repo: gitRepo,
                ref: gitRef,
                build_commands: gitBuildCommands,
                env: parseMapReverse(gitEnv),
                allow_dirty: allowDirty == "true",
                config_files: configFiles,
            })
        });

        if(!res.ok) {
            let errorStr = await res.text()
            error(errorStr)
            return
        }

        success("Git integration created")
    }

    let name: string = service?.ID || "";
    let command: string = service?.Service?.Command || "";
    let directory: string = service?.Service?.Directory || "";
    let target: string = service?.Service?.Target || "ibl-maint";
    let description: string = service?.Service?.Description || "";
    let after: string = service?.Service?.After;
    let brokenValue: string = service?.Service?.Broken ? "0" : "1";

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

    const editService = async () => {
        let editService = await fetch(`/api/createService?update=true`, {
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

<div>
    {#await getMeta()}
        <GreyText>Loading metadata...</GreyText>
    {:then meta}
        <Service service={service} />
        
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
        <ButtonReact
                onclick={() => editService()}
        >
            Edit Service
        </ButtonReact>
    {/await}
</div>

<h2 class="font-semibold text-xl">Git Integration</h2>
{#if service?.Service?.Git}
    <p>Git Integration is correctly configured</p>
{:else}
    <p>Git Integration is not configured</p>
{/if}

<div>
    <InputSm
        id="git-repo"
        label="Git Repo URL"
        placeholder="https://github.com/..."
        bind:value={gitRepo}
        minlength={1}
    />
    <InputSm
        id="git-ref"
        label="Git Ref"
        placeholder="refs/head/master"
        bind:value={gitRef}
        minlength={1}
    />

    <h3 class="font-semibold">Presets</h3>
    {#each Object.entries(gitPresets) as [name, preset]}
        <ButtonReact 
            onclick={() => {
                gitBuildCommands = preset?.buildCmds.map(cmd => cmd.replaceAll("$NAME", service?.ID))
                
                if(preset?.env && preset?.env.length > 0) {
                    gitEnv = preset?.env
                }

                allowDirty = preset?.allowDirty?.toString()

                if(preset?.configFiles && preset?.configFiles.length > 0) {
                    configFiles = preset?.configFiles
                }
            }}
        >
            {name}
        </ButtonReact>
        <span class="ml-2"></span>
    {/each}

    <div class="mb-1"></div>
    
    <MultiInput 
        id="git-build-commands"
        label="Build Commands"
        title="Command"
        placeholder="npm install"
        bind:values={gitBuildCommands}
        minlength={1}
    />

    <MultiInput 
        id="config-files"
        label="Config files to preserve"
        title="Config files"
        placeholder="npm install"
        bind:values={configFiles}
        minlength={1}
    />

    <KvMultiInput
        id="git-env"
        label="Environment Variables"
        title="Key"
        placeholder="KEY"
        bind:values={gitEnv}
        minlength={1}
    />

    <Select
        name="Allow Dirty"
        placeholder="Allow dirty"
        bind:value={allowDirty}
        options={
            new Map([
                ["Yes", "true"],
                ["No", "false"],
            ])
        }
    />
    <GreyText>
        Allow dirty is used to specify whether or not we should always pull new, or if fresh clones are acceptable
    </GreyText>

    <ButtonReact onclick={() => createGit()}>Create/Update</ButtonReact>
</div>