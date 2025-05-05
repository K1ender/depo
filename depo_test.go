package depo_test

import (
	"sync"
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

type mutexDep struct {
	val int
	mu  sync.Mutex
}

func newMutexDep() *mutexDep {
	return &mutexDep{
		val: 0,
	}
}

func TestMutex(t *testing.T) {
	dep := newMutexDep()
	depo := depo.New(dep)
	wg := sync.WaitGroup{}
	depo.Use(func(dep *mutexDep) {
		for range 100 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				dep.mu.Lock()
				dep.val++
				dep.mu.Unlock()
			}()
		}
	})
	wg.Wait()
	depo.Use(func(dep *mutexDep) {
		dep.mu.Lock()
		if dep.val != 100 {
			t.Errorf("expected dep.val to be 100, got %d", dep.val)
		}
		dep.mu.Unlock()
	})
}

type testDep2 struct {
	val string
}

func newTestDep2() *testDep2 {
	return &testDep2{
		val: "test",
	}
}

func TestHas(t *testing.T) {
	dep := newTestDep2()
	depo := depo.New(dep)
	if !depo.Has(dep) {
		t.Errorf("expected depo.Has(dep) to be true")
	}
	if depo.Has(newTestDep()) {
		t.Errorf("expected depo.Has(newTestDep()) to be false")
	}
}
