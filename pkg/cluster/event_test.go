package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/argoproj-labs/argocd-operator/tests/test"
)

func TestRequestEvent(t *testing.T) {
	testClient := fake.NewClientBuilder().Build()

	tests := []struct {
		name         string
		eReq         NamespaceRequest
		desiredEvent *corev1.Event
		wantErr      bool
	}{
		{
			name: "request event",
			eReq: EventRequest{
				InvolvedObjectMeta: metav1.ObjectMeta{
					Name:   test.TestName,
					Labels: test.TestKVP,
				},
				Client: testClient,
			},
			desiredEvent: test.MakeTestEvent(nil, func(ns *corev1.Event) {
				ns.Name = test.TestName
				ns.Labels = test.TestKVP
			}),
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotEvent, err := RequestNamespace(test.eReq)

			if !test.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, test.desiredEvent, gotEvent)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
