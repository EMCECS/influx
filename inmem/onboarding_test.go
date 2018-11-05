package inmem

import (
	"context"
	"testing"

	"github.com/EMCECS/influx"
	platformtesting "github.com/EMCECS/influx/testing"
)

func initOnboardingService(f platformtesting.OnboardingFields, t *testing.T) (platform.OnboardingService, func()) {
	s := NewService()
	s.IDGenerator = f.IDGenerator
	s.TokenGenerator = f.TokenGenerator
	ctx := context.TODO()
	if err := s.PutOnboardingStatus(ctx, !f.IsOnboarding); err != nil {
		t.Fatalf("failed to set new onboarding finished: %v", err)
	}
	return s, func() {}
}

func TestGenerate(t *testing.T) {
	platformtesting.Generate(initOnboardingService, t)
}
