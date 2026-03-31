package simulator

import (
	"context"
	"math/rand"
	domain "micro-vuln-scanner/internal/domain"
	store "micro-vuln-scanner/internal/store"
	"time"
)


type Generator struct {
	store *store.Store
	interval time.Duration
}

func NewGenerator(store *store.Store, interval time.Duration) *Generator {
	if(interval <= 0) {
		interval = time.Second * time.Duration(5)
	}

	return &Generator{
		store: store,
		interval: interval,
	}
}

func generateRandomVulnerability() domain.Vulnerability {
	// This is a placeholder for generating random vulnerabilities.
	// In a real implementation, you would generate realistic vulnerabilities with varying severities and descriptions.
	var severities = []domain.Severity{
    domain.LOW,
    domain.MEDIUM,
    domain.HIGH,
    domain.CRITICAL,
}

	sev := severities[rand.Intn(len(severities))]
	return domain.Vulnerability{
		ID:          time.Now().UTC().Format("20060102150405"), // Unique ID based on timestamp
		Severity:    sev, // Random severity between 0 and 3
		Status: domain.OPEN,
		CreatedAt: time.Now().UTC(),
		ImageName: "example-image" + time.Now().UTC().Format("20060102150405"), // Unique image name based on timestamp
	}
}

func (generator *Generator) Start(ctx context.Context) chan struct{} {
	    done := make(chan struct{})
		tickler := time.NewTicker(generator.interval)

		go func() {
		defer close(done)
			for {
				select {
				case <-tickler.C:
					generator.store.Add(generateRandomVulnerability())
				case <-ctx.Done():
					return
				}
			}
		}()

		return done
}