/*
  WIP code for https://github.com/openshift/enhancements/pull/201
  that will probably end up in the MCO or machineAPI
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openshiftv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/coreos/ign-converter/translate/v24tov31"
	"github.com/coreos/ignition/config/v2_4"
	"github.com/stretchr/objx"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	machineAPINamespace = "openshift-machine-api"
	userDataSuffix      = "-user-data"
	userDataKey         = "userData"
	suffix              = "-ignv3"

	machineAPIGroup       = "machine.openshift.io"
	machineSetOwningLabel = "machine.openshift.io/cluster-api-machineset"
	machineLabelRole      = "machine.openshift.io/cluster-api-machine-role"

	rhcosImages = "https://raw.githubusercontent.com/openshift/installer/release-4.6/data/data/rhcos.json"
)

var roles = []string{"master", "worker"}

func getConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error
	if kubeconfig, ok := os.LookupEnv("KUBECONFIG"); ok {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	return config, nil
}

func updateUserData(ctx context.Context, role string, cs *kubernetes.Clientset) (string, error) {
	name := role + userDataSuffix
	secrets := cs.CoreV1().Secrets(machineAPINamespace)
	targetName := name + suffix
	s, err := secrets.Get(ctx, targetName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", err
		}
	} else {
		return targetName, nil
	}

	s, err = secrets.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	config := s.Data[userDataKey]

	// parse
	cfg, rpt, err := v2_4.Parse(config)
	fmt.Fprintf(os.Stderr, "%s", rpt.String())
	if err != nil || rpt.IsFatal() {
		return "", fmt.Errorf("Error parsing spec v2 config: %w\n%v", err, rpt)
	}

	newCfg, err := v24tov31.Translate(cfg, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to translate config from 2 to 3: %w", err)
	}
	dataOut, err := json.Marshal(newCfg)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal json: %w", err)
	}
	s.Data[userDataKey] = dataOut

	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      targetName,
			Namespace: s.Namespace,
		},
		Data: s.Data,
	}

	if _, err := secrets.Create(context.Background(), newSecret, metav1.CreateOptions{}); err != nil {
		return "", err
	}

	return targetName, nil
}

func machineSetClient(dc dynamic.Interface) dynamic.ResourceInterface {
	machineSetClient := dc.Resource(schema.GroupVersionResource{Group: machineAPIGroup, Resource: "machinesets", Version: "v1beta1"})
	return machineSetClient.Namespace(machineAPINamespace)
}

func objects(from *objx.Value) []objx.Map {
	var values []objx.Map
	switch {
	case from.IsObjxMapSlice():
		return from.ObjxMapSlice()
	case from.IsInterSlice():
		for _, i := range from.InterSlice() {
			if msi, ok := i.(map[string]interface{}); ok {
				values = append(values, objx.Map(msi))
			}
		}
	}
	return values
}

func run(ctx context.Context) error {
	config, err := getConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	configV1Client, err := openshiftv1.NewForConfig(config)
	if err != nil {
		return err
	}
	cvs, err := configV1Client.ClusterVersions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	cv := cvs.Items[0]
	channel := cv.Spec.Channel
	fmt.Printf("OpenShift channel: %s\n", channel)

	newBootimage := bootimageFromChannel(channel)
	if newBootimage == nil {
		return fmt.Errorf("No updated bootimages known for channel %s", channel)
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	targetSecrets := make(map[string]string)
	for _, role := range roles {
		target, err := updateUserData(ctx, role, clientset)
		if err != nil {
			return fmt.Errorf("Failed to generate user-data secret for role %s: %w", role, err)
		}
		targetSecrets[role] = target
	}

	machineSetClient := machineSetClient(dc)
	obj, err := machineSetClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ms := range objects(objx.Map(obj.UnstructuredContent()).Get("items")) {
		udSelector := "spec.template.spec.providerSpec.value.ignition.userDataSecret"
		name := (*ms.Get("metadata.name")).Data().(string)
		curSecret := (*ms.Get(udSelector)).Data().(string)
		labels := (*ms.Get("spec.template.metadata.labels")).Data().(map[string]interface{})
		roleV, ok := labels[machineLabelRole]
		if !ok {
			fmt.Printf("Skipping machineset %s with no role label %s\n", name, machineLabelRole)
			continue
		}
		role := roleV.(string)
		target, ok := targetSecrets[role]
		if !ok {
			fmt.Printf("Skipping machineset %s with unhandled role %s\n", name, role)
			continue
		}
		if curSecret == target {
			fmt.Printf("machineset %s already uses %s\n", name, target)
			continue
		}

		ms.Set(udSelector, target)
		fmt.Printf("Updating machineset %s to use user-data secret %s\n", name, target)
		v := unstructured.Unstructured {
			Object: ms.Value().MSI(),
		}
		machineSetClient.Update(ctx, &v, metav1.UpdateOptions{})
	}

	return nil
}

func main() {
	err := run(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
