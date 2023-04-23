<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
    import Input from "$lib/components/Input.svelte";
	import { error, success } from "$lib/strings";
	import DangerButton from "$lib/components/DangerButton.svelte";
    import ButtonReact from "$lib/components/ButtonReact.svelte";
    
    let publishDomain: string;
    let publishCert: string;
    let publishKey: string;
    let warningNeedsForce: boolean;
    const publishCerts = async (force: boolean) => {
        if(force) {
            let prompt = window.prompt("Are you sure you want to overwrite this domain? Type 'YeS' to continue");

            if(prompt != "YeS") {
                error("Cancelled");
                return;
            }
        }

        let res = await fetch(`/api/nginx/publishCerts?force=${force}`, {
            method: "POST",
            body: JSON.stringify({
                domain: publishDomain,
                cert: publishCert,
                key: publishKey,
            }),
        });

        if(res.ok) {
            success("Successfully published certificates");
            window.location.reload()
        } else {
            let err = await res.text();

            if(err == "ALREADY_EXISTS") {
                warningNeedsForce = true;
                error("This domain already exists. If you want to overwrite it, click the Force Push button below");
                return;
            }

            error(err);
        }
    }
</script>

<h1 class="text-2xl font-semibold">Add NGINX domain</h1>

<h2 class="text-xl font-semibold">Domain Setup</h2>
<GreyText>Follow these steps first to add your domain to Cloudflare</GreyText>
<p class="font-semibold">Note that this is NOT needed if the domain is already previously setup on Nginx</p>
<ol class="list-decimal list-inside">
    <li>Add your domain to Cloudflare normally</li>
    <li>Click SSL/TLS > Overview. Then ensure SSL/TLS encryption mode is set to Full or Full (strict)</li>
    <li>Go to SSL/TLS > Origin Server. Ensure "Authenticated Origin Pulls" is enabled. Then create a new origin certificate</li>
    <li>This will yield two files, a certificate and a private key. Copy the contents of these files and paste them into the fields below</li>
</ol>

<div class="mt-3">
    <InputSm 
        id="publish-domain"
        label="Domain (without any www or http/https)"
        placeholder="infinitybots.gg, botlist.app, narc.live etc."
        bind:value={publishDomain}
        minlength={3}
    />
    <Input
        id="publish-cert"
        label="Certificate (Public Cert)"
        placeholder="-----BEGIN CERTIFICATE-----"
        bind:value={publishCert}
        minlength={256}
    />
    <Input
        id="publish-key"
        label="Certificate (Private Key)"
        placeholder="-----BEGIN PRIVATE KEY-----"
        bind:value={publishKey}
        minlength={256}
    />
    <ButtonReact 
        onclick={() => publishCerts(false)}
    >Publish</ButtonReact>

    {#if warningNeedsForce}
        <h3 class="text-xl font-semibold text-red-400">Force Push</h3>
        <GreyText>Clicking this button will overwrite the existing domain. This is not recommended unless you know what you're doing</GreyText>

        <DangerButton 
            onclick={() => publishCerts(true)}
        >Yes, I'm sure!</DangerButton>
    {/if}
</div>