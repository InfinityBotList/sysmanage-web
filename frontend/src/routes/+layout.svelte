<script>
	import Menubar from '$lib/components/Menubar.svelte';
	import { SvelteToast } from '@zerodevx/svelte-toast'
	import './styles.css';

	const options = {
		
	}
</script>

<!-- Make the div fill whole screen-->
<div class="h-screen w-screen bg-white dark:bg-gray-900">
	<div class="flex flex-col h-full">
		<!--Menubar using tailwind-->
		<Menubar />
		<main class="bg-white dark:bg-gray-900 text-black-700 dark:text-white bg-contain p-8">
			<slot />
		</main>
		<SvelteToast {options} />
	</div>
</div>