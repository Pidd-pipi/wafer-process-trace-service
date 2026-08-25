package main

import "testing"

func TestNormalizeInitializesLabelsMap(t *testing.T) {
	r := normalizeOpsRecord(OpsRecord{ID: "lot-x", Subject: "s", Owner: "o"})
	if r.Labels == nil {
		t.Fatal("normalizeOpsRecord left Labels as nil map")
	}
	r.Labels["site"] = "fab-a"
	if r.Labels["site"] != "fab-a" {
		t.Fatal("Labels map is not writable after normalization")
	}
}

func TestCloneIsolatesLabels(t *testing.T) {
	r := OpsRecord{ID: "lot-x", Labels: map[string]string{"site": "fab-a"}}
	c := r.Clone()
	r.Labels["site"] = "fab-b"
	if c.Labels["site"] != "fab-a" {
		t.Fatalf("clone shares Labels with source: got %q", c.Labels["site"])
	}
}

func TestSnapshotCountsByStatus(t *testing.T) {
	s := newOpsService([]OpsRecord{
		{ID: "a", Owner: "o", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "x"}, Status: OpsStatusActive},
		{ID: "b", Owner: "o", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "x"}, Status: OpsStatusClosed},
	})
	snap := s.Snapshot()
	if snap.ByStatus[OpsStatusActive] != 1 || snap.ByStatus[OpsStatusClosed] != 1 {
		t.Fatalf("ByStatus aggregation wrong: %v", snap.ByStatus)
	}
}

func TestSnapshotCountsByPriority(t *testing.T) {
	s := newOpsService([]OpsRecord{
		{ID: "a", Owner: "o", Priority: OpsPriorityHigh, Labels: map[string]string{"site": "x"}, Status: OpsStatusActive},
		{ID: "b", Owner: "o", Priority: OpsPriorityLow, Labels: map[string]string{"site": "x"}, Status: OpsStatusActive},
	})
	snap := s.Snapshot()
	if snap.ByPriority[OpsPriorityHigh] != 1 || snap.ByPriority[OpsPriorityLow] != 1 {
		t.Fatalf("ByPriority aggregation wrong: %v", snap.ByPriority)
	}
}

func TestSetLabelGuardsNilMap(t *testing.T) {
	r := OpsRecord{ID: "lot-x", Labels: nil}
	r.SetLabel("site", "fab-a")
	if r.Labels == nil || r.Labels["site"] != "fab-a" {
		t.Fatalf("SetLabel did not initialize nil map: %v", r.Labels)
	}
}
