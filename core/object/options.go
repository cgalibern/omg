package object

import (
	"time"

	"github.com/spf13/cobra"
)

type ActionOptionsGlobal struct {
	DryRun bool
}

func (t *ActionOptionsGlobal) init(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&t.DryRun, "dry-run", false, "show the action execution plan")
}

type ActionOptionsLocking struct {
	NoLock      bool
	LockTimeout time.Duration
}

func (t *ActionOptionsLocking) init(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&t.NoLock, "nolock", false, "don't acquire the action lock (danger)")
	cmd.Flags().DurationVar(&t.LockTimeout, "waitlock", 30*time.Second, "Lock acquire timeout")
}

type ActionOptionsResources struct {
	ResourceSelector string
	SubsetSelector   string
	TagSelector      string
}

func (t *ActionOptionsResources) init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&t.ResourceSelector, "rid", "", "resource selector expression (ip#1,app,disk.type=zvol)")
	cmd.Flags().StringVar(&t.SubsetSelector, "subsets", "", "subset selector expression (g1,g2)")
	cmd.Flags().StringVar(&t.TagSelector, "tags", "", "tag selector expression (t1,t2)")
}

type ActionOptionsForce struct {
	Force bool
}

func (t *ActionOptionsForce) init(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&t.Force, "force", false, "allow dangerous operations")
}

type ActionOptionsStart struct {
	ActionOptionsGlobal
	ActionOptionsLocking
	ActionOptionsResources
	ActionOptionsForce
}

func (t *ActionOptionsStart) Init(cmd *cobra.Command) {
	t.ActionOptionsGlobal.init(cmd)
	t.ActionOptionsLocking.init(cmd)
	t.ActionOptionsResources.init(cmd)
	t.ActionOptionsForce.init(cmd)
}