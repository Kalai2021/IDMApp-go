#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Testing Complete User CRUD Operations with Local JWT Token${NC}"

# Base URL
BASE_URL="http://localhost:8080"

# For testing purposes, let's create a simple JWT token
# In a real scenario, you would get this from your local authentication endpoint
echo -e "${YELLOW}Creating test JWT token...${NC}"

# Create a simple JWT token for testing (this is just for testing)
# In production, you would get this from your local authentication endpoint
TOKEN="test-token"

echo -e "${GREEN}Using test token: $TOKEN${NC}"

# Test 1: Health Check (no auth required)
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Login to get authentication token
echo -e "\n2. Login to get authentication token:"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }')

echo $LOGIN_RESPONSE | jq .

# Extract token from response
#TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
export TOKEN=$(cat .auth0_token)
echo "Token: $TOKEN"

# If no token, create a test user first
if [ "$TOKEN" = "null" ] || [ "$TOKEN" = "" ]; then
    echo "No token received. Let's create a test user first..."
    
    # For testing purposes, let's create a simple JWT token
    # In a real scenario, you would get this from Auth0
    # TEST_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXItaWQiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJpYXQiOjE2MzQ1Njc4NzgsImV4cCI6MTYzNDU3MTQ3OH0.test-signature"
    # TOKEN=$TEST_TOKEN
    echo "Using test token: $TOKEN"
fi

# Test 3: Get All Users (with auth)
echo -e "\n3. Get All Users (with auth):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 4: Create a new user (with auth)
echo -e "\n4. Create a new user (with auth):"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "john_doe",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }')

echo $CREATE_RESPONSE | jq .

# Extract user ID from response
USER_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo "Created user ID: $USER_ID"

# Test 5: Get the created user (with auth)
echo -e "\n5. Get the created user (with auth):"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 6: Update the user (with auth)
echo -e "\n6. Update the user (with auth):"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "john_doe_updated",
    "firstName": "John",
    "lastName": "Doe Updated",
    "email": "john.doe.updated@example.com"
  }' | jq .

# Test 7: Get the updated user (with auth)
echo -e "\n7. Get the updated user (with auth):"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 8: Delete the user (with auth)
echo -e "\n8. Delete the user (with auth):"
curl -s -X DELETE $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN"

# Test 9: Try to get the deleted user (should return 404)
echo -e "\n9. Try to get the deleted user (should return 404):"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\nUser CRUD Operations Test Complete!" 
