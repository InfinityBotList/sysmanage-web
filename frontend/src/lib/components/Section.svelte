<script lang="ts">
    export let title: string;

    export let open: boolean = false;
</script>

<details open={open} class="ml-5 mr-5">
    <summary class="font-semibold text-xl opacity-80 hover:cursor-pointer align-middle">
        <span>{title}</span>
    </summary>

    <slot />
</details>