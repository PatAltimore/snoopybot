targetScope = 'subscription'

@minLength(1)
@maxLength(64)
@description('Name of the environment (e.g. dev, prod). Used for resource naming.')
param environmentName string

@minLength(1)
@description('Primary Azure region for all resources.')
param location string

@description('Login server of the existing Azure Container Registry (e.g. myregistry.azurecr.io).')
param containerRegistryLoginServer string

@description('Admin username for the existing Azure Container Registry.')
param containerRegistryUsername string

@description('Admin password for the existing Azure Container Registry.')
@secure()
param containerRegistryPassword string

@description('Mastodon instance URL (e.g. https://mastodon.social).')
param mastodonServer string

@description('Mastodon access token with write:statuses scope.')
@secure()
param mastodonAccessToken string

// ── Naming ────────────────────────────────────────────────────────────────────
var tags = { 'azd-env-name': environmentName }

// uniqueString returns 13 lowercase alphanumeric chars — safe for storage naming
var resourceToken = toLower(uniqueString(subscription().id, environmentName, location))

// ── Resource Group ────────────────────────────────────────────────────────────
resource rg 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: 'rg-${environmentName}'
  location: location
  tags: tags
}

// ── All resources (resource group scope via module) ───────────────────────────
module resources './resources.bicep' = {
  name: 'resources'
  scope: rg
  params: {
    location: location
    tags: tags
    resourceToken: resourceToken
    containerRegistryLoginServer: containerRegistryLoginServer
    containerRegistryUsername: containerRegistryUsername
    containerRegistryPassword: containerRegistryPassword
    mastodonServer: mastodonServer
    mastodonAccessToken: mastodonAccessToken
  }
}

// ── Outputs ───────────────────────────────────────────────────────────────────
// azd stores these in .azure/<env>/.env and injects them into hook scripts

@description('Storage account name — set in .env as AZURE_STORAGE_ACCOUNT')
output AZURE_STORAGE_ACCOUNT string = resources.outputs.storageAccountName

@description('Storage account key — set in .env as AZURE_STORAGE_ACCESS_KEY')
output AZURE_STORAGE_ACCESS_KEY string = resources.outputs.storageAccountKey

@description('Container Apps Job name — set as CONTAINERAPPS_JOB_NAME GitHub secret')
output AZURE_CONTAINER_APPS_JOB_NAME string = resources.outputs.jobName

@description('Resource group name — set as AZURE_RG GitHub secret')
output AZURE_RESOURCE_GROUP string = rg.name
