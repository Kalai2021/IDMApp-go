#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"

echo "Creating User with Auth0 OAuth Token"
echo "===================================="

# Check if we have a token
if [ ! -f ".auth0_token" ]; then
    echo "No Auth0 token found. Getting one first..."
    ./get_auth0_token.sh
fi

# Read the token from file
ACCESS_TOKEN=$(cat .auth0_token)

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
    echo "Error: No valid token found. Please check your Auth0 configuration."
    exit 1
fi

echo "Using Auth0 token: $ACCESS_TOKEN"

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Get All Users (should be empty initially)
echo -e "\n2. Get All Users (initially empty):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 3: Create a new user
echo -e "\n3. Create a new user:"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "john_doe",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }')

echo "Create User Response:"
echo $CREATE_RESPONSE | jq .

# Extract user ID from response
USER_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo -e "\nCreated user ID: $USER_ID"

# Test 4: Get the created user
echo -e "\n4. Get the created user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 5: Get All Users (should now have one user)
echo -e "\n5. Get All Users (should now have one user):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 6: Update the user
echo -e "\n6. Update the user:"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "john_doe_updated",
    "firstName": "John",
    "lastName": "Doe Updated",
    "email": "john.doe.updated@example.com"
  }' | jq .

# Test 7: Get the updated user
echo -e "\n7. Get the updated user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

echo -e "\nUser creation and CRUD operations completed successfully!"
echo "User ID: $USER_ID" 