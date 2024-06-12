package applicationset

import (
	"github.com/argoproj-labs/argocd-operator/tests/test"
	"testing"
)

func TestApplicationSetReconciler_reconcileStatus(t *testing.T) {
	/*resourceName = test.TestArgoCDName
	tests := []struct {
		name            string
		webhook_enabled bool
		wantErr         bool
	}{
		{
			name:            "successful reconcile",
			webhook_enabled: true,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reconciler := makeTestApplicationSetReconciler(
				t, tt.webhook_enabled,
			)
			reconciler.varSetter()
			err := reconciler.ReconcileStatus()
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("Expected error but did not get one")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}*/

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
			name:                      "successful reconcile",
			webhookServerRouteEnabled: true,
			setupClient: func(webhookServerRouteEnabled bool) *ApplicationSetReconciler {
				return makeTestApplicationSetReconciler(t, webhookServerRouteEnabled, ns)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asr := tt.setupClient(tt.webhookServerRouteEnabled)
			asr.varSetter()
			err := asr.ReconcileStatus()
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
