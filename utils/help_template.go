package utils

var HelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{.UsageText}}

   {{.Description}}

CATEGORY: {{.Category}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
`