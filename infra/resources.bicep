targetScope = 'resourceGroup'

param location string
param tags object
param resourceToken string
param containerRegistryLoginServer string
param containerRegistryUsername string

@secure()
param containerRegistryPassword string

@secure()
param twitterConsumerKey string

@secure()
param twitterConsumerSecret string

@secure()
param twitterAccessToken string

@secure()
param twitterAccessTokenSecret string

// ── Storage Account ───────────────────────────────────────────────────────────
// Hosts the 'state' table that tracks the novel line index.
// 'st' prefix + 13-char token = 15 chars (well within the 24-char limit).
resource storageAccount 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'st${resourceToken}'
  location: location
  tags: tags
  sku: { name: 'Standard_LRS' }
  kind: 'StorageV2'
  properties: {
    minimumTlsVersion: 'TLS1_2'
    allowBlobPublicAccess: false
    supportsHttpsTrafficOnly: true
  }
}

// ── Log Analytics Workspace ───────────────────────────────────────────────────
// Required for Container Apps Environment structured logging.
resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  name: 'log-${resourceToken}'
  location: location
  tags: tags
  properties: {
    sku: { name: 'PerGB2018' }
    retentionInDays: 30
  }
}

// ── Container Apps Environment ────────────────────────────────────────────────
resource containerAppsEnv 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: 'cae-${resourceToken}'
  location: location
  tags: tags
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: logAnalytics.properties.customerId
        sharedKey: logAnalytics.listKeys().primarySharedKey
      }
    }
  }
}

// ── Container Apps Job ────────────────────────────────────────────────────────
// Scheduled job: runs daily at 5:00 PM UTC.
// The 'azd-service-name' tag lets azd update the image on `azd deploy`.
resource containerAppsJob 'Microsoft.App/jobs@2024-03-01' = {
  name: 'job-${resourceToken}'
  location: location
  tags: union(tags, { 'azd-service-name': 'snoopybot' })
  properties: {
    environmentId: containerAppsEnv.id
    configuration: {
      triggerType: 'Schedule'
      replicaTimeout: 300
      scheduleTriggerConfig: {
        cronExpression: '0 17 * * *'
        parallelism: 1
        replicaCompletionCount: 1
      }
      secrets: [
        { name: 'twitter-consumer-key', value: twitterConsumerKey }
        { name: 'twitter-consumer-secret', value: twitterConsumerSecret }
        { name: 'twitter-access-token', value: twitterAccessToken }
        { name: 'twitter-access-token-secret', value: twitterAccessTokenSecret }
        { name: 'azure-storage-key', value: storageAccount.listKeys().keys[0].value }
        { name: 'registry-password', value: containerRegistryPassword }
      ]
      registries: [
        {
          server: containerRegistryLoginServer
          username: containerRegistryUsername
          passwordSecretRef: 'registry-password'
        }
      ]
    }
    template: {
      containers: [
        {
          name: 'snoopybot'
          image: '${containerRegistryLoginServer}/snoopybot:latest'
          resources: {
            cpu: json('0.25')
            memory: '0.5Gi'
          }
          env: [
            { name: 'TWITTER_CONSUMER_KEY', secretRef: 'twitter-consumer-key' }
            { name: 'TWITTER_CONSUMER_SECRET', secretRef: 'twitter-consumer-secret' }
            { name: 'TWITTER_ACCESS_TOKEN', secretRef: 'twitter-access-token' }
            { name: 'TWITTER_ACCESS_TOKEN_SECRET', secretRef: 'twitter-access-token-secret' }
            { name: 'AZURE_STORAGE_ACCOUNT', value: storageAccount.name }
            { name: 'AZURE_STORAGE_ACCESS_KEY', secretRef: 'azure-storage-key' }
          ]
        }
      ]
    }
  }
}

// ── Outputs ───────────────────────────────────────────────────────────────────
output storageAccountName string = storageAccount.name
output storageAccountKey string = storageAccount.listKeys().keys[0].value
output jobName string = containerAppsJob.name
