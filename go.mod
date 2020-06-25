module github.com/gearboxworks/bootstrap

go 1.14

replace github.com/newclarity/scribeHelpers/ux => ../scribeHelpers/ux

replace github.com/newclarity/scribeHelpers/toolSelfUpdate => ../scribeHelpers/toolSelfUpdate

replace github.com/newclarity/scribeHelpers/toolRuntime => ../scribeHelpers/toolRuntime

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/google/go-github/v30 v30.1.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/newclarity/scribeHelpers/toolRuntime v0.0.0-20200623081955-45abb1cbefe9
	github.com/newclarity/scribeHelpers/toolSelfUpdate v0.0.0-20200623081955-45abb1cbefe9
	github.com/newclarity/scribeHelpers/ux v0.0.0-20200623081955-45abb1cbefe9
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/tcnksm/go-gitconfig v0.1.2
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
