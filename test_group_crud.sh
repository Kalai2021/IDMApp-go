#!/bin/bash

# Group CRUD Operations Test Script
# This script tests all CRUD operations for Groups with Auth0 authentication

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Testing Complete Group CRUD Operations with Local JWT Token${NC}"

# Base URL
BASE_URL="http://localhost:8080/api/v1"

# Get local JWT token
echo -e "${YELLOW}Getting local JWT token...${NC}"

# For testing purposes, use a simple test token
# In production, you would get this from your local authentication endpoint
TOKEN="test-token"

echo -e "${GREEN}Local JWT token obtained successfully${NC}"

# Test 1: Get all groups (existing data)
echo -e "\n${BLUE}Test 1: Get all groups${NC}"
curl -s -X GET "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 2: Create a new group with unique name
echo -e "\n${BLUE}Test 2: Create a new group${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-developers-2024",
    "displayName": "Test Development Team",
    "description": "Test software development team"
  }')

echo "$CREATE_RESPONSE" | jq '.'

# Extract group ID from response
GROUP_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')
if [ "$GROUP_ID" = "null" ] || [ -z "$GROUP_ID" ]; then
    echo -e "${RED}Failed to create group or extract ID${NC}"
    exit 1
fi

echo -e "${GREEN}Group created with ID: $GROUP_ID${NC}"

# Test 3: Get the created group
echo -e "\n${BLUE}Test 3: Get the created group${NC}"
curl -s -X GET "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 4: Get all groups (should now contain the created group)
echo -e "\n${BLUE}Test 4: Get all groups (after creation)${NC}"
curl -s -X GET "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 5: Update the group
echo -e "\n${BLUE}Test 5: Update the group${NC}"
curl -s -X PUT "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-senior-developers-2024",
    "displayName": "Test Senior Development Team",
    "description": "Test senior software development team",
    "isActive": true
  }' | jq '.'

# Test 6: Get the updated group
echo -e "\n${BLUE}Test 6: Get the updated group${NC}"
curl -s -X GET "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 7: Partial update (only description)
echo -e "\n${BLUE}Test 7: Partial update (only description)${NC}"
curl -s -X PUT "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description for test senior development team"
  }' | jq '.'

# Test 8: Deactivate the group
echo -e "\n${BLUE}Test 8: Deactivate the group${NC}"
curl -s -X PUT "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "isActive": false
  }' | jq '.'

# Test 9: Get the deactivated group
echo -e "\n${BLUE}Test 9: Get the deactivated group${NC}"
curl -s -X GET "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 10: Create another group for testing
echo -e "\n${BLUE}Test 10: Create another group${NC}"
CREATE_RESPONSE2=$(curl -s -X POST "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-qa-team-2024",
    "displayName": "Test QA Team",
    "description": "Test quality assurance team"
  }')

echo "$CREATE_RESPONSE2" | jq '.'

GROUP_ID2=$(echo "$CREATE_RESPONSE2" | jq -r '.id')
echo -e "${GREEN}Second group created with ID: $GROUP_ID2${NC}"

# Test 11: Get all groups (should contain both new groups)
echo -e "\n${BLUE}Test 11: Get all groups (both new groups)${NC}"
curl -s -X GET "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 12: Try to create a group with duplicate name (should fail)
echo -e "\n${BLUE}Test 12: Try to create group with duplicate name${NC}"
curl -s -X POST "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-developers-2024",
    "displayName": "Another Test Development Team",
    "description": "Another test development team"
  }' | jq '.'

# Test 13: Try to get non-existent group
echo -e "\n${BLUE}Test 13: Try to get non-existent group${NC}"
curl -s -X GET "$BASE_URL/groups/00000000-0000-0000-0000-000000000000" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 14: Reactivate the first group
echo -e "\n${BLUE}Test 14: Reactivate the first group${NC}"
curl -s -X PUT "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "isActive": true
  }' | jq '.'

# Test 15: Get the reactivated group
echo -e "\n${BLUE}Test 15: Get the reactivated group${NC}"
curl -s -X GET "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 16: Delete the first group
echo -e "\n${BLUE}Test 16: Delete the first group${NC}"
curl -s -X DELETE "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 17: Try to get the deleted group (should fail)
echo -e "\n${BLUE}Test 17: Try to get the deleted group${NC}"
curl -s -X GET "$BASE_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 18: Get all groups (should only contain the second new group)
echo -e "\n${BLUE}Test 18: Get all groups (after deletion)${NC}"
curl -s -X GET "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 19: Delete the second group
echo -e "\n${BLUE}Test 19: Delete the second group${NC}"
curl -s -X DELETE "$BASE_URL/groups/$GROUP_ID2" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

# Test 20: Get all groups (should be back to original state)
echo -e "\n${BLUE}Test 20: Get all groups (final state)${NC}"
curl -s -X GET "$BASE_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n${GREEN}=== Group CRUD Operations Test Completed ===${NC}" 