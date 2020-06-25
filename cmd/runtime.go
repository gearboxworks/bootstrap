package cmd


//type TypeUpdate struct {
//	//WantVersion string `json:"want_version" mapstructure:"want_version"`
//	//AutoExec    bool   `json:"auto_update" mapstructure:"auto_update"`
//	//IsSymLinked bool   `json:"is_symlinked" mapstructure:"is_symlinked"`
//	//Update      *toolSelfUpdate.TypeSelfUpdate
//
//	*toolRuntime.TypeRuntime
//
//	//CmdName       string    `json:"cmd_name" mapstructure:"cmd_name"`
//	//CmdVersion    string    `json:"cmd_version" mapstructure:"cmd_version"`
//	//CmdSourceRepo string   `json:"cmd_source_repo" mapstructure:"cmd_source_repo"`
//	//CmdBinaryRepo string   `json:"cmd_binary_repo" mapstructure:"cmd_binary_repo"`
//	//Cmd           string    `json:"cmd" mapstructure:"cmd"`
//	//CmdDir        string    `json:"cmd_dir" mapstructure:"cmd_dir"`
//	//CmdFile       string    `json:"cmd_file" mapstructure:"cmd_file"`
//	//
//	//CmdArgs       []string  `json:"cmd_args" mapstructure:"cmd_args"`
//	//
//	//GoRuntime     GoRuntime `json:"go_runtime" mapstructure:"go_runtime"`
//	//
//	//Debug         bool
//	//
//	//
//	//Error         error
//}


//type GoRuntime struct {
//	Os string
//	Arch string
//	Root string
//	Version string
//	Compiler string
//	NumCpus int
//}


//func NewUpdate(BinaryName string, BinaryVersion string, SourceRepo string, BinaryRepo string, debugFlag bool) TypeUpdate {
//
//	r := TypeUpdate {
//		WantVersion: "",
//		AutoExec:    false,
//		IsSymLinked: false,
//		Update:      nil,
//
//		TypeRuntime: toolRuntime.New(BinaryName, BinaryVersion, debugFlag),
//	}
//
//	for range onlyOnce {
//		r.State = r.SetRepos(defaults.SourceRepo, defaults.BinaryRepo)
//
//		//if r.CmdName == "" {
//		//	r.CmdName = r.CmdFile
//		//}
//		//
//		//if r.CmdSourceRepo == "" {
//		//	r.CmdSourceRepo = "github.com/gearboxworks/" + r.CmdName
//		//}
//		//
//		//if r.CmdBinaryRepo == "" {
//		//	r.CmdBinaryRepo = "github.com/gearboxworks/" + r.CmdName
//		//}
//
//		r.Update = toolSelfUpdate.New(r.TypeRuntime)
//		r.Update.LoadCommands(rootCmd, false)
//		if r.Update.State.IsNotOk() {
//			break
//		}
//	}
//
//	return r
//}


//func DecodeToUrl(args ...string) (string, error) {
//	var ret string
//	var err error
//
//	for range onlyOnce {
//		repoPrefix := "github.com"
//		repoString := strings.Join(args, "/")
//
//		switch {
//			case strings.HasPrefix(repoString, "github.com"):
//				repoString = "https://" + repoString
//				fallthrough
//			case strings.HasPrefix(repoString, "http"):
//				// We have a URL
//				u, err := url.Parse(repoString)
//				if err != nil {
//					break
//				}
//				repoString = u.Path
//
//			default:
//				// Leave repoString as is.
//		}
//
//		repoString = strings.ReplaceAll(repoString, "//", "/")
//		u, err := url.Parse(repoString)
//		if err != nil {
//			break
//		}
//		repoString = u.Path
//
//		repoArgs := strings.Split(repoString, "/")
//		switch len(repoArgs) {
//			case 0:
//				err = errors.New(fmt.Sprintf("Url empty"))
//			case 1:
//				// Assume we have been given a binary.
//				u := defaults.Available.GetRepo(repoArgs[0])
//				rt.CmdBinaryRepo = u
//
//			case 2:
//				// Assume we have been given a repo prefix only.
//				rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//
//			case 3:
//				// Assume we have also been given a repo name.
//				rt.Cmd = filepath.Join(rt.CmdDir, rt.CmdName)
//				rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//
//			default:
//				// Assume we have also been given a repo version.
//				rt.Cmd = filepath.Join(rt.CmdDir, rt.CmdName)
//				rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//
//				rt.WantVersion = dropVprefix(repoArgs[3])
//				if rt.WantVersion == LatestVersion {
//					rt.WantVersion = ""
//					break
//				}
//				if rt.WantVersion == "" {
//					break
//				}
//
//				vCheck := toVersionValue(rt.WantVersion)
//				if vCheck == nil {
//					err = errors.New(fmt.Sprintf("Incorrect semver given: '%s'", repoArgs[3]))
//					break
//				}
//		}
//
//		//if len(repoArgs) == 0 {
//		//	break
//		//}
//		//
//		//if len(repoArgs) == 1 {
//		//	break
//		//}
//		//
//		//if len(repoArgs) >= 2 {
//		//	// Assume we have been given a repo prefix only.
//		//	rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//	rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//}
//		//
//		//if len(repoArgs) >= 3 {
//		//	// Assume we have also been given a repo name.
//		//	if repoArgs[2] != "" {
//		//		rt.CmdName = repoArgs[2]
//		//		rt.CmdFile = repoArgs[2]
//		//	}
//		//	rt.Cmd = filepath.Join(rt.CmdDir, rt.CmdName)
//		//	rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//	rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//}
//		//
//		//if len(repoArgs) >= 4 {
//		//	// Assume we have also been given a repo version.
//		//	rt.WantVersion = dropVprefix(repoArgs[3])
//		//	if rt.WantVersion == LatestVersion {
//		//		rt.WantVersion = ""
//		//		break
//		//	}
//		//	if rt.WantVersion == "" {
//		//		break
//		//	}
//		//
//		//	vCheck := toVersionValue(rt.WantVersion)
//		//	if vCheck == nil {
//		//		rt.Error = errors.New(fmt.Sprintf("Incorrect semver given: '%s'", repoArgs[3]))
//		//		break
//		//	}
//		//}
//	}
//
//	return ret, err
//}


