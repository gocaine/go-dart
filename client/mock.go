package client

import (
	"github.com/gocaine/go-dart/hardware"

	log "github.com/Sirupsen/logrus"
)

type MockedClient struct {
}

func NewMockedClient() *MockedClient {
	return &MockedClient{}
}

func (mock *MockedClient) Prepare() error {
	log.Info("Mocked client is ready")
	return nil
}

func (mock *MockedClient) Consume(event hardware.InputEvent) {
	log.WithField("event", event).Info("mocked event handler")
}
func (mock *MockedClient) Shutdown() {
	log.Info("Mocked client is shuting down")
}
