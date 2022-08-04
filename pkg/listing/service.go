package listing

import (
	"errors"
)

var (
	ErrSourceNotFound = errors.New("source not found")
)

type Repository interface {
	FindAllProxies() ([]Proxy, error)
	FindSourceByID(ID string) (Source, error)
}

type Service interface {
	FindAllProxies() ([]Proxy, error)
	FindSourceByID(ID string) (Source, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) FindAllProxies() ([]Proxy, error) {
	proxies, err := s.r.FindAllProxies()

	if err != nil {
		return nil, err
	}

	return proxies, nil
}

func (s *service) FindSourceByID(ID string) (Source, error) {
	var source Source

	storageSource, err := s.r.FindSourceByID(ID)

	if err != nil {
		return source, err
	}

	return Source{
		ID:     storageSource.ID,
		Amount: storageSource.Amount,
		Desc:   storageSource.Desc,
		Date:   storageSource.Date,
	}, nil
}
