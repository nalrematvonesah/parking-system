package observability

import "sync/atomic"

var ActiveAssignments int64

func IncAssign() { atomic.AddInt64(&ActiveAssignments, 1) }
func DecAssign() { atomic.AddInt64(&ActiveAssignments, -1) }
