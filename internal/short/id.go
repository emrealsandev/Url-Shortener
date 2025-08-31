package short

import (
	"context"
	"sync/atomic"
)

var seq uint64 = 100000

func NextID(_ context.Context) uint64 {
	return atomic.AddUint64(&seq, 1)
}
