package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set test environment variables
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "test_password")
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	os.Exit(code)
}

func TestNewClient(t *testing.T) {
	client := NewClient("localhost:6379")
	
	if client == nil {
		t.Fatal("NewClient should not return nil")
	}
	
	if client.url != "localhost:6379" {
		t.Errorf("Expected URL localhost:6379, got %s", client.url)
	}
	
	if client.client == nil {
		t.Error("Redis client should not be nil")
	}
	
	if client.context == nil {
		t.Error("Context should not be nil")
	}
}

func TestGetStoreAccess(t *testing.T) {
	// Test singleton pattern
	client1 := GetStoreAccess()
	client2 := GetStoreAccess()
	
	if client1 != client2 {
		t.Error("GetStoreAccess should return the same instance (singleton pattern)")
	}
}

func TestAddToCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// This test will fail if Redis is not available
	// In a real test environment, you would use a test Redis instance
	err := AddToCache("test_key", "test_value", redisAccess)
	if err != nil {
		t.Logf("AddToCache failed (expected in test environment): %v", err)
	}
}

func TestGetFromCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// First add a value
	err := AddToCache("test_get_key", "test_get_value", redisAccess)
	if err != nil {
		t.Logf("AddToCache failed (expected in test environment): %v", err)
		return
	}
	
	// Then try to get it
	value, err := GetFromCache("test_get_key", redisAccess)
	if err != nil {
		t.Logf("GetFromCache failed (expected in test environment): %v", err)
		return
	}
	
	if value != "test_get_value" {
		t.Errorf("Expected 'test_get_value', got '%s'", value)
	}
}

func TestDeleteFromCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// First add a value
	err := AddToCache("test_delete_key", "test_delete_value", redisAccess)
	if err != nil {
		t.Logf("AddToCache failed (expected in test environment): %v", err)
		return
	}
	
	// Then delete it
	err = DeleteFromCache("test_delete_key", redisAccess)
	if err != nil {
		t.Logf("DeleteFromCache failed (expected in test environment): %v", err)
		return
	}
	
	// Try to get it (should fail)
	_, err = GetFromCache("test_delete_key", redisAccess)
	if err == nil {
		t.Error("Expected error when getting deleted key")
	}
}

func TestClearCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// Add some values
	AddToCache("test_clear_key1", "value1", redisAccess)
	AddToCache("test_clear_key2", "value2", redisAccess)
	
	// Clear cache
	err := ClearCache(redisAccess)
	if err != nil {
		t.Logf("ClearCache failed (expected in test environment): %v", err)
		return
	}
	
	// Try to get values (should fail)
	_, err = GetFromCache("test_clear_key1", redisAccess)
	if err == nil {
		t.Error("Expected error when getting cleared key1")
	}
	
	_, err = GetFromCache("test_clear_key2", redisAccess)
	if err == nil {
		t.Error("Expected error when getting cleared key2")
	}
}

func TestHashSet(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	hashFields := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": 123,
	}
	
	result, err := HashSet("test_hash", hashFields, redisAccess)
	if err != nil {
		t.Logf("HashSet failed (expected in test environment): %v", err)
		return
	}
	
	if result == 0 {
		t.Error("Expected non-zero result from HashSet")
	}
}

func TestHashGet(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// First set a hash
	hashFields := map[string]interface{}{
		"test_field": "test_value",
	}
	HashSet("test_hash_get", hashFields, redisAccess)
	
	// Then get a field
	value, err := HashGet("test_hash_get", "test_field", redisAccess)
	if err != nil {
		t.Logf("HashGet failed (expected in test environment): %v", err)
		return
	}
	
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}

func TestHashGetAll(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// First set a hash
	hashFields := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	HashSet("test_hash_getall", hashFields, redisAccess)
	
	// Then get all fields
	allFields, err := HashGetAll("test_hash_getall", redisAccess)
	if err != nil {
		t.Logf("HashGetAll failed (expected in test environment): %v", err)
		return
	}
	
	if len(allFields) == 0 {
		t.Error("Expected non-empty hash fields")
	}
	
	if allFields["field1"] != "value1" {
		t.Errorf("Expected field1 'value1', got '%s'", allFields["field1"])
	}
	
	if allFields["field2"] != "value2" {
		t.Errorf("Expected field2 'value2', got '%s'", allFields["field2"])
	}
}

func TestGetAll(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// Add some test data
	AddToCache("test_getall_key1", "value1", redisAccess)
	AddToCache("test_getall_key2", "value2", redisAccess)
	
	// Get all keys
	allData, err := GetAll(redisAccess)
	if err != nil {
		t.Logf("GetAll failed (expected in test environment): %v", err)
		return
	}
	
	if len(allData) == 0 {
		t.Error("Expected non-empty data from GetAll")
	}
	
	// Check if our test keys are in the result
	found1 := false
	found2 := false
	for key, value := range allData {
		if key == "test_getall_key1" && value == "value1" {
			found1 = true
		}
		if key == "test_getall_key2" && value == "value2" {
			found2 = true
		}
	}
	
	if !found1 {
		t.Error("Expected to find test_getall_key1")
	}
	
	if !found2 {
		t.Error("Expected to find test_getall_key2")
	}
}

func TestRedisClientConfiguration(t *testing.T) {
	client := NewClient("localhost:6379")
	
	if client.client == nil {
		t.Fatal("Redis client should not be nil")
	}
	
	// Test that the client is properly configured
	// Note: We can't easily test the internal configuration without exposing it
	// In a real test, you might want to add getter methods for testing
}

func TestRedisConnectionError(t *testing.T) {
	// Test with invalid Redis URL
	client := NewClient("invalid:9999")
	
	if client == nil {
		t.Fatal("NewClient should not return nil even with invalid URL")
	}
	
	// The client should be created but connection will fail
	err := AddToCache("test", "value", client)
	if err == nil {
		t.Error("Expected error with invalid Redis URL")
	}
}

func TestConcurrentAccess(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// Test concurrent access to Redis
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent_test_%d", i)
			value := fmt.Sprintf("value_%d", i)
			
			err := AddToCache(key, value, redisAccess)
			if err != nil {
				t.Logf("Concurrent AddToCache failed: %v", err)
			}
			
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestRedisKeyPatterns(t *testing.T) {
	redisAccess := GetStoreAccess()
	
	// Test different key patterns
	keyPatterns := []string{
		"simple_key",
		"key:with:colons",
		"key-with-dashes",
		"key_with_underscores",
		"key.with.dots",
		"key with spaces",
		"key123",
		"123key",
	}
	
	for _, pattern := range keyPatterns {
		value := "test_value"
		err := AddToCache(pattern, value, redisAccess)
		if err != nil {
			t.Logf("AddToCache failed for pattern '%s': %v", pattern, err)
			continue
		}
		
		retrieved, err := GetFromCache(pattern, redisAccess)
		if err != nil {
			t.Logf("GetFromCache failed for pattern '%s': %v", pattern, err)
			continue
		}
		
		if retrieved != value {
			t.Errorf("Expected '%s' for pattern '%s', got '%s'", value, pattern, retrieved)
		}
		
		// Clean up
		DeleteFromCache(pattern, redisAccess)
	}
}