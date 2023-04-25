<script lang="ts">
    export let title: string;

    export let open: boolean = false;
</script>

<details open={open}>
    <summary class="font-semibold text-xl opacity-80 hover:cursor-pointer align-middle">
        <span>{title}</span>
    </summary>

    <article class="ml-5 mb-3">
        <slot />
    </article>
</details>