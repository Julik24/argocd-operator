package applicationset

import (
	"github.com/argoproj-labs/argocd-operator/tests/test"
	"testing"
)

func TestApplicationSetReconciler_reconcileWebhookRoute(t *testing.T) {
	ns := test.MakeTestNamespace(nil)
	//asr := makeTestApplicationSetReconciler(t, true, ns)

	//existingWebhookRoute := asr.getDesiredWebhookRoute()

	tests := []struct {
		name                      string
		webhookServerRouteEnabled bool
		setupClient               func(bool) *ApplicationSetReconciler
		wantErr                   bool
	}{
		{
			name:                      "create a webhookRoute",
			webhookServerRouteEnabled: true,
			setupClient: func(webhookServerRouteEnabled bool) *ApplicationSetReconciler {
				return makeTestApplicationSetReconciler(t, webhookServerRouteEnabled, ns)
			},
			wantErr: false,
		},
		/*{
			name:                      "update a webhookRoute",
			webhookServerRouteEnabled: true,
			setupClient: func(webhookServerRouteEnabled bool) *ApplicationSetReconciler {
				outdatedWebhookRoute := existingWebhookRoute
				outdatedWebhookRoute.ObjectMeta.Labels = test.TestKVP
				return makeTestApplicationSetReconciler(t, webhookServerRouteEnabled, outdatedWebhookRoute, ns)
			},
			wantErr: false,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asr := tt.setupClient(tt.webhookServerRouteEnabled)
			err := asr.reconcileWebhookRoute()
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("Expected error but did not get one")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			/*updatedWebhookRoute := &routev1.Route{}
			err = asr.Client.Get(context.TODO(), types.NamespacedName{Name: AppSetWebhookRouteName, Namespace: test.TestNamespace}, updatedWebhookRoute)
			if err != nil {
				t.Fatalf("Could not get updated WebhookRoute: %v", err)
			}
			assert.Equal(t, testExpectedLabels, updatedWebhookRoute.ObjectMeta.Labels)*/
		})
	}
}

func TestApplicationSetReconciler_deleteWebhookRoute(t *testing.T) {
	ns := test.MakeTestNamespace(nil)
	tests := []struct {
		name                      string
		webhookServerRouteEnabled bool
		routeName                 string
		setupClient               func(bool) *ApplicationSetReconciler
		wantErr                   bool
	}{
		{
			name:                      "successful delete",
			webhookServerRouteEnabled: true,
			routeName:                 "test_name",
			setupClient: func(webhookServerRouteEnabled bool) *ApplicationSetReconciler {
				return makeTestApplicationSetReconciler(t, webhookServerRouteEnabled, ns)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asr := tt.setupClient(tt.webhookServerRouteEnabled)
			if err := asr.deleteRoute(tt.routeName, ns.Name); (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("Expected error but did not get one")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
