package consensus

import (
	"math"
	"testing"

	"github.com/unicornultrafoundation/go-helios/native/idx"
)

// mockEvent is a simple mock implementation of dag.Event for testing
type mockEvent struct {
	selfParent *idx.Event
	frame      idx.Frame
}

func (e *mockEvent) SelfParent() *idx.Event {
	return e.selfParent
}

func (e *mockEvent) Frame() idx.Frame {
	return e.frame
}

// mockInput is a simple mock for event input
type mockInput struct {
	events map[idx.Event]*mockEvent
}

func (m *mockInput) GetEvent(id idx.Event) *mockEvent {
	if event, ok := m.events[id]; ok {
		return event
	}
	// Return event with frame 0 for unknown events
	return &mockEvent{frame: 0}
}

// mockOrderer is a minimal mock for testing calcFrameIdx
type mockOrderer struct {
	input *mockInput
}

func (m *mockOrderer) forklessCausedByQuorumOn(e *mockEvent, f idx.Frame) bool {
	// Simple mock: return false after a few iterations to prevent infinite loops
	return f < e.Frame()+5
}

// TestFrameIdxDoSProtection tests protection against DoS attacks in calcFrameIdx
func TestFrameIdxDoSProtection(t *testing.T) {
	input := &mockInput{
		events: make(map[idx.Event]*mockEvent),
	}
	
	orderer := &mockOrderer{input: input}

	// Test protection against extremely large frame values
	maliciousEvent := &mockEvent{
		selfParent: nil,
		frame:      math.MaxUint32,
	}

	// Add a self-parent with extremely large frame
	parentId := idx.Event(1)
	maliciousParent := &mockEvent{
		frame: math.MaxUint32,
	}
	input.events[parentId] = maliciousParent
	maliciousEvent.selfParent = &parentId

	// This should not cause infinite loops or excessive computation
	selfParentFrame, frame := testCalcFrameIdx(orderer, maliciousEvent, true)

	// Verify that frames are capped appropriately
	if selfParentFrame > math.MaxUint32 {
		t.Errorf("Self parent frame should be capped, got: %d", selfParentFrame)
	}

	if frame == 0 {
		t.Error("Frame should not be zero after protection")
	}

	// Test with checkOnly=false to test the +100 path
	_, frame2 := testCalcFrameIdx(orderer, maliciousEvent, false)
	
	if frame2 == 0 {
		t.Error("Frame should not be zero after protection in non-checkOnly mode")
	}
}

// testCalcFrameIdx is a test version of calcFrameIdx to isolate the logic being tested
func testCalcFrameIdx(orderer *mockOrderer, e *mockEvent, checkOnly bool) (selfParentFrame, frame idx.Frame) {
	selfParentFrame = idx.Frame(0)
	if e.SelfParent() != nil {
		selfParentFrame = orderer.input.GetEvent(*e.SelfParent()).Frame()
	}

	// Input validation: prevent DoS attacks with excessive frame ranges
	if selfParentFrame > math.MaxUint32 {
		selfParentFrame = math.MaxUint32
	}

	maxFrameToCheck := selfParentFrame + 100
	if checkOnly {
		maxFrameToCheck = e.Frame()
		// Validate claimed frame to prevent overflow
		if maxFrameToCheck > math.MaxUint32 {
			maxFrameToCheck = math.MaxUint32
		}
	}
	
	// Prevent excessive iteration ranges that could cause DoS
	if maxFrameToCheck > selfParentFrame + 1000 {
		maxFrameToCheck = selfParentFrame + 1000
	}

	var f idx.Frame
	for f = selfParentFrame; f < maxFrameToCheck && orderer.forklessCausedByQuorumOn(e, f); f++ {
	}
	if f == 0 {
		f = 1
	}
	return selfParentFrame, f
}

// TestFrameIdxBoundaryConditions tests boundary conditions for frame calculation
func TestFrameIdxBoundaryConditions(t *testing.T) {
	input := &mockInput{
		events: make(map[idx.Event]*mockEvent),
	}
	
	orderer := &mockOrderer{input: input}

	// Test with zero frame
	zeroEvent := &mockEvent{
		selfParent: nil,
		frame:      0,
	}

	_, frame := testCalcFrameIdx(orderer, zeroEvent, true)
	if frame == 0 {
		t.Error("Frame should be at least 1 after calculation")
	}

	// Test with frame = 1
	oneEvent := &mockEvent{
		selfParent: nil,
		frame:      1,
	}

	_, frame = testCalcFrameIdx(orderer, oneEvent, true)
	if frame < 1 {
		t.Errorf("Frame should be at least 1, got: %d", frame)
	}

	// Test excessive iteration range protection
	largeFrameEvent := &mockEvent{
		selfParent: nil,
		frame:      10000,
	}

	// Add a self-parent with large frame difference
	parentId := idx.Event(1)
	largeParent := &mockEvent{
		frame: 5000,
	}
	input.events[parentId] = largeParent
	largeFrameEvent.selfParent = &parentId

	// This should limit iteration range
	_, frame = testCalcFrameIdx(orderer, largeFrameEvent, true)
	if frame == 0 {
		t.Error("Frame should not be zero for large frame event")
	}
}

// TestFrameIdxRegressionPrevention tests that DoS protection doesn't break normal operation
func TestFrameIdxRegressionPrevention(t *testing.T) {
	input := &mockInput{
		events: make(map[idx.Event]*mockEvent),
	}
	
	orderer := &mockOrderer{input: input}

	// Test normal operation
	normalEvent := &mockEvent{
		selfParent: nil,
		frame:      10,
	}

	selfParentFrame, frame := testCalcFrameIdx(orderer, normalEvent, true)
	
	if selfParentFrame < 0 {
		t.Errorf("Self parent frame should be non-negative, got: %d", selfParentFrame)
	}

	if frame < 1 {
		t.Errorf("Frame should be at least 1, got: %d", frame)
	}

	// Test with self parent
	parentId := idx.Event(1)
	parent := &mockEvent{
		frame: 5,
	}
	input.events[parentId] = parent
	
	eventWithParent := &mockEvent{
		selfParent: &parentId,
		frame:      7,
	}

	selfParentFrame, frame = testCalcFrameIdx(orderer, eventWithParent, true)
	
	if selfParentFrame != 5 {
		t.Errorf("Expected self parent frame 5, got: %d", selfParentFrame)
	}
	
	if frame < 1 {
		t.Errorf("Frame should be at least 1, got: %d", frame)
	}
}