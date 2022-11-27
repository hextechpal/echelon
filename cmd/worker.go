package cmd

import (
	"fmt"
	"github.com/hashicorp/serf/serf"
	"github.com/hextechpal/echelon/commons"
	"github.com/hextechpal/echelon/config"
	"github.com/hextechpal/echelon/worker"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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
			panic("error parsing config file")
		}

		c, err := config.Load()
		if err != nil {
			panic("error loading config")
		}

		id := commons.GenerateUuid()
		serfCh := make(chan serf.Event)
		cluster, err := initSerf(id, serfCh, c)
		if err != nil {
			log.Fatal(err)
		}

		worker := worker.NewWorker(id, cluster, serfCh)
		go func() {
			addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
			err := worker.Start(addr)
			if err != nil && err != http.ErrServerClosed {
				cluster.Leave()
			}
		}()

		ch := make(chan os.Signal, 2)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		log.Printf("signal received")
		worker.Stop()
	},
}

func initSerf(node string, serfCh chan serf.Event, config *config.Config) (*serf.Serf, error) {
	sc := serf.DefaultConfig()
	sc.Init()
	sc.NodeName = node
	sc.MemberlistConfig.BindAddr = config.Serf.BindAddress
	sc.MemberlistConfig.BindPort = config.Serf.BindPort
	sc.MemberlistConfig.AdvertiseAddr = config.Server.Host
	sc.MemberlistConfig.AdvertisePort = config.Server.Port
	sc.EventCh = serfCh
	cluster, err := serf.Create(sc)
	if err != nil {
		return nil, err
	}
	_, err = cluster.Join([]string{fmt.Sprintf("%s:%d", config.Serf.BindAddress, config.Serf.BindPort)}, true)
	if err != nil {
		return nil, err
	}
	return cluster, err
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().StringP(envArg, "e", "local.env", "Config file for the worker")
}
