package event

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/object88/lighthouse/internal/k8s/predicates"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// blank assignment to verify that ReconcileEvent implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileEvent{}

// ReconcileEvent reconciles a Secret object
type ReconcileEvent struct {
	Client client.Client
	Log    logr.Logger
}

func (r *ReconcileEvent) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		WithLogger(r.Log).
		For(&v1.Event{}).
		WithEventFilter(predicates.ResourceGenerationOrFinalizerChangedPredicate{}).
		Complete(r)
	if err != nil {
		return err
	}

	return nil
}

// Reconcile implements reconcile.Reconciler.
func (r *ReconcileEvent) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	recLogger := r.Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	recLogger.Info("Reconciling ReleaseHistory")

	instance := &v1.Event{}
	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		// There was an error processing the request; requeue
		if !errors.IsNotFound(err) {
			recLogger.Error(err, "Error requesting release history")
			return reconcile.Result{}, err
		}
	}

	// instance.TypeMeta.Kind

	return reconcile.Result{}, nil
}
