#!/bin/bash
# postprovision.sh
# Runs after `azd provision`. Updates the local .env file with the storage
# account values output by Bicep, and prints the GitHub secrets needed for CI/CD.
set -e

ENV_FILE="$(git rev-parse --show-toplevel)/.env"

update_env_var() {
  local key="$1"
  local value="$2"
  if [ -z "$value" ]; then
    echo "  ! $key not found in provisioning outputs, skipping"
    return
  fi
  if grep -q "^${key}=" "$ENV_FILE" 2>/dev/null; then
    sed -i.bak "s|^${key}=.*|${key}=${value}|" "$ENV_FILE"
    rm -f "$ENV_FILE.bak"
  else
    echo "${key}=${value}" >> "$ENV_FILE"
  fi
  echo "  ✓ $key"
}

echo ""
echo "==> Updating .env with provisioned Azure storage credentials..."
update_env_var "AZURE_STORAGE_ACCOUNT"    "$AZURE_STORAGE_ACCOUNT"
update_env_var "AZURE_STORAGE_ACCESS_KEY" "$AZURE_STORAGE_ACCESS_KEY"
echo ""

echo "==> Set these GitHub Actions secrets for CI/CD deployment:"
echo ""
echo "   AZURE_RG                  = $AZURE_RESOURCE_GROUP"
echo "   CONTAINERAPPS_JOB_NAME    = $AZURE_CONTAINER_APPS_JOB_NAME"
echo "   ACR_LOGIN_SERVER          = $AZURE_CONTAINER_REGISTRY_LOGIN_SERVER"
echo "   ACR_USERNAME              = $AZURE_CONTAINER_REGISTRY_USERNAME"
echo "   ACR_PASSWORD              = (your ACR admin password)"
echo ""
echo "   gh secret set AZURE_RG                 --body \"$AZURE_RESOURCE_GROUP\""
echo "   gh secret set CONTAINERAPPS_JOB_NAME   --body \"$AZURE_CONTAINER_APPS_JOB_NAME\""
echo ""
echo "Done."
