<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
    import Input from "$lib/components/Input.svelte";
	import { error, success } from "$lib/strings";
	import DangerButton from "$lib/components/DangerButton.svelte";
    import ButtonReact from "$lib/components/ButtonReact.svelte";
    import Section from "$lib/components/Section.svelte";
	import Select from "$lib/components/Select.svelte";
    
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

    // Derived from types.NginxServerManage but with only the needed field Domain
    interface AvailableDomainListData {
        Domain: string
    }

    interface GetAvailableDomainsResponse {
        Available: string[]
        Claimed: string[]
    }

    const getAvailableDomains = async (): Promise<GetAvailableDomainsResponse> => {
        let domainListRes = await fetch(`/api/nginx/getDomainList`, {
            method: "POST"
        });

        if(!domainListRes.ok) {
            let err = await domainListRes.text();
            error(err);
        }

        let domainListRaw: AvailableDomainListData[] = await domainListRes.json();
        let domainList: string[] = domainListRaw.map((domainObj) => domainObj.Domain.replace(".", "-"));

        let certListRes = await fetch(`/api/nginx/getCertList`, {
            method: "POST"
        });

        if(!certListRes.ok) {
            let err = await certListRes.text();
            error(err);
        }

        let certList: string[] = await certListRes.json();

        let availableDomains = []
        let claimedDomains = []

        // Loop over certList, ensure domain isnt already in domainList, then add it to availableDomains
        for(let cert of certList) {
            let certDomain = cert.replace("cert-", "").replace(".pem", "")

            if(!domainList.includes(certDomain)) {
                availableDomains.push(certDomain);
            } else {
                claimedDomains.push(certDomain);
            }
        }

        return {
            Available: availableDomains,
            Claimed: claimedDomains
        };
    }

    let addDomain: string;

    const addNginxDomain = async () => {
        let res = await fetch(`/api/nginx/addDomain?domain=${addDomain}`, {
            method: "POST",
        });

        if(res.ok) {
            success("Successfully added domain");
            window.location.href = `/plugins/nginx/domain?id=${addDomain}`;
	} else {
            let err = await res.text();
            error(err);
        }
    }
</script>

<h1 class="text-2xl font-semibold">Add NGINX domain</h1>

<Section title="Step 1: Domain Setup">
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
</Section>
<Section title="Step 2: Add Domain">
    {#await getAvailableDomains()}
        <p class="font-semibold">Loading...</p>
    {:then availableDomains}
        <GreyText>
            These are the domains that are available to be added to Nginx. Select one from the dropdown below.

            If you do not see your domain here, ensure that you have followed "Domain Setup" correctly and added the certificate.
        </GreyText>
        <h2 class="text-xl font-semibold">Claimed Domains</h2>
        <ul class="list-disc list-inside">
            {#each availableDomains?.Claimed as domain}
                <li>{domain}</li>
            {/each}
        </ul>
        <h2 class="text-xl font-semibold">Available Domains</h2>
        <ul class="list-disc list-inside">
            {#each availableDomains?.Available as domain}
                <li>{domain}</li>
            {/each}
        </ul>
        <h2 class="text-xl font-semibold">Select A Domain</h2>
        <Select 
            name="domain"
            placeholder="Select a domain"
            bind:value={addDomain}
            options={new Map(availableDomains?.Available?.map(a => [a, a]))}
        />
        <ButtonReact 
            onclick={addNginxDomain}
        >Add Domain</ButtonReact>
    {:catch error}
        <h2 class="text-red-400">{error}</h2>
    {/await}
</Section>
