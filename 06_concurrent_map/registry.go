package concurrent_map

import "sync"

// PackageRegistry tracks installed packages and their versions.
// Multiple goroutines may read and write to it concurrently.
//
// BUG: Concurrent map read/write causes a panic at runtime.
// Fix so all operations are safe for concurrent use.
// You may use sync.Map, sync.RWMutex, or any other approach.

type PackageRegistry struct {
	packages map[string]string
}

func NewPackageRegistry() *PackageRegistry {
	return &PackageRegistry{
		packages: make(map[string]string),
	}
}

func (r *PackageRegistry) Install(name, version string) {
	r.packages[name] = version
}

func (r *PackageRegistry) Get(name string) (string, bool) {
	v, ok := r.packages[name]
	return v, ok
}

func (r *PackageRegistry) List() map[string]string {
	return r.packages
}

func (r *PackageRegistry) Remove(name string) {
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
