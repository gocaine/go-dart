package hardware

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

// A MockedHardware uses keyboard input instead of electronic dart board.
type MockedHardware struct {
	channel chan InputEvent
}

// NewMockedHardware create a new NewMockedHardware.
func NewMockedHardware() *MockedHardware {
	return &MockedHardware{}
}

// Produce is the event-loop responsible of producing the InputEvent based on keyboard input.
func (mock *MockedHardware) Produce(inputEventChannel chan InputEvent) {
	var sector, multiplier int
	mock.channel = inputEventChannel
	for {

		fmt.Print("Enter sector [1-20/25]: ")
		fmt.Scanf("%d", &sector)
		if sector < 0 {
			break
		}

		if sector < 1 || (sector > 20 && sector != 25) {
			fmt.Println("Illegal value (not in range)")
			continue
		}
		fmt.Print("Enter multiplier [1/2/3]: ")
		fmt.Scanf("%d", &multiplier)

		if multiplier < 0 {
			break
		}
		if multiplier < 1 || multiplier > 3 {
			fmt.Println("Illegal value (not in list)")
			continue
		}
		fmt.Println()
		inputEventChannel <- InputEvent{Sector: sector, Multiplier: multiplier}

		// wait a bit during the event processing
		time.Sleep(200 * time.Millisecond)

	}
	mock.Shutdown()
}

// Shutdown release all the resources.
func (mock *MockedHardware) Shutdown() {
	defer close(mock.channel)
	log.Info("Shuting down...")
}
