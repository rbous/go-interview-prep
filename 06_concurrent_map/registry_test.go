package concurrent_map

import (
	"fmt"
	"sync"
	"testing"
)

func TestInstallAndGet(t *testing.T) {
	r := NewPackageRegistry()
	r.Install("nginx", "1.25.0")

	v, ok := r.Get("nginx")
	if !ok || v != "1.25.0" {
		t.Errorf("Get(nginx) = %q, %v; want 1.25.0, true", v, ok)
	}
}

func TestConcurrentInstall(t *testing.T) {
	r := NewPackageRegistry()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.Install(fmt.Sprintf("pkg-%d", i), "1.0.0")
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		v, ok := r.Get(fmt.Sprintf("pkg-%d", i))
		if !ok || v != "1.0.0" {
			t.Errorf("pkg-%d: got %q, %v; want 1.0.0, true", i, v, ok)
		}
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	r := NewPackageRegistry()
	r.Install("base", "1.0.0")

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			r.Install(fmt.Sprintf("pkg-%d", i), "2.0.0")
		}(i)
		go func() {
			defer wg.Done()
			r.Get("base")
		}()
	}
	wg.Wait()
}

func TestBulkInstall(t *testing.T) {
	r := NewPackageRegistry()
	pkgs := map[string]string{
		"curl":    "8.0.0",
		"wget":    "1.21.0",
		"openssl": "3.1.0",
	}
	r.BulkInstall(pkgs)

	for name, version := range pkgs {
		v, ok := r.Get(name)
		if !ok || v != version {
			t.Errorf("Get(%s) = %q, %v; want %s, true", name, v, ok, version)
		}
	}
}

func TestRemoveConcurrent(t *testing.T) {
	r := NewPackageRegistry()
	for i := 0; i < 50; i++ {
		r.Install(fmt.Sprintf("pkg-%d", i), "1.0.0")
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.Remove(fmt.Sprintf("pkg-%d", i))
		}(i)
	}
	wg.Wait()
}

func TestListReturnsCopy(t *testing.T) {
	r := NewPackageRegistry()
	r.Install("a", "1.0")
	r.Install("b", "2.0")

	list := r.List()
	// Mutating the returned map should NOT affect the registry
	delete(list, "a")

	_, ok := r.Get("a")
	if !ok {
		t.Error("deleting from List() result should not affect registry")
	}
}
