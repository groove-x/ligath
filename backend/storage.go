package main

type Storage interface {
	Setup() error
	Close()
	GetPackage(pkg string, ver string) (*Package, error)
	GetParsedPackage(pkg string, ver string) (*Package, error)
	GetNotParsedPackage(pkg string, ver string) (*Package, error)
	GetVerifiedPackage(pkg string, ver string) (*Package, error)
	PutPackage(pkg Package) error
	GetParsedPackages() []PackageListItem
	GetNotParsedPackages() []PackageListItem
	GetVerifiedPackages() []PackageListItem
	GetLicenses() []License
}
