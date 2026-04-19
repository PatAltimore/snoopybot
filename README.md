# snoopybot

Twitter bot that tweets lines from Snoopy's novel and other writings from the Peanuts comic strip. Runs as an Azure Container Apps Job on a daily schedule.

See [@SnoopyAtWork](https://twitter.com/SnoopyAtWork) on Twitter.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Azure Developer CLI (azd)](https://learn.microsoft.com/azure/developer/azure-developer-cli/install-azd)
- [Docker](https://www.docker.com/products/docker-desktop/) (for building the container image)
- An existing [Azure Container Registry](https://learn.microsoft.com/azure/container-registry/)
- Twitter/X API credentials (OAuth 1.0a — consumer key/secret + access token/secret)

## Local development

Copy `.env` and fill in your credentials:

```bash
cp .env .env.local   # or just edit .env directly
```

| Variable | Description |
|---|---|
| `TWITTER_CONSUMER_KEY` | Twitter/X API consumer key |
| `TWITTER_CONSUMER_SECRET` | Twitter/X API consumer secret |
| `TWITTER_ACCESS_TOKEN` | Twitter/X access token |
| `TWITTER_ACCESS_TOKEN_SECRET` | Twitter/X access token secret |
| `AZURE_STORAGE_ACCOUNT` | Azure Storage account name |
| `AZURE_STORAGE_ACCESS_KEY` | Azure Storage account key |
| `DRY_RUN` | Set to `true` to print the tweet without posting it |

Run locally (dry run by default):

```bash
go run .
```

To fire a real tweet, remove `DRY_RUN=true` from `.env` (or set it to `false`) before running.

## Project structure

```
main.go                            # Entry point: validates env vars, calls bot.DoWork()
internal/
  bot/
    data.go                        # Novel lines and miscellaneous quotes
    bot.go                         # DoWork(): coin flip, tweet novel or misc quote
  twitter/
    client.go                      # Minimal OAuth1 Twitter API v2 client
  storage/
    state.go                       # Azure Tables: tracks current novel line index
Dockerfile                         # Multi-stage build → scratch-based image (~8 MB)
azure.yaml                         # Azure Developer CLI template definition
infra/
  main.bicep                       # Subscription-scope: resource group + module
  resources.bicep                  # All Azure resources (storage, container apps env, job)
  main.parameters.bicepparam       # Reads values from azd environment
  scripts/
    postprovision.sh               # Hook: updates .env and prints GitHub secret commands
    postprovision.ps1              # Same, for Windows/PowerShell
.github/
  workflows/
    deploy.yml                     # CI/CD: build image → push to ACR → update job
```

## Provision to Azure

Infrastructure is managed with the [Azure Developer CLI](https://learn.microsoft.com/azure/developer/azure-developer-cli/). A single `azd provision` creates the resource group, storage account, Container Apps Environment, and Container Apps Job.

### 1. Set required environment variables

`azd` handles `AZURE_ENV_NAME` and `AZURE_LOCATION` automatically — it will prompt for the location on first run and remember it. Set the remaining values before provisioning:

```bash
# Existing Azure Container Registry
azd env set AZURE_CONTAINER_REGISTRY_LOGIN_SERVER  <registry>.azurecr.io
azd env set AZURE_CONTAINER_REGISTRY_USERNAME      <acr-admin-username>
azd env set AZURE_CONTAINER_REGISTRY_PASSWORD      <acr-admin-password>

# Twitter/X API credentials
azd env set TWITTER_CONSUMER_KEY        <value>
azd env set TWITTER_CONSUMER_SECRET     <value>
azd env set TWITTER_ACCESS_TOKEN        <value>
azd env set TWITTER_ACCESS_TOKEN_SECRET <value>
```

> ACR admin credentials are found in the Azure Portal under your registry → **Access keys**. Enable the admin user if it is not already on.

### 2. Provision infrastructure

```bash
azd provision
```

This creates the following resources inside a new `rg-<env>` resource group:

| Resource | Purpose |
|---|---|
| Storage Account | Hosts the `state` table tracking the novel line index |
| Log Analytics Workspace | Container Apps structured logging |
| Container Apps Environment | Shared environment for the job |
| Container Apps Job | Runs daily at 5:00 PM UTC; 0.25 vCPU / 0.5 GiB |

After provisioning, the `postprovision` hook automatically updates `AZURE_STORAGE_ACCOUNT` and `AZURE_STORAGE_ACCESS_KEY` in your local `.env` file. It also prints the `gh secret set` commands needed for CI/CD.

### 3. Configure GitHub Actions secrets

The deployment workflow (`.github/workflows/deploy.yml`) requires these repository secrets:

| Secret | Where to get it |
|---|---|
| `ACR_LOGIN_SERVER` | Your ACR login server (e.g. `myregistry.azurecr.io`) |
| `ACR_USERNAME` | ACR admin username |
| `ACR_PASSWORD` | ACR admin password |
| `AZURE_CREDENTIALS` | Service principal JSON — see below |
| `AZURE_RG` | Printed by `postprovision` hook |
| `CONTAINERAPPS_JOB_NAME` | Printed by `postprovision` hook |

Create the `AZURE_CREDENTIALS` service principal:

```bash
az ad sp create-for-rbac \
  --name snoopybot-deploy \
  --role Contributor \
  --scopes /subscriptions/<subscription-id>/resourceGroups/rg-<env> \
  --sdk-auth
```

Copy the JSON output as the `AZURE_CREDENTIALS` secret.

### 4. Deploy the container image

Push to `master` to trigger the GitHub Actions workflow, which builds the Docker image, pushes it to ACR, and updates the Container Apps Job. To manually trigger a job execution after deploying:

```bash
az containerapp job start --name <job-name> --resource-group rg-<env>
```

## Architecture

```
GitHub push → Actions workflow → docker build → ACR → Container Apps Job
                                                           │
                                            (daily 5 PM UTC cron)
                                                           │
                                                    Go binary runs
                                                    ┌──────┴──────┐
                                               50% novel     50% misc
                                                    │
                                            Azure Table Storage
                                            (tracks novel index)
                                                    │
                                              Twitter/X API v2
                                               (posts tweet)
```
