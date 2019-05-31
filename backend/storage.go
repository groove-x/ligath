package main

type Storage interface {
	Setup() error
	Close()
	GetPackage(pkg string, ver string) (*Package, error)
	GetParsedPackages() []PackageListItem
	GetNotParsedPackages() []PackageListItem
	GetManualPackages() []PackageListItem
}
