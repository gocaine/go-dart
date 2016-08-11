package client

import (
	"github.com/gocaine/go-dart/hardware"

	log "github.com/Sirupsen/logrus"
)

// MockedClient testing purpose
type MockedClient struct {
}

// NewMockedClient MockedClient constructor
func NewMockedClient() *MockedClient {
	return &MockedClient{}
}

// Prepare mock implem
func (mock *MockedClient) Prepare() error {
	log.Info("Mocked client is ready")
	return nil
}

// Consume mock implem
func (mock *MockedClient) Consume(event hardware.InputEvent) {
	log.WithField("event", event).Info("mocked event handler")
}

// Shutdown mock implem
func (mock *MockedClient) Shutdown() {
	log.Info("Mocked client is shuting down")
}
