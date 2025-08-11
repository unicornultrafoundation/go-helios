package types

import "github.com/unicornultrafoundation/go-helios/hash"

// Block is a part of an ordered chain of batches of events.
type Block struct {
	Electing hash.Event // Electing event that triggered the consensus decision
	Event    hash.Event // Atropos event containing the final decision
	Cheaters Cheaters   // List of detected cheating validators
}
