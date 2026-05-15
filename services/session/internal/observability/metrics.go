package observability

import "sync/atomic"

var (
	ActiveSessions   int64
	TotalSessions    int64
	TotalPayments    int64
)

func IncActiveSessions()  { atomic.AddInt64(&ActiveSessions, 1) }
func DecActiveSessions()  { atomic.AddInt64(&ActiveSessions, -1) }
func IncTotalSessions()   { atomic.AddInt64(&TotalSessions, 1) }
func IncTotalPayments()   { atomic.AddInt64(&TotalPayments, 1) }
