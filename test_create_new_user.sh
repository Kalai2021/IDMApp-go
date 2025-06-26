#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJCSEZvT0ZKWnJlU1d0X0l0Z0NubCJ9.eyJpc3MiOiJodHRwczovL2Rldi1kZjRsdWQ0bjZ6ejRpNXRnLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJua01zcFd1Y2ZYZXR6T3BpQlhpNXN1bm5tUlFOUDVRWkBjbGllbnRzIiwiYXVkIjoiL2FwaS92MS8iLCJpYXQiOjE3NTA1ODI0NDMsImV4cCI6MTc1MDY2ODg0MywiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIiwiYXpwIjoibmtNc3BXdWNmWGV0ek9waUJYaTVzdW5ubVJRTlA1UVoifQ.V4o6XBV5R-9zhvgmi9V43_hIHYf_qiFD0YM3A16oPxkK1RJuEtw48WPEQX7fcdikoWVDMuVkhd7GNjO62IKPX3EV7yue1LmhdCSGkg8_ErYuzhnjE-rztn6GO35_t7qNTqzHnM-VrgcBwXggKZdcnVBLRJ-ADkrSIAodGSxHEBZ2Rz4Z5NH6OODGT4dQMs5oiCQvCYcG2C46GhtMpN-yRE9rtgatZBNYyEYhZTNjj-ha_6N2wNpoT-lV3nOCxzPXpoNNMaRT5ZdRuNzXajVObggcoqG2SrrgBIN8pSWfahoOSGW-kqzZxKcEjNAbippH0BTlw_QvCMupoAZobIzojQ"

echo "Testing Complete User CRUD Operations with Auth0 Token"
echo "====================================================="

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Create a new user with unique email
echo -e "\n2. Create a new user:"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "jane_smith",
    "firstName": "Jane",
    "lastName": "Smith",
    "email": "jane.smith@example.com",
    "password": "password123"
  }')

echo "Create User Response:"
echo $CREATE_RESPONSE | jq .

# Extract user ID from response
USER_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo -e "\nCreated user ID: $USER_ID"

# Test 3: Get the created user
echo -e "\n3. Get the created user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 4: Update the user
echo -e "\n4. Update the user:"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "jane_smith_updated",
    "firstName": "Jane Updated",
    "lastName": "Smith Updated",
    "email": "jane.smith.updated@example.com"
  }' | jq .

# Test 5: Get the updated user
echo -e "\n5. Get the updated user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 6: Delete the user
echo -e "\n6. Delete the user:"
curl -s -X DELETE $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN"

# Test 7: Try to get the deleted user (should return 404)
echo -e "\n7. Try to get the deleted user (should return 404):"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\nComplete User CRUD Operations Test Completed Successfully! 🎉" 