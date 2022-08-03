package main

type StorageType int

const (
	File StorageType = iota
)

type Storage interface {
	FindAllProxies() ([]Proxy, error)
	FindAllSources() ([]Source, error)
	FindProxyByID(ID string) (Proxy, error)
	FindSourceByID(ID string) (Source, error)
}
