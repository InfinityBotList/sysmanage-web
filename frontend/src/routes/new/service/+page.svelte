<script lang="ts">
	import GreyText from "$lib/components/GreyText.svelte";
	import InputSm from "$lib/components/InputSm.svelte";
	import Select from "$lib/components/Select.svelte";

    let name: string;
    let command: string = "/usr/bin/";
    let directory: string = "/root/";
    let target: string;
    let description: string;
    let after: string = "ibl-maint"; // Usually what you want
    let brokenValue: string = "0";
    let folder: string

    let folderList: string[] = [];

    const getDefinitionFolders = async () => {
        let definitionFolders = await fetch(`/api/getDefinitionFolders`, {
            method: "POST",
        });

        if(!definitionFolders.ok) {
            let error = await definitionFolders.text()

            throw new Error(error)
        }

        folderList = await definitionFolders.json();
        folder = folderList[0];
    }
</script>

<h1 class="text-2xl font-semibold">Create New Service</h1>

<GreyText>If you want to add a build integration or a git deploy hook, you can do so later after creating the service!</GreyText>

<div>
    {#await getDefinitionFolders()}
        <GreyText>Loading folder list...</GreyText>
    {:then fl}
        <div id={JSON.stringify(fl)}></div>
        <Select
            name="folder"
            placeholder="Folder"
            bind:value={folder}
            options={new Map(folderList?.map(folder => [folder, folder]))}
        />
        <div class="mb-3"></div>
    {/await}
    <InputSm 
        id="name"
        label="Service Name"
        placeholder="arcadia, ibl-backup etc."
        bind:value={name}
        minlength={1}
    />
    <InputSm 
        id="command"
        label="Command (must start with /usr/bin/)"
        placeholder="E.g. /usr/bin/arcadia"
        bind:value={command}
        minlength={3}
    />
    <InputSm 
        id="directory"
        label="Directory"
        placeholder="E.g. /root/arcadia"
        bind:value={directory}
        minlength={3}
    />
    <InputSm 
        id="target"
        label="Target"
        placeholder="E.g. ibl"
        bind:value={target}
        minlength={1}
    />
    <InputSm
        id="description"
        label="Description"
        placeholder="E.g. Arcadia"
        bind:value={description}
        minlength={5}
    />
    <InputSm
        id="after"
        label="After"
        placeholder="E.g. ibl-maint"
        bind:value={after}
        minlength={1}
    />
    <Select
        name="broken"
        placeholder="Is the service broken/disabled?"
        bind:value={brokenValue}
        options={new Map([
            ["Yes, it is", "0"],
            ["No, its not", "1"],
        ])}
    />
</div>