package notifications

import (
	"github.com/argoproj-labs/argocd-operator/api/v1alpha1"
	argoproj "github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/argoproj-labs/argocd-operator/pkg/util"
	"github.com/argoproj-labs/argocd-operator/tests/test"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

type SchemeOpt func(*runtime.Scheme) error

func makeTestReconcilerScheme(sOpts ...SchemeOpt) *runtime.Scheme {
	s := scheme.Scheme
	for _, opt := range sOpts {
		_ = opt(s)
	}

	return s
}

func makeTestReconcilerClient(sch *runtime.Scheme, resObjs, subresObjs []client.Object, runtimeObj []runtime.Object) client.Client {
	client := fake.NewClientBuilder().WithScheme(sch)
	if len(resObjs) > 0 {
		client = client.WithObjects(resObjs...)
	}
	if len(subresObjs) > 0 {
		client = client.WithStatusSubresource(subresObjs...)
	}
	if len(runtimeObj) > 0 {
		client = client.WithRuntimeObjects(runtimeObj...)
	}
	return client.Build()
}

func MakeTestNotificationsReconciler(cr *argoproj.ArgoCD, objs ...client.Object) *NotificationsReconciler {
	a := test.MakeTestNotificationsConfiguration()

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	runtimeObjs := []runtime.Object{}
	sch := makeTestReconcilerScheme(v1alpha1.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, runtimeObjs)

	return &NotificationsReconciler{
		Client:   cl,
		Scheme:   sch,
		Logger:   util.NewLogger("notifications-controller"),
		Instance: cr,
	}
}

func TestNotificationsReconciler_Reconcile(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		reconciler   *NotificationsReconciler
		wantErr      bool
	}{
		{
			name:         "successful reconcile",
			resourceName: test.TestArgoCDName,
			reconciler: MakeTestNotificationsReconciler(
				test.MakeTestArgoCD(nil),
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reconciler.VarSetter()
			err := tt.reconciler.Reconcile()
			assert.NoError(t, err)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("Expected error but did not get one")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNotificationsReconciler_DeleteResources(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		reconciler   *NotificationsReconciler
		wantErr      bool
	}{
		{
			name:         "successful reconcile",
			resourceName: test.TestArgoCDName,
			reconciler: MakeTestNotificationsReconciler(
				test.MakeTestArgoCD(nil),
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reconciler.VarSetter()
			err := tt.reconciler.DeleteResources()
			assert.NoError(t, err)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("Expected error but did not get one")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
