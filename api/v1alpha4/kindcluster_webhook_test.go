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

func TestKindClusterUpdateInvalid(t *testing.T) {
	oldCluster := KindCluster{
		ObjectMeta: metav1.ObjectMeta{},
		Spec: KindClusterSpec{
			Replicas:      1,
			Image:         "kindest/node",
			Version:       "v1.21.2",
			FeatureGates:  map[string]bool{},
			RuntimeConfig: map[string]string{},
		},
	}

	tests := []struct {
		name       string
		newCluster *KindCluster
		wantError  bool
	}{
		{
			name: "return no error if no modification",
			newCluster: func() *KindCluster {
				return oldCluster.DeepCopy()
			}(),
			wantError: false,
		},
		{
			name: "don't allow modification of replicas",
			newCluster: func() *KindCluster {
				newCluster := oldCluster.DeepCopy()
				newCluster.Spec.Replicas = 2
				return newCluster
			}(),
			wantError: true,
		},
		{
			name: "don't allow modification of image",
			newCluster: func() *KindCluster {
				newCluster := oldCluster.DeepCopy()
				newCluster.Spec.Image = "newimage"
				return newCluster
			}(),
			wantError: true,
		},
		{
			name: "don't allow modification of version",
			newCluster: func() *KindCluster {
				newCluster := oldCluster.DeepCopy()
				newCluster.Spec.Version = "v1.19.1"
				return newCluster
			}(),
			wantError: true,
		},
		{
			name: "don't allow modification of featureGates",
			newCluster: func() *KindCluster {
				newCluster := oldCluster.DeepCopy()
				newCluster.Spec.FeatureGates = map[string]bool{
					"example": true,
				}
				return newCluster
			}(),
			wantError: true,
		},
		{
			name: "don't allow modification of runtimeConfig",
			newCluster: func() *KindCluster {
				newCluster := oldCluster.DeepCopy()
				newCluster.Spec.RuntimeConfig = map[string]string{
					"example": "true",
				}
				return newCluster
			}(),
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.newCluster.ValidateUpdate(&oldCluster)
			if (err != nil) != tt.wantError {
				t.Errorf("unexpected result - wanted %+v, got %+v", tt.wantError, err)
			}
		})
	}
}
