// Package migration menyimpan skema database beserta file SQL-nya.
package migration

import "embed"

// SQL berisi seluruh file migration. Di-embed supaya binary tidak bergantung
// pada keberadaan folder sql/ saat runtime (mis. di dalam image Docker).
//
//go:embed sql/*.sql
var SQL embed.FS
