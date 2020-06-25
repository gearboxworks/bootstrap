package defaults

import "fmt"

const (
	BinaryName = "bootstrap"
	BinaryVersion = "0.4.6"
	// The version should always be the lowest possible out of all possible binaries.

	RepoPrefix       = "github.com"
	SourceRepoPrefix = "github.com/gearboxworks"
	BinaryRepoPrefix = "github.com/gearboxworks"

	SourceRepo = "github.com/gearboxworks/bootstrap"
	BinaryRepo = "github.com/gearboxworks/bootstrap"
)

type Repo struct {
	Binary string
	Name string
	Owner string
}
type Repos []Repo

var Available = Repos{
	//{Binary: "bootstrap",	Owner: "gearboxworks",	Name: "bootstrap"},
	{Binary: "scribe",		Owner: "gearboxworks",	Name: "scribe"},
	{Binary: "buildtool",	Owner: "gearboxworks",	Name: "buildtool"},
	{Binary: "launch",		Owner: "gearboxworks",	Name: "launch"},
	{Binary: "deploywp",	Owner: "wplib",			Name: "deploywp"},
	{Binary: "squiz",		Owner: "newclarity",	Name: "website-indexer"},
}
// https://github.com/newclarity/website-indexer

func (r *Repos) GetRepo(binary string) string {
	var ret string
	for _, k := range Available {
		if k.Binary == binary {
			ret = fmt.Sprintf("%s", k.Url())
			break
		}
	}
	return ret
}


func (r *Repos) GetRepos() []string {
	var ret []string
	for _, k := range Available {
		ret = append(ret, fmt.Sprintf("%v", k))
	}
	return ret
}


func (r *Repos) GetBinaries() map[string]string {
	ret := make(map[string]string)
	for _, v := range Available {
		ret[v.Binary] = v.Url()
	}
	return ret
}


func (r *Repo) Url() string {
	return fmt.Sprintf("https://%s/%s/%s",
		RepoPrefix,
		r.Owner,
		r.Name,
	)
}


func (r *Repo) String() string {
	return fmt.Sprintf("%s/%s/%s",
		RepoPrefix,
		r.Owner,
		r.Name,
		)
}
