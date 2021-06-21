/*
LICENSE
*/
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/object88/lighthouse/pkg/k8s/apis/engineering.lighthouse/v1alpha1"
	"github.com/object88/lighthouse/pkg/k8s/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type EngineeringV1alpha1Interface interface {
	RESTClient() rest.Interface
	ConfigsGetter
}

// EngineeringV1alpha1Client is used to interact with features provided by the engineering.lighthouse group.
type EngineeringV1alpha1Client struct {
	restClient rest.Interface
}

func (c *EngineeringV1alpha1Client) Configs(namespace string) ConfigInterface {
	return newConfigs(c, namespace)
}

// NewForConfig creates a new EngineeringV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*EngineeringV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &EngineeringV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new EngineeringV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *EngineeringV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new EngineeringV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *EngineeringV1alpha1Client {
	return &EngineeringV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *EngineeringV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
