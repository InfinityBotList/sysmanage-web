interface Link {
    name: string;
    description: string;
    link: string;
    plugin?: string;
}

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
