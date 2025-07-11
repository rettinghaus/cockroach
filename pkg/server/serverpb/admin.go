// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package serverpb

import (
	context "context"
	"sort"

	"github.com/cockroachdb/cockroach/pkg/util/metric"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

// Add adds values from ots to ts.
func (ts *TableStatsResponse) Add(ots *TableStatsResponse) {
	ts.RangeCount += ots.RangeCount
	ts.ReplicaCount += ots.ReplicaCount
	ts.ApproximateDiskBytes += ots.ApproximateDiskBytes
	ts.Stats.Add(ots.Stats)

	// The stats in TableStatsResponse were generated by getting separate stats
	// for each node, then aggregating them into TableStatsResponse.
	// So resulting NodeCount should be the same, unless ots contains nodeData
	// in MissingNodes that isn't already tracked in ts.MissingNodes.
	// Note: when comparing missingNode objects, there's a chance that the nodeId
	// could be the same, but that the error messages differ. Keeping the first
	// and dropping subsequent ones seems reasonable to do, and is what is done
	// here.
	missingNodeIds := make(map[string]struct{})
	for _, nodeData := range ts.MissingNodes {
		missingNodeIds[nodeData.NodeID] = struct{}{}
	}
	for _, nodeData := range ots.MissingNodes {
		if _, found := missingNodeIds[nodeData.NodeID]; !found {
			ts.MissingNodes = append(ts.MissingNodes, nodeData)
			ts.NodeCount--
		}
	}
}

func (r DecommissionPreCheckResponse_NodeReadiness) String() string {
	switch r {
	case DecommissionPreCheckResponse_UNKNOWN:
		return "unknown"
	case DecommissionPreCheckResponse_READY:
		return "ready"
	case DecommissionPreCheckResponse_ALREADY_DECOMMISSIONED:
		return "already decommissioned"
	case DecommissionPreCheckResponse_ALLOCATION_ERRORS:
		return "allocation errors"
	default:
		panic("unknown decommission node readiness")
	}
}

type TenantAdminServer interface {
	Liveness(context.Context, *LivenessRequest) (*LivenessResponse, error)
}

// Empty is true if there are no unavailable ranges and no error performing
// healthcheck.
func (r *RecoveryVerifyResponse_UnavailableRanges) Empty() bool {
	return len(r.Ranges) == 0 && len(r.Error) == 0
}

// GetInternalTimeseriesNamesFromServer is a helper that uses the provided
// ClientConn to query the AllMetricMetadata endpoint, and returns the set of
// all possible internal metric names as a sorted slice.
//
// This is *not* the list of timeseries names. Instead, it is that list but
// adding `cr.node.` and `cr.store.` prefixes (both copies are emitted, since we
// can't tell what the true prefix for each metric is). Additionally, for histograms
// we generate the names for the quantiles that are exported (internal TSDB does
// not support full histograms).
func GetInternalTimeseriesNamesFromServer(
	ctx context.Context, ac RPCAdminClient,
) ([]string, error) {
	resp, err := ac.AllMetricMetadata(ctx, &MetricMetadataRequest{})
	if err != nil {
		return nil, err
	}
	var sl []string
	for name, meta := range resp.Metadata {
		if meta.MetricType == io_prometheus_client.MetricType_HISTOGRAM {
			// See usage of HistogramMetricComputers in pkg/server/status/recorder.go.
			for _, q := range metric.HistogramMetricComputers {
				sl = append(sl, name+q.Suffix)
			}
		} else {
			sl = append(sl, name)
		}
	}
	out := make([]string, 0, 2*len(sl))
	for _, prefix := range []string{"cr.node.", "cr.store."} {
		for _, name := range sl {
			out = append(out, prefix+name)
		}
	}
	sort.Strings(out)
	return out, nil
}
