package notifications

import (
	"context"
	"github.com/argoproj-labs/argocd-operator/api/v1alpha1"
	argoproj "github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/argoproj-labs/argocd-operator/pkg/util"
	"github.com/argoproj-labs/argocd-operator/tests/test"
	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	//subresObjs := []client.Object{a}
	runtimeObjs := []runtime.Object{}
	sch := makeTestReconcilerScheme(v1alpha1.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, objs, runtimeObjs)

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
		/*{
			name:         "not successful reconcile",
			resourceName: test.TestArgoCDName,
			reconciler: MakeTestNotificationsReconciler(
				test.MakeTestArgoCD(nil),
			),
			wantErr: true,
		},*/
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

/*func TestNotificationsReconciler_ReconcileBad(t *testing.T) {
	a := makeTestNotificationsConfiguration(func(a *v1alpha1.NotificationsConfiguration) {
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	runtimeObjs := []runtime.Object{}
	sch := makeTestReconcilerScheme(v1alpha1.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, runtimeObjs)
	r := makeTestReconciler(cl, sch)

	a.Spec = v1alpha1.NotificationsConfigurationSpec{
		// Add a default template for test
		Templates: map[string]string{
			"template.app-created": `email:
			subject: Application {{.app.metadata.name}} has been created.
		  message: Application {{.app.metadata.name}} has been created.
		  teams:
			title: Application {{.app.metadata.name}} has been created.`,
		},
		// Add a default template for test
		Triggers: map[string]string{
			"trigger.on-created": `- description: Application is created.
			oncePer: app.metadata.name
			send:
			- app-created
			when: "true"`,
		},
	}

	err := r.reconcileNotificationsConfigmap(a)
	assert.NoError(t, err)

	// Verify if the ConfigMap is created
	testCM := &corev1.ConfigMap{}
	assert.NoError(t, r.Client.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      ArgoCDNotificationsConfigMap,
			Namespace: a.Namespace,
		},
		testCM))

	// Verify that the configmap has the default template
	assert.NotEqual(t, testCM.Data["template.app-created"], "")

	// Verify that the configmap has the default trigger
	assert.NotEqual(t, testCM.Data["trigger.on-created"], "")
}*/

func TestReconcileNotifications_CreateRoles(t *testing.T) {
	a := test.MakeTestArgoCD(nil)

	r := MakeTestNotificationsReconciler(a)
	err := r.reconcileRole()
	assert.NoError(t, err)

	err = r.reconcileDeployment()
	assert.NoError(t, err)

	testRole := &rbacv1.Role{}

	assert.NoError(t, r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      "test-argocd",
		Namespace: a.Namespace,
	}, testRole))

	/*desiredPolicyRules := policyRuleForNotificationsController()

	assert.Equal(t, desiredPolicyRules, testRole.Rules)

	a.Spec.Notifications.Enabled = false
	_, err = r.reconcileNotificationsRole(a)
	assert.NoError(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      generateResourceName(common.ArgoCDNotificationsControllerComponent, a),
		Namespace: a.Namespace,
	}, testRole)
	assert.True(t, errors.IsNotFound(err))*/
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