//func (rt *TypeUpdate) SetApp(runtime *TypeUpdate, args ...string) *ux.State {
//	return rt.State
//	for range onlyOnce {
//		if rt.WantVersion == LatestVersion {
//			rt.WantVersion = ""
//		}
//		if rt.CmdVersion == LatestVersion {
//			rt.CmdVersion = runtime.CmdVersion
//		} else if rt.CmdVersion == "" {
//			rt.CmdVersion = runtime.CmdVersion
//		}
//
//		rt.CmdSourceRepo = rt.CmdSourceRepo + "/" + rt.CmdName
//		rt.CmdBinaryRepo = rt.CmdBinaryRepo + "/" + rt.CmdName
//
//		if rt.Cmd == runtime.Cmd {
//			rt.Cmd = filepath.Join(filepath.Dir(rt.Cmd), rt.CmdName)
//		}
//
//		if len(args) == 0 {
//			break
//		}
//
//		//repoPrefix := "github.com"
//		repoString := strings.Join(args, "/")
//		switch {
//			case strings.HasPrefix(repoString, "github.com"):
//				repoString = "https://" + repoString
//				fallthrough
//			case strings.HasPrefix(repoString, "http"):
//				// We have a URL
//				u, err := url.Parse(repoString)
//				if err != nil {
//					break
//				}
//				repoString = u.Path
//
//			default:
//				// Leave repoString as is.
//		}
//
//		repoString = strings.ReplaceAll(repoString, "//", "/")
//		u, err := url.Parse(repoString)
//		if err != nil {
//			break
//		}
//		repoString = u.Path
//
//		//repoArgs := strings.Split(repoString, "/")
//		//switch len(repoArgs) {
//		//	case 0:
//		//		rt.Error = errors.New(fmt.Sprintf("Url empty"))
//		//	case 1:
//		//		// Assume we have been given a binary.
//		//		u := defaults.Available.GetRepo(repoArgs[0])
//		//		rt.CmdName = repoArgs[0]
//		//		rt.CmdFile = repoArgs[0]
//		//		rt.CmdBinaryRepo = u
//		//		rt.CmdSourceRepo = u
//		//		//rt.Error = errors.New(fmt.Sprintf("Incorrect url given"))
//		//
//		//	case 2:
//		//		// Assume we have been given a repo prefix only.
//		//		rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//		rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//
//		//	case 3:
//		//		// Assume we have also been given a repo name.
//		//		if repoArgs[2] != "" {
//		//			rt.CmdName = repoArgs[2]
//		//			rt.CmdFile = repoArgs[2]
//		//		}
//		//		rt.Cmd = filepath.Join(rt.CmdDir, rt.CmdName)
//		//		rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//		rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//
//		//	default:
//		//		// Assume we have also been given a repo version.
//		//		if repoArgs[2] != "" {
//		//			rt.CmdName = repoArgs[2]
//		//			rt.CmdFile = repoArgs[2]
//		//		}
//		//		rt.Cmd = filepath.Join(rt.CmdDir, rt.CmdName)
//		//		rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//		rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//
//  		//		rt.WantVersion = dropVprefix(repoArgs[3])
//		//		if rt.WantVersion == LatestVersion {
//		//			rt.WantVersion = ""
//		//			break
//		//		}
//		//		if rt.WantVersion == "" {
//		//			break
//		//		}
//		//
//		//		vCheck := toVersionValue(rt.WantVersion)
//		//		if vCheck == nil {
//		//			rt.Error = errors.New(fmt.Sprintf("Incorrect semver given: '%s'", repoArgs[3]))
//		//			break
//		//		}
//		//}
//
//		//if len(repoArgs) == 0 {
//		//	break
//		//}
//		//
//		//if len(repoArgs) == 1 {
//		//	break
//		//}
//		//
//		//if len(repoArgs) >= 2 {
//		//	// Assume we have been given a repo prefix only.
//		//	rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//	rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//}
//		//
//		//if len(repoArgs) >= 3 {
//		//	// Assume we have also been given a repo name.
//		//	if repoArgs[2] != "" {
//		//		rt.CmdName = repoArgs[2]
//		//		rt.CmdFile = repoArgs[2]
//		//	}
//		//	rt.Cmd = filepath.Join(rt.CmdDir, rt.CmdName)
//		//	rt.CmdBinaryRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//	rt.CmdSourceRepo = repoPrefix + "/" + repoArgs[1] + "/" + rt.CmdName
//		//}
//		//
//		//if len(repoArgs) >= 4 {
//		//	// Assume we have also been given a repo version.
//		//	rt.WantVersion = dropVprefix(repoArgs[3])
//		//	if rt.WantVersion == LatestVersion {
//		//		rt.WantVersion = ""
//		//		break
//		//	}
//		//	if rt.WantVersion == "" {
//		//		break
//		//	}
//		//
//		//	vCheck := toVersionValue(rt.WantVersion)
//		//	if vCheck == nil {
//		//		rt.Error = errors.New(fmt.Sprintf("Incorrect semver given: '%s'", repoArgs[3]))
//		//		break
//		//	}
//		//}
//	}
//
//	return rt.State
//}


