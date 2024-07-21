// Package migrations contains migrations for data base.
package migrations

import "embed"

//go:embed *.sql
var Content embed.FS
