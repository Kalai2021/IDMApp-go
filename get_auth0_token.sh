#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    source .env
else
    echo "Error: .env file not found. Please create one based on env.example"
    exit 1
fi

# Auth0 Configuration from environment variables
AUTH0_DOMAIN="${AUTH0_DOMAIN:-your-domain.auth0.com}"
AUTH0_CLIENT_ID="${AUTH0_CLIENT_ID:-your-client-id}"
AUTH0_CLIENT_SECRET="${AUTH0_CLIENT_SECRET:-your-client-secret}"
AUTH0_AUDIENCE="${AUTH0_AUDIENCE:-your-audience}"

# Check if credentials are properly configured
if [ "$AUTH0_DOMAIN" = "your-domain.auth0.com" ] || [ "$AUTH0_CLIENT_ID" = "your-client-id" ]; then
    echo "Error: Please configure your Auth0 credentials in the .env file"
    echo "Required variables: AUTH0_DOMAIN, AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_AUDIENCE"
    exit 1
fi

echo "Getting OAuth Token from Auth0"
echo "=============================="

# Get client credentials token
echo "1. Getting client credentials token..."
TOKEN_RESPONSE=$(curl -s -X POST "https://$AUTH0_DOMAIN/oauth/token" \
  -H "Content-Type: application/json" \
  -d "{
    \"client_id\": \"$AUTH0_CLIENT_ID\",
    \"client_secret\": \"$AUTH0_CLIENT_SECRET\",
    \"audience\": \"$AUTH0_AUDIENCE\",
    \"grant_type\": \"client_credentials\"
  }")

echo "Token Response:"
echo $TOKEN_RESPONSE | jq .

# Extract access token
ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.access_token')
echo -e "\nAccess Token:"
echo $ACCESS_TOKEN

# Save token to file for use in other scripts
echo $ACCESS_TOKEN > .auth0_token
echo -e "\nToken saved to .auth0_token file"

echo -e "\nOAuth Token obtained successfully!" 