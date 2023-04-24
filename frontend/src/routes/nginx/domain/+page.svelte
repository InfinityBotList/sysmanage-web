<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import ObjectRender from "$lib/components/ObjectRender.svelte";
	import NgServer from "./NGServer.svelte";

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
        <h1 class="text-2xl font-semibold">Viewing {domain?.Domain}</h1>

        <h2 class="text-xl font-semibold">Server List</h2>

        <div class="flex flex-col space-y-2">
            {#each domain?.Server?.Servers as server, i}
                <NgServer bind:server={server} i={i} />
            {/each}
        </div>
        <small>
            <ObjectRender object={domain} />
        </small>
    {:catch err}
        <p class="text-red-500">{err}</p>
    {/await}
</div>