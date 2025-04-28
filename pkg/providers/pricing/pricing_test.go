/*
Copyright 2025 The CloudPilot AI Authors.

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

package pricing_test

import (
	"context"
	"net/http"
	"slices"
	"testing"

	"github.com/cloudpilot-ai/karpenter-provider-gcp/pkg/providers/pricing"
)

var testInstanceTypes = []string{
	"e2-standard-32",
	"n2-standard-16",
	"c2-standard-8",
	"m1-ultramem-40",
	"a2-highgpu-1g",
}

func TestDefaultProvider_LivenessProbe(t *testing.T) {
	provider := pricing.NewDefaultProvider(context.Background(), "europe-west4")
	if err := provider.LivenessProbe(&http.Request{}); err != nil {
		t.Errorf("LivenessProbe failed: %v", err)
	}
}

func TestDefaultProvider_InitialPrices(t *testing.T) {
	// Initialize the provider with europe-west4 region
	provider := pricing.NewDefaultProvider(context.Background(), "europe-west4")

	// Get all instance types
	instanceTypes := provider.InstanceTypes()
	if len(instanceTypes) == 0 {
		t.Fatal("No instance types found")
	}

	// Test getting ondemand prices for various instance types
	for _, instanceType := range testInstanceTypes {
		_, found := provider.OnDemandPrice(instanceType)
		if !found {
			t.Errorf("Failed to find on-demand price for %s", instanceType)
			continue
		}
	}

	// Test getting spot prices for various instance types
	// Should fail, cause inital prices doesn't contain spot prices
	for _, instanceType := range testInstanceTypes {
		_, found := provider.SpotPrice(instanceType, "europe-west4-a")
		if found {
			t.Errorf("Expected to not find spot price for %s instance type, while prices wasn't updated", instanceType)
			continue
		}
	}

	// Test getting price for a non-existent instance type
	_, found := provider.OnDemandPrice("non-existent-type")
	if found {
		t.Error("Expected to not find price for non-existent instance type")
	}
}

func TestDefaultProvider_OnDemandPrice(t *testing.T) {
	// Initialize the provider with europe-west4 region
	provider := pricing.NewDefaultProvider(context.Background(), "europe-west4")

	// Test price update
	if err := provider.UpdateOnDemandPricing(context.Background()); err != nil {
		t.Fatalf("Failed to update on-demand pricing: %v", err)
	}

	// Test getting prices for various instance types
	for _, instanceType := range testInstanceTypes {
		_, found := provider.OnDemandPrice(instanceType)
		if !found {
			t.Errorf("Failed to find on-demand price for %s", instanceType)
			continue
		}
	}

	// Test getting price for a non-existent instance type
	_, found := provider.OnDemandPrice("non-existent-type")
	if found {
		t.Error("Expected to not find price for non-existent instance type")
	}
}

func TestDefaultProvider_SpotPrice(t *testing.T) {
	// Initialize the provider with europe-west4 region
	provider := pricing.NewDefaultProvider(context.Background(), "europe-west4")

	// Test price update
	if err := provider.UpdateSpotPricing(context.Background()); err != nil {
		t.Fatalf("Failed to update spot pricing: %v", err)
	}

	// Test getting spot prices for various instance types
	for _, instanceType := range testInstanceTypes {
		_, found := provider.SpotPrice(instanceType, "europe-west4-a")
		if !found {
			t.Errorf("Failed to find spot price for %s", instanceType)
			continue
		}
	}

	// Test getting price for a non-existent instance type
	_, found := provider.SpotPrice("non-existent-type", "europe-west4-a")
	if found {
		t.Error("Expected to not find price for non-existent instance type")
	}
}

func TestDefaultProvider_InstanceTypes(t *testing.T) {
	// Initialize the provider with europe-west4 region
	provider := pricing.NewDefaultProvider(context.Background(), "europe-west4")

	// Updating prices to retrieve a instance types from runtime prices
	if err := provider.UpdateOnDemandPricing(context.Background()); err != nil {
		t.Fatalf("Failed to update on-demand pricing: %v", err)
	}

	// Get all instance types
	instanceTypes := provider.InstanceTypes()
	if len(instanceTypes) == 0 {
		t.Fatal("No instance types found")
	}

	// Check if all test instance types are in the list
	for _, testInstanceType := range testInstanceTypes {
		if !slices.Contains(instanceTypes, testInstanceType) {
			t.Errorf("%s not found in instance types list", testInstanceType)
		}
	}
}
