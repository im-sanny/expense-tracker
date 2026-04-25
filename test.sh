#!/bin/bash
# test.sh - Full integration test suite for expense-tracker

set -e  # Exit on first error

BASE_URL="http://localhost:3000"
EMAIL="test-$(date +%s)@example.com"  # Unique email per test run
PASSWORD="SecurePass123!"
COOKIES="cookies.txt"

echo "🧪 Starting integration tests..."
echo "Test user: $EMAIL"

# ─── AUTH FLOW ─────────────────────────────────────────────
echo -e "\n🔐 Testing Auth Flow..."

# 1. Register
echo "→ Registering user..."
REGISTER_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
REGISTER_CODE=$(echo "$REGISTER_RESP" | tail -n1)
if [ "$REGISTER_CODE" != "201" ]; then
  echo "❌ Register failed: $REGISTER_RESP"
  exit 1
fi
echo "✅ Register: $REGISTER_CODE"

# 2. Login
echo "→ Logging in..."
curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" \
  -c "$COOKIES" > /dev/null
if [ ! -f "$COOKIES" ] || ! grep -q "access_token" "$COOKIES"; then
  echo "❌ Login failed: cookies not set"
  exit 1
fi
echo "✅ Login: cookies saved"

# 3. Test protected route
echo "→ Testing protected /track endpoint..."
TRACK_RESP=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/track" -b "$COOKIES")
TRACK_CODE=$(echo "$TRACK_RESP" | tail -n1)
if [ "$TRACK_CODE" != "200" ]; then
  echo "❌ Protected route failed: $TRACK_RESP"
  exit 1
fi
echo "✅ Protected route: $TRACK_CODE"

# ─── CRUD OPERATIONS ───────────────────────────────────────
echo -e "\n📝 Testing CRUD Operations..."

# 4. Create expense
echo "→ Creating expense..."
CREATE_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/track" \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-04-23T00:00:00Z","amount":890,"note":"Coffee + croissant"}' \
  -b "$COOKIES")
CREATE_CODE=$(echo "$CREATE_RESP" | tail -n1)
if [ "$CREATE_CODE" != "201" ]; then  # ← Expect 201, not 200
  echo "❌ Create failed: $CREATE_RESP"
  exit 1
fi
EXPENSE_ID=$(echo "$CREATE_RESP" | head -n1 | jq -r '.id // empty')
echo "✅ Create: $CREATE_CODE (id: $EXPENSE_ID)"

# 5. Get all expenses
echo "→ Getting all expenses..."
curl -s -X GET "$BASE_URL/track" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Get all: OK"

# 6. Get by ID
if [ -n "$EXPENSE_ID" ]; then
  echo "→ Getting expense by ID..."
  curl -s -X GET "$BASE_URL/track/$EXPENSE_ID" -b "$COOKIES" | jq -c '.' > /dev/null
  echo "✅ Get by ID: OK"
fi

# 7. Update (PUT)
if [ -n "$EXPENSE_ID" ]; then
  echo "→ Updating expense (PUT)..."
  curl -s -X PUT "$BASE_URL/track/$EXPENSE_ID" \
    -H "Content-Type: application/json" \
    -d '{"date":"2026-04-23T00:00:00Z","amount":1200,"note":"Updated: Coffee + pastry"}' \
    -b "$COOKIES" | jq -c '.' > /dev/null
  echo "✅ Update (PUT): OK"
fi

# 8. Patch (partial update)
if [ -n "$EXPENSE_ID" ]; then
  echo "→ Patching expense..."
  curl -s -X PATCH "$BASE_URL/track/$EXPENSE_ID" \
    -H "Content-Type: application/json" \
    -d '{"note":"Patched note"}' \
    -b "$COOKIES" | jq -c '.' > /dev/null
  echo "✅ Patch: OK"
fi

# ─── QUERY PARAMETERS ──────────────────────────────────────
echo -e "\n🔍 Testing Query Parameters..."

# 9. Pagination
echo "→ Testing pagination..."
curl -s -X GET "$BASE_URL/track?page=1&limit=10" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Pagination: OK"

# 10. Date range
echo "→ Testing date range..."
curl -s -X GET "$BASE_URL/track?from=2026-01-01&to=2026-12-31" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Date range: OK"

# 11. Amount filter
echo "→ Testing amount filter..."
curl -s -X GET "$BASE_URL/track?min=100&max=500" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Amount filter: OK"

# 12. Search
echo "→ Testing search..."
curl -s -X GET "$BASE_URL/track?q=coffee" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Search: OK"

# 13. Combined filters
echo "→ Testing combined filters..."
curl -s -X GET "$BASE_URL/track?q=coffee&min=100&from=2026-01-01&page=1&limit=10" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Combined filters: OK"

# ─── TOKEN MANAGEMENT ──────────────────────────────────────
echo -e "\n🔄 Testing Token Management..."

# 14. Refresh token
echo "→ Refreshing access token..."
curl -s -X POST "$BASE_URL/auth/refresh" -b "$COOKIES" | jq -c '.' > /dev/null
echo "✅ Refresh: OK"

# 15. Logout
echo "→ Logging out..."
LOGOUT_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/logout" -b "$COOKIES")
LOGOUT_CODE=$(echo "$LOGOUT_RESP" | tail -n1)
if [ "$LOGOUT_CODE" != "200" ]; then
  echo "❌ Logout failed: $LOGOUT_RESP"
  exit 1
fi
echo "✅ Logout: $LOGOUT_CODE"

# ✅ Clear cookies file after logout (simulates client deleting cookies)
> "$COOKIES"  # Truncate the file

# 16. Verify logout (should fail - no cookies sent)
echo "→ Verifying logout (should be unauthorized)..."
UNAUTH_RESP=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/track" -b "$COOKIES")
UNAUTH_CODE=$(echo "$UNAUTH_RESP" | tail -n1)
if [ "$UNAUTH_CODE" != "401" ]; then
  echo "❌ Logout verification failed: expected 401, got $UNAUTH_CODE"
  exit 1
fi
echo "✅ Logout verified: $UNAUTH_CODE"

# ─── NEGATIVE TESTS ────────────────────────────────────────
echo -e "\n🚫 Testing Negative Cases..."

# 17. Invalid login
echo "→ Testing invalid credentials..."
INVALID_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"wrongpass\"}")
INVALID_CODE=$(echo "$INVALID_RESP" | tail -n1)
if [ "$INVALID_CODE" != "401" ]; then
  echo "❌ Invalid login should return 401, got $INVALID_CODE"
  exit 1
fi
echo "✅ Invalid credentials: $INVALID_CODE"

# 18. Missing auth on protected route
echo "→ Testing unprotected access..."
NOAUTH_RESP=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/track")
NOAUTH_CODE=$(echo "$NOAUTH_RESP" | tail -n1)
if [ "$NOAUTH_CODE" != "401" ]; then
  echo "❌ Protected route should return 401 without auth, got $NOAUTH_CODE"
  exit 1
fi
echo "✅ Protected route without auth: $NOAUTH_CODE"

# ─── CLEANUP ───────────────────────────────────────────────
echo -e "\n🧹 Cleanup..."
rm -f "$COOKIES"
echo "✅ Test cookies removed"

echo -e "\n🎉 All tests passed!"
