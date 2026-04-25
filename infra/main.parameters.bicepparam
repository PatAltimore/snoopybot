using './main.bicep'

// `environmentName` and `location` are injected automatically by azd — no entry needed here.
// azd will prompt for `location` once on first provision and remember it.
//
// All other variables must be set before running `azd provision`:
//
//   azd env set AZURE_CONTAINER_REGISTRY_LOGIN_SERVER  <registry>.azurecr.io
//   azd env set AZURE_CONTAINER_REGISTRY_USERNAME      <acr-admin-username>
//   azd env set AZURE_CONTAINER_REGISTRY_PASSWORD      <acr-admin-password>
//   azd env set MASTODON_SERVER                        https://mastodon.social
//   azd env set MASTODON_ACCESS_TOKEN                  <value>

param containerRegistryLoginServer = readEnvironmentVariable('AZURE_CONTAINER_REGISTRY_LOGIN_SERVER')
param containerRegistryUsername    = readEnvironmentVariable('AZURE_CONTAINER_REGISTRY_USERNAME')
param containerRegistryPassword    = readEnvironmentVariable('AZURE_CONTAINER_REGISTRY_PASSWORD')

param mastodonServer      = readEnvironmentVariable('MASTODON_SERVER')
param mastodonAccessToken = readEnvironmentVariable('MASTODON_ACCESS_TOKEN')
