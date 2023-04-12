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

    interface Preset {
        [key: string]: {
            git: string[],
            env: [string, string][],
            allowDirty: boolean,
            configFiles: string[],
        }
    }

    const gitPresets: Preset = {
        "NPM": {
            git: [
                "npm install",
                "npm run build",
            ],
            env: [],
            allowDirty: true,
            configFiles: []
        },
        "Yarn": {
            git: [
                "yarn install",
                "yarn install --dev",
                "yarn run build"
            ],
            env: [],
            allowDirty: true,
            configFiles: []
        },
        "Go": {
            git: [
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
    let configFiles: string[] = service?.Service?.ConfigFiles || [];
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
</script>

<DangerButton 
    onclick={() => deleteService()}
>
    <Icon icon="material-symbols:delete-outline-sharp" color="white" />
    <span class="ml-2">Delete</span>
</DangerButton>

{#if deleteServiceTaskId != ""}
    <h2 class="text-red-500">Delete service log ID: {deleteServiceTaskId}</h2>
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
    <h2 class="text-red-500">Deploy service log ID: {deployTaskId}</h2>
    <TaskWindow 
        output={deployTaskOutput}
    />
{/if}

<h2 class="font-semibold">Git Integration</h2>
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
                gitBuildCommands = preset?.git
                
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