// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/detach/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	//Check if default values must be used
	if o.applierScenariosOptions.ValuesPath == "" {
		o.values = make(map[string]interface{})
		mc := make(map[string]interface{})
		o.values["managedCluster"] = mc
		if o.clusterName != "" {
			mc["name"] = o.clusterName
		} else {
			return fmt.Errorf("values or name are missing")
		}
		if o.replicas != -1 {
			mc["replicas"] = o.replicas
		} else {
			return fmt.Errorf("values or replicas are missing")
		}
	} else {
		//Read values
		o.values, err = appliercmd.ConvertValuesFileToValuesMap(o.applierScenariosOptions.ValuesPath, "")
		if err != nil {
			return err
		}
	}

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	return nil
}

func (o *Options) validate() error {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	imc, ok := o.values["managedCluster"]
	if !ok || imc == nil {
		return fmt.Errorf("managedCluster is missing")
	}
	mc := imc.(map[string]interface{})

	if o.clusterName == "" {
		iname, ok := mc["name"]
		if !ok || iname == nil {
			return fmt.Errorf("cluster name is missing")
		}
		o.clusterName = iname.(string)
		if len(o.clusterName) == 0 {
			return fmt.Errorf("managedCluster.name not specified")
		}
	}

	mc["name"] = o.clusterName

	if o.replicas == -1 {
		ireplicas, ok := mc["replicas"]
		if !ok || ireplicas == nil {
			return fmt.Errorf("replicas number is missing")
		}
		o.replicas = ireplicas.(int)
		if o.replicas < 0 {
			return fmt.Errorf("replicas must greather or equal to zero")
		}
	}

	mc["replicas"] = o.replicas

	return nil
}

func (o *Options) run() error {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		reader := scenario.GetApplierScenarioResourcesReader()
		return reader.ExtractAssets(scenarioDirectory, o.applierScenariosOptions.OutTemplatesDir)
	}
	client, err := helpers.GetControllerRuntimeClientFromFlags(o.applierScenariosOptions.ConfigFlags)
	if err != nil {
		return err
	}
	return o.runWithClient(client)
}

func (o *Options) runWithClient(client crclient.Client) error {
	mp := &unstructured.Unstructured{}
	mp.SetKind("MachinePool")
	mp.SetAPIVersion("hive.openshift.io/v1")
	err := client.Get(context.TODO(),
		crclient.ObjectKey{
			Name:      o.clusterName + "-worker",
			Namespace: o.clusterName}, mp)
	if err != nil {
		return err
	}
	patch := crclient.MergeFrom(mp.DeepCopyObject())
	spec := mp.Object["spec"].(map[string]interface{})
	spec["replicas"] = o.replicas

	return client.Patch(context.TODO(), mp, patch)

}
