package cmd

import (
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/toolRuntime"
	"github.com/newclarity/scribeHelpers/toolSelfUpdate"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"strings"
)


const onlyOnce = "1"
var onlyTwice = []string{"", ""}
var RunTime *toolRuntime.TypeRuntime
var SelfUpdate *toolSelfUpdate.TypeSelfUpdate


func init() {
	for range onlyOnce {
		if RunTime == nil {
			RunTime = toolRuntime.New(defaults.BinaryName, defaults.BinaryVersion, false)
			RunTime.SetRepos(defaults.SourceRepo, defaults.BinaryRepo)
		}

		if SelfUpdate == nil {
			SelfUpdate = toolSelfUpdate.New(RunTime)
			//SelfUpdate.LoadCommands(rootCmd, false)
			if SelfUpdate.State.IsNotOk() {
				break
			}

			rootCmd.Flags().BoolP(FlagVersion, "v", false, ux.SprintfBlue("Display version of %s", SelfUpdate.Runtime.CmdName))
		}
	}
}


var rootCmd = &cobra.Command{
	Use:   defaults.BinaryName,
	Short: "Bootstrap is the automatic app downloader.",
	Long: "Bootstrap is the automatic app downloader.",
	SilenceErrors: true,
	Run: gbRootFunc,
}
func gbRootFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		if SelfUpdate.FlagCheckVersion(cmd) {
			SelfUpdate.State.SetOk()
			break
		}

		if !SelfUpdate.IsBootstrapBinary() {
			VersionUpdate(cmd, nil)
			break
		}

		// Show help if no commands specified.
		if len(args) == 0 {
			PrintHelp()
			_ = cmd.Help()
			break
		}
	}
}


func Execute() *ux.State {
	for range onlyOnce {
		SetHelp(rootCmd)

		err := rootCmd.Execute()
		if err == nil {
			break
		}

		if !strings.HasPrefix(err.Error(), "unknown command") {
			SelfUpdate.State.SetError(err)
			break
		}

		// Assume a 'version update'
		//if !SelfUpdate.IsBootstrapBinary() {
			VersionUpdate(rootCmd, RunTime.FullArgs)
		//}
	}

	return SelfUpdate.State
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


func PrintHelp() {
	ux.PrintfCyan("bootstrap")
	ux.PrintflnBlue(" is intended to automatically download the correct binary from a GitHub repository.\n")

	if RunTime.CmdFile == defaults.BinaryName {

	} else {
		ux.PrintflnBlue("The '%s' executable is running this bootstrap code.", RunTime.CmdName)

		ux.PrintflnBlue("To be able to use the real '%s' executable, replace this bootstrap binary using this command:", RunTime.CmdName)
		ux.PrintflnCyan("%s version update", RunTime.CmdName)
		ux.PrintflnBlue("\nThis will update and replace the current '%s' file with the correct executable.\n\n", RunTime.CmdName)
	}
}


func SetHelp(c *cobra.Command) {
	var tmplHelp string
	var tmplUsage string

	cobra.AddTemplateFunc("GetUsage", _GetUsage)
	cobra.AddTemplateFunc("GetVersion", _GetVersion)
	cobra.AddTemplateFunc("HelpExamples", HelpExamples)

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

{{- if .HasAvailableLocalFlags }}
{{ SprintfBlue "\nFlags:" }}
{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}

{{- if .HasAvailableInheritedFlags }}
{{ SprintfBlue "\nGlobal Flags:" }}
{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
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


func HelpExamples() string {
	var ret string

	ret += ux.SprintfWhite("\nExamples for: %s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionUpdate)
	ret += ux.SprintfBlue(" - List all available versions of the '%s' binary.\n", defaults.BinaryName)
	ret += ux.SprintfMagenta("\t%s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionUpdate)
	ret += ux.SprintfBlue(" - Update to the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool\n", defaults.BinaryName, CmdVersion, CmdVersionUpdate)
	ret += ux.SprintfBlue(" - Update to the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool/latest\n", defaults.BinaryName, CmdVersion, CmdVersionUpdate)
	ret += ux.SprintfBlue(" - Update to version 1.1.3 within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool/1.1.3\n", defaults.BinaryName, CmdVersion, CmdVersionUpdate)

	ret += ux.SprintfWhite("\nExamples for: %s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionCheck)
	ret += ux.SprintfBlue(" - Check the latest version of the '%s' binary.\n", defaults.BinaryName)
	ret += ux.SprintfMagenta("\t%s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionCheck)
	ret += ux.SprintfBlue(" - Check the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool\n", defaults.BinaryName, CmdVersion, CmdVersionCheck)
	ret += ux.SprintfBlue(" - Check the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool/latest\n", defaults.BinaryName, CmdVersion, CmdVersionCheck)
	ret += ux.SprintfBlue(" - Check version 1.1.3 within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool/1.1.3\n", defaults.BinaryName, CmdVersion, CmdVersionCheck)

	ret += ux.SprintfWhite("\nExamples for: %s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionList)
	ret += ux.SprintfBlue(" - List all available versions of the '%s' binary.\n", defaults.BinaryName)
	ret += ux.SprintfMagenta("\t%s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionList)
	ret += ux.SprintfBlue(" - List all available versions within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool\n", defaults.BinaryName, CmdVersion, CmdVersionList)

	ret += ux.SprintfWhite("\nExamples for: %s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionInfo)
	ret += ux.SprintfBlue(" - Show info on the current version of the '%s' binary.\n", defaults.BinaryName)
	ret += ux.SprintfMagenta("\t%s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionInfo)
	ret += ux.SprintfBlue(" - Show info on the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool\n", defaults.BinaryName, CmdVersion, CmdVersionInfo)
	ret += ux.SprintfBlue(" - Show info on the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool/latest\n", defaults.BinaryName, CmdVersion, CmdVersionInfo)
	ret += ux.SprintfBlue(" - Show info on version 1.1.3 within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool/1.1.3\n", defaults.BinaryName, CmdVersion, CmdVersionInfo)

	ret += ux.SprintfWhite("\nExamples for: %s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionLatest)
	ret += ux.SprintfBlue(" - Show the latest version of the '%s' binary.\n", defaults.BinaryName)
	ret += ux.SprintfMagenta("\t%s %s %s\n", defaults.BinaryName, CmdVersion, CmdVersionLatest)
	ret += ux.SprintfBlue(" - Show the latest version within the buildtool repo.\n")
	ret += ux.SprintfMagenta("\t%s %s %s gearboxworks/buildtool\n", defaults.BinaryName, CmdVersion, CmdVersionLatest)

	ret += ux.SprintfWhite("\nSymlinking methods:\n")
	ret += ux.SprintfBlue(" - Show the latest version of buildtool.\n")
	ret += ux.SprintfMagenta("\tln -s %s ./buildtool\n", RunTime.Cmd)
	ret += ux.SprintfMagenta("\t./buildtool %s %s\n", CmdVersion, CmdVersionInfo)
	ret += ux.SprintfBlue(" - Update to the latest version of buildtool, (no args will automatically update).\n")
	ret += ux.SprintfMagenta("\tln -s %s ./buildtool\n", RunTime.Cmd)
	ret += ux.SprintfMagenta("\t./buildtool\n")

	ret += ux.SprintfWhite("\n")

	return ret
}
