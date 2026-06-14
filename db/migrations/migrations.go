// Package migrations holds the goose migrations applied at startup.
// SQL migrations are embedded; Go migrations register themselves via init().
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
