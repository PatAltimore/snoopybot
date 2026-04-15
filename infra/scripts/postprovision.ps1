# postprovision.ps1
# Runs after `azd provision`. Updates the local .env file with the storage
# account values output by Bicep, and prints the GitHub secrets needed for CI/CD.

$ErrorActionPreference = 'Stop'

$repoRoot = git rev-parse --show-toplevel
$envFile  = Join-Path $repoRoot '.env'

function Update-EnvVar {
    param([string]$Key, [string]$Value)

    if (-not $Value) {
        Write-Host "  ! $Key not found in provisioning outputs, skipping"
        return
    }

    $content = if (Test-Path $envFile) { Get-Content $envFile } else { @() }

    if ($content -match "^${Key}=") {
        $content = $content -replace "^${Key}=.*", "${Key}=${Value}"
    } else {
        $content += "${Key}=${Value}"
    }

    $content | Set-Content $envFile -Encoding UTF8
    Write-Host "  ✓ $Key"
}

Write-Host ""
Write-Host "==> Updating .env with provisioned Azure storage credentials..."
Update-EnvVar 'AZURE_STORAGE_ACCOUNT'    $env:AZURE_STORAGE_ACCOUNT
Update-EnvVar 'AZURE_STORAGE_ACCESS_KEY' $env:AZURE_STORAGE_ACCESS_KEY
Write-Host ""

Write-Host "==> Set these GitHub Actions secrets for CI/CD deployment:"
Write-Host ""
Write-Host "   AZURE_RG               = $env:AZURE_RESOURCE_GROUP"
Write-Host "   CONTAINERAPPS_JOB_NAME = $env:AZURE_CONTAINER_APPS_JOB_NAME"
Write-Host "   ACR_LOGIN_SERVER       = $env:AZURE_CONTAINER_REGISTRY_LOGIN_SERVER"
Write-Host "   ACR_USERNAME           = $env:AZURE_CONTAINER_REGISTRY_USERNAME"
Write-Host "   ACR_PASSWORD           = (your ACR admin password)"
Write-Host ""
Write-Host "   gh secret set AZURE_RG               --body `"$env:AZURE_RESOURCE_GROUP`""
Write-Host "   gh secret set CONTAINERAPPS_JOB_NAME --body `"$env:AZURE_CONTAINER_APPS_JOB_NAME`""
Write-Host ""
Write-Host "Done."
