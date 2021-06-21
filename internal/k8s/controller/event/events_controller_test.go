package event

import (
	"context"
	"testing"

	"github.com/object88/lighthouse/internal/logging/testlogger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Test_Controller_Event(t *testing.T) {
	namespace := "foo"

	e := v1.Event{}

	rs := &ReconcileEvent{
		Client: fake.NewFakeClient(&e),
		Log:    testlogger.TestLogger{T: t},
		// VersionedClient: fake.NewSimpleClientset(rel),
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "secretname",
			Namespace: namespace,
		},
	}

	if result, err := rs.Reconcile(context.TODO(), req); err != nil {
		t.Fatalf("Unexpected error while reconciling: %s", err.Error())
	} else if result.Requeue {
		t.Errorf("Unexpectedly set to requeue")
	}
}
