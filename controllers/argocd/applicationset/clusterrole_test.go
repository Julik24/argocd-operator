package applicationset

import (
	"github.com/argoproj-labs/argocd-operator/tests/test"
	"testing"
)

func TestApplicationSetReconciler_ReconcileClusterRole(t *testing.T) {
	resourceName = test.TestArgoCDName
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
			err := reconciler.reconcileClusterRole()
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

func TestApplicationSetReconciler_DeleteClusterRole(t *testing.T) {
	resourceName = test.TestArgoCDName
	ns := test.MakeTestNamespace(nil)
	tests := []struct {
		name            string
		webhook_enabled bool
		role_name       string
		wantErr         bool
	}{
		{
			name:            "successful delete",
			webhook_enabled: true,
			role_name:       ns.Name,
			wantErr:         false,
		},
		{
			name:            "successful reconcile",
			webhook_enabled: true,
			role_name:       "some unexisting role",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reconciler := makeTestApplicationSetReconciler(
				t, tt.webhook_enabled,
			)
			reconciler.varSetter()
			err := reconciler.deleteClusterRole(tt.role_name)
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
