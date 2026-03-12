# snoopybot
Twitter bot that tweets lines from Snoopy's novel and other writings from the Peanuts comic strip. Runs as an Azure Functions Timer Trigger.

See [@SnoopyAtWork](https://twitter.com/SnoopyAtWork) on Twitter.

## Prerequisites

- Node.js 18+
- [Azure Functions Core Tools v4](https://learn.microsoft.com/azure/azure-functions/functions-run-local)
- An Azure Storage account (used by both the bot state and Azure Functions runtime)
- Twitter/X API credentials

## Local development

```bash
npm install
```

Fill in your secrets in `local.settings.json`:

| Setting | Description |
|---|---|
| `TWITTER_CONSUMER_KEY` | Twitter/X API consumer key |
| `TWITTER_CONSUMER_SECRET` | Twitter/X API consumer secret |
| `TWITTER_ACCESS_TOKEN` | Twitter/X access token |
| `TWITTER_ACCESS_TOKEN_SECRET` | Twitter/X access token secret |
| `AZURE_STORAGE_ACCOUNT` | Azure Storage account name |
| `AZURE_STORAGE_ACCESS_KEY` | Azure Storage account key |

Start the function locally:

```bash
npm start
```

The timer trigger is set to run daily at 5:00 PM UTC. To test immediately, use the Azure Functions Core Tools HTTP endpoint:

```bash
curl -X POST http://localhost:7071/admin/functions/snoopyTrigger
```

## Deploy to Azure

1. Create a Function App in the Azure Portal (Node.js 18+, Consumption plan)
2. Add the settings from the table above as **Application Settings**
3. Deploy using one of:
   - VS Code Azure Functions extension
   - `func azure functionapp publish <your-function-app-name>`
   - GitHub Actions

## Project structure

```
host.json                      # Azure Functions host config
local.settings.json            # Local dev settings (gitignored)
package.json
src/
  bot.js                       # Bot logic (tweets, Azure Table state)
  functions/
    snoopyTrigger.js           # Timer trigger (daily at 5 PM UTC)
```
