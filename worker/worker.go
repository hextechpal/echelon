package worker

import (
	"context"
	"fmt"
	"github.com/hashicorp/serf/serf"
	"log"
	"net/http"
	"time"
)

type Worker struct {
	id     string
	doneCh chan bool

	cluster *serf.Serf
	serfCh  chan serf.Event

	server *http.Server
}

func NewWorker(id string, cluster *serf.Serf, serfCh chan serf.Event) *Worker {
	w := &Worker{
		id:      id,
		doneCh:  make(chan bool),
		cluster: cluster,
		serfCh:  serfCh,
	}
	return w
}

func (w *Worker) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	w.server = server
	go w.monitor()
	return w.server.ListenAndServe()
}

func (w *Worker) Stop() {
	log.Printf("leaving cluster")
	err := w.cluster.Leave()
	if err != nil {
		log.Printf("err leaving cluster err=%v\n", err)
	}

	log.Printf("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = w.server.Shutdown(ctx)
	cancel()
}

func (w *Worker) monitor() {
	for {
		select {
		case event := <-w.serfCh:
			switch t := event.(type) {
			case serf.MemberEvent:
				fmt.Printf("MemberEvent %T\n", t)
			case serf.UserEvent:
				fmt.Printf("UserEvent %T\n", t)
			}
		}
	}
}
