#!/bin/bash

# Organization CRUD Operations Test Script
# This script tests all CRUD operations for Organizations with Auth0 authentication

BASE_URL="http://localhost:8090"
API_BASE="$BASE_URL/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Organization CRUD Operations Test ===${NC}"

# Get Auth0 token
echo -e "${YELLOW}Getting Auth0 token...${NC}"

# Load environment variables
if [ -f .env ]; then
    source .env
else
    echo -e "${RED}Error: .env file not found. Please create one based on env.example${NC}"
    exit 1
fi

# Use environment variables for Auth0 configuration
AUTH0_DOMAIN="${AUTH0_DOMAIN:-your-domain.auth0.com}"
AUTH0_CLIENT_ID="${AUTH0_CLIENT_ID:-your-client-id}"
AUTH0_CLIENT_SECRET="${AUTH0_CLIENT_SECRET:-your-client-secret}"
AUTH0_AUDIENCE="${AUTH0_AUDIENCE:-your-audience}"

# Check if credentials are properly configured
if [ "$AUTH0_DOMAIN" = "your-domain.auth0.com" ] || [ "$AUTH0_CLIENT_ID" = "your-client-id" ]; then
    echo -e "${RED}Error: Please configure your Auth0 credentials in the .env file${NC}"
    echo -e "${RED}Required variables: AUTH0_DOMAIN, AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_AUDIENCE${NC}"
    exit 1
fi

TOKEN=$(curl -s -X POST "https://$AUTH0_DOMAIN/oauth/token" \
  -H "Content-Type: application/json" \
  -d "{
    \"client_id\": \"$AUTH0_CLIENT_ID\",
    \"client_secret\": \"$AUTH0_CLIENT_SECRET\",
    \"audience\": \"$AUTH0_AUDIENCE\",
    \"grant_type\": \"client_credentials\"
  }" | jq -r '.access_token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get Auth0 token${NC}"
    exit 1
fi

echo -e "${GREEN}Auth0 token obtained successfully${NC}"

# Test 1: Get all organizations (should be empty initially)
echo -e "\n${BLUE}Test 1: Get all organizations${NC}"
curl -s -X GET "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 2: Create a new organization
echo -e "\n${BLUE}Test 2: Create a new organization${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "acme-corp",
    "displayName": "ACME Corporation",
    "description": "A leading technology company"
  }')

echo "$CREATE_RESPONSE" | jq '.'

# Extract organization ID from response
ORG_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')
if [ "$ORG_ID" = "null" ] || [ -z "$ORG_ID" ]; then
    echo -e "${RED}Failed to create organization or extract ID${NC}"
    exit 1
fi

echo -e "${GREEN}Organization created with ID: $ORG_ID${NC}"

# Test 3: Get the created organization
echo -e "\n${BLUE}Test 3: Get the created organization${NC}"
curl -s -X GET "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 4: Get all organizations (should now contain the created organization)
echo -e "\n${BLUE}Test 4: Get all organizations (after creation)${NC}"
curl -s -X GET "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 5: Update the organization
echo -e "\n${BLUE}Test 5: Update the organization${NC}"
curl -s -X PUT "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "acme-enterprise",
    "displayName": "ACME Enterprise Solutions",
    "description": "Enterprise technology solutions provider",
    "isActive": true
  }' | jq '.'

# Test 6: Get the updated organization
echo -e "\n${BLUE}Test 6: Get the updated organization${NC}"
curl -s -X GET "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 7: Partial update (only description)
echo -e "\n${BLUE}Test 7: Partial update (only description)${NC}"
curl -s -X PUT "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description for ACME Enterprise Solutions"
  }' | jq '.'

# Test 8: Deactivate the organization
echo -e "\n${BLUE}Test 8: Deactivate the organization${NC}"
curl -s -X PUT "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "isActive": false
  }' | jq '.'

# Test 9: Get the deactivated organization
echo -e "\n${BLUE}Test 9: Get the deactivated organization${NC}"
curl -s -X GET "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 10: Create another organization for testing
echo -e "\n${BLUE}Test 10: Create another organization${NC}"
CREATE_RESPONSE2=$(curl -s -X POST "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "tech-startup",
    "displayName": "Tech Startup Inc",
    "description": "Innovative startup company"
  }')

echo "$CREATE_RESPONSE2" | jq '.'

ORG_ID2=$(echo "$CREATE_RESPONSE2" | jq -r '.id')
echo -e "${GREEN}Second organization created with ID: $ORG_ID2${NC}"

# Test 11: Get all organizations (should contain both organizations)
echo -e "\n${BLUE}Test 11: Get all organizations (both organizations)${NC}"
curl -s -X GET "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 12: Try to create an organization with duplicate name (should fail)
echo -e "\n${BLUE}Test 12: Try to create organization with duplicate name${NC}"
curl -s -X POST "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "acme-corp",
    "displayName": "Another ACME Corp",
    "description": "Another ACME corporation"
  }' | jq '.'

# Test 13: Try to get non-existent organization
echo -e "\n${BLUE}Test 13: Try to get non-existent organization${NC}"
curl -s -X GET "$API_BASE/orgs/00000000-0000-0000-0000-000000000000" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 14: Reactivate the first organization
echo -e "\n${BLUE}Test 14: Reactivate the first organization${NC}"
curl -s -X PUT "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "isActive": true
  }' | jq '.'

# Test 15: Get the reactivated organization
echo -e "\n${BLUE}Test 15: Get the reactivated organization${NC}"
curl -s -X GET "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 16: Delete the first organization
echo -e "\n${BLUE}Test 16: Delete the first organization${NC}"
curl -s -X DELETE "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 17: Try to get the deleted organization (should fail)
echo -e "\n${BLUE}Test 17: Try to get the deleted organization${NC}"
curl -s -X GET "$API_BASE/orgs/$ORG_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 18: Get all organizations (should only contain the second organization)
echo -e "\n${BLUE}Test 18: Get all organizations (after deletion)${NC}"
curl -s -X GET "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 19: Delete the second organization
echo -e "\n${BLUE}Test 19: Delete the second organization${NC}"
curl -s -X DELETE "$API_BASE/orgs/$ORG_ID2" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 20: Get all organizations (should be empty again)
echo -e "\n${BLUE}Test 20: Get all organizations (final state)${NC}"
curl -s -X GET "$API_BASE/orgs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n${GREEN}=== Organization CRUD Operations Test Completed ===${NC}" 