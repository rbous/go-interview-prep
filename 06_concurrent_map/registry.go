package concurrent_map

import "sync"

// PackageRegistry tracks installed packages and their versions.
// Multiple goroutines may read and write to it concurrently.

type PackageRegistry struct {
	packages map[string]string
	mu sync.RWMutex
}

func NewPackageRegistry() *PackageRegistry {
	return &PackageRegistry{
		packages: make(map[string]string),
	}
}

func (r *PackageRegistry) Install(name, version string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.packages[name] = version
}

func (r *PackageRegistry) Get(name string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.packages[name]
	return v, ok
}

func (r *PackageRegistry) List() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tmp := make(map[string]string, len(r.packages))
	for k, v := range r.packages {
		tmp[k] = v
	}
	return tmp
}

func (r *PackageRegistry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.packages, name)
}

// BulkInstall installs packages concurrently from a map of name->version.
func (r *PackageRegistry) BulkInstall(pkgs map[string]string) {
	var wg sync.WaitGroup
	for name, version := range pkgs {
		wg.Add(1)
		go func(n, v string) {
			defer wg.Done()
			r.Install(n, v)
		}(name, version)
	}
	wg.Wait()
}
