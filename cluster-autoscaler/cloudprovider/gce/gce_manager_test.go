/*
Copyright 2016 The Kubernetes Authors.

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

package gce

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/kubernetes/pkg/api/v1"

	"github.com/stretchr/testify/assert"
)

func TestBuildGenericLabels(t *testing.T) {
	labels, err := buildGenericLabels(GceRef{
		Name:    "kubernetes-minion-group",
		Project: "mwielgus-proj",
		Zone:    "us-central1-b"},
		"n1-standard-8", "sillyname")
	assert.Nil(t, err)
	assert.Equal(t, "us-central1", labels[metav1.LabelZoneRegion])
	assert.Equal(t, "us-central1-b", labels[metav1.LabelZoneFailureDomain])
	assert.Equal(t, "sillyname", labels[metav1.LabelHostname])
	assert.Equal(t, "n1-standard-8", labels[metav1.LabelInstanceType])
	assert.Equal(t, defaultArch, labels[metav1.LabelArch])
	assert.Equal(t, defaultOS, labels[metav1.LabelOS])
}

func TestExtractLabelsFromKubeEnv(t *testing.T) {
	kubeenv := "ENABLE_NODE_PROBLEM_DETECTOR: 'daemonset'\n" +
		"NODE_LABELS: a=b,c=d,cloud.google.com/gke-nodepool=pool-3,cloud.google.com/gke-preemptible=true\n" +
		"DNS_SERVER_IP: '10.0.0.10'\n"

	labels, err := extractLabelsFromKubeEnv(kubeenv)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(labels))
	assert.Equal(t, "b", labels["a"])
	assert.Equal(t, "d", labels["c"])
	assert.Equal(t, "pool-3", labels["cloud.google.com/gke-nodepool"])
	assert.Equal(t, "true", labels["cloud.google.com/gke-preemptible"])
}

func TestBuildReadyConditions(t *testing.T) {
	conditions := buildReadyConditions()
	foundReady := false
	for _, condition := range conditions {
		if condition.Type == apiv1.NodeReady && condition.Status == apiv1.ConditionTrue {
			foundReady = true
		}
	}
	assert.True(t, foundReady)
}
