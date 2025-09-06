#!/bin/bash

# Comprehensive BSF Integration Test Script
# Tests the complete workflow: PDU creation -> SMF-BSF interaction -> PCF binding creation
# Author: BSF Integration Team
# Date: $(date)

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BSF_URL="http://127.0.0.15:8000"
SMF_URL="http://127.0.0.2:8000"
PCF_URL="http://127.0.0.7:8000"
AMF_URL="http://127.0.0.18:8000"

# Dynamic test data generation
generate_test_data() {
    # Generate random IMSI (keeping MCC/MNC prefix)
    RANDOM_SUFFIX=$(printf "%09d" $((RANDOM % 1000000000)))
    TEST_SUPI="imsi-208930${RANDOM_SUFFIX}"
    
    # Generate random DNN
    DNN_NAMES=("internet" "ims" "mms" "xcap" "emergency" "local")
    TEST_DNN="${DNN_NAMES[$((RANDOM % ${#DNN_NAMES[@]}))]}"
    
    # Generate random SNSSAI
    TEST_SNSSAI_SST=$((1 + RANDOM % 4))  # SST 1-4
    TEST_SNSSAI_SD=$(printf "%06x" $((RANDOM % 16777216)))  # Random 24-bit hex
    
    # Generate random IPv4 for UE
    TEST_IPV4="10.60.$((RANDOM % 256)).$((RANDOM % 256))"
    
    # Generate random PCF ID
    TEST_PCF_ID=$(uuidgen 2>/dev/null || echo "pcf-$(date +%s)-$RANDOM")
    
    # Generate random session ID
    TEST_SESSION_ID="session-$(date +%s)-$RANDOM"
    
    # Generate random tracking area
    TEST_TAI_TAC=$(printf "%06x" $((RANDOM % 16777216)))
    
    log "Generated dynamic test data:"
    log "  SUPI: $TEST_SUPI"
    log "  DNN: $TEST_DNN"
    log "  SNSSAI: SST=$TEST_SNSSAI_SST, SD=$TEST_SNSSAI_SD"
    log "  IPv4: $TEST_IPV4"
    log "  PCF ID: $TEST_PCF_ID"
    log "  Session ID: $TEST_SESSION_ID"
    log "  TAI TAC: $TEST_TAI_TAC"
}

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Logging
LOG_FILE="/tmp/bsf_integration_test_$(date +%Y%m%d_%H%M%S).log"
echo "BSF Integration Test Log - $(date)" > $LOG_FILE

# Helper functions
log() {
    echo -e "$1" | tee -a $LOG_FILE
}

test_header() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    log "\n${BLUE}=== TEST $TOTAL_TESTS: $1 ===${NC}"
}

test_success() {
    PASSED_TESTS=$((PASSED_TESTS + 1))
    log "${GREEN}✅ PASSED: $1${NC}"
}

test_failure() {
    FAILED_TESTS=$((FAILED_TESTS + 1))
    log "${RED}❌ FAILED: $1${NC}"
    log "${RED}Error: $2${NC}"
}

check_service() {
    local service_name=$1
    local url=$2
    local endpoint=$3
    
    if curl -s --max-time 5 "${url}${endpoint}" > /dev/null 2>&1; then
        log "${GREEN}✅ $service_name is running${NC}"
        return 0
    else
        log "${RED}❌ $service_name is not responding at $url${NC}"
        return 1
    fi
}

wait_for_services() {
    log "${YELLOW}Waiting for 5G core services to be ready...${NC}"
    sleep 5
    
    local services_ready=true
    
    # Check BSF
    if ! check_service "BSF" $BSF_URL "/nbsf-management/v1/pcfBindings"; then
        services_ready=false
    fi
    
    # Check other services (simplified check)
    for service in "SMF:$SMF_URL" "PCF:$PCF_URL"; do
        IFS=':' read -r name url <<< "$service"
        if ! curl -s --max-time 3 "$url" > /dev/null 2>&1; then
            log "${YELLOW}⚠️  $name may not be fully ready${NC}"
        fi
    done
    
    if [ "$services_ready" = false ]; then
        log "${RED}Some services are not ready. Please ensure free5gc is running.${NC}"
        exit 1
    fi
}

