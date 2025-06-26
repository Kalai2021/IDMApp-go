#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"

echo "Testing User CRUD Operations (No Auth - Development)"
echo "==================================================="

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s -X GET http://localhost:8090/health | jq .

# Test 2: Try to create a user without auth (should fail with auth error)
echo -e "\n2. Try to create user without auth (should fail with auth error):"
curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john_doe",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }' | jq .

# Test 3: Try with a simple test token
echo -e "\n3. Try with a simple test token:"
curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "name": "john_doe",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }' | jq .

# Test 4: Get all users with test token
echo -e "\n4. Get all users with test token:"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer test-token" | jq .

echo -e "\nTest completed!" 