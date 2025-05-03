package depo_test

import (
	"testing"

	"github.com/k1ender/depo"
)

type testDep struct {
	val string
}

func newTestDep() *testDep {
	return &testDep{
		val: "test",
	}
}

func TestInjector(t *testing.T) {
	dep := newTestDep()
	depo := depo.New(dep)
	depo.Use(func(dep *testDep) {
		if dep.val != "test" {
			t.Errorf("expected dep.val to be 'test', got '%s'", dep.val)
		}
	})
}

type kvstore struct {
	values map[string]string
}

func newKVStore() *kvstore {
	return &kvstore{
		values: map[string]string{},
	}
}

func TestPtr(t *testing.T) {
	kv := newKVStore()
	depo := depo.New(kv)
	depo.Use(func(kv *kvstore) {
		kv.values["test"] = "test"
	})
	depo.Use(func(kv *kvstore) {
		if kv.values["test"] != "test" {
			t.Errorf("expected kv.values['test'] to be 'test', got '%s'", kv.values["test"])
		}
	})
	if kv.values["test"] != "test" {
		t.Errorf("expected kv.values['test'] to be 'test', got '%s'", kv.values["test"])
	}
}
