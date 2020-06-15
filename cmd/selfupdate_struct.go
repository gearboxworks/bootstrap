package cmd

import (
	"errors"
	"fmt"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"path/filepath"
	"runtime"
	"strings"
)


type SelfUpdateGetter interface {
}

type SelfUpdateArgs struct {
	name       *string
	version    *string
	sourceRepo *string
	binaryRepo *string

	logging    *bool
}

type TypeSelfUpdate struct {
	name       *StringValue
	version    *VersionValue
	sourceRepo *StringValue
	binaryRepo *StringValue
	logging    *FlagValue

	config     *selfupdate.Config
	ref        *selfupdate.Updater

	useRepo    string

	runtime    *TypeRuntime
	Error      error
}


func (su *TypeSelfUpdate) IsNil() error {
	if su == nil {
		return errors.New("TypeSelfUpdate is nil")
	}
	return nil
}


func New(rt *TypeRuntime) *TypeSelfUpdate {
	//rt = rt.EnsureNotNil()

	te := TypeSelfUpdate{
		name:       toStringValue(rt.CmdName),
		version:    toVersionValue(rt.CmdVersion),
		sourceRepo: toStringValue(stripUrlPrefix(rt.CmdSourceRepo)),
		binaryRepo: toStringValue(stripUrlPrefix(rt.CmdBinaryRepo)),
		logging:    toBoolValue(rt.Debug),

		config:     &selfupdate.Config{
			APIToken:            "",
			EnterpriseBaseURL:   "",
			EnterpriseUploadURL: "",
			Validator:           nil, 	// &MyValidator{},
			Filters:             []string{},
		},

		useRepo:    "",

		runtime: rt,
		Error:   nil,
		//State:   ux.NewState(rt.CmdName, rt.Debug),
	}
	//te.State.SetPackage("")
	//te.State.SetFunctionCaller()

	// Workaround for selfupdate not being flexible enough to support variable asset names
	// Should enable a template similar to GoReleaser.
	// EG: {{ .ProjectName }}-{{ .Os }}_{{ .Arch }}
	//var asset string
	//asset, te.State = toolGhr.GetAsset(rt.CmdBinaryRepo, "latest")
	//te.config.Filters = append(te.config.Filters, asset)

	// Ignore the above and just make sure all filenames are lowercase.
	te.config.Filters = append(te.config.Filters, addFilters(rt.CmdName, runtime.GOOS, runtime.GOARCH)...)
	te.ref, _ = selfupdate.NewUpdater(*te.config)
	if *te.logging {
		selfupdate.EnableLog()
	}

	return &te
}


func addFilters(Binary string, Os string, Arch string) []string {
	var ret []string
	ret = append(ret, fmt.Sprintf("(?i)%s_.*_%s_%s.*", Binary, Os, Arch))
	ret = append(ret, fmt.Sprintf("(?i)%s_%s_%s.*", Binary, Os, Arch))
	ret = append(ret, fmt.Sprintf("(?i)%s-.*_%s_%s.*", Binary, Os, Arch))
	ret = append(ret, fmt.Sprintf("(?i)%s-%s_%s.*", Binary, Os, Arch))
	if Arch == "amd64" {
		// This is recursive - so be careful what you place in the "Arch" argument.
		ret = append(ret, addFilters(Binary, Os, "x86_64.*")...)
		ret = append(ret, addFilters(Binary, Os, "64.*")...)
		ret = append(ret, addFilters(Binary, Os, "64bit.*")...)
	}
	return ret
}


func (su *TypeSelfUpdate) IsValid() error {
	for range onlyOnce {
		if su.name.IsNotValid() {
			su.Error = errors.New("binary name is not defined - selfupdate disabled")
			break
		}

		if su.version.IsNotValid() {
			su.Error = errors.New("binary version is not defined - selfupdate disabled")
			break
		}

		// Refer to binary repo definition first.
		if su.binaryRepo.IsValid() {
			su.useRepo = su.binaryRepo.String()
			su.Error = nil
			break
		}

		// If binary repo is not set, use source repo.
		if su.sourceRepo.IsValid() {
			su.useRepo = su.sourceRepo.String()
			su.Error = nil
			break
		}

		su.Error = errors.New(errorNoRepo)
	}

	return su.Error
}


func (su *TypeSelfUpdate) getRepo() string {
	var ret string

	for range onlyOnce {
		if su.binaryRepo.IsValid() {
			ret = su.binaryRepo.String()
			break
		}
		if su.sourceRepo.IsValid() {
			ret = su.sourceRepo.String()
			break
		}
	}

	return ret
}


