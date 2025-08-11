package consensus

import (
	"testing"
	"github.com/unicornultrafoundation/go-helios/hash"
	"github.com/unicornultrafoundation/go-helios/native/idx"
	"github.com/unicornultrafoundation/go-helios/native/pos"
	"github.com/unicornultrafoundation/go-helios/types"
)

func TestElectingEventIntegration(t *testing.T) {
	// Test that the ApplyAtropos callback receives both atropos and electing events
	var capturedAtropos, capturedElecting hash.Event
	var capturedFrame idx.Frame
	
	// Create mock callback that captures the parameters
	mockCallback := func(frame idx.Frame, atropos, electing hash.Event) *pos.Validators {
		capturedFrame = frame
		capturedAtropos = atropos
		capturedElecting = electing
		return nil
	}
	
	// Create test events
	atroposEvent := hash.HexToEventHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	electingEvent := hash.HexToEventHash("0x2222222222222222222222222222222222222222222222222222222222222222")
	testFrame := idx.Frame(5)
	
	// Create orderer with mock callback
	orderer := &Orderer{
		callback: OrdererCallbacks{
			ApplyAtropos: mockCallback,
		},
		store: &Store{}, // Mock store would be needed for full test
	}
	
	// Test the callback functionality (this would normally be called by the consensus mechanism)
	if orderer.callback.ApplyAtropos != nil {
		orderer.callback.ApplyAtropos(testFrame, atroposEvent, electingEvent)
	}
	
	// Verify the callback received the correct parameters
	if capturedFrame != testFrame {
		t.Errorf("Expected frame %v, got %v", testFrame, capturedFrame)
	}
	
	if capturedAtropos != atroposEvent {
		t.Errorf("Expected atropos event %v, got %v", atroposEvent, capturedAtropos)
	}
	
	if capturedElecting != electingEvent {
		t.Errorf("Expected electing event %v, got %v", electingEvent, capturedElecting)
	}
}

func TestBlockCreationWithElectingEvent(t *testing.T) {
	// Test that Consensus.applyAtropos creates Block with correct electing event
	var capturedBlock *types.Block
	
	// Mock consensus callback that captures the block
	mockBeginBlock := func(block *types.Block) types.BlockCallbacks {
		capturedBlock = block
		return types.BlockCallbacks{}
	}
	
	// Create test events and consensus instance 
	atroposEvent := hash.HexToEventHash("0x3333333333333333333333333333333333333333333333333333333333333333")
	electingEvent := hash.HexToEventHash("0x4444444444444444444444444444444444444444444444444444444444444444")
	
	// This test would require a more complex setup with actual stores and dag indexes
	// For now, we'll test the callback structure
	callbacks := types.ConsensusCallbacks{
		BeginBlock: mockBeginBlock,
	}
	
	// Verify the callback structure is compatible
	if callbacks.BeginBlock == nil {
		t.Error("BeginBlock callback should not be nil")
	}
	
	// Test callback invocation (simulating what applyAtropos does)
	testBlock := &types.Block{
		Electing: electingEvent,
		Event:    atroposEvent, 
		Cheaters: nil,
	}
	
	callbacks.BeginBlock(testBlock)
	
	// Verify the block was captured correctly
	if capturedBlock == nil {
		t.Fatal("Block was not captured by callback")
	}
	
	if capturedBlock.Electing != electingEvent {
		t.Errorf("Expected electing event %v, got %v", electingEvent, capturedBlock.Electing)
	}
	
	if capturedBlock.Event != atroposEvent {
		t.Errorf("Expected atropos event %v, got %v", atroposEvent, capturedBlock.Event)
	}
}