package cmd

import (
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
)


const onlyOnce = "1"
var onlyTwice = []string{"", ""}


func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)

	versionCmd.AddCommand(versionCheckCmd)
	versionCmd.AddCommand(versionListCmd)
	versionCmd.AddCommand(versionInfoCmd)
	versionCmd.AddCommand(versionLatestCmd)
	versionCmd.AddCommand(versionUpdateCmd)
}


var versionCmd = &cobra.Command{
	Use:   CmdVersion,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Self-manage executable."),
	Long:  ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Self-manage executable."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = Version(cmd, args...)
	},
}
var selfUpdateCmd = &cobra.Command{
	Use:   CmdSelfUpdate,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Update version of executable."),
	Long: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Check and update the latest version."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = VersionUpdate()
	},
}


var versionCheckCmd = &cobra.Command{
	Use:   CmdVersionCheck,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Show any version updates."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Show any version updates."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = VersionCheck()
	},
}
var versionListCmd = &cobra.Command{
	Use:   CmdVersionList,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Self-manage executable."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Self-manage executable."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = VersionList(args...)
	},
}
var versionInfoCmd = &cobra.Command{
	Use:   CmdVersionInfo,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on current version."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on current version."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = VersionInfo(args...)
	},
}
var versionLatestCmd = &cobra.Command{
	Use:   CmdVersionLatest,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
	Long:  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = VersionInfo(CmdVersionLatest)
	},
}
var versionUpdateCmd = &cobra.Command{
	Use:   CmdVersionUpdate,
	Short: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Update version of executable."),
	Long: ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Check and update the latest version."),
	Run: func(cmd *cobra.Command, args []string) {
		Runtime.Error = SetApp(&Runtime, &Target)
		Runtime.Error = VersionUpdate()
	},
}
