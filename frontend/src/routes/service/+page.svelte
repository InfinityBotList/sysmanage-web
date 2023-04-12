<script lang="ts">
    import GreyText from "$lib/components/GreyText.svelte";
	import EditService from "./EditService.svelte";

    let service: any;

    const getServiceId = (): string => {
        let searchParams = new URLSearchParams(window.location.search);

        return searchParams.get("id") || "";
    }

    const getService = async () => {
        if(!getServiceId()) {
            throw new Error("No service id provided");
        }

		let serviceList = await fetch(`/api/getServiceList`, {
			method: "POST",
		});

		if(!serviceList.ok) {
			let error = await serviceList.text()

			throw new Error(error)
		} 

		let list = await serviceList.json();

        service = list.find((service: any) => service?.ID == getServiceId());

        if(!service) {
            throw new Error("Service not found");
        }

        return service;
    }
</script>

<div>
    {#await getService()}
        <GreyText>Loading metadata...</GreyText>
    {:then service}
        <h1 class="text-2xl font-semibold">Editting {service?.ID}</h1>
        <EditService service={service} />
    {/await}
</div>