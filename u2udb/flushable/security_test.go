package flushable

import (
	"math"
	"testing"

	"github.com/unicornultrafoundation/go-helios/u2udb/devnulldb"
)

// TestMaliciousInputHandling tests handling of malicious input patterns
func TestMaliciousInputHandling(t *testing.T) {
	parent := devnulldb.New()
	db := Wrap(parent)
	defer db.Close()

	// Test nil key/value protection in Put
	err := db.Put(nil, []byte("value"))
	if err == nil {
		t.Error("Expected error for nil key in Put")
	}

	err = db.Put([]byte("key"), nil)
	if err == nil {
		t.Error("Expected error for nil value in Put")
	}

	// Test nil key protection in Delete
	err = db.Delete(nil)
	if err == nil {
		t.Error("Expected error for nil key in Delete")
	}

	// Test extremely large key protection
	largeKey := make([]byte, math.MaxUint16+1)
	err = db.Put(largeKey, []byte("value"))
	if err == nil {
		t.Error("Expected error for extremely large key")
	}

	// Test extremely large value protection
	largeValue := make([]byte, math.MaxInt32/2+1)
	err = db.Put([]byte("key"), largeValue)
	if err == nil {
		t.Error("Expected error for extremely large value")
	}

	// Test large key in Delete
	err = db.Delete(largeKey)
	if err == nil {
		t.Error("Expected error for extremely large key in Delete")
	}
}

// TestBatchMaliciousInputHandling tests batch operations against malicious input
func TestBatchMaliciousInputHandling(t *testing.T) {
	parent := devnulldb.New()
	db := Wrap(parent)
	defer db.Close()

	batch := db.NewBatch()

	// Test nil key/value protection in batch Put
	err := batch.Put(nil, []byte("value"))
	if err == nil {
		t.Error("Expected error for nil key in batch Put")
	}

	err = batch.Put([]byte("key"), nil)
	if err == nil {
		t.Error("Expected error for nil value in batch Put")
	}

	// Test nil key protection in batch Delete
	err = batch.Delete(nil)
	if err == nil {
		t.Error("Expected error for nil key in batch Delete")
	}

	// Test extremely large key protection in batch
	largeKey := make([]byte, math.MaxUint16+1)
	err = batch.Put(largeKey, []byte("value"))
	if err == nil {
		t.Error("Expected error for extremely large key in batch")
	}

	// Test extremely large value protection in batch
	largeValue := make([]byte, math.MaxInt32/2+1)
	err = batch.Put([]byte("key"), largeValue)
	if err == nil {
		t.Error("Expected error for extremely large value in batch")
	}

	// Test large key in batch Delete
	err = batch.Delete(largeKey)
	if err == nil {
		t.Error("Expected error for extremely large key in batch Delete")
	}

	batch.Reset()
}

// TestSizeEstimationOverflowProtection tests protection against size estimation overflow
func TestSizeEstimationOverflowProtection(t *testing.T) {
	parent := devnulldb.New()
	db := Wrap(parent)
	defer db.Close()

	// Fill up to near max size
	bigKey := make([]byte, 1000)
	bigValue := make([]byte, math.MaxInt32/4)
	
	// This should succeed without overflow
	err := db.Put(bigKey, bigValue)
	if err != nil {
		t.Fatalf("Unexpected error for large but valid input: %v", err)
	}

	// Check that size estimation is reasonable
	size := db.NotFlushedSizeEst()
	if size < 0 {
		t.Error("Size estimation should never be negative")
	}

	// Adding another large item should cap the size, not overflow
	err = db.Put(bigKey, bigValue)
	if err != nil {
		t.Fatalf("Unexpected error for second large input: %v", err)
	}

	newSize := db.NotFlushedSizeEst()
	if newSize < size {
		t.Error("Size should not decrease without explicit operations")
	}

	// Test batch size overflow protection
	batch := db.NewBatch()
	for i := 0; i < 100; i++ {
		key := make([]byte, 10000)
		value := make([]byte, 10000)
		err := batch.Put(key, value)
		if err != nil {
			break // Expected for very large batches
		}
	}

	batchSize := batch.ValueSize()
	if batchSize < 0 {
		t.Error("Batch size should never be negative")
	}

	batch.Reset()
}

// TestSecurityRegressionPrevention tests that security fixes don't break normal operation
func TestSecurityRegressionPrevention(t *testing.T) {
	parent := devnulldb.New()
	db := Wrap(parent)
	defer db.Close()

	// Test normal operation still works
	key := []byte("testkey")
	value := []byte("testvalue")

	err := db.Put(key, value)
	if err != nil {
		t.Fatalf("Normal Put should work: %v", err)
	}

	has, err := db.Has(key)
	if err != nil {
		t.Fatalf("Normal Has should work: %v", err)
	}
	if !has {
		t.Error("Key should exist after Put")
	}

	retrieved, err := db.Get(key)
	if err != nil {
		t.Fatalf("Normal Get should work: %v", err)
	}
	if string(retrieved) != string(value) {
		t.Error("Retrieved value should match put value")
	}

	err = db.Delete(key)
	if err != nil {
		t.Fatalf("Normal Delete should work: %v", err)
	}

	has, err = db.Has(key)
	if err != nil {
		t.Fatalf("Has should work after Delete: %v", err)
	}
	if has {
		t.Error("Key should not exist after Delete")
	}

	// Test batch operations
	batch := db.NewBatch()
	for i := 0; i < 10; i++ {
		batchKey := []byte("batchkey" + string(rune('0'+i)))
		batchValue := []byte("batchvalue" + string(rune('0'+i)))
		err := batch.Put(batchKey, batchValue)
		if err != nil {
			t.Fatalf("Normal batch Put should work: %v", err)
		}
	}

	err = batch.Write()
	if err != nil {
		t.Fatalf("Normal batch Write should work: %v", err)
	}

	batch.Reset()
}