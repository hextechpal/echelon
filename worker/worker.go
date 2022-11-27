package worker

import (
	"context"
	"fmt"
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
	name       string
	doneCh     chan bool
	stopLock   sync.Mutex
	membership *discovery.Membership
	server     *http.Server
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
	if err := w.initServer(fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) initServer(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hell0"))
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
		BindAddr:  fmt.Sprintf("%s:%d", c.Serf.BindAddress, c.Serf.BindPort),
		Tags:      nil,
		JoinAddrs: c.Serf.JoinAddrs,
	})
	if err != nil {
		return err
	}
	w.membership = membership
	return nil

}

func (w *Worker) Stop() {
	w.stopLock.Lock()
	defer w.stopLock.Unlock()

	log.Printf("leaving cluster")
	err := w.membership.Leave()
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
		case e := <-w.membership.EventCh():
			switch e.EventType() {
			case serf.EventMemberJoin:
				for _, member := range e.(serf.MemberEvent).Members {
					if w.membership.IsLocal(member) {
						continue
					}
					if err := w.Join(member.Name, member.Tags["rpc_addr"]); err != nil {
						log.Printf("failed to join %v, %v", member, err)
					}
				}
			case serf.EventMemberLeave, serf.EventMemberFailed:
				for _, member := range e.(serf.MemberEvent).Members {
					if w.membership.IsLocal(member) {
						return
					}
					if err := w.Leave(member.Name); err != nil {
						log.Printf("failed to join %v, %v", member, err)
					}
				}
			}
		}
	}
}

func (w *Worker) Join(name, addr string) error {
	//TODO implement me
	panic("implement me")
}

func (w *Worker) Leave(name string) error {
	//TODO implement me
	panic("implement me")
}
