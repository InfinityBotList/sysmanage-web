<script lang="ts">
    export let title: string;

    export let open: boolean = false;
</script>

<details open={open} class="hover:cursor-pointer">
    <summary>
        <h4 class="font-semibold text-xl opacity-80">{title}</h4>

        <div class="flex flex-row-reverse">
            <svg
                class="w-6 h-6 text-gray-500"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor">
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 9l-7 7-7-7" />
            </svg>
        </div>
    </summary>

    <slot />
</details>