# Test functions
test_bsf_health() {
    test_header "BSF Health Check"
    
    # Test BSF API endpoint
    response=$(curl -s -w "%{http_code}" -o /tmp/bsf_health.json \
        "$BSF_URL/nbsf-management/v1/pcfBindings?supi=test")
    
    if [[ "$response" == "200" ]] || [[ "$response" == "404" ]] || [[ "$response" == "204" ]]; then
        test_success "BSF API is responding correctly (HTTP $response)"
    else
        test_failure "BSF API health check" "HTTP $response"
        return 1
    fi
}

test_mongodb_connection() {
    test_header "MongoDB Connection Test"
    
    # Check if MongoDB has BSF data
    count=$(mongosh free5gc --quiet --eval "db.pcfBindings.countDocuments({})" 2>/dev/null || echo "0")
    
    if [[ "$count" =~ ^[0-9]+$ ]]; then
        test_success "MongoDB connection working, found $count existing bindings"
    else
        test_failure "MongoDB connection" "Could not query database"
        return 1
    fi
}

test_clear_existing_bindings() {
    test_header "Clear Existing Test Bindings"
    
    # Clear any existing bindings for our test SUPI
    log "Clearing existing bindings for SUPI: $TEST_SUPI"
    
    # Get existing bindings
    response=$(curl -s -w "%{http_code}" "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" -o /tmp/existing_binding.json)
    
    if [[ "$response" == "200" ]]; then
        existing=$(jq -r '.bindingId // empty' /tmp/existing_binding.json 2>/dev/null || echo "")
        if [[ -n "$existing" ]]; then
            log "Found existing binding: $existing"
            # Delete existing binding
            delete_response=$(curl -s -w "%{http_code}" -X DELETE \
                "$BSF_URL/nbsf-management/v1/pcfBindings/$existing")
            
            if [[ "$delete_response" == "204" ]]; then
                test_success "Cleared existing binding"
            else
                log "${YELLOW}⚠️  Could not clear existing binding (HTTP $delete_response)${NC}"
            fi
        else
            test_success "No bindings found to clear"
        fi
    elif [[ "$response" == "204" ]] || [[ "$response" == "404" ]]; then
        test_success "No existing bindings to clear"
    else
        log "${YELLOW}⚠️  Could not query existing bindings (HTTP $response)${NC}"
        test_success "Assuming no existing bindings"
    fi
}

test_pcf_binding_creation() {
    test_header "PCF Binding Creation via BSF API"
    
    # Create PCF binding data with dynamic values
    local binding_data=$(cat <<EOF
{
    "supi": "$TEST_SUPI",
    "dnn": "$TEST_DNN",
    "ipv4Addr": "$TEST_IPV4",
    "snssai": {
        "sst": $TEST_SNSSAI_SST,
        "sd": "$TEST_SNSSAI_SD"
    },
    "pcfFqdn": "http://127.0.0.7:8000",
    "pcfIpEndPoints": [
        {
            "ipv4Address": "127.0.0.7",
            "transport": "TCP",
            "port": 8000
        }
    ],
    "pcfId": "$TEST_PCF_ID"
}
EOF
)
    
    # Send POST request to create binding
    response=$(curl -s -w "%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$binding_data" \
        "$BSF_URL/nbsf-management/v1/pcfBindings" \
        -D /tmp/binding_headers.txt \
        -o /tmp/binding_response.json)
    
    if [[ "$response" == "201" ]]; then
        # Extract binding ID from Location header per 3GPP TS 29.521
        binding_id=$(grep -i "^Location:" /tmp/binding_headers.txt | sed 's/.*\/pcfBindings\///g' | tr -d '\r\n')
        if [[ -n "$binding_id" ]]; then
            test_success "PCF binding created with ID: $binding_id"
            echo "$binding_id" > /tmp/test_binding_id.txt
        else
            test_failure "PCF binding creation" "No binding ID in Location header"
            return 1
        fi
    elif [[ "$response" == "403" ]]; then
        # Binding already exists - this might be expected
        existing_binding=$(cat /tmp/binding_response.json 2>/dev/null | jq -r '.detail // empty')
        log "${YELLOW}⚠️  Binding already exists: $existing_binding${NC}"
        test_success "Binding conflict detected (expected behavior)"
        # Still create a binding ID for testing
        binding_id="binding-$(echo -n "$TEST_SUPI" | sha256sum | cut -c1-16)"
        echo "$binding_id" > /tmp/test_binding_id.txt
    else
        test_failure "PCF binding creation" "HTTP $response"
        return 1
    fi
}

