#!/bin/bash

# API Testing Script

BASE_URL="http://localhost:8080"

echo "🚀 User Management API Test Suite"
echo "=================================="

# Test 1: Health Check
echo "📋 Testing health check..."
curl -s "${BASE_URL}/health" | jq . || echo "❌ Health check failed"
echo ""

# Test 2: Create User
echo "👤 Creating a new user..."
USER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "phone": "1234567890",
    "address": "123 Main St"
  }')

echo "$USER_RESPONSE" | jq . || echo "❌ Create user failed"

# Extract user ID for subsequent tests
USER_ID=$(echo "$USER_RESPONSE" | jq -r '.data.id' 2>/dev/null)
echo "Created user with ID: $USER_ID"
echo ""

# Test 3: Get All Users
echo "📋 Getting all users..."
curl -s "${BASE_URL}/api/v1/users?page=1&page_size=10" | jq . || echo "❌ Get all users failed"
echo ""

# Test 4: Get User by ID
if [ "$USER_ID" != "null" ] && [ "$USER_ID" != "" ]; then
  echo "🔍 Getting user by ID: $USER_ID"
  curl -s "${BASE_URL}/api/v1/users/${USER_ID}" | jq . || echo "❌ Get user by ID failed"
  echo ""

  # Test 5: Update User
  echo "✏️ Updating user..."
  curl -s -X PUT "${BASE_URL}/api/v1/users/${USER_ID}" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "John Updated",
      "email": "john.updated@example.com",
      "age": 31,
      "phone": "0987654321",
      "address": "456 Oak Ave"
    }' | jq . || echo "❌ Update user failed"
  echo ""

  # Test 6: Delete User
  echo "🗑️ Deleting user..."
  curl -s -X DELETE "${BASE_URL}/api/v1/users/${USER_ID}" | jq . || echo "❌ Delete user failed"
  echo ""
fi

echo "✅ API test suite completed!"
