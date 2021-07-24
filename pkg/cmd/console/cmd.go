// Copyright Contributors to the Open Cluster Management project
package console

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/console/clusterclaim"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/console/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "console",
		Short: "open a console",
	}

	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))

	return cmd
}