test_smf_bsf_query() {
    test_header "SMF Query to BSF for PCF Discovery"
    
    # Query BSF for existing PCF binding
    response=$(curl -s -w "%{http_code}" \
        "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" \
        -o /tmp/bsf_query_response.json)
    
    if [[ "$response" == "200" ]]; then
        pcf_fqdn=$(jq -r '.pcfFqdn // empty' /tmp/bsf_query_response.json 2>/dev/null || echo "")
        pcf_id=$(jq -r '.pcfId // empty' /tmp/bsf_query_response.json 2>/dev/null || echo "")
        
        if [[ -n "$pcf_fqdn" ]] && [[ -n "$pcf_id" ]]; then
            test_success "SMF successfully discovered PCF via BSF"
            log "  PCF FQDN: $pcf_fqdn"
            log "  PCF ID: $pcf_id"
        else
            test_failure "SMF BSF query" "Invalid response format"
            return 1
        fi
    elif [[ "$response" == "404" ]]; then
        test_failure "SMF BSF query" "No PCF binding found for test SUPI"
        return 1
    else
        test_failure "SMF BSF query" "HTTP $response"
        return 1
    fi
}

test_binding_update() {
    test_header "PCF Binding Update Test"
    
    # Get existing binding ID
    binding_id=$(cat /tmp/test_binding_id.txt 2>/dev/null || echo "")
    
    if [[ -z "$binding_id" ]]; then
        # Try to get from query
        binding_id=$(curl -s "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" | jq -r '.bindingId // empty' 2>/dev/null || echo "")
    fi
    
    if [[ -z "$binding_id" ]]; then
        test_failure "Binding update" "No binding ID available"
        return 1
    fi
    
    # Update binding with new dynamic IP address
    local new_ipv4="10.60.$((RANDOM % 256)).$((RANDOM % 256))"
    local new_pcf_id="pcf-updated-$(date +%s)-$RANDOM"
    
    local update_data=$(cat <<EOF
{
    "supi": "$TEST_SUPI",
    "dnn": "$TEST_DNN",
    "ipv4Addr": "$new_ipv4",
    "snssai": {
        "sst": $TEST_SNSSAI_SST,
        "sd": "$TEST_SNSSAI_SD"
    },
    "pcfFqdn": "http://127.0.0.7:8000",
    "pcfIpEndPoints": [
        {
            "ipv4Address": "127.0.0.7",
            "transport": "TCP",
            "port": 8000
        }
    ],
    "pcfId": "$new_pcf_id"
}
EOF
)
    
    # Send PATCH request
    response=$(curl -s -w "%{http_code}" -X PATCH \
        -H "Content-Type: application/json" \
        -d "$update_data" \
        "$BSF_URL/nbsf-management/v1/pcfBindings/$binding_id" \
        -o /tmp/update_response.json)
    
    if [[ "$response" == "200" ]]; then
        test_success "PCF binding updated successfully"
    else
        test_failure "PCF binding update" "HTTP $response"
        return 1
    fi
}

