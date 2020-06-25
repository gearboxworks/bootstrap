package cmd

import (
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
)


func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
	rootCmd.AddCommand(linksCmd)

	versionCmd.AddCommand(versionCheckCmd)
	versionCmd.AddCommand(versionListCmd)
	versionCmd.AddCommand(versionInfoCmd)
	versionCmd.AddCommand(versionLatestCmd)
	versionCmd.AddCommand(versionUpdateCmd)
	versionCmd.AddCommand(versionExamplesCmd)
}


var versionCmd = &cobra.Command{
	Use:   CmdVersion,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Self-manage executable."),
	Long:  ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Self-manage executable."),
	Run: Version,
}
var selfUpdateCmd = &cobra.Command{
	Use:   CmdSelfUpdate,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Update version of executable."),
	Long:  ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Check and update the latest version."),
	Run: VersionUpdate,
}
var linksCmd = &cobra.Command{
	Use:   CmdLinks,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Create symlinks for all supported apps."),
	Long:  ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Create symlinks for all supported apps."),
	Run: VersionLinks,
}

var versionUpdateCmd = &cobra.Command{
	Use:   CmdVersionUpdate,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Update version of executable."),
	Long: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Check and update the latest version."),
	Run: VersionUpdate,
	Example: ux.SprintfMagenta("%s %s", CmdVersion, CmdVersionUpdate) + ux.SprintfBlue(" - List all available versions of this binary.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool", CmdVersion, CmdVersionUpdate) + ux.SprintfBlue(" - Update to the latest version within the buildtool repo.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool/latest", CmdVersion, CmdVersionUpdate) + ux.SprintfBlue(" - Update to the latest version within the buildtool repo.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool/1.1.3", CmdVersion, CmdVersionUpdate) + ux.SprintfBlue(" - Update to version 1.1.3 within the buildtool repo.\n"),
}
var versionCheckCmd = &cobra.Command{
	Use:   CmdVersionCheck,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Show any version updates."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Show any version updates."),
	Run: VersionCheck,
	Example: ux.SprintfMagenta("%s %s", CmdVersion, CmdVersionCheck) + ux.SprintfBlue(" - Check the latest version for this binary.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool", CmdVersion, CmdVersionCheck) + ux.SprintfBlue(" - Check the latest version within the buildtool repo.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool/latest", CmdVersion, CmdVersionCheck) + ux.SprintfBlue(" - Check the latest version within the buildtool repo.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool/1.1.3", CmdVersion, CmdVersionCheck) + ux.SprintfBlue(" - Check version 1.1.3 within the buildtool repo.\n"),
}
var versionListCmd = &cobra.Command{
	Use:   CmdVersionList,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Self-manage executable."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Self-manage executable."),
	Run: VersionList,
	Example: ux.SprintfMagenta("%s %s", CmdVersion, CmdVersionList) + ux.SprintfBlue(" - List all available versions of this binary.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool", CmdVersion, CmdVersionList) + ux.SprintfBlue(" - List all available versions within the buildtool repo.\n"),
}
var versionInfoCmd = &cobra.Command{
	Use:   CmdVersionInfo,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on current version."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on current version."),
	Run: VersionInfo,
	Example: ux.SprintfMagenta("%s %s", CmdVersion, CmdVersionInfo) + ux.SprintfBlue(" - Show info on the current version of this binary.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool", CmdVersion, CmdVersionInfo) + ux.SprintfBlue(" - Show info on the latest version within the buildtool repo.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool/latest", CmdVersion, CmdVersionInfo) + ux.SprintfBlue(" - Show info on the latest version within the buildtool repo.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool/1.1.3", CmdVersion, CmdVersionInfo) + ux.SprintfBlue(" - Show info on version 1.1.3 within the buildtool repo.\n"),
}
var versionLatestCmd = &cobra.Command{
	Use:   CmdVersionLatest,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
	Run: VersionLatest,
	Example: ux.SprintfMagenta("%s %s", CmdVersion, CmdVersionLatest) + ux.SprintfBlue(" - Show the latest version of this binary.\n") +
		ux.SprintfMagenta("%s %s gearboxworks/buildtool", CmdVersion, CmdVersionLatest) + ux.SprintfBlue(" - Show the latest version within the buildtool repo.\n"),
}
var versionExamplesCmd = &cobra.Command{
	Use:   CmdVersionExamples,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Show examples."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Show examples."),
	Run: VersionExamples,
}
