// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package tenantcapabilitiesauthorizer

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv/kvpb"
	"github.com/cockroachdb/cockroach/pkg/multitenant/tenantcapabilities"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// AllowEverythingAuthorizer is a tenantcapabilities.Authorizer that
// allows all operations
type AllowEverythingAuthorizer struct{}

var _ tenantcapabilities.Authorizer = &AllowEverythingAuthorizer{}

// NewAllowEverythingAuthorizer constructs and returns a AllowEverythingAuthorizer.
func NewAllowEverythingAuthorizer() *AllowEverythingAuthorizer {
	return &AllowEverythingAuthorizer{}
}

// HasCrossTenantRead returns true if a tenant can read from other tenants.
func (n *AllowEverythingAuthorizer) HasCrossTenantRead(
	ctx context.Context, tenID roachpb.TenantID,
) bool {
	return true
}

// HasCapabilityForBatch implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) HasCapabilityForBatch(
	context.Context, roachpb.TenantID, *kvpb.BatchRequest,
) error {
	return nil
}

// BindReader implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) BindReader(tenantcapabilities.Reader) {}

// HasNodeStatusCapability implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) HasNodeStatusCapability(
	ctx context.Context, tenID roachpb.TenantID,
) error {
	return nil
}

// HasTSDBQueryCapability implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) HasTSDBQueryCapability(
	ctx context.Context, tenID roachpb.TenantID,
) error {
	return nil
}

// HasNodelocalStorageCapability implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) HasNodelocalStorageCapability(
	ctx context.Context, tenID roachpb.TenantID,
) error {
	return nil
}

// IsExemptFromRateLimiting implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) IsExemptFromRateLimiting(
	context.Context, roachpb.TenantID,
) bool {
	return true
}

// HasProcessDebugCapability implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) HasProcessDebugCapability(
	ctx context.Context, tenID roachpb.TenantID,
) error {
	return nil
}

// HasTSDBAllMetricsCapability implements the tenantcapabilities.Authorizer interface.
func (n *AllowEverythingAuthorizer) HasTSDBAllMetricsCapability(
	ctx context.Context, tenID roachpb.TenantID,
) error {
	return nil
}
