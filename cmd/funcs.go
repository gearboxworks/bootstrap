package cmd

import (
	"context"
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/google/go-github/v30/github"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"os"
)


func Version(cmd *cobra.Command) error {
	err := VersionShow()
	SetHelp(cmd)
	PrintHelp()
	err = cmd.Help()
	return err
}


func VersionShow() error {
	ux.PrintfBlue("%s ", defaults.BinaryName)
	ux.PrintflnCyan("v%s", defaults.BinaryVersion)
	return nil
}


func VersionInfo() error {
	update := New(&Target)

	for range onlyOnce {
		if update.Error != nil {
			break
		}

		update.Error = update.PrintVersion(update.version.String())
		if update.Error != nil {
			break
		}
	}

	return update.Error
}


func VersionList() error {
	update := New(&Target)

	for range onlyOnce {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			token, _ = gitconfig.GithubToken()
		}
		gh := github.NewClient(newHTTPClient(context.Background(), token))
		var rels []*github.RepositoryRelease
		rels, _, update.Error = gh.Repositories.ListReleases(context.Background(), update.owner.String(), update.name.String(), nil)
		if update.Error != nil {
			break
		}

		for _, rel := range rels {
			update.Error = update.PrintVersionSummary(*rel.TagName)
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

		update.Error = update.UpdateTo()
		if update.Error != nil {
			break
		}

		if !Target.AutoExec {
			break
		}

		// AutoExec will execute the new binary with the same args as given.
		update.Error = Run(Target.Cmd, Target.CmdArgs...)
		if update.Error != nil {
			break
		}
	}

	return update.Error
}


func VersionExamples() string {
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
	ret += ux.SprintfMagenta("\tln -s %s ./buildtool\n", Runtime.Cmd)
	ret += ux.SprintfMagenta("\t./buildtool %s %s\n", CmdVersion, CmdVersionInfo)
	ret += ux.SprintfBlue(" - Update to the latest version of buildtool, (no args will automatically update).\n")
	ret += ux.SprintfMagenta("\tln -s %s ./buildtool\n", Runtime.Cmd)
	ret += ux.SprintfMagenta("\t./buildtool\n")

	ret += ux.SprintfWhite("\n")

	return ret
}
