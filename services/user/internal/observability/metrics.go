package observability

import "sync/atomic"

var (
	TotalRegistrations int64
	ActiveSessions     int64
)

func IncRegistrations()  { atomic.AddInt64(&TotalRegistrations, 1) }
func IncActiveSessions() { atomic.AddInt64(&ActiveSessions, 1) }
func DecActiveSessions() { atomic.AddInt64(&ActiveSessions, -1) }