func (su *TypeSelfUpdate) Update() error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		ux.PrintflnBlue("Checking '%s' for version greater than v%s", su.useRepo, su.version.String())
		previous := su.version.ToSemVer()
		// up.UpdateCommand(cmdPath, current, slug)
		//latest, err := su.ref.UpdateSelf(previous, su.useRepo)
		latest, err := su.ref.UpdateCommand(filepath.Join(su.runtime.CmdDir, su.runtime.CmdName), previous, su.useRepo)
		if err != nil {
			su.Error = err
			break
		}

		if previous.Equals(latest.Version) {
			ux.PrintflnOk("%s is up to date: v%s", su.name.String(), su.version.String())
		} else {
			ux.PrintflnOk("%s updated to v%s", su.name.String(), latest.Version)
			if latest.ReleaseNotes != "" {
				ux.PrintflnOk("%s %s Release Notes:\n%s", su.name.String(), latest.Version, latest.ReleaseNotes)
			}
		}
	}

	return su.Error
}


func (su *TypeSelfUpdate) GetVersion(version *VersionValue) *selfupdate.Release {
	var release *selfupdate.Release

	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		var ok bool
		var err error

		switch {
		case version.IsNotValid():
			fallthrough
		case version.IsLatest():
			release, ok, err = su.ref.DetectLatest(su.useRepo)

		default:
			v := version.String()
			if !strings.HasPrefix(v, "v") {
				v = "v" + v
			}
			release, ok, err = su.ref.DetectVersion(su.useRepo, v)
		}

		if !ok {
			su.Error = errors.New(errorNoVersion)
			break
		}
		if err != nil {
			su.Error = errors.New(fmt.Sprintf("%s - %s", errorNoVersion, err))
			break
		}

		//su.State.SetOutput(release)
	}

	return release
}


func (su *TypeSelfUpdate) PrintVersion(version *VersionValue) error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		release := su.GetVersion(version)
		if su.Error != nil {
			break
		}

		fmt.Printf(printVersion(release))
	}

	return su.Error
}


func (su *TypeSelfUpdate) IsUpdated() error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		latest := su.GetVersion(nil)
		if su.Error != nil {
			break
		}

		current := su.GetVersion(su.version)

		if current == nil {
			ux.PrintflnWarning("%s can be updated to v%s.",
				su.name.String(),
				latest.Version.String())
			ux.PrintflnWarning("Current version info unknown.")
			ux.PrintflnBlue("Current version (v%s)\n", su.version.String())
			ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
			fmt.Printf(printVersion(latest))

			su.Error = nil
			break
		}

		if current.Version.Equals(latest.Version) {
			ux.PrintflnOk("%s is up to date at v%s.",
				su.name.String(),
				su.version.String())
			fmt.Printf(printVersion(current))
			break
		}

		if current.Version.LE(latest.Version) {
			ux.PrintflnOk("%s can be updated to v%s.",
				su.name.String(),
				su.version.String())

			ux.PrintflnBlue("Current version (v%s)", current.Version.String())
			fmt.Printf(printVersion(current))
			ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
			fmt.Printf(printVersion(latest))
			break
		}

		if current.Version.GT(latest.Version) {
			ux.PrintflnWarning("%s is more recent at v%s, (latest is %s).",
				su.name.String(),
				su.version.String(),
				latest.Version.String(),
				)

			ux.PrintflnBlue("Current version (v%s)", current.Version.String())
			fmt.Printf(printVersion(current))
			ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
			fmt.Printf(printVersion(latest))
			break
		}
	}

	return su.Error
}


func (su *TypeSelfUpdate) Set(s SelfUpdateArgs) error {
	if s.name != nil {
		su.name = (*StringValue)(s.name)
	}

	if s.version != nil {
		su.version = toVersionValue(*s.version)
	}

	if s.binaryRepo != nil {
		su.binaryRepo = (*StringValue)(s.binaryRepo)
	}

	if s.sourceRepo != nil {
		su.sourceRepo = (*StringValue)(s.sourceRepo)
	}

	if s.logging != nil {
		su.logging = (*FlagValue)(s.logging)
	} else {
		su.logging = &defaultFalse
	}

	su.Error = su.IsValid()

	return su.Error
}


func (su *TypeSelfUpdate) SetDebug(value bool) error {
	su.logging = (*FlagValue)(&value)
	su.Error = su.IsValid()
	return su.Error
}


func (su *TypeSelfUpdate) SetName(value string) error {
	su.name = (*StringValue)(&value)
	su.Error = su.IsValid()
	return su.Error
}


func (su *TypeSelfUpdate) SetVersion(value string) error {
	su.version = toVersionValue(value)
	su.Error = su.IsValid()
	return su.Error
}


func (su *TypeSelfUpdate) SetSourceRepo(value string) error {
	su.sourceRepo = (*StringValue)(&value)
	su.Error = su.IsValid()
	return su.Error
}


func (su *TypeSelfUpdate) SetBinaryRepo(value string) error {
	su.binaryRepo = (*StringValue)(&value)
	su.Error = su.IsValid()
	return su.Error
}
