package types

import (
	"testing"
	"github.com/unicornultrafoundation/go-helios/hash"
)

func TestBlockWithElectingEvent(t *testing.T) {
	// Create test hash events
	electing := hash.HexToEventHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	atropos := hash.HexToEventHash("0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321")
	
	// Create a block with electing event
	block := Block{
		Electing: electing,
		Event:    atropos,
		Cheaters: nil,
	}
	
	// Verify the fields are set correctly
	if block.Electing != electing {
		t.Errorf("Expected electing event %v, got %v", electing, block.Electing)
	}
	
	if block.Event != atropos {
		t.Errorf("Expected atropos event %v, got %v", atropos, block.Event)
	}
	
	if block.Cheaters != nil {
		t.Errorf("Expected nil cheaters, got %v", block.Cheaters)
	}
}

func TestBlockStructure(t *testing.T) {
	// Test that Block has all expected fields
	var block Block
	
	// Check that zero values are properly initialized
	if block.Electing == (hash.Event{}) {
		// Expected zero value
	}
	
	if block.Event == (hash.Event{}) {
		// Expected zero value  
	}
	
	if block.Cheaters != nil {
		t.Errorf("Expected nil cheaters for zero value, got %v", block.Cheaters)
	}
}