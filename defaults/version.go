package defaults

import "fmt"

const (
	BinaryName = "bootstrap"
	BinaryVersion = "0.4.3"
	// The version should always be the lowest possible out of all possible binaries.

	RepoPrefix       = "github.com"
	SourceRepoPrefix = RepoPrefix + "/gearboxworks"
	BinaryRepoPrefix = RepoPrefix + "/gearboxworks"

	SourceRepo = SourceRepoPrefix + "/" + BinaryName
	BinaryRepo = BinaryRepoPrefix + "/" + BinaryName
)

type Repo struct {
	Binary string
	Name string
	Owner string
}
type Repos []Repo

var available = Repos{
	{Binary: "scribe",		Owner: "gearboxworks", Name: "scribe"},
	{Binary: "bootstrap",	Owner: "gearboxworks", Name: "bootstrap"},
	{Binary: "buildtool",	Owner: "gearboxworks", Name: "buildtool"},
	{Binary: "launch",		Owner: "gearboxworks", Name: "launch"},
}


func (r *Repos) Get(binary string) string {
	var ret string
	for _, k := range available {
		if k.Binary == binary {
			ret = fmt.Sprintf("%v", k)
			break
		}
	}
	return ret
}


func (r *Repos) GetAll() []string {
	var ret []string
	for _, k := range available {
		ret = append(ret, fmt.Sprintf("%v", k))
	}
	return ret
}


func (r *Repo) String() string {
	return fmt.Sprintf("%s/%s/%s",
		RepoPrefix,
		r.Owner,
		r.Name,
		)
}
