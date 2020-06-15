package cmd

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"strings"
)


func GetSemVer(v string) *VersionValue {
	sver := (VersionValue)(semver.MustParse(v))
	return &sver
}


func printVersion(release *selfupdate.Release) string {
	var ret string

	for range onlyOnce {
		ret += ux.SprintfBlue("Repository release information:\n")
		ret += fmt.Sprintf("Executable: %s v%s\n",
			ux.SprintfBlue(release.RepoName),
			ux.SprintfWhite(release.Version.String()),
		)

		ret += fmt.Sprintf("Url: %s\n", ux.SprintfBlue(release.URL))

		//ret += fmt.Sprintf("TypeRepo Owner: %s\n", ux.SprintfBlue(release.RepoOwner))
		//ret += fmt.Sprintf("TypeRepo Name: %s\n", ux.SprintfBlue(release.RepoName))

		ret += fmt.Sprintf("Size: %s\n", ux.SprintfBlue("%d", release.AssetByteSize))

		ret += fmt.Sprintf("Published Date: %s\n", ux.SprintfBlue(release.PublishedAt.String()))

		if release.ReleaseNotes != "" {
			ret += fmt.Sprintf("Release Notes: %s\n", ux.SprintfBlue(release.ReleaseNotes))
		}
	}

	return ret
}


func stripUrlPrefix(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSpace(url)

	return url
}
