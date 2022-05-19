// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"

	"github.com/stolostron/cm-cli/pkg/cmd/create/hypershiftdeployment/scenario"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	scenarioDirectory = "create"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")

var example = `
# Create a hypershiftdeployment
%[1]s create hypershiftdeployment --values values.yaml

# Create a hypershiftdeployment with cluster name overwrite by args
%[1]s create hypershiftdeployment mycluster --values values.yaml
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "hypershiftdeployment",
		Aliases:      []string{"hypershiftdeployments", "hd", "hds"},
		Short:        "hypershiftdeployment a cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			isSupported, err := helpers.IsSupported(o.CMFlags)
			if err != nil {
				return err
			}
			if !isSupported {
				return fmt.Errorf("this command '%s %s' is only available on %s or %s",
					helpers.GetExampleHeader(),
					strings.Join(os.Args[1:], " "),
					helpers.RHACM,
					helpers.MCE)
			}
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
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

	cmd.SetUsageTemplate(clusteradmhelpers.UsageTempate(cmd, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().StringVarP(&o.clusterNamespace, "namespace", "n", "", "Name of the cluster")
	cmd.Flags().StringVar(&o.valuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")

	return cmd
}
