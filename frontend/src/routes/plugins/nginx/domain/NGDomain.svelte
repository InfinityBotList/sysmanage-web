<script lang="ts">
	import NgServer from "./NGServer.svelte";
    import ObjectRender from "$lib/components/ObjectRender.svelte";
    import Section from "$lib/components/Section.svelte";
	import DangerButton from "$lib/components/DangerButton.svelte";
	import ButtonReact from "$lib/components/ButtonReact.svelte";
	import { error, success } from "$lib/strings";
	import { goto } from "$app/navigation";
	import { newTask } from "$lib/tasks";
	import TaskWindow from "$lib/components/TaskWindow.svelte";

    interface NgDomain {
        Domain: string,
        Server: NgServerList
    }

    interface NGServerInterface {
        ID: string,
        Names: string[],
        Comment: string,
        Broken: boolean,
        Locations: NGLocation[],
    }

    interface NGLocation {
        Path: string,
        Proxy?: string,
        Opts?: string[],
    }

    interface NgServerList {
        Servers: NGServerInterface[]
    }

    export let domain: NgDomain;

    const saveChanges = async () => {
        let res = await fetch(`/api/nginx/updateDomain`, {
            method: "POST",
            body: JSON.stringify(domain)
        });

        if (!res.ok) {
            let err = await res.text();
            error(err);
            return;
        }

        success("Successfully updated domain!");
    }

    let deleteDomainTaskId: string = "";
    let deleteDomainTaskOutput: string[] = [];
    const deleteDomain = async () => {
        // Send prompt
        let p = prompt("Are you sure you want to delete this domain? Type 'YES' to confirm.");

        if (p !== "YES") {
            return;
        }

        let res = await fetch(`/api/nginx/deleteDomain?domainName=${domain}`, {
            method: "POST",
        });

		if(!res.ok) {
			let errorStr = await res.text()
			error(errorStr)
			return
		}

		deleteDomainTaskId = await res.text()

		newTask(deleteDomainTaskId, (output: string[]) => {
			deleteDomainTaskOutput = output
		})
    }
</script>

<h1 class="text-2xl font-semibold">Viewing {domain?.Domain}</h1>

<h2 class="text-xl font-semibold">Server List</h2>

<div class="flex flex-col space-y-2">
    {#each domain?.Server?.Servers as server, i}
        <Section title={server.ID}>
            <DangerButton 
                onclick={() => {
                    // Delete the server
                    domain.Server.Servers = domain.Server.Servers.filter((_, index) => index !== i);      
                }}
            >
                Delete Server
            </DangerButton>
            <ButtonReact
                onclick={() => {
                    // Add a new server
                    domain.Server.Servers = [...domain.Server.Servers, {
                        ID: "Not Specified",
                        Names: [],
                        Comment: "",
                        Broken: false,
                        Locations: []
                    }];
                }}
            >
                Add Server Below
            </ButtonReact>    
            <ButtonReact
                onclick={() => {
                    // Move the server up
                    let newIndex = prompt("Which index would you like to move this server to? (0 is the top)");
                    if (newIndex === null) return;

                    let newIndexNum = parseInt(newIndex);

                    if (isNaN(newIndexNum)) {
                        alert("That is not a number!");
                        return;
                    }

                    if (newIndexNum < 0 || newIndexNum > domain?.Server?.Servers.length - 1) {
                        alert("That is not a valid index!");
                        return;
                    }

                    // Get the location
                    let location = domain?.Server?.Servers[i];

                    // Get the location we are moving it to
                    let newLocation = domain?.Server?.Servers[newIndexNum];

                    // Set the location we are moving it to to the location we are moving
                    domain.Server.Servers[newIndexNum] = location;
                    domain.Server.Servers[i] = newLocation;
                }}
            >
                Move Server
            </ButtonReact>
                
            <NgServer bind:server={server} i={i} />
        </Section>
    {/each}
    <ButtonReact
        onclick={() => {
            // Add a new server
            domain.Server.Servers = [...domain.Server.Servers, {
                ID: "Not Specified",
                Names: [],
                Comment: "",
                Broken: false,
                Locations: [
                    {
                        Path: "/"
                    }
                ]
            }];
        }}
    >
        Add Server
    </ButtonReact>
    <hr class="mt-2 mb-2" />
    <Section title="Tree View">
        <ObjectRender object={domain} />
    </Section>
</div>

<ButtonReact
    onclick={saveChanges}
>Save changes</ButtonReact>
<DangerButton
    onclick={deleteDomain}
>Delete Domain</DangerButton>

{#if deleteDomainTaskId != ""}
    <TaskWindow 
        output={deleteDomainTaskOutput}
    />
{/if}