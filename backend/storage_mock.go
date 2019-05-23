package main

type MockStorage struct {
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (s *MockStorage) GetPackages() []Package {
	return []Package{
		{
			ID:      "vim",
			Version: "2.:8.0.0197-4+deb9u1",
			Copyrights: []Copyright{
				{
					Notice: "Copyright (c) 1988-2003 Bram Moolenaar",
					License: License{
						Name: "Charityware (GPL compatible)",
						Body: "foobar",
					},
				},
			},
		},
	}
}