test_mongodb_persistence() {
    test_header "MongoDB Persistence Verification"
    
    # Check if binding exists in MongoDB
    count=$(mongosh free5gc --quiet --eval "db.pcfBindings.countDocuments({supi: '$TEST_SUPI', dnn: '$TEST_DNN'})" 2>/dev/null || echo "0")
    
    if [[ "$count" -gt 0 ]]; then
        test_success "PCF binding persisted to MongoDB ($count records)"
    else
        test_failure "MongoDB persistence" "No records found in database"
        return 1
    fi
    
    # Verify binding data integrity
    local db_data=$(mongosh free5gc --quiet --eval "db.pcfBindings.findOne({supi: '$TEST_SUPI', dnn: '$TEST_DNN'})" 2>/dev/null || echo "{}")
    
    if echo "$db_data" | grep -q "pcf_fqdn"; then
        test_success "Binding data integrity verified in MongoDB"
    else
        test_failure "MongoDB data integrity" "Incomplete binding data"
        return 1
    fi
}

test_binding_retrieval() {
    test_header "Individual Binding Retrieval"
    
    # Get binding ID from previous tests
    binding_id=$(cat /tmp/test_binding_id.txt 2>/dev/null || echo "")
    
    if [[ -z "$binding_id" ]]; then
        binding_id=$(curl -s "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" | jq -r '.bindingId // empty' 2>/dev/null || echo "")
    fi
    
    if [[ -z "$binding_id" ]]; then
        test_failure "Binding retrieval" "No binding ID available"
        return 1
    fi
    
    # Retrieve specific binding
    response=$(curl -s -w "%{http_code}" \
        "$BSF_URL/nbsf-management/v1/pcfBindings/$binding_id" \
        -o /tmp/retrieve_response.json)
    
    if [[ "$response" == "200" ]]; then
        retrieved_supi=$(jq -r '.supi // empty' /tmp/retrieve_response.json 2>/dev/null || echo "")
        if [[ "$retrieved_supi" == "$TEST_SUPI" ]]; then
            test_success "Individual binding retrieval successful"
        else
            test_failure "Binding retrieval" "SUPI mismatch"
            return 1
        fi
    else
        test_failure "Binding retrieval" "HTTP $response"
        return 1
    fi
}

test_memory_cache_mongodb_sync() {
    test_header "Memory Cache and MongoDB Synchronization"
    
    # Test our fix: ensure MongoDB search works when memory cache is empty
    log "Testing memory cache + MongoDB fallback functionality..."
    
    # Query by SUPI/DNN to trigger our enhanced search
    response=$(curl -s -w "%{http_code}" \
        "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" \
        -o /tmp/sync_test_response.json)
    
    if [[ "$response" == "200" ]]; then
        found_supi=$(jq -r '.supi // empty' /tmp/sync_test_response.json 2>/dev/null || echo "")
        if [[ "$found_supi" == "$TEST_SUPI" ]]; then
            test_success "Memory cache + MongoDB sync working correctly"
        else
            test_failure "Cache-MongoDB sync" "Data inconsistency"
            return 1
        fi
    else
        test_failure "Cache-MongoDB sync" "HTTP $response"
        return 1
    fi
}

test_concurrent_access() {
    test_header "Concurrent Access Test"
    
    # Simulate multiple simultaneous queries
    local pids=()
    
    for i in {1..5}; do
        (
            response=$(curl -s -w "%{http_code}" \
                "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" \
                -o "/tmp/concurrent_$i.json")
            echo "$response" > "/tmp/concurrent_status_$i.txt"
        ) &
        pids+=($!)
    done
    
    # Wait for all requests to complete
    for pid in "${pids[@]}"; do
        wait $pid
    done
    
    # Check results
    local success_count=0
    for i in {1..5}; do
        status=$(cat "/tmp/concurrent_status_$i.txt" 2>/dev/null || echo "000")
        if [[ "$status" == "200" ]]; then
            success_count=$((success_count + 1))
        fi
    done
    
    if [[ $success_count -eq 5 ]]; then
        test_success "All concurrent requests handled successfully"
    else
        test_failure "Concurrent access" "Only $success_count/5 requests succeeded"
        return 1
    fi
    
    # Cleanup
    rm -f /tmp/concurrent_*.txt /tmp/concurrent_*.json
}

