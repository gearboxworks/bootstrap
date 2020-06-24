package cmd

import (
	"context"
	"errors"
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/google/go-github/v30/github"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"os"
	"path/filepath"
)


func Links() error {
	update := New(&Target)

	for range onlyOnce {
		ux.PrintflnBlue("Creating all supported application links.")

		bins := defaults.Available.GetBinaries()
		update.Error = os.Chdir(Runtime.CmdDir)

		links := make(map[string]string)
		var failed bool
		for _, binaryFile := range bins {
			var err error
			var linkStat os.FileInfo

			linkStat, err = os.Lstat(binaryFile)
			if linkStat == nil {
				// Symlink doesn't exist - create.
				err = os.Symlink(Runtime.CmdFile, binaryFile)
				if err != nil {
					continue
				}

				linkStat, err = os.Lstat(binaryFile)
				if linkStat == nil {
					continue
				}

				links[binaryFile] = "created"
			}
			//fmt.Printf("'%s' (%s) => '%s'\n", k, binaryFile, binaryFile)
			//fmt.Printf("\tReadlink() => %s\n", l)
			//fmt.Printf("\tLstat() => %s	%s	%s	%s	%d\n",
			//	linkStat.Name(),
			//	linkStat.IsDir(),
			//	linkStat.Mode().String(),
			//	linkStat.ModTime().String(),
			//	linkStat.Size(),
			//)

			// Symlink exists - validate.
			l, _ := os.Readlink(binaryFile)
			if l == defaults.BinaryName {
			}

			fpel, err := filepath.EvalSymlinks(binaryFile)
			//fmt.Printf("%s\n", fpel)
			if fpel != Runtime.CmdFile {
				links[binaryFile] = "incorrect link"
				failed = true
				continue
			}

			links[binaryFile] = "exists"
		}

		if failed {
			update.Error = errors.New("failed to add all application links")
			ux.PrintflnWarning("Failed to add all application links.")
			//break
		}
		for k, v := range links {
			ux.PrintflnCyan("%s    \t- %s - from repo %s/%s", k, v, defaults.BinaryRepoPrefix, k)
		}
	}

	return update.Error
}


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
