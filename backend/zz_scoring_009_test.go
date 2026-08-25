package main

import "testing"

func TestCacheSizeNeverExceedsLimit(t *testing.T) {
	c := newOpsCache(4)
	for i := 0; i < 20; i++ {
		c.Put(string(rune('a'+i%26))+string(rune('0'+i)), OpsRecord{ID: string(rune('a' + i%26))})
	}
	if got := c.Size(); got > 4 {
		t.Fatalf("cache exceeded limit: size %d", got)
	}
}

func TestCacheEvictsOldestEntry(t *testing.T) {
	c := newOpsCache(3)
	c.Put("first", OpsRecord{ID: "first"})
	c.Put("second", OpsRecord{ID: "second"})
	c.Put("third", OpsRecord{ID: "third"})
	c.Put("fourth", OpsRecord{ID: "fourth"})
	if _, ok := c.Get("first"); ok {
		t.Fatal("oldest entry was not evicted")
	}
	if _, ok := c.Get("fourth"); !ok {
		t.Fatal("newest entry should still be present")
	}
}

func TestCacheGetReturnsCopy(t *testing.T) {
	c := newOpsCache(4)
	c.Put("lot-x", OpsRecord{ID: "lot-x", Labels: map[string]string{"site": "fab-a"}})
	got, ok := c.Get("lot-x")
	if !ok {
		t.Fatal("expected cached entry")
	}
	got.Labels["site"] = "mutated"
	again, _ := c.Get("lot-x")
	if again.Labels["site"] == "mutated" {
		t.Fatal("Get returned a shared reference into the cache")
	}
}
