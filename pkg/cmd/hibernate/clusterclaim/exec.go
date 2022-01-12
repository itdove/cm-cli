// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterclaim name is missing")
	}
	o.ClusterClaims = args[0]
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.HibernateClusterClaims(o.ClusterClaims, o.SkipSchedule, o.CMFlags.DryRun, o.outputFile)
}
