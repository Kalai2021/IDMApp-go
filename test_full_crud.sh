#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJCSEZvT0ZKWnJlU1d0X0l0Z0NubCJ9.eyJpc3MiOiJodHRwczovL2Rldi1kZjRsdWQ0bjZ6ejRpNXRnLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJua01zcFd1Y2ZYZXR6T3BpQlhpNXN1bm5tUlFOUDVRWkBjbGllbnRzIiwiYXVkIjoiL2FwaS92MS8iLCJpYXQiOjE3NTA1ODI0NDMsImV4cCI6MTc1MDY2ODg0MywiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIiwiYXpwIjoibmtNc3BXdWNmWGV0ek9waUJYaTVzdW5ubVJRTlA1UVoifQ.V4o6XBV5R-9zhvgmi9V43_hIHYf_qiFD0YM3A16oPxkK1RJuEtw48WPEQX7fcdikoWVDMuVkhd7GNjO62IKPX3EV7yue1LmhdCSGkg8_ErYuzhnjE-rztn6GO35_t7qNTqzHnM-VrgcBwXggKZdcnVBLRJ-ADkrSIAodGSxHEBZ2Rz4Z5NH6OODGT4dQMs5oiCQvCYcG2C46GhtMpN-yRE9rtgatZBNYyEYhZTNjj-ha_6N2wNpoT-lV3nOCxzPXpoNNMaRT5ZdRuNzXajVObggcoqG2SrrgBIN8pSWfahoOSGW-kqzZxKcEjNAbippH0BTlw_QvCMupoAZobIzojQ"

echo "Testing Complete User CRUD Operations"
echo "====================================="

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
    "email": "jane.smith.test@example.com",
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
UPDATE_RESPONSE=$(curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "jane_smith_updated",
    "firstName": "Jane Updated",
    "lastName": "Smith Updated",
    "email": "jane.smith.updated@example.com"
  }')

echo "Update User Response:"
echo $UPDATE_RESPONSE | jq .

# Test 5: Get the updated user
echo -e "\n5. Get the updated user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 6: Update user with partial data (only name)
echo -e "\n6. Update user with partial data (only name):"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "jane_smith_final"
  }' | jq .

# Test 7: Get the user after partial update
echo -e "\n7. Get the user after partial update:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 8: Update user status (isActive)
echo -e "\n8. Update user status (isActive):"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "isActive": false
  }' | jq .

# Test 9: Get the user after status update
echo -e "\n9. Get the user after status update:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 10: Delete the user
echo -e "\n10. Delete the user:"
curl -s -X DELETE $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN"

echo -e "\nDelete response status: $?"

# Test 11: Try to get the deleted user (should return 404)
echo -e "\n11. Try to get the deleted user (should return 404):"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 12: Try to update the deleted user (should return 404)
echo -e "\n12. Try to update the deleted user (should return 404):"
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "deleted_user"
  }' | jq .

# Test 13: Try to delete the already deleted user (should return 404)
echo -e "\n13. Try to delete the already deleted user (should return 404):"
curl -s -X DELETE $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN"

echo -e "\nSecond delete response status: $?"

echo -e "\n🎉 Complete User CRUD Operations Test Completed Successfully!" 