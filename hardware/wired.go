package hardware

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/kidoman/embd"
	"io/ioutil"
	"time"
	// Configure for rpi
	_ "github.com/kidoman/embd/host/rpi"
)

// WiredHardware is actually listening for hardware events.
type WiredHardware struct {
	over          chan bool
	eventReciever chan InputEvent
	inputs        []embd.DigitalPin
	outputs       []embd.DigitalPin
	board         *board
}

type board struct {
	sectors [][]common.Sector
}

// NewWiredHardware create a new NewWiredHardware.
func NewWiredHardware(runCalibration bool) *WiredHardware {
	hardware := &WiredHardware{}
	if !runCalibration {
		// Not running a calibration, load the board dataset
		content, err := ioutil.ReadFile("boards/default.json")
		if err != nil {
			panic("Missing board configuration")
		}
		hardware.board = &board{}
		err = json.Unmarshal(content, &hardware.board.sectors)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("board: %v", hardware.board.sectors)
	}
	hardware.over = make(chan bool)
	return hardware
}

func (hardware *WiredHardware) bootstrap() {
	hardware.outputs = make([]embd.DigitalPin, 7)
	hardware.inputs = make([]embd.DigitalPin, 10)

	err := embd.InitGPIO()
	if err != nil {
		log.Fatalf("oops %v", err)
	}
	for i, n := range []int{17, 22, 4, 10, 9, 11, 5} {
		log.Infof("preparing output #%d GPIO_%d", i, n)
		hardware.outputs[i], err = embd.NewDigitalPin(n)
		if err != nil {
			log.Fatalf("oops %v", err)
		}
		hardware.outputs[i].SetDirection(embd.Out)
		hardware.outputs[i].ActiveLow(false)
		hardware.outputs[i].Write(0)
	}

	for i, n := range []int{18, 23, 24, 25, 8, 7, 12, 16, 20, 21} {
		log.Infof("preparing input #%d GPIO_%d", i, n)
		hardware.inputs[i], err = embd.NewDigitalPin(n)
		if err != nil {
			log.Fatalf("oops %v", err)
		}
		if n == 8 || n == 7 {
			hardware.inputs[i].ActiveLow(true)
			hardware.inputs[i].PullDown()
		} else {
			hardware.inputs[i].PullUp()
		}
		hardware.inputs[i].SetDirection(embd.In)
	}

}

// Produce is the main event-loop reading the hardware matrix and producing InputEvent.
func (hardware *WiredHardware) Produce(inputEventChannel chan InputEvent) {
	hardware.eventReciever = inputEventChannel

	hardware.bootstrap()

	defer hardware.releaseGPIO()

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

		var i int
		for i = 0; i < 10; i++ {
			v, err := hardware.inputs[i].Read()
			if err != nil {
				log.Fatalf("oops %v", err)
			}
			if v == 1 {
				hardware.notify(i, down)
			}
		}
	}
}

func (hardware *WiredHardware) notify(in int, out int) {
	log.WithFields(log.Fields{"in": in, "out": out}).Warn("hit")
	if hardware.board != nil {
		sector := hardware.board.sectors[out][in]
		log.WithFields(log.Fields{"sector": sector.Pos, "multiplier": sector.Val}).Warn("hit")
		hardware.eventReciever <- InputEvent{sector.Pos, sector.Val}
	} else {
		hardware.eventReciever <- InputEvent{out, in}
	}
	time.Sleep(500 * time.Millisecond)
}

func (hardware *WiredHardware) releaseGPIO() {
	log.Warn("end of event-loop, releasing")
	for _, pin := range hardware.inputs {
		log.Infof("Releasing %d", pin.N())
		pin.Close()
	}
	for _, pin := range hardware.outputs {
		log.Infof("Releasing %d", pin.N())
		pin.Close()
	}
}

// Shutdown releases all the GPIO pins.
func (hardware *WiredHardware) Shutdown() {
	log.Warnln("shuting down hardware...")
	hardware.over <- true
	log.Debug("Notified game is over")
	close(hardware.eventReciever)
	log.Debug("eventReciever has been closed")

}