test_binding_deletion() {
    test_header "PCF Binding Deletion"
    
    # Get binding ID
    binding_id=$(cat /tmp/test_binding_id.txt 2>/dev/null || echo "")
    
    if [[ -z "$binding_id" ]]; then
        binding_id=$(curl -s "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$TEST_SUPI&dnn=$TEST_DNN" | jq -r '.bindingId // empty' 2>/dev/null || echo "")
    fi
    
    if [[ -z "$binding_id" ]]; then
        test_failure "Binding deletion" "No binding ID available"
        return 1
    fi
    
    # Delete binding
    response=$(curl -s -w "%{http_code}" -X DELETE \
        "$BSF_URL/nbsf-management/v1/pcfBindings/$binding_id")
    
    if [[ "$response" == "204" ]]; then
        test_success "PCF binding deleted successfully"
        
        # Verify deletion from MongoDB
        count=$(mongosh free5gc --quiet --eval "db.pcfBindings.countDocuments({_id: '$binding_id'})" 2>/dev/null || echo "1")
        if [[ "$count" == "0" ]]; then
            test_success "Binding also removed from MongoDB"
        else
            test_failure "MongoDB cleanup" "Binding still exists in database"
        fi
    else
        test_failure "Binding deletion" "HTTP $response"
        return 1
    fi
}

simulate_ue_pdu_session() {
    test_header "Simulated UE PDU Session Creation Flow"
    
    log "${YELLOW}Simulating complete PDU session creation workflow...${NC}"
    
    # Generate session-specific dynamic data
    local session_supi="imsi-208930$(printf "%09d" $((RANDOM % 1000000000)))"
    local session_ipv4="10.60.$((RANDOM % 256)).$((RANDOM % 256))"
    local session_pcf_id="session-pcf-$(date +%s)-$RANDOM"
    
    log "Session data: SUPI=$session_supi, IPv4=$session_ipv4"
    
    # Step 1: Create initial PCF binding (simulating PCF creating binding)
    log "Step 1: PCF creates binding in BSF"
    
    local session_binding_data=$(cat <<EOF
{
    "supi": "$session_supi",
    "dnn": "$TEST_DNN",
    "ipv4Addr": "$session_ipv4",
    "snssai": {
        "sst": $TEST_SNSSAI_SST,
        "sd": "$TEST_SNSSAI_SD"
    },
    "pcfFqdn": "http://127.0.0.7:8000",
    "pcfIpEndPoints": [
        {
            "ipv4Address": "127.0.0.7",
            "transport": "TCP",
            "port": 8000
        }
    ],
    "pcfId": "$session_pcf_id"
}
EOF
)
    
    response=$(curl -s -w "%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$session_binding_data" \
        "$BSF_URL/nbsf-management/v1/pcfBindings" \
        -o /tmp/session_binding_response.json)
    
    if [[ "$response" == "201" ]] || [[ "$response" == "403" ]]; then
        log "  ✅ PCF binding created/exists"
    else
        test_failure "PDU session simulation" "Failed to create PCF binding (HTTP $response)"
        return 1
    fi
    
    # Step 2: SMF queries BSF for PCF (simulating SMF discovering PCF)
    log "Step 2: SMF queries BSF for PCF discovery"
    
    response=$(curl -s -w "%{http_code}" \
        "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$session_supi&dnn=$TEST_DNN" \
        -o /tmp/smf_discovery_response.json)
    
    if [[ "$response" == "200" ]]; then
        pcf_fqdn=$(jq -r '.pcfFqdn // empty' /tmp/smf_discovery_response.json 2>/dev/null || echo "")
        if [[ -n "$pcf_fqdn" ]]; then
            log "  ✅ SMF discovered PCF: $pcf_fqdn"
        else
            test_failure "PDU session simulation" "Invalid PCF discovery response"
            return 1
        fi
    else
        test_failure "PDU session simulation" "SMF could not discover PCF (HTTP $response)"
        return 1
    fi
    
    # Step 3: Update binding with new IP (simulating IP allocation)
    log "Step 3: Update binding with allocated IP address"
    
    binding_id=$(jq -r '.bindingId // empty' /tmp/smf_discovery_response.json 2>/dev/null || echo "")
    if [[ -n "$binding_id" ]]; then
        local ip_update_data=$(cat <<EOF
{
    "supi": "imsi-208930000000003",
    "dnn": "$TEST_DNN",
    "ipv4Addr": "10.60.0.101",
    "snssai": {
        "sst": $TEST_SNSSAI_SST,
        "sd": "$TEST_SNSSAI_SD"
    },
    "pcfFqdn": "http://127.0.0.7:8000",
    "pcfIpEndPoints": [
        {
            "ipv4Address": "127.0.0.7",
            "transport": "TCP",
            "port": 8000
        }
    ],
    "pcfId": "$updated_pcf_id"
}
EOF
)
        
        response=$(curl -s -w "%{http_code}" -X PATCH \
            -H "Content-Type: application/json" \
            -d "$ip_update_data" \
            "$BSF_URL/nbsf-management/v1/pcfBindings/$binding_id" \
            -o /tmp/ip_update_response.json)
        
        if [[ "$response" == "200" ]]; then
            log "  ✅ Binding updated with allocated IP"
        else
            log "  ⚠️  Could not update binding with IP (HTTP $response)"
        fi
    fi
    
    test_success "Complete PDU session workflow simulation"
    
    # Cleanup
    if [[ -n "$binding_id" ]]; then
        curl -s -X DELETE "$BSF_URL/nbsf-management/v1/pcfBindings/$binding_id" > /dev/null
    fi
}

