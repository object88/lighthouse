/*
LICENSE
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/object88/lighthouse/pkg/k8s/client/clientset/versioned/typed/engineering.lighthouse/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeEngineeringV1alpha1 struct {
	*testing.Fake
}

func (c *FakeEngineeringV1alpha1) Configs(namespace string) v1alpha1.ConfigInterface {
	return &FakeConfigs{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeEngineeringV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
