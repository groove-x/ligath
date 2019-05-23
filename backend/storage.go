package main

type Storage interface {
	GetPackages() []Package
}

type BoltStorage struct {
}

func NewBoltStorage() *BoltStorage {
	return &BoltStorage{}
}

func (s *BoltStorage) GetPackages() []Package {
	return []Package{}
}
