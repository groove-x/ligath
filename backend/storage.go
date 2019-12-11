package main

type Storage interface {
	Close()
	GetPackage(pkg string, ver string) (*Package, error)
	GetParsedPackage(pkg string, ver string) (*Package, error)
	GetNotParsedPackage(pkg string, ver string) (*Package, error)
	GetVerifiedPackage(pkg string, ver string) (*Package, error)
	PutPackage(pkg Package) error
	GetPackagesWithLicense(license string) []PackageListItem
	GetParsedPackages() []PackageListItem
	GetNotParsedPackages() []PackageListItem
	GetVerifiedPackages() []PackageListItem
	GetEmptyCopyrightPackages() []PackageListItem
	GetLicenses() []License
}
