package commands

import (
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/objectaction"
	"opensvc.com/opensvc/core/path"
)

type (
	// CmdObjectUnfreeze is the cobra flag set of the unfreeze command.
	CmdObjectUnfreeze struct {
		Global object.OptsGlobal
		Async  object.OptsAsync
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *CmdObjectUnfreeze) Init(kind string, parent *cobra.Command, selector *string) {
	cmd := t.cmd(kind, selector)
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *CmdObjectUnfreeze) cmd(kind string, selector *string) *cobra.Command {
	return &cobra.Command{
		Use:   "unfreeze",
		Short: "unfreeze the selected objects",
		Run: func(cmd *cobra.Command, args []string) {
			t.run(selector, kind)
		},
	}
}

func (t *CmdObjectUnfreeze) run(selector *string, kind string) {
	mergedSelector := mergeSelector(*selector, t.Global.ObjectSelector, kind, "")
	objectaction.New(
		objectaction.WithLocal(t.Global.Local),
		objectaction.WithObjectSelector(mergedSelector),
		objectaction.WithFormat(t.Global.Format),
		objectaction.WithColor(t.Global.Color),
		objectaction.WithServer(t.Global.Server),
		objectaction.WithAsyncTarget("thawed"),
		objectaction.WithAsyncWatch(t.Async.Watch),
		objectaction.WithRemoteNodes(t.Global.NodeSelector),
		objectaction.WithRemoteAction("unfreeze"),
		objectaction.WithLocalRun(func(p path.T) (interface{}, error) {
			return nil, object.NewActorFromPath(p).Unfreeze()
		}),
	).Do()
}
