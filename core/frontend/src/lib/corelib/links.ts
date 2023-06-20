interface Link {
    Title:       string
	Description: string
	LinkText:    string
	Href:        string
}

/*
export const links: Link[] = [
    {
        name: "Service Management",
        description: "Systemd service management",
        link: "/plugins/systemd",
        plugin: "systemd"
    },
    {
        name: "Nginx Management",
        description: "Add, update, remove and manage nginx-proxied domains",
        link: "/plugins/nginx",
        plugin: "nginx"
    },
]
*/

let links: Link[] = []

export const getLinks = async (): Promise<Link[]> => {
    if(links.length) {
        return links // return cached links
    }

    let res = await fetch("/api/frontend/getRegisteredLinks", {
        method: "POST"
    })

    if(!res.ok) {
        throw new Error("Failed to fetch links")
    }

    links = await res.json()

    return links
}