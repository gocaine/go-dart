package client

import (
	"github.com/gocaine/go-dart/hardware"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"os"
)

// CalibrationClient can be used to generate the dart matrix
type CalibrationClient struct {
	board      string
	hits       [][]common.Sector
	sector     int
	multiplier int
}

// NewCalibrationClient prepare a client for the calibration
func NewCalibrationClient(board string) *CalibrationClient {
	log.Info("initializing the calibration process")
	calibration := &CalibrationClient{board: board, sector: 1, multiplier: 1}
	calibration.hits = make([][]common.Sector, 7)
	for i := range calibration.hits {
		calibration.hits[i] = make([]common.Sector, 10)
	}
	calibration.output()
	return calibration
}

// Prepare prints the first output
func (calibration *CalibrationClient) Prepare() error {
	calibration.output()
	return nil
}

func (calibration *CalibrationClient) output() {
	log.Infof("please hit %d * %d...", calibration.sector, calibration.multiplier)
}

func (calibration *CalibrationClient) next() {
	if calibration.multiplier == 3 {
		if calibration.sector == 20 {
			calibration.sector = 25
		} else {
			calibration.sector++
		}
		calibration.multiplier = 1
	} else {
		if calibration.multiplier == 2 && calibration.sector == 25 {
			result, _ := json.Marshal(calibration.hits)
			log.Infof("%s", result)
			os.Exit(0)
			// stop here
		} else {
			calibration.multiplier++
		}
	}
	calibration.output()
}

// Consume handles next input
func (calibration *CalibrationClient) Consume(event hardware.InputEvent) {
	log.Warnf("%d * %d @ [%d][%d]", calibration.sector, calibration.multiplier, event.Sector, event.Multiplier)
	calibration.hits[event.Sector][event.Multiplier] = common.Sector{Pos: calibration.sector, Val: calibration.multiplier}
	calibration.next()
}

// Shutdown prints shutdown message
func (calibration *CalibrationClient) Shutdown() {
	log.Info("Calibration client is shuting down")
}
