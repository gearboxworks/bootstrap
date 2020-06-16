package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/kardianos/osext"
	"github.com/newclarity/scribeHelpers/ux"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)


type TypeRuntime struct {
	WantVersion string `json:"want_version" mapstructure:"want_version"`
	AutoExec    bool   `json:"auto_update" mapstructure:"auto_update"`
	IsSymLinked bool   `json:"is_symlinked" mapstructure:"is_symlinked"`

	CmdName       string    `json:"cmd_name" mapstructure:"cmd_name"`
	CmdVersion    string    `json:"cmd_version" mapstructure:"cmd_version"`
	CmdSourceRepo string    `json:"cmd_source_repo" mapstructure:"cmd_source_repo"`
	CmdBinaryRepo string    `json:"cmd_binary_repo" mapstructure:"cmd_binary_repo"`
	Cmd           string    `json:"cmd" mapstructure:"cmd"`
	CmdDir        string    `json:"cmd_dir" mapstructure:"cmd_dir"`
	CmdFile       string    `json:"cmd_file" mapstructure:"cmd_file"`

	CmdArgs       []string  `json:"cmd_args" mapstructure:"cmd_args"`

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


func NewRuntime(BinaryName string, BinaryVersion string, SourceRepo string, BinaryRepo string, debugFlag bool) TypeRuntime {
	r := TypeRuntime {
		WantVersion: "",
		AutoExec:    false,
		IsSymLinked: false,

		CmdName:       BinaryName,
		CmdVersion:    BinaryVersion,
		CmdSourceRepo: SourceRepo,
		CmdBinaryRepo: BinaryRepo,
		Cmd:           "",
		CmdDir:        "",
		CmdFile:       "",

		CmdArgs:       os.Args[1:],

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

		if r.CmdFile != defaults.BinaryName {
			r.IsSymLinked = true
		}

		if r.CmdName == "" {
			r.CmdName = r.CmdFile
		}

		if r.CmdSourceRepo == "" {
			r.CmdSourceRepo = "github.com/gearboxworks/" + r.CmdName
		}

		if r.CmdBinaryRepo == "" {
			r.CmdBinaryRepo = "github.com/gearboxworks/" + r.CmdName
		}
	}

	return r
}


func (target *TypeRuntime) SetApp(runtime *TypeRuntime, args ...string) error {
	for range onlyOnce {
		if target.WantVersion == LatestVersion {
			target.WantVersion = ""
		}
		if target.CmdVersion == LatestVersion {
			target.CmdVersion = runtime.CmdVersion
		} else if target.CmdVersion == "" {
			target.CmdVersion = runtime.CmdVersion
		}

		target.CmdSourceRepo = target.CmdSourceRepo + "/" + target.CmdName
		target.CmdBinaryRepo = target.CmdBinaryRepo + "/" + target.CmdName

		if target.Cmd == runtime.Cmd {
			target.Cmd = filepath.Join(filepath.Dir(target.Cmd), target.CmdName)
		}

		if len(args) == 0 {
			break
		}

		repoPrefix := "github.com"
		repoString := strings.Join(args, "/")
		switch {
		case strings.HasPrefix(repoString, "github.com"):
			repoString = "https://" + repoString
			fallthrough
		case strings.HasPrefix(repoString, "http"):
			// We have a URL
			u, err := url.Parse(repoString)
			if err != nil {
				break
			}
			repoString = u.Path

		default:
			// Leave repoString as is.
		}

		repoArgs := strings.Split(repoString, "/")
		if len(repoArgs) == 0 {
			break
		}

		if len(repoArgs) >= 1 {
			// Assume we have been given a repo prefix only.
			target.CmdBinaryRepo = repoPrefix + "/" + repoArgs[0] + "/" + target.CmdName
			target.CmdSourceRepo = repoPrefix + "/" + repoArgs[0] + "/" + target.CmdName
		}

		if len(repoArgs) >= 2 {
			// Assume we have also been given a repo name.
			if repoArgs[1] != "" {
				target.CmdName = repoArgs[1]
				target.CmdFile = repoArgs[1]
			}
			target.Cmd = filepath.Join(target.CmdDir, target.CmdName)
			target.CmdBinaryRepo = repoPrefix + "/" + repoArgs[0] + "/" + target.CmdName
			target.CmdSourceRepo = repoPrefix + "/" + repoArgs[0] + "/" + target.CmdName
		}

		if len(repoArgs) >= 3 {
			// Assume we have also been given a repo version.
			target.WantVersion = dropVprefix(repoArgs[2])
			if target.WantVersion == LatestVersion {
				target.WantVersion = ""
				break
			}
			if target.WantVersion == "" {
				break
			}

			vCheck := toVersionValue(target.WantVersion)
			if vCheck == nil {
				target.Error = errors.New(fmt.Sprintf("Incorrect semver given: '%s'", repoArgs[2]))
				break
			}
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

		if filepath.Base(link) != defaults.BinaryName {
			err = errors.New("symlink not pointing to bootstrap")
			break
		}

		runtimeBin = filepath.Join(filepath.Dir(targetBin), filepath.Base(link))
		ux.PrintflnOk("Removing symlink %s (%s)", targetBin, filepath.Base(link))
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


func Run(exe string, args ...string) error {
	var err error

	for range onlyOnce {
		ux.PrintflnWhite("Executing the real binary: '%s'", exe)
		c := exec.Command(exe, args...)

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		c.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		err = c.Run()
		waitStatus := c.ProcessState.Sys().(syscall.WaitStatus)
		waitStatus.ExitStatus()
	}

	return err
}