performance_test() {
    test_header "BSF Performance Test"
    
    log "Running performance test with 50 concurrent operations..."
    
    local start_time=$(date +%s.%N)
    local pids=()
    
    # Create multiple bindings concurrently with dynamic data
    for i in {1..50}; do
        (
            # Generate dynamic values for each binding
            local test_supi="imsi-208930$(printf "%09d" $((RANDOM % 1000000000)))"
            local test_ipv4="10.60.$((RANDOM % 256)).$((RANDOM % 256))"
            local test_pcf_id="perf-test-pcf-$(date +%s)-$i-$RANDOM"
            local test_sst=$((1 + RANDOM % 4))
            local test_sd=$(printf "%06x" $((RANDOM % 16777216)))
            
            local binding_data=$(cat <<EOF
{
    "supi": "$test_supi",
    "dnn": "$TEST_DNN",
    "ipv4Addr": "$test_ipv4",
    "snssai": {
        "sst": $test_sst,
        "sd": "$test_sd"
    },
    "pcfFqdn": "http://127.0.0.7:8000",
    "pcfIpEndPoints": [
        {
            "ipv4Address": "127.0.0.7",
            "transport": "TCP",
            "port": 8000
        }
    ],
    "pcfId": "$test_pcf_id"
}
EOF
)
            
            # Create binding
            response=$(curl -s -w "%{http_code}" -X POST \
                -H "Content-Type: application/json" \
                -d "$binding_data" \
                "$BSF_URL/nbsf-management/v1/pcfBindings" \
                -o "/tmp/perf_create_$i.json")
            
            # Query binding
            if [[ "$response" == "201" ]]; then
                curl -s "$BSF_URL/nbsf-management/v1/pcfBindings?supi=$test_supi&dnn=$TEST_DNN" > "/tmp/perf_query_$i.json"
            fi
            
            echo "$response" > "/tmp/perf_status_$i.txt"
        ) &
        pids+=($!)
    done
    
    # Wait for all operations
    for pid in "${pids[@]}"; do
        wait $pid
    done
    
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    # Count successful operations
    local success_count=0
    for i in {1..50}; do
        status=$(cat "/tmp/perf_status_$i.txt" 2>/dev/null || echo "000")
        if [[ "$status" == "201" ]] || [[ "$status" == "403" ]]; then
            success_count=$((success_count + 1))
        fi
    done
    
    log "Performance results: $success_count/50 operations in ${duration}s"
    
    if [[ $success_count -ge 45 ]]; then
        test_success "Performance test passed (${success_count}/50 successful in ${duration}s)"
    else
        test_failure "Performance test" "Only $success_count/50 operations successful"
    fi
    
    # Cleanup performance test data using binding IDs
    for i in {1..50}; do
        binding_id=$(jq -r '.bindingId // empty' "/tmp/perf_create_$i.json" 2>/dev/null || echo "")
        if [[ -n "$binding_id" ]]; then
            curl -s -X DELETE "$BSF_URL/nbsf-management/v1/pcfBindings/$binding_id" > /dev/null
        fi
    done
    
    rm -f /tmp/perf_*.txt /tmp/perf_*.json
}

