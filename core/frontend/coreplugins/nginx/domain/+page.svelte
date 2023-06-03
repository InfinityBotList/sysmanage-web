<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import NgDomain from "./NGDomain.svelte";

    const getDomainId = (): string => {
        let searchParams = new URLSearchParams(window.location.search);

        return searchParams.get("id") || "";
    }

    const getDomain = async () => {
        if(!getDomainId()) {
            throw new Error("No domain name provided in query");
        }

		let domainList = await fetch(`/api/nginx/getDomainList`, {
			method: "POST",
		});

		if(!domainList.ok) {
			let error = await domainList.text()

			throw new Error(error)
		} 

		let list = await domainList.json();

        let domain = list.find((domain: any) => domain?.Domain == getDomainId());

        if(!domain) {
            throw new Error("Domain not found");
        }

        return domain;
    }
</script>

<div>
    {#await getDomain()}
        <GreyText>Loading metadata...</GreyText>
    {:then domain}
        <NgDomain domain={domain} />
    {:catch err}
        <p class="text-red-500">{err}</p>
    {/await}
</div>