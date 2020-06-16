package cmd

//noinspection ALL
const (
	CmdSelfUpdate		= "selfupdate"
	CmdVersion 			= "version"
	CmdVersionInfo		= "info"
	CmdVersionList		= "list"
	CmdVersionLatest	= "latest"
	CmdVersionCheck		= "check"
	CmdVersionUpdate	= "update"
	CmdVersionExamples	= "examples"

	FlagVersion 		= "version"
)


const (
	errorNoRepo = "repo is not defined"
	errorNoVersion = "specified version '%s' not found in repo"
	LatestVersion = "latest"
)


var defaultFalse = FlagValue(false)
