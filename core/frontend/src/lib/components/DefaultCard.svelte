<!--Card made using tailwind-->

<script lang="ts">
    export let title: string;

    export let onclick = () => {};
</script>


<article class="w-full block p-6 m-3 md:w-1/2 sm:w-96 bg-white rounded-lg shadow-md border border-gray-200 dark:bg-gray-800 dark:border-gray-700">
    <h5 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white hover:cursor-pointer" on:click={onclick} on:keypress={onclick}>{@html title}</h5>
    <p class="font-normal text-gray-700 dark:text-gray-400"><slot /></p>
</article>
