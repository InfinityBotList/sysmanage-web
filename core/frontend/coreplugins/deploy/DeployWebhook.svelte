<script lang="ts">
	import { page } from "$app/stores";
	import InputSm from "$lib/components/InputSm.svelte";
    import Select from "$lib/components/Select.svelte";
    import type { DeployWebhook } from "./dpsettings";

    const getDeployWebhookSourceTypes = async () => {
        let res = await fetch(`/api/deploy/getDeployWebhookSourceTypes`, {
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

    export let id: string
    export let webhook: DeployWebhook
</script>

<div>
    {#await getDeployWebhookSourceTypes()}
        <h2 class="text-xl">Loading deploy webhook list</h2>
    {:then srcs}
        <Select
            name="Webhook Type"
            placeholder="Choose webhook type"
            bind:value={webhook.type}
            options={
                new Map([
                    ...parseSrc(srcs)
                ])
            }
        />
        <InputSm
            id="id"
            label="id"
            placeholder="github-deploy etc."
            bind:value={webhook.id}
            minlength={1}
        />
        {#if webhook?.id && webhook?.token}
            <p>URL: {$page.url.origin}/api/deploy/createDeploy?id={id}&wid={webhook?.id}&token={webhook?.token}</p>
        {/if}
    {:catch err}
        <h2 class="text-red-500">{err}</h2> 
    {/await}
</div>