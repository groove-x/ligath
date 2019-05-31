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

func (s *MockStorage) Setup() error {
	return nil
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

func (s *MockStorage) GetParsedPackages() []PackageListItem {
	return []PackageListItem{
		{Name: vim.Name, Version: vim.Version},
	}
}

func (s *MockStorage) GetNotParsedPackages() []PackageListItem {
	return []PackageListItem{}
}

func (s *MockStorage) GetManualPackages() []PackageListItem {
	return []PackageListItem{}
}
