package main

var vim = Package{
	Name:    "vim",
	Version: "2.:8.0.0197-4+deb9u1",
	Copyrights: []Copyright{
		{
			FileRange: []string{"*"},
			Copyright: "Copyright (c) 1988-2003 Bram Moolenaar",
			License: License{
				Name: "Charityware (GPL compatible)",
				Body: "foobar",
			},
		},
	},
	RawCopyright: "raw copyright",
}

type MockStorage struct {
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (s *MockStorage) Close() {
}

func (s *MockStorage) GetPackage(pkg, ver string) (*Package, error) {
	if pkg == "vim" {
		return &vim, nil
	} else {
		return nil, nil
	}
}

func (s *MockStorage) GetParsedPackage(pkg, ver string) (*Package, error) {
	return &Package{}, nil
}

func (s *MockStorage) GetNotParsedPackage(pkg, ver string) (*Package, error) {
	return &Package{}, nil
}

func (s *MockStorage) GetVerifiedPackage(pkg, ver string) (*Package, error) {
	return &Package{}, nil
}

func (s *MockStorage) PutPackage(pkg Package) error {
	return nil
}

func (s *MockStorage) GetParsedPackages() []PackageListItem {
	return []PackageListItem{
		{Name: vim.Name, Version: vim.Version},
	}
}

func (s *MockStorage) GetNotParsedPackages() []PackageListItem {
	return []PackageListItem{}
}

func (s *MockStorage) GetVerifiedPackages() []PackageListItem {
	return []PackageListItem{}
}

func (s *MockStorage) GetLicenses() []License {
	return []License{
		GPLv2(),
	}
}
