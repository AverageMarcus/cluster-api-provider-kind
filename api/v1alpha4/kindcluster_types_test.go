/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha4

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKindClusterNamespacedName(t *testing.T) {

	tests := []struct {
		cluster *KindCluster
		want    string
	}{
		{
			cluster: &KindCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
			},
			want: "default-test-cluster",
		},
		{
			cluster: &KindCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-cluster",
					Namespace: "example-namespace",
				},
			},
			want: "example-namespace-my-cluster",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if tt.cluster.NamespacedName() != tt.want {
				t.Errorf("unexpected result - wanted %+v, got %+v", tt.want, tt.cluster.NamespacedName())
			}
		})
	}
}
