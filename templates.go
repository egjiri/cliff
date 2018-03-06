package cli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func updateTemplates() {
	addTemplateFuncs()
	rootCmd.SetHelpTemplate(colorTemplate)
}

func addTemplateFuncs() {
	cobra.AddTemplateFunc("colorYellow", color.YellowString)
	cobra.AddTemplateFunc("colorGreen", color.GreenString)
	cobra.AddTemplateFunc("colorBold", func(format string) string {
		return color.New(color.Bold).Sprint(format)
	})
	cobra.AddTemplateFunc("colorGrey", func(format string) string {
		return color.New(color.FgWhite).Sprint(format)
	})
}

var colorTemplate = `
{{- if .Long}}
	{{- .Long}}

{{end -}}

{{"Usage: " | colorGreen}}{{if .Runnable}}
{{- if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"| colorYellow}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
{{- .CommandPath | colorYellow}} {{"COMMAND" | colorBold}}{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{if .HasAvailableSubCommands}}

{{"Available Commands:" | colorGreen}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding | colorBold }} {{.Short | colorGrey}}{{end}}{{end}}{{end}}

{{- if .HasAvailableLocalFlags}}

{{"Flags:" | colorGreen}}
{{.LocalFlags.FlagUsages | trimRightSpace -}}
{{end}}

{{- if .HasAvailableInheritedFlags}}

{{"Global Flags:" | colorGreen}}
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
