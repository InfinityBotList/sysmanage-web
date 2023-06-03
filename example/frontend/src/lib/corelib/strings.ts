import { toast } from'@zerodevx/svelte-toast'

export function title(str: string) {
    return str.replace(/(^|\s)\S/g, function(t: string) { return t.toUpperCase() });
}

export function cast(a: any): string {
    return a.toString()
}

export function castToObj(a: any): any {
    if(typeof a === "object") {
        return a
    }
    return {}
}

export const success = (m: string) => toast.push(m, {
  theme: {
    '--toastBackground': 'black',
    '--toastColor': 'white',
    '--toastBarBackground': 'olive'
  }
})

export const warning = (m: string) => toast.push(m, { theme: { 
    '--toastBackground': 'black',
    '--toastColor': 'white',
    '--toastBarBackground': 'orange'
 } 
})

export const error = (m: string) => toast.push(m, { theme: { 
    '--toastBackground': 'black',
    '--toastColor': 'white',
    '--toastBarBackground': 'maroon'
 } })
