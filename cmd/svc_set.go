package cmd

import (
	"opensvc.com/opensvc/core/commands"
)

var (
	svcSet commands.CmdObjectSet
)

func init() {
	svcSet.Init("svc", svcCmd, &selectorFlag)
}