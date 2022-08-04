package exporting

import (
	"testing"
	"time"

	"github.com/AdiPP/reconciliation/pkg/listing"
	"github.com/stretchr/testify/assert"
)

func TestFindSourceByID(t *testing.T) {
	mr := newMockStorage()
	s := NewService(mr)

	_, err := s.FindSourceByID("zoUr")

	assert.Nil(t, err)
}

type mockStorage struct {
	proxies []Proxy
	sources []Source
}

func (m *mockStorage) FindAllProxies() ([]listing.Proxy, error) {
	proxies := []listing.Proxy{}

	for _, pp := range m.proxies {
		p := listing.Proxy{
			ID:     pp.ID,
			Amount: pp.Amount,
			Desc:   pp.Desc,
			Date:   pp.Date,
		}
		proxies = append(proxies, p)
	}

	return proxies, nil
}

func (m *mockStorage) FindSourceByID(ID string) (listing.Source, error) {
	var source listing.Source

	for _, ss := range m.sources {
		if ss.ID == ID {
			return listing.Source{
				ID:     ss.ID,
				Amount: ss.Amount,
				Desc:   ss.Desc,
				Date:   ss.Date,
			}, nil
		}
	}

	return source, listing.ErrSourceNotFound
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		[]Proxy{
			{
				ID:     "zoUr",
				Amount: 24,
				Desc:   "A",
				Date:   getDate("2021-06-30"),
			},
		},
		[]Source{
			{
				ID:     "zoUr",
				Amount: 24,
				Desc:   "A",
				Date:   getDate("2021-06-30"),
			},
		},
	}
}

func getDate(dateString string) time.Time {
	date, err := time.Parse("2006-01-02", dateString)

	if err != nil {
		panic(err.Error())
	}

	return date
}
