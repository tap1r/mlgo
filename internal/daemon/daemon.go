// Package daemon starts a mongod or mongos that cannot use --fork.
//
// MongoDB 8.2 and newer reject --fork on macOS. The replacement is a new
// session plus a wait until that process is listening, which is the barrier
// --fork used to provide. The caller must not use exec.CommandContext for
// the server: cancelling that context would kill it.
package daemon

import "time"

// WaitTimeout is how long StartAndWait blocks before giving up.
// The process is left running so a slow startup can still finish.
const WaitTimeout = 120 * time.Second
