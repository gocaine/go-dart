package hardware

import (
	"time"

	"github.com/gocaine/go-dart/common"

	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/embd"
	// Configure for rpi
	_ "github.com/kidoman/embd/host/rpi"
)

// WiredHardware is actually listening for hardware events.
type WiredHardware struct {
	over          chan bool
	eventReciever chan InputEvent
	inputs        []embd.DigitalPin
	outputs       []embd.DigitalPin
}

var sectors = [7][10]common.Sector{
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
	{common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}, common.Sector{Val: 1, Pos: 1}}, /*  initializers for row indexed by 0 */
}

// NewWiredHardware create a new NewWiredHardware.
func NewWiredHardware() *WiredHardware {
	hardware := &WiredHardware{}
	hardware.init()
	return hardware
}

func (hardware *WiredHardware) init() {
	log.Infoln("hello from on board led")

	defer embd.CloseLED()
	for index := 0; index < 8; index++ {
		embd.LEDToggle("LED0")
		time.Sleep(250 * time.Millisecond)
	}
	embd.LEDToggle("LED0")
	hardware.over = make(chan bool)
}

// Produce is the main event-loop reading the hardware matrix and producing InputEvent.
func (hardware *WiredHardware) Produce(inputEventChannel chan InputEvent) {
	hardware.eventReciever = inputEventChannel
	defer hardware.releaseGPIO()
	hardware.outputs = make([]embd.DigitalPin, 7)
	hardware.inputs = make([]embd.DigitalPin, 10)
	var err error

	err = embd.InitGPIO()
	if err != nil {
		log.Fatalf("oops %v", err)
	}
	for i, n := range []int{2, 3, 4, 5, 6, 7, 8} {
		log.Infof("preparing output #%d GPIO_%d", i, n)
		hardware.outputs[i], err = embd.NewDigitalPin(n)
		if err != nil {
			log.Fatalf("oops %v", err)
		}
		hardware.outputs[i].SetDirection(embd.Out)
		hardware.outputs[i].ActiveLow(false)
		hardware.outputs[i].Write(0)
	}

	for i, n := range []int{9, 10, 11, 12, 13, 16, 17, 18, 19, 20} {
		log.Infof("preparing input #%d GPIO_%d", i, n)
		hardware.inputs[i], err = embd.NewDigitalPin(n)
		if err != nil {
			log.Fatalf("oops %v", err)
		}
		hardware.inputs[i].PullUp()
		hardware.inputs[i].ActiveLow(true)
		hardware.inputs[i].SetDirection(embd.In)
	}

	down := 0
	//var total uint = 0
	for {
		select {
		case <-hardware.over:
			log.Warn("Game is over")
			return
		default:
			// proceed
		}
		hardware.outputs[down].Write(0)
		down = (down + 1) % 7
		hardware.outputs[down].Write(1)
		time.Sleep(5 * time.Microsecond)

		var i uint
		for i = 0; i < 10; i++ {
			v, err := hardware.inputs[i].Read()
			if err != nil {
				log.Fatalf("oops %v", err)
			}
			if v == 0 {
				log.Infof("hit %d %d", i, down)
			}
		}
	}
}

func (hardware *WiredHardware) notify(in int, out byte) {
	log.WithFields(log.Fields{"in": in, "out": out}).Warn("hit !")
}

func (hardware *WiredHardware) releaseGPIO() {
	log.Warn("end of event-loop, releasing")
	for _, pin := range hardware.inputs {
		log.Info("Releasing %d", pin.N())
		pin.Close()
	}
	for _, pin := range hardware.outputs {
		log.Info("Releasing %d", pin.N())
		pin.Close()
	}
	log.Info("Cleared")
}

// Shutdown releases all the GPIO pins.
func (hardware *WiredHardware) Shutdown() {
	log.Warnln("shuting down hardware...")
	hardware.over <- true
	close(hardware.eventReciever)
}
