package store

import (
	domain "micro-vuln-scanner/internal/domain"
	"sync"
)

type Store struct {
	mutex sync.RWMutex
	vulnerabilities []domain.Vulnerability
	maxSize int
}


func (store *Store) Add(vuln domain.Vulnerability)  {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if len(store.vulnerabilities) >= store.maxSize {
		store.vulnerabilities = store.vulnerabilities[1:]
	}
	store.vulnerabilities = append(store.vulnerabilities, vuln)
}

func (store *Store) GetAll() []domain.Vulnerability {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	result:= make([]domain.Vulnerability, len(store.vulnerabilities))
	copy(result, store.vulnerabilities)
	return result
}

func (store *Store) GetBySeverity(severity domain.Severity) []domain.Vulnerability {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	var result []domain.Vulnerability = []domain.Vulnerability{}
	for _,vuln := range store.vulnerabilities {
		if vuln.Severity == severity {
			result = append(result, vuln)
		}
	}
	return result
}

func NewStore(maxSize int) *Store {
	return &Store{
		vulnerabilities: make([]domain.Vulnerability, 0),
		maxSize: maxSize,
	}
}