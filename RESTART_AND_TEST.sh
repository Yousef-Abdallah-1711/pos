#!/bin/bash

echo "========================================="
echo "Restaurant-Website Restart & Test Script"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}STEP 1: Kill existing pnpm dev process${NC}"
echo "You MUST do this manually:
  1. Find the terminal running 'pnpm dev'
  2. Press Ctrl+C to stop it
  3. Come back and run this script again"
echo ""

echo -e "${YELLOW}STEP 2: Clear Next.js cache${NC}"
cd "c:\Users\OPT\Desktop\POS\frontend"
rm -rf "apps/restaurant-website/.next"
echo -e "${GREEN}✓ Cache cleared${NC}"
echo ""

echo -e "${YELLOW}STEP 3: Waiting for you to restart pnpm dev...${NC}"
echo "Run this command in a NEW terminal:"
echo -e "${GREEN}cd c:\\Users\\OPT\\Desktop\\POS\\frontend && pnpm dev${NC}"
echo ""
echo "Wait for output like:"
echo "  @pos-saas/restaurant-website:dev: ✓ Starting..."
echo ""

echo -e "${YELLOW}STEP 4: Testing routes...${NC}"
echo "Press Enter when dev server is running and shows ✓ Starting..."
read

# Wait a few seconds for server to be ready
echo "Waiting for server to be ready..."
sleep 3

echo ""
echo "========================================="
echo "TESTING RESTAURANT-WEBSITE ROUTES"
echo "========================================="
echo ""

# Test 1: Root path
echo -e "${YELLOW}Test 1: Root path (http://localhost:3003/)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3003/" | tail -1)
if [ "$RESPONSE" = "200" ]; then
  echo -e "${GREEN}✓ PASS - Returns 200 OK${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE (should be 200)${NC}"
fi
echo ""

# Test 2: English homepage
echo -e "${YELLOW}Test 2: English homepage (http://localhost:3003/en)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3003/en" | tail -1)
if [ "$RESPONSE" = "200" ]; then
  echo -e "${GREEN}✓ PASS - Returns 200 OK${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE (should be 200)${NC}"
fi
echo ""

# Test 3: Arabic homepage
echo -e "${YELLOW}Test 3: Arabic homepage (http://localhost:3003/ar)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3003/ar" | tail -1)
if [ "$RESPONSE" = "200" ]; then
  echo -e "${GREEN}✓ PASS - Returns 200 OK${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE (should be 200)${NC}"
fi
echo ""

# Test 4: Menu page
echo -e "${YELLOW}Test 4: Menu page (http://localhost:3003/en/menu)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3003/en/menu" | tail -1)
if [ "$RESPONSE" = "200" ]; then
  echo -e "${GREEN}✓ PASS - Returns 200 OK${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE (should be 200)${NC}"
fi
echo ""

# Test 5: Maintenance page (should exist but return specific content)
echo -e "${YELLOW}Test 5: Maintenance page (http://localhost:3003/maintenance)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3003/maintenance" | tail -1)
if [ "$RESPONSE" = "200" ]; then
  echo -e "${GREEN}✓ PASS - Returns 200 OK${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE (should be 200)${NC}"
fi
echo ""

# Test 6: Invalid path (should return 404)
echo -e "${YELLOW}Test 6: Invalid path (http://localhost:3003/invalid)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3003/invalid" | tail -1)
if [ "$RESPONSE" = "404" ]; then
  echo -e "${GREEN}✓ PASS - Returns 404 Not Found${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE (should be 404)${NC}"
fi
echo ""

# Test 7: Dashboard still works
echo -e "${YELLOW}Test 7: Dashboard homepage (http://localhost:3002/en/dashboard)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3002/en/dashboard" | tail -1)
if [ "$RESPONSE" = "200" ] || [ "$RESPONSE" = "307" ]; then
  echo -e "${GREEN}✓ PASS - Returns $RESPONSE (dashboard working)${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE${NC}"
fi
echo ""

# Test 8: Orders page (should work now with hook fix)
echo -e "${YELLOW}Test 8: Orders page (http://localhost:3002/en/dashboard/orders)${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" "http://localhost:3002/en/dashboard/orders" | tail -1)
if [ "$RESPONSE" = "200" ] || [ "$RESPONSE" = "307" ]; then
  echo -e "${GREEN}✓ PASS - Returns $RESPONSE (orders page working)${NC}"
else
  echo -e "${RED}✗ FAIL - Returns $RESPONSE${NC}"
fi
echo ""

echo "========================================="
echo "TEST SUMMARY"
echo "========================================="
echo -e "${GREEN}All tests completed!${NC}"
echo ""
echo "If all tests PASS:"
echo "  ✓ Restaurant-website is working"
echo "  ✓ Dashboard is working"
echo "  ✓ No more redirect loops"
echo ""
echo "Next: Open browser and visit:"
echo "  - http://localhost:3002/en/dashboard (Dashboard)"
echo "  - http://localhost:3003/en (Restaurant Website)"
echo ""
