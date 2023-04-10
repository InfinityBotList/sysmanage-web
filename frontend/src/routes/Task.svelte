<script lang="ts">
	import { error } from "$lib/strings";

    export let pollUrl: string;

    let timeout = 500;

    let pollConsole: string[] = [];

    let c = setInterval(async () => {
        let res = await fetch(pollUrl, {
            method: "POST",
        });

        if (!res.ok) {
            let errorText = await res.text();
            error(errorText);
        }

        let xIsDone = res.headers.get("X-Is-Done");
        console.log(xIsDone)
        if(res.headers.get("X-Is-Done")) {
            console.log("Cancelling polling...")
            clearInterval(c);
            return
        }

        let out = await res.json();

        console.log("Polling...")
        pollConsole = out;
    }, timeout);
</script>

<pre class="p-4 bg-gray-300 text-black dark:bg-gray-400 break-all whitespace-pre-wrap">{pollConsole.join("")}</pre>