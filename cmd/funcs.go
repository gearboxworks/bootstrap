package cmd

import (
	"errors"
	"fmt"
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/toolSelfUpdate"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)


func CheckArgs(self *toolSelfUpdate.TypeSelfUpdate, args ...string) *ux.State {
	for range onlyOnce {
		//ux.PrintflnBlue("Creating all supported application links.")

		ux.PrintflnBlue("IsRunningAsFile: %v", self.Runtime.IsRunningAsFile())
		ux.PrintflnBlue("IsRunningAsLink: %v", self.Runtime.IsRunningAsLink())
		ux.PrintflnBlue("IsRunningAs: %v", self.Runtime.IsRunningAs(defaults.BinaryName))

		if self.IsBootstrapBinary() {
			// Exit as we are the bootstrap binary.
			self.State.SetOk()
			break
		}

		if self.Runtime.IsRunningAsLink() {
			self.CreateDummyBinary()
			if self.State.IsNotOk() {
				break
			}
		}

		// Lookup binary name...
		repoUrl := defaults.Available.GetRepo(self.Runtime.CmdFile)
		if repoUrl == "" {
			self.State.SetError("Binary '%s' has no known repo.", self.Runtime.CmdFile)
			break
		}

		self.SetBinaryRepo(repoUrl)
		if self.State.IsNotOk() {
			break
		}
	}

	return self.State
}


func CheckRunTime(self *toolSelfUpdate.TypeSelfUpdate, args ...string) *ux.State {
	for range onlyOnce {
		if self.Runtime.Debug {
			ux.PrintflnBlue("IsRunningAsFile: %v", self.Runtime.IsRunningAsFile())
			ux.PrintflnBlue("IsRunningAsLink: %v", self.Runtime.IsRunningAsLink())
			ux.PrintflnBlue("IsRunningAs: %v", self.Runtime.IsRunningAs(defaults.BinaryName))
			ux.PrintflnBlue("IsBootstrapBinary: %v", self.IsBootstrapBinary())
		}

		if self.IsBootstrapBinary() {
			// Exit as we are the bootstrap binary.
			self.State.SetOk()
			break
		}

		if self.Runtime.IsRunningAsLink() {
			self.AutoExec = true
			self.CreateDummyBinary()
			if self.State.IsNotOk() {
				break
			}
		}

		// Lookup binary name...
		repoUrl := defaults.Available.GetRepo(self.Runtime.CmdFile)
		if repoUrl == "" {
			self.State.SetError("Binary '%s' has no known repo.", self.Runtime.CmdFile)
			break
		}

		self.SetSourceRepo(repoUrl)
		if self.State.IsNotOk() {
			break
		}

		self.SetBinaryRepo(repoUrl)
		if self.State.IsNotOk() {
			break
		}
	}

	return self.State
}


//func (rt *TypeUpdate) Version(cmd *cobra.Command) error {
//	err := rt.VersionShow()
//	SetHelp(cmd)
//	PrintHelp()
//	err = cmd.Help()
//	return err
//}
//
//
//func (rt *TypeUpdate) VersionShow() error {
//	ux.PrintfBlue("%s ", defaults.BinaryName)
//	ux.PrintflnCyan("v%s", defaults.BinaryVersion)
//	return nil
//}
//
//
//func (rt *TypeUpdate) VersionInfo() error {
//	var err error
//	for range onlyOnce {
//		err = rt.Update.PrintVersion(CurrentVersion)
//		if err != nil {
//			break
//		}
//	}
//	return err
//}
//
//
//func (rt *TypeUpdate) VersionList() error {
//	var err error
//	for range onlyOnce {
//		token := os.Getenv("GITHUB_TOKEN")
//		if token == "" {
//			token, _ = gitconfig.GithubToken()
//		}
//
//		gh := github.NewClient(newHTTPClient(context.Background(), token))
//		var rels []*github.RepositoryRelease
//		rels, _, err = gh.Repositories.ListReleases(context.Background(), rt.Update.owner.String(), rt.Update.name.String(), nil)
//		if err != nil {
//			break
//		}
//
//		for _, rel := range rels {
//			err = rt.Update.PrintVersionSummary(*rel.TagName)
//			if err != nil {
//				break
//			}
//		}
//	}
//	return err
//}
//
//
//func (rt *TypeUpdate) VersionCheck() error {
//	var err error
//	for range onlyOnce {
//		err = rt.Update.IsUpdated()
//		if err != nil {
//			break
//		}
//	}
//	return err
//}
//
//
//func (rt *TypeUpdate) VersionUpdate() error {
//	var err error
//	for range onlyOnce {
//		err = CreateDummyBinary(SelfUpdate.Cmd, SelfUpdate.Cmd)
//		if err != nil {
//			break
//		}
//
//		err = rt.Update.IsUpdated()
//		if err != nil {
//			break
//		}
//
//		err = rt.Update.UpdateTo()
//		if err != nil {
//			break
//		}
//
//		if !SelfUpdate.AutoExec {
//			break
//		}
//
//		// AutoExec will execute the new binary with the same args as given.
//		err = Run(SelfUpdate.Cmd, SelfUpdate.CmdArgs...)
//		if err != nil {
//			break
//		}
//	}
//	return err
//}


func Version(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		SelfUpdate.State = CheckRunTime(SelfUpdate, args...)
		if SelfUpdate.State.IsNotOk() {
			break
		}

		SelfUpdate.State = SelfUpdate.Version(cmd, args...)
		if SelfUpdate.State.IsNotOk() {
			return
		}
	}
}


