package listing

import (
	"errors"

	"github.com/AdiPP/reconciliation/pkg/storage/memory"
)

var (
	ErrSourceNotFound = errors.New("source not found")
)

type Repository interface {
	FindAllProxies() ([]memory.Proxy, error)
	FindSourceByID(ID string) (memory.Source, error)
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
	storageProxies, err := s.r.FindAllProxies()

	if err != nil {
		return nil, err
	}

	var proxies []Proxy

	for _, v := range storageProxies {
		proxies = append(
			proxies,
			Proxy{
				ID:     v.ID,
				Amount: v.Amount,
				Desc:   v.Desc,
				Date:   v.Date,
			},
		)
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
