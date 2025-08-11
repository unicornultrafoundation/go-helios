package wlru

import (
	"math"
	"testing"
)

// TestNegativeInputValidation tests validation for negative input values
func TestNegativeInputValidation(t *testing.T) {
	// Test negative maxSize in constructor - should return error
	cache, err := New(100, -10)
	if err == nil {
		t.Fatal("Expected error for negative maxSize")
	}
	if cache != nil {
		t.Fatal("Expected nil cache for negative maxSize")
	}

	// Test negative maxSize in Resize with valid cache
	validCache, err := New(100, 10)
	if err != nil {
		t.Fatalf("Failed to create valid cache: %v", err)
	}
	
	evicted := validCache.Resize(100, -5)
	if evicted < 0 {
		t.Errorf("Expected non-negative evicted count, got: %d", evicted)
	}
}

// TestOverflowProtection tests protection against integer overflow conditions
func TestOverflowProtection(t *testing.T) {
	// Test maxSize overflow protection in constructor
	cache, err := NewWithEvict(math.MaxUint32, math.MaxInt32+1, nil)
	if err != nil {
		t.Fatalf("Expected no error for large maxSize, got: %v", err)
	}
	if cache == nil {
		t.Fatal("Expected cache to be created with capped maxSize")
	}

	// Test maxWeight overflow protection in constructor
	cache2, err := NewWithEvict(math.MaxUint32+1, 100, nil)
	if err != nil {
		t.Fatalf("Expected no error for large maxWeight, got: %v", err)
	}
	if cache2 == nil {
		t.Fatal("Expected cache to be created with capped maxWeight")
	}

	// Test overflow protection in Resize
	cache.Resize(math.MaxUint32+1, math.MaxInt32+1)
	// Should not crash or cause overflow
}

// TestMaliciousInputHandling tests handling of malicious input patterns
func TestMaliciousInputHandling(t *testing.T) {
	cache, err := New(1000, 10)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test nil key/value protection in Add
	evicted := cache.Add(nil, "value", 1)
	if evicted != 0 {
		t.Errorf("Expected 0 evicted for nil key, got: %d", evicted)
	}

	evicted = cache.Add("key", nil, 1)
	if evicted != 0 {
		t.Errorf("Expected 0 evicted for nil value, got: %d", evicted)
	}

	// Test excessive weight protection
	cache.Add("key1", "value1", math.MaxUint32+1)
	// Should not crash

	// Test nil key protection in other methods
	if cache.Contains(nil) {
		t.Error("Expected false for nil key in Contains")
	}

	if _, ok := cache.Get(nil); ok {
		t.Error("Expected false for nil key in Get")
	}

	if _, ok := cache.Peek(nil); ok {
		t.Error("Expected false for nil key in Peek")
	}

	if cache.Remove(nil) {
		t.Error("Expected false for nil key in Remove")
	}

	// Test nil protection in ContainsOrAdd
	if ok, _ := cache.ContainsOrAdd(nil, "value", 1); ok {
		t.Error("Expected false for nil key in ContainsOrAdd")
	}

	if ok, _ := cache.ContainsOrAdd("key", nil, 1); ok {
		t.Error("Expected false for nil value in ContainsOrAdd")
	}

	// Test nil protection in PeekOrAdd
	if _, ok, _ := cache.PeekOrAdd(nil, "value", 1); ok {
		t.Error("Expected false for nil key in PeekOrAdd")
	}

	if _, ok, _ := cache.PeekOrAdd("key", nil, 1); ok {
		t.Error("Expected false for nil value in PeekOrAdd")
	}
}

// TestSecurityRegressionPrevention tests that security fixes don't break existing functionality
func TestSecurityRegressionPrevention(t *testing.T) {
	cache, err := New(1000, 10)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test normal operation still works
	cache.Add("key1", "value1", 10)
	if !cache.Contains("key1") {
		t.Error("Expected key1 to be present")
	}

	if value, ok := cache.Get("key1"); !ok || value != "value1" {
		t.Error("Expected to retrieve value1 for key1")
	}

	// Test resize still works
	evicted := cache.Resize(500, 5)
	if evicted < 0 {
		t.Errorf("Expected non-negative evicted count, got: %d", evicted)
	}
}