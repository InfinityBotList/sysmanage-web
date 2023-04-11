import { error } from "./strings";

const timeout = 500;

export const newTask = (logId: string, callback: (outp: string[]) => void) => {
    let c = setInterval(async () => {
        let res = await fetch(`/api/getLogEntry?id=${logId}`, {
            method: "POST",
        });

        if (!res.ok) {
            let errorText = await res.text();
            error(errorText);
        }

        let xIsDone = res.headers.get("X-Is-Done");
        console.log(xIsDone)

        let out = await res.json();
        callback(out)

        if(res.headers.get("X-Is-Done")) {
            console.log("Cancelling polling...")
            clearInterval(c);
            return
        }

        console.log("Polling...")
    }, timeout);
}