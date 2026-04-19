# snoopybot

Bot that posts lines from Snoopy's novel and other writings from the Peanuts comic strip to Mastodon, Threads, or both. Runs as an Azure Container Apps Job on a daily schedule.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Azure Developer CLI (azd)](https://learn.microsoft.com/azure/developer/azure-developer-cli/install-azd)
- [Docker](https://www.docker.com/products/docker-desktop/) (for building the container image)
- An existing [Azure Container Registry](https://learn.microsoft.com/azure/container-registry/)
- A Mastodon account and/or Threads account (at least one required)

## Local development

Fill in your credentials in `.env`. At least one posting platform must be configured.

**Mastodon** (optional):

| Variable | Description |
|---|---|
| `MASTODON_SERVER` | Your Mastodon instance URL (e.g. `https://mastodon.social`) |
| `MASTODON_ACCESS_TOKEN` | Access token with `write:statuses` scope |

To get a Mastodon access token: go to your instance → **Preferences** → **Development** → **New Application** → enable `write:statuses` → copy the access token.

**Threads** (optional):

| Variable | Description |
|---|---|
| `THREADS_USER_ID` | Your Threads user ID |
| `THREADS_ACCESS_TOKEN` | Long-lived access token with `threads_content_publish` scope |

To get Threads credentials: go to the [Meta developer portal](https://developers.facebook.com) → create an app → add the Threads product → generate a long-lived access token. Your user ID is available from `GET https://graph.threads.net/v1.0/me?access_token=<token>`.

**Always required**:

| Variable | Description |
|---|---|
| `AZURE_STORAGE_ACCOUNT` | Azure Storage account name |
| `AZURE_STORAGE_ACCESS_KEY` | Azure Storage account key |
| `DRY_RUN` | Set to `true` to print the post without sending it |

Run locally (dry run by default):

```bash
go run .
```

To fire a real post, set `DRY_RUN=false` in `.env` before running. Both platforms post the same content simultaneously if both are configured.

## Project structure

```
main.go                            # Entry point: validates env vars, calls bot.DoWork()
internal/
  bot/
    data.go                        # Novel lines and miscellaneous quotes
    bot.go                         # DoWork(): coin flip, post novel or misc quote
  mastodon/
    client.go                      # Mastodon API client (POST /api/v1/statuses)
  threads/
    client.go                      # Threads API client (two-step container + publish)
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

# Platform credentials — set one or both
azd env set MASTODON_SERVER        https://mastodon.social
azd env set MASTODON_ACCESS_TOKEN  <value>

azd env set THREADS_USER_ID        <value>
azd env set THREADS_ACCESS_TOKEN   <value>
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
                                       ┌─────────────────────┐
                                       │  Mastodon API v1    │
                                       │  Threads API v1.0   │
                                       │  (any/all enabled)  │
                                       └─────────────────────┘
```
