interface Data {
    allowedIps: string[],
    type: string,
    url: string,
    token: string,
    ref: string,
    brokenValue: string,
    outputPath: string,
    commands: string[],
    webhooks: DeployWebhook[],
    timeout: number,
    env: [string, string][],
    configFiles: string[],
}

interface DeployWebhook {
    id: string,
    token: string,
    type: string,
}

export {
    Data,
    DeployWebhook,
}