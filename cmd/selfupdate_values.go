package cmd

import "github.com/blang/semver"


type StringValue string
type VersionValue semver.Version
type FlagValue bool


func (v *VersionValue) String() string {
	return (semver.Version)(*v).String()
}
func toVersionValue(version string) *VersionValue {
	v := VersionValue(semver.MustParse(version))
	return &v
}
func toSemVer(version string) semver.Version {
	return semver.MustParse(version)
}
func toStringValue(s string) *StringValue {
	v := StringValue(s)
	return &v
}
func toBoolValue(b bool) *FlagValue {
	v := FlagValue(b)
	return &v
}


func (v *VersionValue) ToSemVer() semver.Version {
	return semver.Version(*v)
}
func (v *VersionValue) IsValid() bool {
	var ok bool
	for range onlyOnce {
		if v == nil {
			break
		}

		err := (semver.Version)(*v).Validate()
		if err != nil {
			break
		}

		ok = true
	}
	return ok
}
func (v *VersionValue) IsNotValid() bool {
	return !v.IsValid()
}
func (v *VersionValue) IsLatest() bool {
	return (semver.Version)(*v).String() == LatestVersion
}


func (v *StringValue) String() string {
	return string(*v)
}


func (v *StringValue) IsValid() bool {
	var ok bool
	for range onlyOnce {
		if v == nil {
			break
		}
		if *v == "" {
			break
		}
		ok = true
	}
	return ok
}
func (v *StringValue) IsNotValid() bool {
	return !v.IsValid()
}


func (v *StringValue) IsNil() bool {
	if v == nil {
		return true
	}
	return false
}
func (v *StringValue) IsNotNil() bool {
	return !v.IsNil()
}


func (v *StringValue) IsEmpty() bool {
	if v == nil {
		return true
	}
	return false
}
func (v *StringValue) IsNotEmpty() bool {
	return !v.IsEmpty()
}
