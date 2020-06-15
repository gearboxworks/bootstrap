package cmd

import (
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
)


func Version(cmd *cobra.Command, args ...string) error {
	err := VersionShow()
	SetHelp(cmd)
	PrintHelp(cmd)
	err = cmd.Help()
	return err
}


func VersionShow() error {
	ux.PrintfBlue("%s ", defaults.BinaryName)
	ux.PrintflnCyan("v%s", defaults.BinaryVersion)
	return nil
}


func VersionInfo(args ...string) error {
	update := New(&Target)

	for range onlyOnce {
		if len(args) == 0 {
			args = []string{CmdVersionLatest}
		}

		if update.Error != nil {
			break
		}

		for _, v := range args {
			update.Error = update.PrintVersion(GetSemVer(v))
			if update.Error != nil {
				break
			}
		}
	}

	return update.Error
}


func VersionList(args ...string) error {
	update := New(&Target)

	for range onlyOnce {
		if len(args) == 0 {
			// @TODO = Obtain full list of versions.
			args = []string{CmdVersionLatest}
		}

		if update.Error != nil {
			break
		}

		for _, v := range args {
			update.Error = update.PrintVersion(GetSemVer(v))
			if update.Error != nil {
				break
			}
		}
	}

	return update.Error
}


func VersionCheck() error {
	update := New(&Target)

	for range onlyOnce {
		if update.Error != nil {
			break
		}

		update.Error = update.IsUpdated()
		if update.Error != nil {
			break
		}
	}

	return update.Error
}


func VersionUpdate() error {
	update := New(&Target)

	for range onlyOnce {
		if update.Error != nil {
			break
		}

		update.Error = CreateDummyBinary(Runtime.Cmd, Target.Cmd)
		if update.Error != nil {
			break
		}

		update.Error = update.IsUpdated()
		if update.Error != nil {
			break
		}

		update.Error = update.Update()
		if update.Error != nil {
			break
		}
	}

	return update.Error
}


func _GetUsage(c *cobra.Command) string {
	var str string

	if c.Parent() == nil {
		str += ux.SprintfCyan("%s ", c.Name())
	} else {
		str += ux.SprintfCyan("%s ", c.Parent().Name())
		str += ux.SprintfGreen("%s ", c.Use)
	}

	if c.HasAvailableSubCommands() {
		str += ux.SprintfGreen("[command] ")
		str += ux.SprintfCyan("<args> ")
	}

	return str
}


func _GetVersion(c *cobra.Command) string {
	var str string

	if c.Parent() == nil {
		str = ux.SprintfBlue("%s ", defaults.BinaryName)
		str += ux.SprintfCyan("v%s", defaults.BinaryVersion)
	}

	return str
}


func PrintHelp(c *cobra.Command) {
	ux.PrintfCyan("bootstrap")
	ux.PrintflnBlue(" is intended to automatically download the correct binary from a GitHub repository.\n")

	if Runtime.CmdFile == defaults.BinaryName {

	} else {
		ux.PrintflnBlue("The '%s' executable is running this bootstrap code.", Runtime.CmdName)

		ux.PrintflnBlue("To be able to use the real '%s' executable, replace this bootstrap binary using this command:", Runtime.CmdName)
		ux.PrintflnCyan("%s version update", Runtime.CmdName)
		ux.PrintflnBlue("\nThis will update and replace the current '%s' file with the correct executable.\n\n", Runtime.CmdName)
	}
}


func SetHelp(c *cobra.Command) {
	var tmplHelp string
	var tmplUsage string

	cobra.AddTemplateFunc("GetUsage", _GetUsage)
	cobra.AddTemplateFunc("GetVersion", _GetVersion)

	cobra.AddTemplateFunc("SprintfBlue", ux.SprintfBlue)
	cobra.AddTemplateFunc("SprintfCyan", ux.SprintfCyan)
	cobra.AddTemplateFunc("SprintfGreen", ux.SprintfGreen)
	cobra.AddTemplateFunc("SprintfMagenta", ux.SprintfMagenta)
	cobra.AddTemplateFunc("SprintfRed", ux.SprintfRed)
	cobra.AddTemplateFunc("SprintfWhite", ux.SprintfWhite)
	cobra.AddTemplateFunc("SprintfYellow", ux.SprintfYellow)

	tmplUsage += `
{{ SprintfBlue "Usage: " }}
	{{ GetUsage . }}

{{- if gt (len .Aliases) 0 }}
{{ SprintfBlue "\nAliases:" }} {{ .NameAndAliases }}
{{- end }}

{{- if .HasExample }}
{{ SprintfBlue "\nExamples:" }}
	{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
{{- range .Commands }}
{{- if (or .IsAvailableCommand (eq .Name "help")) }}
	{{ rpad (SprintfGreen .Name) .NamePadding}}	- {{ .Short }}{{ end }}
{{- end }}
{{- end }}

{{- if .HasHelpSubCommands }}
{{- SprintfBlue "\nAdditional help topics:" }}
{{- range .Commands }}
{{- if .IsAdditionalHelpTopicCommand }}
	{{ rpad (SprintfGreen .CommandPath) .CommandPathPadding }} {{ .Short }}
{{- end }}
{{- end }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nUse" }} {{ SprintfCyan .CommandPath }} {{ SprintfCyan "help" }} {{ SprintfGreen "[command]" }} {{ SprintfBlue "for more information about a command." }}
{{- end }}
`

	tmplHelp = `{{ GetVersion . }}

{{ SprintfBlue "Commmand:" }} {{ SprintfCyan .Use }}

{{ SprintfBlue "Description:" }} 
	{{ with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`

	//c.SetHelpCommand(c)
	//c.SetHelpFunc(PrintHelp)
	c.SetHelpTemplate(tmplHelp)
	c.SetUsageTemplate(tmplUsage)
}


//func ProcessArgs(toolArgs *TypeRuntime, cmd *cobra.Command, args []string) error {
//	state := Runtime.State
//
//	for range onlyOnce {
//		//err := toolArgs.Runtime.SetArgs(cmd.Use)
//		//if err != nil {
//		//	state.SetError(err)
//		//	break
//		//}
//		//
//		//err = toolArgs.Runtime.AddArgs(args...)
//		//if err != nil {
//		//	state.SetError(err)
//		//	break
//		//}
//
//		SetApp(Runtime.CmdName)
//	}
//
//	return state
//}


// 	BinaryRepo = "github.com/gearboxworks/scribe"
//	BinaryRepo = "github.com/gearboxworks/bootstrap"
//	BinaryRepo = "github.com/gearboxworks/buildtool"
//	BinaryRepo = "github.com/wplib/deploywp"
//	BinaryRepo = "github.com/gearboxworks/launch"
