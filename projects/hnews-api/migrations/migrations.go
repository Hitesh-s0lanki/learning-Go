// Package migrations embeds the SQL migration files so they ship inside the
// compiled binary — no separate files to deploy.
package migrations

import "embed"

// FS holds every .sql migration in this directory, applied in filename order.
//
//go:embed *.sql
var FS embed.FS
