<!--Tailwind menubar-->

<script lang="ts">
  import { getLinks } from "$lib/corelib/links";

  function onClickMenu() {
    console.log("clicked");

    const menu = document.getElementById("menu");

    if (menu) {
      menu.classList.toggle("hidden");
    }
  }
</script>

<nav class="bg-white border-gray-200 px-2 sm:px-4 rounded dark:bg-gray-900">
  <div class="flex flex-wrap justify-between items-center mx-auto">
    <a href="/" class="flex items-center text-black">
        <img src="https://cdn.infinitybots.gg/images/png/Infinity.png" class="mr-3 h-6 sm:h-9" alt="IBL Logo" />
        <span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">Infinity Portal</span>
    </a>
    <button on:click={onClickMenu} type="button" class="inline-flex items-center p-2 ml-3 text-sm text-gray-500 rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600" aria-controls="navbar-default" aria-expanded="false">
      <span class="sr-only">Open main menu</span>
      <svg class="w-6 h-6" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd"></path></svg>
    </button>
  </div>
</nav>

<!--Menu for mobile users-->

<div id="menu" class="hidden">
  <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
    <a href="/" class="block px-3 py-2 text-base font-medium text-violet-600 rounded-md hover:text-gray-900 hover:bg-gray-50 dark:text-amber-400 dark:hover:text-white dark:hover:bg-gray-700">Home</a>

    {#await getLinks() then links}
      {#each links as link}
        <a href={link.Href} class="block px-3 py-2 text-base font-medium text-violet-600 rounded-md hover:text-gray-900 hover:bg-gray-50 dark:text-amber-400 dark:hover:text-white dark:hover:bg-gray-700">{link.Title}</a>
      {/each}
    {:catch error}
      <h2 class="text-xl">Error: {error?.message ? error?.message : error}</h2>
    {/await}
  </div>
</div>