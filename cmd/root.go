package cmd

import (
	"errors"
	"fmt"
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/kardianos/osext"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)


type TypeRuntime struct {
	CmdName       string    `json:"cmd_name" mapstructure:"cmd_name"`
	CmdVersion    string    `json:"cmd_version" mapstructure:"cmd_version"`
	CmdSourceRepo string    `json:"cmd_source_repo" mapstructure:"cmd_source_repo"`
	CmdBinaryRepo string    `json:"cmd_binary_repo" mapstructure:"cmd_binary_repo"`
	Cmd           string    `json:"cmd" mapstructure:"cmd"`
	CmdDir        string    `json:"cmd_dir" mapstructure:"cmd_dir"`
	CmdFile       string    `json:"cmd_file" mapstructure:"cmd_file"`

	GoRuntime     GoRuntime `json:"go_runtime" mapstructure:"go_runtime"`

	Debug         bool

	Error         error
}
type GoRuntime struct {
	Os string
	Arch string
	Root string
	Version string
	Compiler string
	NumCpus int
}

var Runtime TypeRuntime
var Target TypeRuntime


func init() {
	Runtime = NewRuntime(defaults.BinaryName, defaults.BinaryVersion, defaults.SourceRepo, defaults.BinaryRepo, false)
	Target = NewRuntime("", "", "", "", false)

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

	rootCmd.PersistentFlags().StringVarP(&Target.CmdName, "bin" ,"b", Target.CmdName, ux.SprintfBlue("Name of binary to download."))
	rootCmd.PersistentFlags().StringVarP(&Target.CmdBinaryRepo, "repo" ,"r", Target.CmdBinaryRepo, ux.SprintfBlue("Url of binary repo to download."))

	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpAll, loadTools.FlagHelpAll, "", false, ux.SprintfBlue("Show all help."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpVariables, loadTools.FlagHelpVariables, "", false, ux.SprintfBlue("Help on template variables."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpFunctions, loadTools.FlagHelpFunctions, "", false, ux.SprintfBlue("Help on template functions."))
	//rootCmd.PersistentFlags().BoolVarP(&Runtime.HelpExamples, loadTools.FlagHelpExamples, "", false, ux.SprintfBlue("Help on template examples."))

	rootCmd.Flags().BoolP(FlagVersion, "v", false, ux.SprintfBlue("Display version of " + defaults.BinaryName))
}


func NewRuntime(BinaryName string, BinaryVersion string, SourceRepo string, BinaryRepo string, debugFlag bool) TypeRuntime {
	r := TypeRuntime {
		//RunAsBootStrap: false,

		CmdName:       BinaryName,
		CmdVersion:    BinaryVersion,
		CmdSourceRepo: SourceRepo,
		CmdBinaryRepo: BinaryRepo,
		Cmd:           "",
		CmdDir:        "",
		CmdFile:       "",

		GoRuntime: GoRuntime{
			Os:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Root:     runtime.GOROOT(),
			Version:  runtime.Version(),
			Compiler: runtime.Compiler,
			NumCpus:  runtime.NumCPU(),
		},

		Debug: debugFlag,
		Error: nil,
	}

	for range onlyOnce {
		exe, err := osext.Executable()
		if err != nil{
			r.Error = err
			break
		}
		r.Cmd = exe
		r.CmdDir = path.Dir(exe)
		r.CmdFile = path.Base(exe)
		r.Debug = debugFlag

		if r.CmdName == "" {
			r.CmdName = r.CmdFile
		}

		if r.CmdSourceRepo == "" {
			r.CmdSourceRepo = "github.com/gearboxworks/" + r.CmdName
		}

		if r.CmdBinaryRepo == "" {
			r.CmdBinaryRepo = "github.com/gearboxworks/" + r.CmdName
		}

		if r.CmdFile == defaults.BinaryName {
			//r.RunAsBootStrap = true
			break
		}
	}

	return r
}


func SetApp(runtime *TypeRuntime, target *TypeRuntime) error {
	for range onlyOnce {
		if target.CmdVersion == "" {
			target.CmdVersion = runtime.CmdVersion
		}

		if target.CmdSourceRepo == runtime.CmdSourceRepo {
			target.CmdSourceRepo = "github.com/gearboxworks/" + target.CmdName
		}

		if target.CmdBinaryRepo == runtime.CmdBinaryRepo {
			target.CmdBinaryRepo = "github.com/gearboxworks/" + target.CmdName
		}

		if target.Cmd == runtime.Cmd {
			target.Cmd = filepath.Join(filepath.Dir(target.Cmd), target.CmdName)
		}

		if runtime.CmdFile != defaults.BinaryName {
			// We are running the binary either as a symlink or filename other than 'bootstrap'.
		}
	}

	return target.Error
}


func CreateDummyBinary(runtimeBin string, targetBin string) error {
	var err error

	for range onlyOnce {
		var link string

		link, err = os.Readlink(targetBin)
		if link == "" {
			if targetBin == runtimeBin {
				err = nil
				break
			}
			_, err = os.Stat(targetBin)
			if os.IsNotExist(err) {
				// File doesn't exist - need to create it.
				err = CopyFile(runtimeBin, targetBin)
			}

			break
		}

		_, err = os.Stat(targetBin)
		if os.IsNotExist(err) {
			if link != "" {
				// File is a dud - no action.
				err = errors.New(fmt.Sprintf("file '%s' is a symlink pointing to non-existant file '%s'", targetBin, link))
				break
			}

			// File doesn't exist - need to create it.
			err = CopyFile(runtimeBin, targetBin)
			break
		}

		if link != defaults.BinaryName {
			err = errors.New("symlink not pointing to bootstrap")
			break
		}

		runtimeBin = filepath.Join(filepath.Dir(targetBin), link)
		ux.PrintflnOk("Removing symlink %s (%s)", targetBin, link)
		err = os.Remove(targetBin)
		if err != nil {
			break
		}

		err = CopyFile(runtimeBin, targetBin)
	}

	return err
}


func CopyFile(runtimeBin string, targetBin string) error {
	var err error

	for range onlyOnce {
		var input []byte
		input, err = ioutil.ReadFile(runtimeBin)
		if err != nil {
			break
		}

		err = ioutil.WriteFile(targetBin, input, 0755)
		if err != nil {
			fmt.Println("Error creating", targetBin)
			break
		}
	}

	return err
}


func CompareBinary(runtimeBin string, newBin string) error {
	var err error

	for range onlyOnce {
		var srcBin []byte
		srcBin, err = ioutil.ReadFile(runtimeBin)
		if err != nil {
			break
		}
		if srcBin == nil {
			break
		}

		var targetBin []byte
		targetBin, err = ioutil.ReadFile(newBin)
		if err != nil {
			break
		}
		if targetBin == nil {
			break
		}

		if len(srcBin) != len(targetBin) {
			break
		}

		err = errors.New("binary files differ")
		for i := range srcBin {
			if srcBin[i] != targetBin[i] {
				err = nil
				break
			}
		}
	}

	return err
}


var rootCmd = &cobra.Command{
	Use:   defaults.BinaryName,
	Short: "Bootstrap is the automatic app downloader.",
	Long: "Bootstrap is the automatic app downloader.",
	Run: gbRootFunc,
}
func gbRootFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ok bool
		fl := cmd.Flags()

		// ////////////////////////////////
		// Show version.
		ok, _ = fl.GetBool(FlagVersion)
		if ok {
			Runtime.Error = VersionShow()
			break
		}

		// Show help if no commands specified.
		if len(args) == 0 {
			PrintHelp(cmd)
			Runtime.Error = cmd.Help()
			break
		}
	}
}


func Execute() error {
	for range onlyOnce {
		SetHelp(rootCmd)

		err := rootCmd.Execute()

		if err != nil {
			if strings.HasPrefix(err.Error(), "unknown command") {
				PrintHelp(rootCmd)
				err = nil
			}

			_ = rootCmd.Help()
			Runtime.Error = err
			break
		}
	}

	return Runtime.Error
}
