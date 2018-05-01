package cliff

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func updateTemplates(cmd *cobra.Command) {
	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)
}

func init() {
	addTemplateFuncs()
}

func addTemplateFuncs() {
	cobra.AddTemplateFunc("colorCyan", color.CyanString)
	cobra.AddTemplateFunc("colorGreen", color.GreenString)
	cobra.AddTemplateFunc("colorYellow", color.YellowString)
	cobra.AddTemplateFunc("colorBold", func(format string) string {
		return color.New(color.Bold).Sprint(format)
	})
	cobra.AddTemplateFunc("colorGrey", func(format string) string {
		return color.New(color.FgWhite).Sprint(format)
	})
}

const helpTemplate = `{{with (or .Long .Short)}}{{. | colorCyan | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

const usageTemplate = `
{{"Usage:" | colorGreen}}{{if .Runnable}}
  {{.UseLine | colorYellow}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath | colorYellow}} {{"[command]" | colorYellow}}{{end}}{{if gt (len .Aliases) 0}}

{{"Aliases: " | colorGreen}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{"Examples: " | colorGreen}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

{{"Available Commands: " | colorGreen}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding | colorBold }}   {{.Short | colorGrey}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{"Flags: " | colorGreen}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{"Global Flags: " | colorGreen}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

{{"Additional help topics: " | colorGreen}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