//func CreateDummyBinary(runtimeBin string, targetBin string) error {
//	var err error
//
//	for range onlyOnce {
//		var link string
//
//		link, err = os.Readlink(targetBin)
//		if link == "" {
//			if targetBin == runtimeBin {
//				err = nil
//				break
//			}
//			_, err = os.Stat(targetBin)
//			if os.IsNotExist(err) {
//				// File doesn't exist - need to create it.
//				err = CopyFile(runtimeBin, targetBin)
//			}
//
//			break
//		}
//
//		_, err = os.Stat(targetBin)
//		if os.IsNotExist(err) {
//			if link != "" {
//				// File is a dud - no action.
//				err = errors.New(fmt.Sprintf("file '%s' is a symlink pointing to non-existant file '%s'", targetBin, link))
//				break
//			}
//
//			// File doesn't exist - need to create it.
//			err = CopyFile(runtimeBin, targetBin)
//			break
//		}
//
//		if filepath.Base(link) != defaults.BinaryName {
//			err = errors.New("symlink not pointing to bootstrap")
//			break
//		}
//
//		if filepath.IsAbs(link) {
//			runtimeBin = link
//		} else {
//			runtimeBin = filepath.Join(filepath.Dir(targetBin), filepath.Base(link))
//		}
//		ux.PrintflnOk("Removing symlink %s (%s)", targetBin, filepath.Base(link))
//		err = os.Remove(targetBin)
//		if err != nil {
//			break
//		}
//
//		err = CopyFile(runtimeBin, targetBin)
//	}
//
//	return err
//}
//
//
//func CopyFile(runtimeBin string, targetBin string) error {
//	var err error
//
//	for range onlyOnce {
//		var input []byte
//		input, err = ioutil.ReadFile(runtimeBin)
//		if err != nil {
//			break
//		}
//
//		err = ioutil.WriteFile(targetBin, input, 0755)
//		if err != nil {
//			fmt.Println("Error creating", targetBin)
//			break
//		}
//	}
//
//	return err
//}
//
//
//func CompareBinary(runtimeBin string, newBin string) error {
//	var err error
//
//	for range onlyOnce {
//		var srcBin []byte
//		srcBin, err = ioutil.ReadFile(runtimeBin)
//		if err != nil {
//			break
//		}
//		if srcBin == nil {
//			break
//		}
//
//		var targetBin []byte
//		targetBin, err = ioutil.ReadFile(newBin)
//		if err != nil {
//			break
//		}
//		if targetBin == nil {
//			break
//		}
//
//		if len(srcBin) != len(targetBin) {
//			break
//		}
//
//		err = errors.New("binary files differ")
//		for i := range srcBin {
//			if srcBin[i] != targetBin[i] {
//				err = nil
//				break
//			}
//		}
//	}
//
//	return err
//}


//func Run(exe string, args ...string) error {
//	var err error
//
//	for range onlyOnce {
//		ux.PrintflnWhite("Executing the real binary: '%s'", exe)
//		c := exec.Command(exe, args...)
//
//		var stdoutBuf, stderrBuf bytes.Buffer
//		c.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
//		c.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
//		err = c.Run()
//		waitStatus := c.ProcessState.Sys().(syscall.WaitStatus)
//		waitStatus.ExitStatus()
//	}
//
//	return err
//}
