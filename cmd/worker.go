package cmd

import (
	"github.com/hextechpal/echelon/config"
	"github.com/hextechpal/echelon/worker"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	envArg = "env-file"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Starts a echelon worker",
	Long:  "Starts an echelon worker with the given config",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString(envArg)
		err := godotenv.Load(path)
		if err != nil {
			log.Printf("cannot parse config file")
		}

		c, err := config.Load()
		if err != nil {
			panic("error loading config")
		}

		log.Printf("Inittiating worker with config=%v\n", c)
		worker, err := worker.NewWorker(c)
		if err != nil {
			panic("error initiating worker")
		}

		ch := make(chan os.Signal, 2)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		log.Printf("signal received")
		worker.Stop()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().StringP(envArg, "e", "local.env", "Config file for the worker")
}
