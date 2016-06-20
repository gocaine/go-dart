package cmd

import (
	"go-dart/client"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sampleCmd = &cobra.Command{
	Use:   "sample <gamestyle>",
	Short: "Sample",
	Run: func(cmd *cobra.Command, arg []string) {

		if len(arg) < 1 {
			log.Panic("Missing argument. Retry with a gamestyle : 101, 301...")
		}

		client := client.NewClient()

		response, err := client.CreateGame(arg[0])
		if err != nil {
			log.Error(err)
		}
		log.Println(response)

		response, err = client.GetState(1)
		if err != nil {
			log.Error(err)
		}
		log.Println(response)

		resource, err := client.CreatePlayer(1, "erwann")
		if err != nil {
			log.Error(err)
		}
		log.Println(resource)

		state, err := client.FireDart(1, 20, 3)
		if err != nil {
			log.Error(err)
		}
		log.Println(state)

		response, err = client.GetState(1)
		if err != nil {
			log.Error(err)
		}
		log.Println(response)

	},
}
