package migration

import "embed"

//go:embed *.sql
var MigrateFS embed.FS
