package applicationset

import (
	"github.com/argoproj-labs/argocd-operator/pkg/resource"
	"github.com/argoproj-labs/argocd-operator/tests/mock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	argoproj "github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/argoproj-labs/argocd-operator/common"
	"github.com/argoproj-labs/argocd-operator/pkg/util"
	"github.com/argoproj-labs/argocd-operator/tests/test"
)

var testExpectedLabels = common.DefaultResourceLabels(test.TestArgoCDName, test.TestNamespace, common.AppSetControllerComponent)

func makeTestApplicationSetReconciler(t *testing.T, webhookServerRouteEnabled bool, objs ...runtime.Object) *ApplicationSetReconciler {
	s := scheme.Scheme

	assert.NoError(t, routev1.Install(s))
	assert.NoError(t, argoproj.AddToScheme(s))

	cl := fake.NewClientBuilder().WithScheme(s).WithRuntimeObjects(objs...).Build()

	return &ApplicationSetReconciler{
		Client: cl,
		Scheme: s,
		Logger: util.NewLogger("appset-controller"),
		Instance: test.MakeTestArgoCD(nil, func(a *argoproj.ArgoCD) {
			a.Spec.ApplicationSet = &argoproj.ArgoCDApplicationSet{
				WebhookServer: argoproj.WebhookServerSpec{
					Route: argoproj.ArgoCDRouteSpec{
						Enabled: webhookServerRouteEnabled,
					},
				},
			}
		}),
	}
}

func TestApplicationSetReconciler_Reconcile(t *testing.T) {
	/*expectedResources := []client.Object{
		test.MakeTestServiceAccount(
			func(sa *corev1.ServiceAccount) {
				sa.Name = resourceName
			},
		),
		test.MakeTestService(nil,
			func(s *corev1.Service) {
				s.Name = resourceName
			},
		),
		test.MakeTestDeployment(nil,
			func(d *appsv1.Deployment) {
				d.Name = resourceName
			},
		),
	}*/

	tests := []struct {
		name              string
		webhookEnabled    bool
		expectedResources []client.Object
	}{
		{
			name:           "no webhook enabled",
			webhookEnabled: false,
			expectedResources: []client.Object{
				test.MakeTestRole(nil,
					func(r *rbacv1.Role) {
						r.Name = "test-argocd-appset"
					},
				),
				test.MakeTestServiceAccount(
					func(sa *corev1.ServiceAccount) {
						sa.Name = "test-argocd-appset"
					},
				),
				test.MakeTestRoleBinding(nil,
					func(rb *rbacv1.RoleBinding) {
						rb.Name = "test-argocd-appset"
					},
				),
				test.MakeTestService(nil,
					func(s *corev1.Service) {
						s.Name = "test-argocd-appset"
					},
				),
				test.MakeTestDeployment(nil,
					func(d *appsv1.Deployment) {
						d.Name = "test-argocd-appset"
					},
				),
			},
		},
		{
			name:           "webhook enabled",
			webhookEnabled: true,
			expectedResources: []client.Object{
				test.MakeTestRole(nil,
					func(r *rbacv1.Role) {
						r.Name = "test-argocd-redis-ha"
					},
				),
				test.MakeTestServiceAccount(
					func(sa *corev1.ServiceAccount) {
						sa.Name = "test-argocd-redis"
					},
				),
				test.MakeTestRoleBinding(nil,
					func(rb *rbacv1.RoleBinding) {
						rb.Name = "test-argocd-redis"
					},
				),
				test.MakeTestService(nil,
					func(s *corev1.Service) {
						s.Name = "test-argocd-redis"
					},
				),
				test.MakeTestDeployment(nil,
					func(d *appsv1.Deployment) {
						d.Name = "test-argocd-redis-ha-haproxy"
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoServerName := "test-argocd-reposerver"
			reconciler := makeTestApplicationSetReconciler(
				t, tt.webhookEnabled,
			)

			reconciler.varSetter()
			reconciler.RepoServer = mock.NewRepoServer(mockRepoServerName, test.TestNamespace, reconciler.Client)
			err := reconciler.Reconcile()
			assert.NoError(t, err)

			/*for _, obj := range tt.expectedResources {
				_, err := resource.GetObject(obj.GetName(), test.TestNamespace, obj, reconciler.Client)
				assert.NoError(t, err)
			}*/
		})
	}
}

func TestApplicationSetReconciler_DeleteResources(t *testing.T) {
	tests := []struct {
		name              string
		webhookEnabled    bool
		expectedResources []client.Object
	}{
		{
			name:           "no webhook enabled",
			webhookEnabled: false,
			expectedResources: []client.Object{
				test.MakeTestRole(nil,
					func(r *rbacv1.Role) {
						r.Name = "test-argocd-appset"
					},
				),
				test.MakeTestServiceAccount(
					func(sa *corev1.ServiceAccount) {
						sa.Name = "test-argocd-appset"
					},
				),
				test.MakeTestRoleBinding(nil,
					func(rb *rbacv1.RoleBinding) {
						rb.Name = "test-argocd-appset"
					},
				),
				test.MakeTestService(nil,
					func(s *corev1.Service) {
						s.Name = "test-argocd-appset"
					},
				),
				test.MakeTestDeployment(nil,
					func(d *appsv1.Deployment) {
						d.Name = "test-argocd-appset"
					},
				),
			},
		},
		{
			name:           "webhook enabled",
			webhookEnabled: true,
			expectedResources: []client.Object{
				test.MakeTestRole(nil,
					func(r *rbacv1.Role) {
						r.Name = "test-argocd-redis-ha"
					},
				),
				test.MakeTestServiceAccount(
					func(sa *corev1.ServiceAccount) {
						sa.Name = "test-argocd-redis"
					},
				),
				test.MakeTestRoleBinding(nil,
					func(rb *rbacv1.RoleBinding) {
						rb.Name = "test-argocd-redis"
					},
				),
				test.MakeTestService(nil,
					func(s *corev1.Service) {
						s.Name = "test-argocd-redis"
					},
				),
				test.MakeTestDeployment(nil,
					func(d *appsv1.Deployment) {
						d.Name = "test-argocd-redis-ha-haproxy"
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoServerName := "test-argocd-reposerver"
			reconciler := makeTestApplicationSetReconciler(
				t, tt.webhookEnabled,
			)

			reconciler.varSetter()
			reconciler.RepoServer = mock.NewRepoServer(mockRepoServerName, test.TestNamespace, reconciler.Client)
			err := reconciler.DeleteResources()
			assert.NoError(t, err)

			for _, obj := range tt.expectedResources {
				_, err := resource.GetObject(obj.GetName(), test.TestNamespace, obj, reconciler.Client)
				assert.True(t, apierrors.IsNotFound(err))
			}
		})
	}
}