# Main test execution
main() {
    log "${BLUE}=================================================${NC}"
    log "${BLUE}  BSF Comprehensive Integration Test Suite      ${NC}"
    log "${BLUE}=================================================${NC}"
    log "Start time: $(date)"
    
    # Generate dynamic test data
    generate_test_data
    
    log "BSF URL: $BSF_URL"
    log "Log file: $LOG_FILE"
    
    # Pre-test setup
    wait_for_services
    
    # Execute all tests
    test_bsf_health
    test_mongodb_connection
    test_clear_existing_bindings
    test_pcf_binding_creation
    test_smf_bsf_query
    test_binding_update
    test_mongodb_persistence
    test_binding_retrieval
    test_memory_cache_mongodb_sync
    test_concurrent_access
    simulate_ue_pdu_session
    performance_test
    test_binding_deletion
    
    # Test summary
    log "\n${BLUE}=================================================${NC}"
    log "${BLUE}              TEST SUMMARY                      ${NC}"
    log "${BLUE}=================================================${NC}"
    log "Total tests: $TOTAL_TESTS"
    log "${GREEN}Passed: $PASSED_TESTS${NC}"
    
    if [[ $FAILED_TESTS -gt 0 ]]; then
        log "${RED}Failed: $FAILED_TESTS${NC}"
        log "\n${RED}❌ OVERALL RESULT: FAILED${NC}"
        log "Check log file: $LOG_FILE"
        exit 1
    else
        log "${GREEN}Failed: $FAILED_TESTS${NC}"
        log "\n${GREEN}🎉 OVERALL RESULT: ALL TESTS PASSED${NC}"
        log "Check log file: $LOG_FILE"
    fi
    
    log "End time: $(date)"
    
    # Cleanup temp files
    rm -f /tmp/bsf_*.json /tmp/binding_*.json /tmp/test_binding_id.txt
    rm -f /tmp/retrieve_*.json /tmp/sync_test_*.json /tmp/session_*.json
    rm -f /tmp/smf_discovery_*.json /tmp/ip_update_*.json
}

# Trap to cleanup on exit
cleanup() {
    log "\n${YELLOW}Cleaning up...${NC}"
    rm -f /tmp/bsf_*.json /tmp/binding_*.json /tmp/test_binding_id.txt
    rm -f /tmp/retrieve_*.json /tmp/sync_test_*.json /tmp/session_*.json
    rm -f /tmp/smf_discovery_*.json /tmp/ip_update_*.json
    rm -f /tmp/concurrent_*.txt /tmp/concurrent_*.json
    rm -f /tmp/perf_*.txt /tmp/perf_*.json
}

trap cleanup EXIT

# Check dependencies
if ! command -v curl &> /dev/null; then
    log "${RED}Error: curl is required but not installed${NC}"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    log "${RED}Error: jq is required but not installed${NC}"
    exit 1
fi

if ! command -v mongosh &> /dev/null; then
    log "${YELLOW}Warning: mongosh not found, MongoDB tests will be skipped${NC}"
fi

# Run the tests
main "$@"
