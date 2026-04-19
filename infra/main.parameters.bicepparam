using './main.bicep'

// `environmentName` and `location` are injected automatically by azd — no entry needed here.
// azd will prompt for `location` once on first provision and remember it.
//
// All other variables must be set before running `azd provision`:
//
//   azd env set AZURE_CONTAINER_REGISTRY_LOGIN_SERVER  <registry>.azurecr.io
//   azd env set AZURE_CONTAINER_REGISTRY_USERNAME      <acr-admin-username>
//   azd env set AZURE_CONTAINER_REGISTRY_PASSWORD      <acr-admin-password>
//   azd env set TWITTER_CONSUMER_KEY                   <value>
//   azd env set TWITTER_CONSUMER_SECRET                <value>
//   azd env set TWITTER_ACCESS_TOKEN                   <value>
//   azd env set TWITTER_ACCESS_TOKEN_SECRET            <value>

param containerRegistryLoginServer = readEnvironmentVariable('AZURE_CONTAINER_REGISTRY_LOGIN_SERVER')
param containerRegistryUsername    = readEnvironmentVariable('AZURE_CONTAINER_REGISTRY_USERNAME')
param containerRegistryPassword    = readEnvironmentVariable('AZURE_CONTAINER_REGISTRY_PASSWORD')

param twitterConsumerKey       = readEnvironmentVariable('TWITTER_CONSUMER_KEY')
param twitterConsumerSecret    = readEnvironmentVariable('TWITTER_CONSUMER_SECRET')
param twitterAccessToken       = readEnvironmentVariable('TWITTER_ACCESS_TOKEN')
param twitterAccessTokenSecret = readEnvironmentVariable('TWITTER_ACCESS_TOKEN_SECRET')
