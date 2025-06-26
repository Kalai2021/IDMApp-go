#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"
echo "Testing User CRUD Operations - Simple Approach"
echo "=============================================="

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Try to login with non-existent user
echo -e "\n2. Try to login with non-existent user:"
curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq .

# Test 3: Let's test the endpoints directly with a mock token
# First, let's create a proper JWT token for testing
echo -e "\n3. Testing with a mock JWT token:"

# Create a simple JWT token (this is just for testing)
# In production, you would get this from Auth0
MOCK_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXItaWQiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJpYXQiOjE2MzQ1Njc4NzgsImV4cCI6MTYzNDU3MTQ3OH0.signature"

echo "Using mock token: $MOCK_TOKEN"

# Test 4: Get All Users
echo -e "\n4. Get All Users:"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $MOCK_TOKEN" | jq .

# Test 5: Create a new user
echo -e "\n5. Create a new user:"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MOCK_TOKEN" \
  -d '{
    "name": "test_user",
    "firstName": "Test",
    "lastName": "User",
    "email": "test@example.com",
    "password": "password123"
  }')

echo $CREATE_RESPONSE | jq .

# Extract user ID from response
USER_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo "Created user ID: $USER_ID"

# Test 6: Get the created user
echo -e "\n6. Get the created user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $MOCK_TOKEN" | jq .

# Test 7: Update the user
echo -e "\n7. Update the user:"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MOCK_TOKEN" \
  -d '{
    "name": "test_user_updated",
    "firstName": "Test Updated",
    "lastName": "User Updated",
    "email": "test.updated@example.com"
  }' | jq .

# Test 8: Get the updated user
echo -e "\n8. Get the updated user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $MOCK_TOKEN" | jq .

# Test 9: Delete the user
echo -e "\n9. Delete the user:"
curl -s -X DELETE $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $MOCK_TOKEN"

# Test 10: Try to get the deleted user
echo -e "\n10. Try to get the deleted user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $MOCK_TOKEN" | jq .

echo -e "\nUser CRUD Operations Test Complete!" 