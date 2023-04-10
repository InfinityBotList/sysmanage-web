<!--Card made using tailwind-->

<script lang="ts">
    export let title: string;

    export let onclick = () => {};
</script>


<article on:keypress={onclick} on:click={onclick} class="block p-6 m-3 max-w-sm bg-white shadow-md rounded-lg border border-gray-200 dark:bg-gray-800 dark:border-gray-700">
    <a href={"javascript:void(0)"}>
        <h5 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">{@html title}</h5>
    </a>
    <p class="font-normal text-gray-700 dark:text-gray-400"><slot /></p>
</article>
