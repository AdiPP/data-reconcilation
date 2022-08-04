package listing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindAllProxies(t *testing.T) {
	mr := newMockStorage()
	s := NewService(mr)

	_, err := s.FindAllProxies()

	assert.Nil(t, err)
}

type mockStorage struct {
	proxies []Proxy
	sources []Source
}

func (m *mockStorage) FindAllProxies() ([]Proxy, error) {
	return m.proxies, nil
}

func (m *mockStorage) FindSourceByID(ID string) (Source, error) {
	var source Source

	for _, ss := range m.sources {
		if ss.ID == ID {
			return ss, nil
		}
	}

	return source, ErrSourceNotFound
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
