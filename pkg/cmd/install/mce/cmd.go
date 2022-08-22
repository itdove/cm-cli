// Copyright Contributors to the Open Cluster Management project
package mce

import (
	"fmt"
	"os"
	"strings"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# install Advanced cluster management
%[1]s install mce --namespace <namespace> --channel <channel> [--manual-approval]
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "mce",
		Short:        "install mce",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			isSupported, err := helpers.IsOpenshift(o.CMFlags)
			if err != nil {
				return err
			}
			if !isSupported {
				return fmt.Errorf("this command '%s %s' is only available on Openshift",
					helpers.GetExampleHeader(),
					strings.Join(os.Args[1:], " "))
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.complete(c, args); err != nil {
				return err
			}
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.channel, "channel", "", "The channel to use")
	cmd.Flags().StringVar(&o.namespace, "namespace", "multicluster-engine", "The namespace where to install MCE")
	cmd.Flags().StringVar(&o.operatorGroup, "operatorGroup", "multicluster-engine", "The operator group")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().BoolVar(&o.wait, "wait", false, "Wait until MCE installed is completed")
	cmd.Flags().BoolVar(&o.manualApproval, "manual-approval", false, "Set for manual approval otherwize automatic")
	cmd.Flags().IntVar(&o.timeout, "timeout", 30, "Timeout to get MCE installed in minutes")
	return cmd
}
