<!--Tailwind menubar-->

<script lang="ts">
	import { goto } from "$app/navigation";

  function onClickMenu() {
    console.log("clicked");

    const menu = document.getElementById("menu");

    if (menu) {
      menu.classList.toggle("hidden");
    }
  }

  function logoutUser() {
    // Clear user_id, username and api_token cookies
    document.cookie = "user_id=; max-age=0; path=/;";
    document.cookie = "username=; max-age=0; path=/;";
    document.cookie = "api_token=; max-age=0; path=/;";
    window.location.href = "/";
  }

  function loginUser() {
    localStorage.setItem("redirect", window.location.href);
    goto("/login");
  }
</script>

<nav class="bg-white border-gray-200 px-2 sm:px-4 rounded dark:bg-gray-900">
  <div class="flex flex-wrap justify-between items-center mx-auto">
    <a href="/" class="flex items-center text-black">
        <img src="https://cdn.infinitybots.xyz/images/png/Infinity.png" class="mr-3 h-6 sm:h-9" alt="IBL Logo" />
        <span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">Infinity Portal</span>
    </a>
    <button on:click={onClickMenu} type="button" class="inline-flex items-center p-2 ml-3 text-sm text-gray-500 rounded-lg md:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600" aria-controls="navbar-default" aria-expanded="false">
      <span class="sr-only">Open main menu</span>
      <svg class="w-6 h-6" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd"></path></svg>
    </button>
    <div class="hidden w-full md:block md:w-auto" id="navbar-default">
      <ul class="flex flex-col p-4 mt-4 bg-gray-50 rounded-lg border border-gray-100 md:flex-row md:space-x-8 md:mt-0 md:text-sm md:font-medium md:border-0 md:bg-white dark:bg-gray-800 md:dark:bg-gray-900 dark:border-gray-700">
        <li>
          <a href={"/about"} class="block py-2 pr-4 pl-3 text-gray-700 rounded hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent">About</a>
        </li>
      </ul>
    </div>
  </div>
</nav>

<!--Menu for mobile users-->

<div id="menu" class="hidden md:hidden">
  <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
    <a href="/" class="block px-3 py-2 text-base font-medium text-violet-600 rounded-md hover:text-gray-900 hover:bg-gray-50 dark:text-amber-400 dark:hover:text-white dark:hover:bg-gray-700">Home</a>
  </div>
</div>