package marketplace

import (
	"sync"
	"time"
)

type TokenManager struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time

	mu sync.Mutex
}

func (t *TokenManager) Set(access, refresh string, expiresIn int) {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.AccessToken = access
	t.RefreshToken = refresh
	t.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (t *TokenManager) IsExpired() bool {

	t.mu.Lock()
	defer t.mu.Unlock()

	return time.Now().After(t.ExpiresAt)
}

func (t *TokenManager) Get() string {

	t.mu.Lock()
	defer t.mu.Unlock()

	return t.AccessToken
}
