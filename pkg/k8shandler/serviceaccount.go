package k8shandler

import (
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	v1alpha1 "github.com/t0ffel/elasticsearch-operator/pkg/apis/elasticsearch/v1alpha1"
)

// CreateOrUpdateServiceAccount ensures the existence of the serviceaccount for Elasticsearch cluster
func CreateOrUpdateServiceAccount(dpl *v1alpha1.Elasticsearch) error {
	// In case no serviceaccount is specified in the spec, we'll use the namespace's default service account
	if dpl.Spec.ServiceAccountName == "" {
		return nil
	}

	owner := asOwner(dpl)

	err := createOrUpdateServiceAccount(dpl.Spec.ServiceAccountName, dpl.Namespace, owner)
	if err != nil {
		return fmt.Errorf("Failure creating ServiceAccount %v", err)
	}

	return nil
}

func createOrUpdateServiceAccount(serviceAccountName, namespace string, owner metav1.OwnerReference) error {
	elasticsearchSA := serviceAccount(serviceAccountName, namespace)
	addOwnerRefToObject(elasticsearchSA, owner)
	err := sdk.Create(elasticsearchSA)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("Failure constructing serviceaccount for the Elasticsearch cluster: %v", err)
	}
	return nil
}

// serviceAccount returns a v1.ServiceAccount object
func serviceAccount(serviceAccountName string, namespace string) *v1.ServiceAccount {
	return &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
	}
}