package cmd

import (
	"errors"
	"fmt"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"runtime"
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
	owner       *StringValue
	name       *StringValue
	version    *VersionValue
	sourceRepo *StringValue
	binaryRepo *StringValue
	logging    *FlagValue

	config     *selfupdate.Config
	ref        *selfupdate.Updater

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
		owner:      toOwnerValue(rt.CmdSourceRepo),
		name:       toStringValue(rt.CmdName),
		version:    toVersionValue(rt.WantVersion),
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

		runtime: rt,
		Error:   nil,
	}


	// Workaround for selfupdate not being flexible enough to support variable asset names
	// Should enable a template similar to GoReleaser.
	// EG: {{ .ProjectName }}-{{ .Os }}_{{ .Arch }}
	//var asset string
	//asset, te.State = toolGhr.GetAsset(rt.CmdBinaryRepo, LatestVersion)
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

		if su.binaryRepo.IsValid() {
			su.Error = nil
			break
		}

		if su.sourceRepo.IsValid() {
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


func (su *TypeSelfUpdate) UpdateTo() error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		wantVersion := su.version.ToSemVer()
		newRelease := su.GetVersion(su.version)
		if newRelease == nil {
			ux.PrintflnError("Version v%s doesn't exist for '%s'", su.version.String(), su.name.String())
			break
		}
		su.Error = su.ref.UpdateTo(newRelease, su.runtime.Cmd)
		if su.Error != nil {
			break
		}
		if wantVersion.Equals(newRelease.Version) {
			ux.PrintflnOk("%s updated to v%s", su.name.String(), su.version.String())
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
				release, ok, err = su.ref.DetectLatest(su.binaryRepo.String())

			default:
				v := addVprefix(version.String())
				release, ok, err = su.ref.DetectVersion(su.binaryRepo.String(), v)
		}

		if !ok {
			su.Error = errors.New(fmt.Sprintf(errorNoVersion, version.String()))
			break
		}
		if err != nil {
			su.Error = errors.New(fmt.Sprintf("%s - %s", fmt.Sprintf(errorNoVersion, version.String()), err))
			break
		}

		su.version = toVersionValue(release.Name)
	}

	return release
}


func (su *TypeSelfUpdate) PrintVersion(version string) error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		var release *selfupdate.Release
		switch version {
			case LatestVersion:
				fallthrough
			case "":
				release = su.GetVersion(nil)

			default:
				release = su.GetVersion(toVersionValue(version))
		}
		if su.Error != nil {
			break
		}

		fmt.Print(printVersion(release))
	}

	return su.Error
}


func (su *TypeSelfUpdate) PrintVersionSummary(version string) error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		var release *selfupdate.Release
		if version == LatestVersion {
			release = su.GetVersion(nil)
		} else {
			release = su.GetVersion(toVersionValue(version))
		}
		if su.Error != nil {
			break
		}

		fmt.Print(printVersionSummary(release))
	}

	return su.Error
}


func (su *TypeSelfUpdate) IsUpdated() error {
	for range onlyOnce {
		su.Error = su.IsValid()
		if su.Error != nil {
			break
		}

		newRelease := su.GetVersion(su.version)
		if newRelease == nil {
			break
		}

		currentRelease := su.GetVersion(toVersionValue(su.runtime.CmdVersion))

		if currentRelease == nil {
			ux.PrintflnOk("%s can be updated to v%s.",
				su.name.String(),
				newRelease.Version.String())
			ux.PrintflnBlue("Current version (v%s)", su.version.String())
			ux.PrintflnWarning("Info for current version unknown.\n")

			ux.PrintflnBlue("Updated version (v%s)", newRelease.Version.String())
			fmt.Printf("%s\n", printVersion(newRelease))

			su.Error = nil
			break
		}

		if currentRelease.Version.Equals(newRelease.Version) {
			ux.PrintflnOk("%s is up to date at v%s.",
				su.name.String(),
				su.version.String())
			fmt.Printf("%s\n", printVersion(currentRelease))
			break
		}

		if currentRelease.Version.LE(newRelease.Version) {
			ux.PrintflnOk("%s can be updated to v%s.",
				su.name.String(),
				su.version.String())

			ux.PrintflnBlue("Current version (v%s)", currentRelease.Version.String())
			fmt.Printf("%s\n", printVersion(currentRelease))

			ux.PrintflnBlue("Updated version (v%s)", newRelease.Version.String())
			fmt.Printf("%s\n", printVersion(newRelease))
			break
		}

		if currentRelease.Version.GT(newRelease.Version) {
			ux.PrintflnWarning("%s is more recent at v%s, (latest is %s).",
				su.name.String(),
				su.version.String(),
				newRelease.Version.String(),
				)

			ux.PrintflnBlue("Current version (v%s)", currentRelease.Version.String())
			fmt.Printf("%s\n", printVersion(currentRelease))

			ux.PrintflnBlue("Updated version (v%s)", newRelease.Version.String())
			fmt.Printf("%s\n", printVersion(newRelease))
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
