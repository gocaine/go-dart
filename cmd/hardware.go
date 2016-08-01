package cmd

import (
	"github.com/gocaine/go-dart/client"
	"github.com/gocaine/go-dart/hardware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
)

func hardwareCmd() *cobra.Command {

	var hardwareCmd = &cobra.Command{
		Use:   "hardware",
		Short: "Start wired client",
		Long:  "Start a client fully wired to the electronic dartboard",
		Run: func(cmd *cobra.Command, arg []string) {
			log.Info("wiring...")

			var producer hardware.InputEventProducer
			inputEventChannel := make(chan hardware.InputEvent)
			noWire, _ := cmd.Flags().GetBool("no-wire")
			if noWire {
				log.Info("well, in fact let's use the keyboard...")
				producer = hardware.NewMockedHardware()
			} else {
				producer = hardware.NewWiredHardware()

			}

			noServer, _ := cmd.Flags().GetBool("no-server")
			var consumer hardware.InputEventConsumer
			if noServer {
				log.Info("well, in fact let's print events...")
				consumer = client.NewMockedClient()
			} else {
				consumer = client.NewWrappedClient()
			}

			c := make(chan os.Signal, 1)
			// trap ctrl-c
			signal.Notify(c, os.Interrupt)
			go func() {
				sig := <-c
				log.Warnf("Caught Ctrl-C (%v)", sig)
				producer.Shutdown()
			}()

			// event-loop, wait for input and push to the server
			go producer.Produce(inputEventChannel)
			for {
				select {
				case event, more := <-inputEventChannel:
					if !more {
						// channel has been closed
						return
					}
					consumer.Consume(event)
				}
			}

		},
	}

	hardwareCmd.Flags().BoolP("no-wire", "m", false, "mock the hardware (for dev pupose only)")
	hardwareCmd.Flags().Bool("no-server", false, "mock the server (for dev pupose only)")

	hardwareCmd.Flags().StringP("board", "b", "test", "name of the board")
	viper.BindPFlag("board", hardwareCmd.Flags().Lookup("board"))

	return hardwareCmd
}
