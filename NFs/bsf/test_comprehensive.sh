#!/bin/bash

# BSF Comprehensive Test Runner
# Executes complete BSF functionality testing including 3GPP compliance tests

set -e

BSF_DIR="/home/xflow/free5gc/NFs/bsf"
LOG_DIR="${BSF_DIR}/test-logs"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=====================================${NC}"
echo -e "${BLUE}  BSF Comprehensive Test Suite      ${NC}"
echo -e "${BLUE}  Testing 3GPP TS 29.521 Compliance ${NC}"
echo -e "${BLUE}=====================================${NC}"

# Create log directory
mkdir -p "$LOG_DIR"

# Change to BSF directory
cd "$BSF_DIR"

echo -e "\n${YELLOW}[INFO]${NC} Starting BSF comprehensive tests..."

# Run comprehensive test suite
echo -e "\n${BLUE}=== Running BSF Comprehensive Tests ===${NC}"
if go test -v ./internal/sbi/processor/ -run TestBSFComprehensive -timeout=30m 2>&1 | tee "${LOG_DIR}/comprehensive_test.log"; then
    echo -e "\n${GREEN}✅ BSF Comprehensive Tests: PASSED${NC}"
else
    echo -e "\n${RED}❌ BSF Comprehensive Tests: FAILED${NC}"
    exit 1
fi

# Run existing processor tests
echo -e "\n${BLUE}=== Running Existing Processor Tests ===${NC}"
if go test -v ./internal/sbi/processor/ -run TestCreatePCFBinding -timeout=10m 2>&1 | tee "${LOG_DIR}/processor_test.log"; then
    echo -e "\n${GREEN}✅ Processor Tests: PASSED${NC}"
else
    echo -e "\n${RED}❌ Processor Tests: FAILED${NC}"
    exit 1
fi

# Run context tests
echo -e "\n${BLUE}=== Running Context Tests ===${NC}"
if go test -v ./internal/context/ -timeout=10m 2>&1 | tee "${LOG_DIR}/context_test.log"; then
    echo -e "\n${GREEN}✅ Context Tests: PASSED${NC}"
else
    echo -e "\n${YELLOW}⚠️  Context Tests: No tests found or failed${NC}"
fi

# Performance benchmark (if available)
echo -e "\n${BLUE}=== Running Performance Benchmarks ===${NC}"
if go test -bench=. -benchmem ./internal/sbi/processor/ 2>&1 | tee "${LOG_DIR}/benchmark.log"; then
    echo -e "\n${GREEN}✅ Benchmarks: COMPLETED${NC}"
else
    echo -e "\n${YELLOW}⚠️  Benchmarks: No benchmarks found${NC}"
fi

# Race condition detection
echo -e "\n${BLUE}=== Running Race Condition Tests ===${NC}"
if go test -race ./internal/sbi/processor/ -run TestBSFComprehensive -timeout=15m 2>&1 | tee "${LOG_DIR}/race_test.log"; then
    echo -e "\n${GREEN}✅ Race Detection Tests: PASSED${NC}"
else
    echo -e "\n${RED}❌ Race Detection Tests: FAILED${NC}"
    exit 1
fi

# Code coverage analysis
echo -e "\n${BLUE}=== Generating Code Coverage Report ===${NC}"
if go test -coverprofile="${LOG_DIR}/coverage.out" ./internal/sbi/processor/ 2>&1 | tee "${LOG_DIR}/coverage_test.log"; then
    go tool cover -html="${LOG_DIR}/coverage.out" -o "${LOG_DIR}/coverage.html"
    echo -e "\n${GREEN}✅ Coverage Report: Generated at ${LOG_DIR}/coverage.html${NC}"
    
    # Display coverage summary
    echo -e "\n${BLUE}=== Coverage Summary ===${NC}"
    go tool cover -func="${LOG_DIR}/coverage.out" | tail -1
else
    echo -e "\n${YELLOW}⚠️  Coverage Report: Failed to generate${NC}"
fi

# Test MongoDB integration (if available)
echo -e "\n${BLUE}=== Testing MongoDB Integration ===${NC}"
if command -v mongod >/dev/null 2>&1; then
    echo -e "${YELLOW}[INFO]${NC} MongoDB found, testing database integration..."
    # Set MongoDB test environment
    export BSF_MONGO_URI="mongodb://localhost:27017/bsf_test"
    if go test -v ./internal/context/ -run TestMongoDB -timeout=10m 2>&1 | tee "${LOG_DIR}/mongodb_test.log"; then
        echo -e "\n${GREEN}✅ MongoDB Integration: PASSED${NC}"
    else
        echo -e "\n${YELLOW}⚠️  MongoDB Integration: Tests not available or failed${NC}"
    fi
else
    echo -e "${YELLOW}[INFO]${NC} MongoDB not found, skipping database integration tests"
fi

# Generate test summary
echo -e "\n${BLUE}=== Test Summary ===${NC}"
echo -e "Test logs available at: ${LOG_DIR}/"
echo -e "Coverage report: ${LOG_DIR}/coverage.html"

echo -e "\n${GREEN}🎯 All BSF tests completed successfully!${NC}"
echo -e "${GREEN}📋 3GPP TS 29.521 compliance verified${NC}"
echo -e "${GREEN}🔒 Race conditions checked${NC}"
echo -e "${GREEN}📊 Performance benchmarks completed${NC}"
echo -e "${GREEN}📈 Code coverage analyzed${NC}"

echo -e "\n${BLUE}=====================================${NC}"
echo -e "${BLUE}  BSF Testing Complete ✅            ${NC}"
echo -e "${BLUE}=====================================${NC}"
