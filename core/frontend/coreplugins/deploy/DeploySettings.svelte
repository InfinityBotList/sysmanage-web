<script lang="ts">
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import GreyText from "$lib/components/GreyText.svelte";
	import InputNumber from "$lib/components/InputNumber.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
	import MultiInput from "$lib/components/MultiInput.svelte";
    import KvMultiInput from "$lib/components/KVMultiInput.svelte";
	import Select from "$lib/components/Select.svelte";
	import DeployWebhook from "./DeployWebhook.svelte";
    import type { Data } from "./dpsettings";
    import Section from "$lib/components/Section.svelte";

    const getDeploySourceTypes = async () => {
        let res = await fetch(`/api/deploy/getDeploySourceTypes`, {
            method: "POST",
        });

        if(!res.ok) {
            let error = await res.text()

            throw new Error(error)
        }

        return await res.json() as string[]
    }

    const parseSrc = (srcs: string[]): [string, string][] => {
        return srcs?.map(src => [src, src])
    }

    interface Preset {
        [key: string]: {
            buildCmds: string[],
            env: [string, string][],
            allowDirty: boolean,
            configFiles: string[],
        }
    }

    const presets: Preset = {
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

    export let id: string
    export let data: Data;
</script>

<InputSm
    id="id"
    label="Deployment ID"
    placeholder="infinity-next-deploy"
    bind:value={id}
    minlength={1}
/>

<h2 class="text-xl font-semibold">IP Whitelist</h2>

<GreyText>
    If this has more than one IP, then only IPs in this list will be able to use this API
</GreyText>

<MultiInput 
    id="allowed-ips"
    label="Allowed IP whitelist"
    title="IP Whitelist"
    placeholder="X.X.X.X"
    bind:values={data.allowedIps}
    minlength={1}
/>

<h2 class="text-xl font-semibold">Deploy Source</h2>

<div>
    {#await getDeploySourceTypes()}
        <h2 class="text-xl">Loading deploy source list</h2>
    {:then srcs}
        <Select
            name="Source Type"
            placeholder="Choose source type"
            bind:value={data.type}
            options={
                new Map([
                    ...parseSrc(srcs)
                ])
            }
        />
        <InputSm
            id="url"
            label="URL"
            placeholder="https://..."
            bind:value={data.url}
            minlength={1}
        />
        <InputSm
            id="token"
            label="Token"
            placeholder="Any API token you need to access the source"
            bind:value={data.token}
            minlength={1}
        />
        <InputSm
            id="ref"
            label="Reference"
            placeholder="refs/head/master"
            bind:value={data.ref}
            minlength={1}
        />
    {:catch err}
        <h2 class="text-red-500">{err}</h2>  
    {/await}
</div>

<Select
    name="broken"
    placeholder="Is the service broken/disabled?"
    bind:value={data.brokenValue}
    options={new Map([
        ["Yes, it is", "0"],
        ["No, its not", "1"],
    ])}
/>

<InputSm
    id="output-path"
    label="Output Path"
    placeholder="/var/www/html etc."
    bind:value={data.outputPath}
    minlength={1}
/>

<h3 class="font-semibold">Presets</h3>

<div>
    {#each Object.entries(presets) as [name, preset]}
        <ButtonReact 
            onclick={() => {
                data.commands = preset?.buildCmds
                
                if(preset?.env && preset?.env.length > 0) {
                    data.env = preset?.env
                }

                if(preset?.configFiles && preset?.configFiles.length > 0) {
                    data.configFiles = preset?.configFiles
                }
            }}
        >
            {name}
        </ButtonReact>
        <span class="ml-2"></span>
    {/each}
</div>

<div class="mb-1"></div>

<MultiInput 
    id="commands"
    label="Build Commands"
    title="Command"
    placeholder="npm install..."
    bind:values={data.commands}
    minlength={1}
/>

<h3 class="text-xl font-semibold">Webhooks</h3>

{#each data.webhooks as webh}
    <Section title={webh?.id || "Not Specified"}>
        <DeployWebhook 
            id={id}
            bind:webhook={webh}
        />
    </Section>
{/each}
<ButtonReact onclick={() => {
    data.webhooks.push({
        id: "git",
        token: "",
        type: ""
    })
    data.webhooks = data.webhooks
}}>
    New Webhook
</ButtonReact>

<h3 class="text-xl font-semibold">Misc.</h3>

<InputNumber
    id="timeout"
    label="Timeout"
    placeholder="1234"
    bind:value={data.timeout}
/>

<KvMultiInput
    id="git-env"
    label="Environment Variables"
    title="Key"
    placeholder="KEY"
    bind:values={data.env}
    minlength={1}
/>

<MultiInput 
    id="config-files"
    label="Config files to preserve"
    title="Config files"
    placeholder="npm install"
    bind:values={data.configFiles}
    minlength={1}
/>

<!-- <ButtonReact onclick={() => createGit()}>Create/Update</ButtonReact> -->