package crds

import (
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/cache"
)

type Watcher struct {
	log       logr.Logger
	clientset *clientset.Clientset
}

func New(log logr.Logger, clientset *clientset.Clientset) *Watcher {
	// s := runtime.NewScheme()
	// v1.AddToScheme(s)

	// codecs := serializer.NewCodecFactory(s)
	// enc := unstructured.NewJSONFallbackEncoder(codecs.LegacyCodec(s.PrioritizedVersionsAllGroups()...))
	return &Watcher{
		log:       log,
		clientset: clientset,
		// encoder:   enc,
	}
}

func (w *Watcher) GetInformer() cache.SharedIndexInformer {
	factory := externalversions.NewSharedInformerFactoryWithOptions(w.clientset, 10*time.Second)

	// factory := informers.NewSharedInformerFactory(w.clientset, 10*time.Second)

	// factory.Discovery().V1beta1().EndpointSlices().Informer()
	informer := factory.Apiextensions().V1().CustomResourceDefinitions().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: w.added,
		// UpdateFunc: w.updated,
		// DeleteFunc: w.deleted,
	})

	return informer
}

func (w *Watcher) added(obj interface{}) {
	if p, ok := obj.(*v1.CustomResourceDefinition); ok {
		w.log.Info("added CRD", "name", p.Name)
	}

}