func VersionLinks(cmd *cobra.Command, args []string) {
	var err error

	for range onlyOnce {
		ux.PrintflnBlue("Creating all supported application links.")

		bins := defaults.Available.GetBinaries()
		err = os.Chdir(RunTime.CmdDir)

		links := make(map[string]string)
		var failed bool
		for _, binaryFile := range bins {
			var err error
			var linkStat os.FileInfo

			linkStat, err = os.Lstat(binaryFile)
			if linkStat == nil {
				// Symlink doesn't exist - create.
				err = os.Symlink(RunTime.CmdFile, binaryFile)
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
			if fpel != RunTime.CmdFile {
				links[binaryFile] = "incorrect link"
				failed = true
				continue
			}

			links[binaryFile] = "exists"
		}

		if failed {
			err = errors.New("failed to add all application links")
			ux.PrintflnWarning("Failed to add all application links.")
			//break
		}
		for k, v := range links {
			ux.PrintflnCyan("%s    \t- %s - from repo %s/%s", k, v, defaults.BinaryRepoPrefix, k)
		}
	}

	if err != nil {
		SelfUpdate.State.SetError(err)
	}
}


func VersionUpdate(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		SelfUpdate.State = CheckRunTime(SelfUpdate, args...)
		if SelfUpdate.State.IsNotOk() {
			break
		}

		ux.PrintflnWarning("The binary '%s' will be installed from the '%s' repo...", SelfUpdate.Runtime.CmdFile, SelfUpdate.GetBinaryRepo())
		SelfUpdate.State = SelfUpdate.VersionUpdate()
		if SelfUpdate.State.IsNotOk() {
			break
		}

		if !SelfUpdate.AutoExec {
			break
		}

		// AutoExec will execute the new binary with the same args as given.
		SelfUpdate.State = SelfUpdate.AutoRun()
		if SelfUpdate.State.IsNotOk() {
			break
		}
	}
}


func VersionCheck(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		SelfUpdate.State = CheckRunTime(SelfUpdate, args...)
		if SelfUpdate.State.IsNotOk() {
			break
		}

		SelfUpdate.State = SelfUpdate.VersionCheck()
		if SelfUpdate.State.IsNotOk() {
			return
		}
	}
}


func VersionList(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		SelfUpdate.State = CheckRunTime(SelfUpdate, args...)
		if SelfUpdate.State.IsNotOk() {
			break
		}

		SelfUpdate.State = SelfUpdate.VersionList(args...)
		if SelfUpdate.State.IsNotOk() {
			return
		}
	}
}


func VersionInfo(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		SelfUpdate.State = CheckRunTime(SelfUpdate, args...)
		if SelfUpdate.State.IsNotOk() {
			break
		}

		if len(args) == 0 {
			args = []string{SelfUpdate.GetVersion()}
		}

		SelfUpdate.State = SelfUpdate.VersionInfo(args...)
		if SelfUpdate.State.IsNotOk() {
			return
		}
	}
}


func VersionLatest(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		SelfUpdate.State = CheckRunTime(SelfUpdate, args...)
		if SelfUpdate.State.IsNotOk() {
			break
		}

		SelfUpdate.State = SelfUpdate.VersionInfo(LatestVersion)
		if SelfUpdate.State.IsNotOk() {
			return
		}
	}
}


func VersionExamples(cmd *cobra.Command, args []string) {
	fmt.Print(HelpExamples())
}
