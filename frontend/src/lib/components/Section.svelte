<script lang="ts">
    export let title: string;

    export let open: boolean = false;
</script>

<details open={open} class="hover:cursor-pointer">
    <summary>
        <span class="font-semibold text-xl opacity-80">{title}</span>
    </summary>

    <slot />
</details>