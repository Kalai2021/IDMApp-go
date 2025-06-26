#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJCSEZvT0ZKWnJlU1d0X0l0Z0NubCJ9.eyJpc3MiOiJodHRwczovL2Rldi1kZjRsdWQ0bjZ6ejRpNXRnLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJua01zcFd1Y2ZYZXR6T3BpQlhpNXN1bm5tUlFOUDVRWkBjbGllbnRzIiwiYXVkIjoiL2FwaS92MS8iLCJpYXQiOjE3NTA1ODI0NDMsImV4cCI6MTc1MDY2ODg0MywiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIiwiYXpwIjoibmtNc3BXdWNmWGV0ek9waUJYaTVzdW5ubVJRTlA1UVoifQ.V4o6XBV5R-9zhvgmi9V43_hIHYf_qiFD0YM3A16oPxkK1RJuEtw48WPEQX7fcdikoWVDMuVkhd7GNjO62IKPX3EV7yue1LmhdCSGkg8_ErYuzhnjE-rztn6GO35_t7qNTqzHnM-VrgcBwXggKZdcnVBLRJ-ADkrSIAodGSxHEBZ2Rz4Z5NH6OODGT4dQMs5oiCQvCYcG2C46GhtMpN-yRE9rtgatZBNYyEYhZTNjj-ha_6N2wNpoT-lV3nOCxzPXpoNNMaRT5ZdRuNzXajVObggcoqG2SrrgBIN8pSWfahoOSGW-kqzZxKcEjNAbippH0BTlw_QvCMupoAZobIzojQ"

echo "Testing User Creation with Your Auth0 Token"
echo "=========================================="

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Get All Users (initially empty)
echo -e "\n2. Get All Users (initially empty):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 3: Create a new user
echo -e "\n3. Create a new user:"
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

echo "Create User Response:"
echo $CREATE_RESPONSE | jq .

# Extract user ID from response
USER_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo -e "\nCreated user ID: $USER_ID"

# Test 4: Get the created user
echo -e "\n4. Get the created user:"
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# Test 5: Get All Users (should now have one user)
echo -e "\n5. Get All Users (should now have one user):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\nUser creation test completed!" 