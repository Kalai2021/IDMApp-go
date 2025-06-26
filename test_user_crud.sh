#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"
echo "Testing User CRUD Operations"
echo "============================"

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Get All Users (should be empty initially)
echo -e "\n2. Get All Users (initially empty):"
curl -s -X GET $BASE_URL/users | jq .

# Test 3: Create a new user
echo -e "\n3. Create a new user:"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
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

# Test 4: Get the created user
echo -e "\n4. Get the created user:"
curl -s -X GET $BASE_URL/users/$USER_ID | jq .

# Test 5: Get All Users (should now have one user)
echo -e "\n5. Get All Users (should now have one user):"
curl -s -X GET $BASE_URL/users | jq .

# Test 6: Update the user
echo -e "\n6. Update the user:"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john_doe_updated",
    "firstName": "John",
    "lastName": "Doe Updated",
    "email": "john.doe.updated@example.com"
  }' | jq .

# Test 7: Get the updated user
echo -e "\n7. Get the updated user:"
curl -s -X GET $BASE_URL/users/$USER_ID | jq .

# Test 8: Create another user
echo -e "\n8. Create another user:"
curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "jane_smith",
    "firstName": "Jane",
    "lastName": "Smith",
    "email": "jane.smith@example.com",
    "password": "password456"
  }' | jq .

# Test 9: Get All Users (should now have two users)
echo -e "\n9. Get All Users (should now have two users):"
curl -s -X GET $BASE_URL/users | jq .

# Test 10: Delete the first user
echo -e "\n10. Delete the first user:"
curl -s -X DELETE $BASE_URL/users/$USER_ID

# Test 11: Try to get the deleted user (should return 404)
echo -e "\n11. Try to get the deleted user (should return 404):"
curl -s -X GET $BASE_URL/users/$USER_ID | jq .

# Test 12: Get All Users (should now have one user again)
echo -e "\n12. Get All Users (should now have one user again):"
curl -s -X GET $BASE_URL/users | jq .

echo -e "\nUser CRUD Operations Test Complete!" 