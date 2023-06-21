interface Link {
    Title:       string
	Description: string
	LinkText:    string
	Href:        string
}

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