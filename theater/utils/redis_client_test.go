package utils

import (
	"fmt"
	"testing"
)

func TestAddToCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	err := AddToCache("foo", "bar", redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetFromCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	err := AddToCache("foo", "bar", redisAccess)
	if err != nil {
		t.Fatalf("Failed to add to cache: %v", err)
	}

	value, err := GetFromCache("foo", redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "bar" {
		t.Errorf("Expected value 'bar', got %s", value)
	}
}

func TestDeleteFromCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	err := AddToCache("foo", "bar", redisAccess)
	if err != nil {
		t.Fatalf("Failed to add to cache: %v", err)
	}

	err = DeleteFromCache("foo", redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	value, err := GetFromCache("foo", redisAccess)
	if err == nil {
		t.Errorf("Expected error, got value %s", value)
	}
}

func TestClearCache(t *testing.T) {
	redisAccess := GetStoreAccess()
	err := AddToCache("foo", "bar", redisAccess)
	if err != nil {
		t.Fatalf("Failed to add to cache: %v", err)
	}

	err = ClearCache(redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	value, err := GetFromCache("foo", redisAccess)
	if err == nil {
		t.Errorf("Expected error, got value %s", value)
	}
}
func TestHashSet(t *testing.T) {
	redisAccess := GetStoreAccess()
	hashFields := map[string]interface{}{
		"model": "Deimos",
		"brand": "Ergonom",
		"type":  "Enduro bikes",
		"price": "4972",
	}

	_, err := HashSet("bike:1", hashFields, redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHashGet(t *testing.T) {
	redisAccess := GetStoreAccess()
	hashFields := map[string]interface{}{
		"model": "Deimos",
		"brand": "Ergonom",
		"type":  "Enduro bikes",
		"price": "4972",
	}
	_, err := HashSet("bike:1", hashFields, redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	res1, err := HashGet("bike:1", "model", redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if res1 != "Deimos" {
		t.Errorf("Expected value 'Deimos', got %s", res1)
	}

	res2, err := HashGet("bike:1", "price", redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if res2 != "4972" {
		t.Errorf("Expected value '4972', got %s", res2)
	}
}

func TestHashGetAll(t *testing.T) {
	redisAccess := GetStoreAccess()
	hashFields := map[string]interface{}{
		"model": "Deimos",
		"brand": "Ergonom",
		"type":  "Enduro bikes",
		"price": "4972",
	}
	_, err := HashSet("bike:1", hashFields, redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	res1, err := HashGetAll("bike:1", redisAccess)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(res1) == 0 {
		t.Error("Expected hash to have fields")
	}
	tmp := fmt.Sprintf("%v", res1)
	if tmp != "map[brand:Ergonom model:Deimos price:4972 type:Enduro bikes]" {
		t.Errorf("Expected value 'map[brand:Ergonom model:Deimos price:4972 type:Enduro bikes]', got %s", tmp)
	}
}
