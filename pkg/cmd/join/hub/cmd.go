// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"fmt"
	"path/filepath"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/applierscenarios"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/join/hub/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

var example = `
# Init hub
%[1]s join hub
`

const (
	scenarioDirectory = "join"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")

// NewCmd ...
func NewCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(f, streams)

	cmd := &cobra.Command{
		Use:          "hub",
		Short:        "join a hub",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
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

	cmd.SetUsageTemplate(applierscenarios.UsageTempate(cmd, scenario.GetApplierScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().StringVar(&o.token, "hub-token", "", "The token to access the hub")
	cmd.Flags().StringVar(&o.hubServerInternal, "hub-server-internal", "", "The internal api server url to the hub")
	cmd.Flags().StringVar(&o.hubServerExternal, "hub-server-external", "", "The external api server url to the hub")
	cmd.Flags().StringVar(&o.clusterName, "name", "", "The name of the joining cluster")

	o.applierScenariosOptions.AddFlags(cmd.Flags())
	o.applierScenariosOptions.ConfigFlags.AddFlags(cmd.Flags())

	return cmd
}
