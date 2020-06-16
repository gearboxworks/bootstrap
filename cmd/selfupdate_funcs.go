package cmd

import (
	"context"
	"fmt"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)


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


func printVersionSummary(release *selfupdate.Release) string {
	var ret string

	for range onlyOnce {
		ret += ux.SprintfBlue("\n")
		ret += fmt.Sprintf("Executable: %s v%s\n",
			ux.SprintfBlue(release.RepoName),
			ux.SprintfWhite(release.Version.String()),
		)
		ret += fmt.Sprintf("Url: %s\n", ux.SprintfBlue(release.URL))
		ret += fmt.Sprintf("Size: %s\n", ux.SprintfBlue("%d", release.AssetByteSize))
		ret += fmt.Sprintf("Published Date: %s\n", ux.SprintfBlue(release.PublishedAt.String()))
	}

	return ret
}


func dropVprefix(v string) string {
	return strings.TrimPrefix(v, "v")
}


func addVprefix(v string) string {
	return "v" + strings.TrimPrefix(v, "v")
}


func stripUrlPrefix(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSpace(url)

	return url
}

func newHTTPClient(ctx context.Context, token string) *http.Client {
	if token == "" {
		return http.DefaultClient
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, src)
}
