package notifications

import (
	"github.com/argoproj-labs/argocd-operator/tests/test"
	"testing"
)

func TestNotificationsReconciler_reconcileStatus(t *testing.T) {
	resourceName = test.TestArgoCDName
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
			err := tt.reconciler.ReconcileStatus()
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
