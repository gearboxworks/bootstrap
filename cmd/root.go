package cmd

import (
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"strings"
)


const onlyOnce = "1"
//var onlyTwice = []string{"", ""}
var Runtime TypeRuntime
var Target TypeRuntime


func init() {
	Runtime = NewRuntime(defaults.BinaryName, defaults.BinaryVersion, defaults.SourceRepo, defaults.BinaryRepo, false)
	Target = NewRuntime("", "", defaults.SourceRepoPrefix, defaults.BinaryRepoPrefix, false)

	//rootCmd.PersistentFlags().StringVarP(&Runtime.Json.Filename, loadTools.FlagJsonFile, "j", loadTools.DefaultJsonFile, ux.SprintfBlue("Alternative JSON file."))
	//rootCmd.PersistentFlags().StringVarP(&Runtime.Template.Filename, loadTools.FlagTemplateFile, "t", loadTools.DefaultTemplateFile, ux.SprintfBlue("Alternative template file."))
	//rootCmd.PersistentFlags().StringVarP(&Runtime.Output.Filename, loadTools.FlagOutputFile, "o", loadTools.DefaultOutFile, ux.SprintfBlue("Output file."))
	//rootCmd.PersistentFlags().StringVarP(&Runtime.WorkingPath.Filename, loadTools.FlagWorkingPath, "p", loadTools.DefaultWorkingPath, ux.SprintfBlue("Set working path."))
	//
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.Chdir, loadTools.FlagChdir, "c", false, ux.SprintfBlue("Change to directory containing %s", loadTools.DefaultJsonFile))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.RemoveTemplate, loadTools.FlagRemoveTemplate, "", false, ux.SprintfBlue("Remove template file afterwards."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.ForceOverwrite, loadTools.FlagForce, "f", false, ux.SprintfBlue("Force overwrite of output files."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.RemoveOutput, loadTools.FlagRemoveOutput, "", false, ux.SprintfBlue("Remove output file afterwards."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.QuietProgress, loadTools.FlagQuiet, "q", false, ux.SprintfBlue("Silence progress in shell scripts."))
	//
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.Debug, loadTools.FlagDebug ,"d", false, ux.SprintfBlue("DEBUG mode."))

	rootCmd.PersistentFlags().StringVarP(&Target.CmdName, "bin" ,"b", Target.CmdName, ux.SprintfBlue("Name of target binary to download."))
	_ = rootCmd.PersistentFlags().MarkHidden("bin")
	rootCmd.PersistentFlags().StringVarP(&Target.CmdBinaryRepo, "repo" ,"r", Target.CmdBinaryRepo, ux.SprintfBlue("Url of target binary repo to download."))
	_ = rootCmd.PersistentFlags().MarkHidden("repo")
	rootCmd.PersistentFlags().StringVarP(&Target.CmdVersion, "ver" ,"", Target.CmdVersion, ux.SprintfBlue("Version of target binary to download."))
	_ = rootCmd.PersistentFlags().MarkHidden("ver")
	rootCmd.PersistentFlags().BoolVarP(&Target.AutoExec, "auto" ,"", false, ux.SprintfBlue("Auto-update when symlinked."))
	_ = rootCmd.PersistentFlags().MarkHidden("ver")

	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpAll, loadTools.FlagHelpAll, "", false, ux.SprintfBlue("Show all help."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpVariables, loadTools.FlagHelpVariables, "", false, ux.SprintfBlue("Help on template variables."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpFunctions, loadTools.FlagHelpFunctions, "", false, ux.SprintfBlue("Help on template functions."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpExamples, loadTools.FlagHelpExamples, "", false, ux.SprintfBlue("Help on template examples."))

	rootCmd.Flags().BoolP(FlagVersion, "v", false, ux.SprintfBlue("Display version of " + defaults.BinaryName))
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
		var ok bool
		fl := cmd.Flags()

		// Show version.
		ok, _ = fl.GetBool(FlagVersion)
		if ok {
			Runtime.Error = VersionShow()
			break
		}

		if Runtime.IsSymLinked {
			ux.PrintflnWarning("This binary will be auto-updated from the '%s' repo...", Target.CmdBinaryRepo)
			// Assume a 'version update'
			Runtime.Error = Target.SetApp(&Runtime, args...)
			if Runtime.Error != nil {
				return
			}
			Runtime.AutoExec = true
			Target.AutoExec = true
			Runtime.Error = VersionUpdate()
			break
		}

		// Show help if no commands specified.
		if len(args) == 0 {
			PrintHelp()
			Runtime.Error = cmd.Help()
			break
		}
	}
}


func Execute() error {
	for range onlyOnce {
		SetHelp(rootCmd)

		err := rootCmd.Execute()
		if err == nil {
			break
		}

		if !strings.HasPrefix(err.Error(), "unknown command") {
			//PrintHelp()
			//Runtime.Error = rootCmd.Help()
			break
		}

		gbRootFunc(rootCmd, []string{})

		//// Assume a 'version update'
		//Runtime.AutoExec = true
		//Target.AutoExec = true
		//
		//Runtime.Error = Target.SetApp(&Runtime)
		//if Runtime.Error == nil {
		//	Runtime.Error = VersionUpdate()
		//}
		//Runtime.Error = nil
		//break
	}

	return Runtime.Error
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
	cobra.AddTemplateFunc("VersionExamples", VersionExamples)

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
