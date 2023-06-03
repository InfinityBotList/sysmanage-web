<script lang="ts">
	import InputSm from "$lib/components/InputSm.svelte";
    import MultiInput from "$lib/components/MultiInput.svelte";
	import Section from "$lib/components/Section.svelte";
    import Select from "$lib/components/Select.svelte";
	import NgLocation from "./NGLocation.svelte";
	import DangerButton from "$lib/components/DangerButton.svelte";
	import ButtonReact from "$lib/components/ButtonReact.svelte";

    interface NGServer {
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

    export let server: NGServer;
    export let i: number;

    let broken = server.Broken ? "true" : "false";

    $: server.Broken = broken === "true";
</script>

<div>
    <h3 class="text-xl font-semibold">Editting {server.ID}</h3>

    <InputSm
        id={`s-id-${i}`}
        label="ID"
        placeholder="E.g. popplio, arcadia-rpc etc."
        minlength={1}
        bind:value={server.ID}
    />

    <MultiInput 
        id={`s-names-${i}`} 
        title="Subdomain"
        label="Subdomain (@root for root domain)"
        bind:values={server.Names} 

        placeholder="example.com, www.example.com etc."
        minlength={3}
    />

    <div class="mb-2"></div>

    <InputSm
        id={`s-comment-${i}`}
        label="Comment"
        placeholder="E.g. Popplio Web API"
        minlength={1}
        bind:value={server.Comment}
    />
    <Select
        name="broken"
        placeholder="Is the server broken/disabled?"
        bind:value={broken}
        options={new Map([
            ["Yes, it is", "true"],
            ["No, its not", "false"],
        ])}
    />

    <h3 class="text-xl font-semibold">Locations</h3>

    {#each server.Locations as loc, i}
        <Section title={loc.Path}>
            <DangerButton 
                onclick={() => {
                    // Delete the location
                    server.Locations = server.Locations.filter((_, index) => index !== i);      
                }}
            >
                Delete Location
            </DangerButton>
            <ButtonReact
                onclick={() => {
                    // Add a new location
                    server.Locations = [...server.Locations.slice(0, i + 1), {
                        Path: "Not Specified",
                    }, ...server.Locations.slice(i + 1)];
                }}
            >
                Add Location Below
            </ButtonReact>
            <ButtonReact 
                onclick={() => {
                    // Move the location up
                    let newIndex = prompt("Which index would you like to move this location to? (0 is the top)");
                    if (newIndex === null) return;

                    let newIndexNum = parseInt(newIndex);

                    if (isNaN(newIndexNum)) {
                        alert("That is not a number!");
                        return;
                    }

                    if (newIndexNum < 0 || newIndexNum > server.Locations.length - 1) {
                        alert("That is not a valid index!");
                        return;
                    }

                    // Get the location
                    let location = server.Locations[i];

                    // Get the location we are moving it to
                    let newLocation = server.Locations[newIndexNum];

                    // Set the location we are moving it to to the location we are moving
                    server.Locations[newIndexNum] = location;
                    server.Locations[i] = newLocation;

                }}
            >
                Move Location
            </ButtonReact>
            <NgLocation bind:location={loc} i={i} />
        </Section>
    {/each}
    <ButtonReact
        onclick={() => {
            // Add a new location
            server.Locations = [...server.Locations, {
                Path: "Not Specified",
            }];
        }}
    >
        Add Location
    </ButtonReact>
</div>