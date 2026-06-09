package public

import "embed"

//go:embed all:dist
var Public embed.FS

//go:embed all:builtin_static
var BuiltinStatic embed.FS
