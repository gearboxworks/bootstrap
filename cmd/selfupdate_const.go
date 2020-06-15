package cmd


const (
	CmdSelfUpdate		= "selfupdate"
	CmdVersion 			= "version"
	CmdVersionInfo		= "info"
	CmdVersionList		= "list"
	CmdVersionLatest	= "latest"
	CmdVersionCheck		= "check"
	CmdVersionUpdate	= "update"

	FlagVersion 		= "version"
)


const (
	errorNoRepo = "repo is not defined - selfupdate disabled"
	errorNoVersion = "no versions in repo - selfupdate disabled"
	LatestVersion = "latest"
)


var defaultFalse = FlagValue(false)
