package worker

import (
	"context"
	"github.com/goombaio/namegenerator"
	"github.com/hashicorp/serf/serf"
	"github.com/hextechpal/echelon/config"
	"github.com/hextechpal/echelon/discovery"
	"log"
	"net/http"
	"sync"
	"time"
)

type Worker struct {
	name     string
	doneCh   chan bool
	stopLock sync.Mutex
	cluster  *discovery.Cluster
	server   *http.Server
}

func NewWorker(c *config.Config) (*Worker, error) {
	ng := namegenerator.NewNameGenerator(time.Now().UnixMilli())
	w := &Worker{
		name:   ng.Generate(),
		doneCh: make(chan bool),
	}
	if err := w.initMemberShip(c); err != nil {
		return nil, err
	}
	if err := w.initServer(c.ServerAddress); err != nil {
		return nil, err
	}
	go w.monitor()
	return w, nil
}

func (w *Worker) initServer(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hell0"))
	})
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	w.server = server
	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			w.Stop()
		}
	}()
	return nil
}

func (w *Worker) initMemberShip(c *config.Config) error {
	membership, err := discovery.NewMembership(discovery.Config{
		NodeName:  w.name,
		BindAddr:  c.BindAddress,
		Tags:      nil,
		JoinAddrs: c.JoinAddresses,
	})
	if err != nil {
		return err
	}
	w.cluster = membership
	return nil

}

func (w *Worker) Stop() {
	w.stopLock.Lock()
	defer w.stopLock.Unlock()

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
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case e := <-w.cluster.EventCh():
			switch e.EventType() {
			case serf.EventMemberJoin:
				for _, member := range e.(serf.MemberEvent).Members {
					if w.cluster.IsLocal(member) {
						continue
					}
					if err := w.Join(member.Name, member.Tags["rpc_addr"]); err != nil {
						log.Printf("failed to join %v, %v", member, err)
					}
				}
			case serf.EventMemberLeave, serf.EventMemberFailed:
				for _, member := range e.(serf.MemberEvent).Members {
					if w.cluster.IsLocal(member) {
						return
					}
					if err := w.Leave(member.Name); err != nil {
						log.Printf("failed to join %v, %v", member, err)
					}
				}

			}
		case <-ticker.C:
			log.Printf("tick happened")
			for _, m := range w.cluster.Members() {
				log.Printf("name=%s, addr=%s, status=%v\n", m.Name, m.Addr, m.Status)
			}
		}
	}
}

func (w *Worker) Join(name, addr string) error {
	log.Printf("member joined name=%s, addr=%s\n", name, addr)
	return nil
}

func (w *Worker) Leave(name string) error {
	log.Printf("member left name=%s\n", name)
	return nil
}
