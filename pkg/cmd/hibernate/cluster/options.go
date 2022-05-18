// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags *genericclioptionscm.CMFlags
	// list of cluster names (comma-separated)
	Clusters string
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
