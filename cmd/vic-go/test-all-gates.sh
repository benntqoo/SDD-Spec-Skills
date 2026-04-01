#!/bin/bash
# Test script to demonstrate all gate checks

echo "=== Testing VIBE-SDD Gate Checks ==="
echo ""

# Test 1: Gate 0 - Requirements
echo "1. Testing Gate 0: Requirements Completeness"
./vic.exe spec gate 0
echo ""

# Test 2: Gate 1 - Architecture
echo "2. Testing Gate 1: Architecture Completeness"
./vic.exe spec gate 1
echo ""

# Test 3: Gate 2 - Code Alignment
echo "3. Testing Gate 2: Code Alignment"
./vic.exe spec gate 2
echo ""

# Test 4: Gate 3 - Test Coverage
echo "4. Testing Gate 3: Test Coverage"
./vic.exe spec gate 3
echo ""

# Test 5: Smart Gate Selection
echo "5. Testing Smart Gate Selection"
./vic.exe gate smart --output json
echo ""

# Test 6: Gate Status
echo "6. Current Gate Status"
./vic.exe gate status
echo ""

# Test 7: Blocking Check
echo "7. Blocking Gate Check"
./vic.exe gate check --blocking
echo ""

echo "=== All tests completed ==="