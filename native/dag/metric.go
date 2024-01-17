package dag

import (
	"fmt"

	"github.com/unicornultrafoundation/go-helios/native/idx"
)

type Metric struct {
	Num  idx.Event
	Size uint64
}

func (m Metric) String() string {
	return fmt.Sprintf("{Num=%d,Size=%d}", m.Num, m.Size)
}
