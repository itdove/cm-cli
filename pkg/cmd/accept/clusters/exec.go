// Copyright Contributors to the Open Cluster Management project
package clusters

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/detach/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spf13/cobra"
	certificatesv1 "k8s.io/api/certificates/v1"
)

const (
	groupName               = "system:bootstrappers:managedcluster"
	userNameSignaturePrefix = "system:bootstrap:"
	clusterLabel            = "open-cluster-management.io/cluster-name"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	//Check if default values must be used
	if o.applierScenariosOptions.ValuesPath == "" {
		o.values = make(map[string]interface{})
		clusters := make([]string, 0)
		if o.clusters != "" {
			cs := strings.Split(o.clusters, ",")
			for _, c := range cs {
				clusters = append(clusters, strings.TrimSpace(c))
			}
			o.values["clusters"] = clusters
		} else {
			return fmt.Errorf("values or name are missing")
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
	iclusters, ok := o.values["clusters"]
	if !ok || iclusters == nil {
		return fmt.Errorf("clusters is missing")
	}
	clusters := iclusters.([]string)

	if o.clusters == "" {
		if len(clusters) == 0 {
			return fmt.Errorf("cluster names are missing")
		}
	}

	alreadyProvidedCluster := make(map[string]bool)
	for _, c := range clusters {
		alreadyProvidedCluster[c] = true
	}

	cs := strings.Split(o.clusters, ",")
	for _, c := range cs {
		if _, ok := alreadyProvidedCluster[c]; !ok {
			clusters = append(clusters, strings.TrimSpace(c))
		}
	}
	o.values["clusters"] = clusters

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
	for _, clusterName := range o.values["clusters"].([]string) {

		csrs := &certificatesv1.CertificateSigningRequestList{}
		ls := labels.SelectorFromSet(labels.Set{
			clusterLabel: clusterName,
		})
		err := client.List(context.TODO(),
			csrs,
			&crclient.ListOptions{
				LabelSelector: ls,
			})
		if err != nil {
			return err
		}
		var csr *certificatesv1.CertificateSigningRequest
		for _, item := range csrs.Items {
			//Does not have the correct name prefix
			if !strings.HasPrefix(item.Spec.Username, userNameSignaturePrefix) {
				continue
			}
			//Check groups
			var group string
			for _, g := range item.Spec.Groups {
				if g == groupName {
					group = g
					break
				}
			}
			//Does not contain the correct group
			if len(group) == 0 {
				continue
			}
			//Check if already approved or denied
			done := false
			for _, c := range item.Status.Conditions {
				if c.Type == certificatesv1.CertificateApproved || c.Type == certificatesv1.CertificateDenied {
					done = true
					break
				}
			}
			if done {
				continue
			}
			csr = &item
			break
		}

		if csr != nil {
			if csr.Status.Conditions == nil {
				csr.Status.Conditions = make([]certificatesv1.CertificateSigningRequestCondition, 0)
			}

			csr.Status.Conditions = append(csr.Status.Conditions, certificatesv1.CertificateSigningRequestCondition{
				Status:         corev1.ConditionTrue,
				Type:           certificatesv1.CertificateApproved,
				Reason:         fmt.Sprintf("%sApprove", helpers.GetExampleHeader()),
				Message:        fmt.Sprintf("This CSR was approved by %s certificate approve.", helpers.GetExampleHeader()),
				LastUpdateTime: metav1.Now(),
			})

			kubeClient, err := o.factory.KubernetesClientSet()
			if err != nil {
				return err
			}
			signingRequest := kubeClient.CertificatesV1().CertificateSigningRequests()
			if _, err := signingRequest.UpdateApproval(context.TODO(), csr.Name, csr, metav1.UpdateOptions{}); err != nil {
				return err
			}

			fmt.Printf("CSR %s approved\n", csr.Name)
		} else {
			fmt.Printf("no CSR to approve for cluster %s\n", clusterName)
		}

		mc := &unstructured.Unstructured{}
		mc.SetKind("ManagedCluster")
		mc.SetAPIVersion("cluster.open-cluster-management.io/v1")
		err = client.Get(context.TODO(),
			crclient.ObjectKey{
				Name: clusterName,
			},
			mc)
		if err != nil {
			return err
		}
		spec := mc.Object["spec"].(map[string]interface{})
		hubAcceptsClient, ok := spec["hubAcceptsClient"]
		if !ok || !hubAcceptsClient.(bool) {
			patch := crclient.MergeFrom(mc.DeepCopyObject())
			spec["hubAcceptsClient"] = true

			err = client.Patch(context.TODO(), mc, patch)
			if err != nil {
				return err
			}
			fmt.Printf("set httpAcceptsClient to true for cluster %s\n", clusterName)
		} else {
			fmt.Printf("httpAcceptsClient already set for cluster %s\n", clusterName)
		}
	}
	return nil

}
