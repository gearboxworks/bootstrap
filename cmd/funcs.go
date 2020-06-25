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
	"strings"
)


func CheckRunTime(self *toolSelfUpdate.TypeSelfUpdate, args ...string) *ux.State {
	for range onlyOnce {
		self.SetOldVersion(toolSelfUpdate.EarliestSemVer)
		name := strings.Join(args, "")

		if !self.IsBootstrapBinary() {
			// Lookup binary name...
			repoUrl := defaults.Available.GetRepo(self.Runtime.CmdFile)
			if repoUrl != "" {
				self.State = self.SetSourceRepo(repoUrl)
				if self.State.IsNotOk() {
					break
				}

				self.State = self.SetBinaryRepo(repoUrl)
				if self.State.IsNotOk() {
					break
				}

				self.State = self.SetVersion(name)
				if self.State.IsNotOk() {
					break
				}
			}
			break
		}

		self.State = self.SetRepo(args...)
		if self.State.IsNotOk() {
			// Lookup binary name...
			repoUrl := defaults.Available.GetRepo(name)
			self.State = self.SetRepo(repoUrl)
			if self.State.IsNotOk() {
				self.State.SetError("No known repo for binary '%s' at repo '%s'.", name, self.GetRepo())
				break
			}
		}

		//ux.PrintflnBlue("IsRunningAsFile: %v", self.Runtime.IsRunningAsFile())
		//ux.PrintflnBlue("IsRunningAsLink: %v", self.Runtime.IsRunningAsLink())
		//ux.PrintflnBlue("IsRunningAs: %v", self.Runtime.IsRunningAs(defaults.BinaryName))
		//ux.PrintflnBlue("Repo: %v", self.GetRepo())
		self.State.SetOk()
	}

	return self.State
}


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
		ux.PrintflnBlue("Installing placeholder application links.")

		bins := defaults.Available.GetBinaries()
		err = os.Chdir(RunTime.CmdDir)

		links := make(map[string]string)
		var failed bool
		for binaryFile, binaryUrl := range bins {
			var err error
			var linkStat os.FileInfo

			binaryFile = filepath.Join(RunTime.CmdDir, binaryFile)
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
				ux.PrintflnOk("%s    \t- %s - (repository %s)", filepath.Base(binaryFile), links[binaryFile], binaryUrl)
				continue
			}
			//fmt.Printf("'%s' (%s) => '%s'\n", k, binaryFile, binaryFile)

			// Symlink exists - validate.
			l, _ := os.Readlink(binaryFile)
			if l == defaults.BinaryName {
			}

			var link string
			link, err = filepath.EvalSymlinks(binaryFile)
			//fmt.Printf("%s\n", link)
			if link != RunTime.Cmd {
				links[binaryFile] = "incorrect link"
				ux.PrintflnError("%s    \t- %s - (repository %s)", filepath.Base(binaryFile), links[binaryFile], binaryUrl)
				failed = true
				continue
			}

			links[binaryFile] = "exists"
			ux.PrintflnOk("%s    \t- %s - (repository %s)", filepath.Base(binaryFile), links[binaryFile], binaryUrl)
			fmt.Printf("")
		}

		if failed {
			err = errors.New("Failed to install some applications")
			ux.PrintflnWarning("%s", err)
		}

		//for k, v := range links {
		//	ux.PrintflnCyan("%s    \t- %s - from repo %s/%s", k, v, defaults.BinaryRepoPrefix, k)
		//}
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

		//SelfUpdate.State = SelfUpdate.CreateDummyBinary()
		//if SelfUpdate.State.IsNotOk() {
		//	break
		//}

		ux.PrintflnWarning("The binary '%s' will be installed from the '%s' repo...", SelfUpdate.GetName(), SelfUpdate.GetBinaryRepo())
		SelfUpdate.State = SelfUpdate.VersionUpdate()
		if SelfUpdate.State.IsNotOk() {
			break
		}

		//if !SelfUpdate.AutoExec {
		//	break
		//}
		//
		//// AutoExec will execute the new binary with the same args as given.
		//SelfUpdate.State = SelfUpdate.AutoRun()
		//if SelfUpdate.State.IsNotOk() {
		//	break
		//}
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

		SelfUpdate.State = SelfUpdate.VersionInfo(SelfUpdate.GetVersion())
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
	_ = rootCmd.Help()
	fmt.Print(HelpExamples())
}
