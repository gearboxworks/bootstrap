package cmd

import (
	"github.com/blang/semver"
	"strings"
)


type StringValue string
type VersionValue semver.Version
type FlagValue bool


func (v *VersionValue) String() string {
	if v == nil {
		return ""
	}
	vers := (semver.Version)(*v).String()
	return dropVprefix(vers)
}
func toVersionValue(version string) *VersionValue {
	if version == LatestVersion {
		return nil
	}
	if version == "" {
		return nil
	}
	version = dropVprefix(version)
	v, err := semver.Parse(version)
	if err != nil {
		return nil
	}
	ret := VersionValue(v)
	return &ret
}
func toStringValue(s string) *StringValue {
	v := StringValue(s)
	return &v
}
func toOwnerValue(s string) *StringValue {
	s = stripUrlPrefix(s)
	if strings.Contains(s, "/") {
		sa := strings.Split(s, "/")
		switch {
			case len(sa) == 0:
				// Nada
			default:
				s = sa[0]
		}

	}
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
			ok = true	// Assume "latest"
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
	if v == nil {
		return true
	}
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
