import { error } from "./strings";

const timeout = 500;

export const newTask = (logId: string, callback: (outp: string[]) => void) => {
    let hasStarted = false;
    let c = setInterval(async () => {
        let res = await fetch(`/api/logger/getLogEntry?id=${logId}`, {
            method: "POST",
        });

        if (!res.ok) {
            let errorText = await res.text();
            error(errorText);
        }

        let xIsDone = res.headers.get("X-Is-Done");
        console.log(xIsDone)

        let out = await res.json();

        if(out?.length == 0) {
            if(!hasStarted) {
                console.log("No output yet...")
                return
            } else {
                console.log("Task is done.")
                console.log("Cancelling polling...")
                clearInterval(c);
                return    
            }
        }

        hasStarted = true;

        callback(out)

        if(res.headers.get("X-Is-Done")) {
            console.log("Cancelling polling...")
            clearInterval(c);
            return
        }

        console.log("Polling...")
    }, timeout);
}