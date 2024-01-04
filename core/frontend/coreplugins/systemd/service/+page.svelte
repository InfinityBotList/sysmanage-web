<script lang="ts">
    import GreyText from "$lib/components/GreyText.svelte";
	import EditService from "./EditService.svelte";

    let serviceId: string | null = null
    const getServiceId = (): string => {
        if(serviceId) return serviceId;
        let searchParams = new URLSearchParams(window.location.search);
        serviceId = searchParams.get("id") || "";
        return serviceId
    }

    const getService = async () => {
        if(!getServiceId()) {
            throw new Error("No service id provided");
        }

		let serviceList = await fetch(`/api/systemd/getServiceList`, {
			method: "POST",
		});

		if(!serviceList.ok) {
			let error = await serviceList.text()

			throw new Error(error)
		} 

		let list = await serviceList.json();

        let service = list.find((service: any) => service?.ID == getServiceId());

        if(!service) {
            throw new Error("Service not found");
        }

        return service;
    }
</script>

<div>
    {#await getService()}
        <GreyText>Loading service...</GreyText>
    {:then service}
        <h1 class="text-2xl font-semibold">Viewing {service?.ID}</h1>
        <EditService service={service} />
    {:catch err}
        <p class="text-red-500">{err}</p>
    {/await}
</div>