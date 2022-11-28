package discovery

import (
	"github.com/hashicorp/serf/serf"
	"net"
)

type Cluster struct {
	serf   *serf.Serf
	serfCh chan serf.Event
}

func NewMembership(c Config) (*Cluster, error) {
	m := &Cluster{}
	if err := m.setupSerf(c); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Cluster) setupSerf(c Config) error {
	m.serfCh = make(chan serf.Event)
	addr, err := net.ResolveTCPAddr("tcp", c.BindAddr)
	if err != nil {
		return err
	}

	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port

	config.EventCh = m.serfCh
	config.Tags = c.Tags
	config.NodeName = c.NodeName

	m.serf, err = serf.Create(config)
	if err != nil {
		return err
	}

	if c.JoinAddrs != nil && len(c.JoinAddrs) > 0 {
		_, err = m.serf.Join(c.JoinAddrs, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Cluster) IsLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

func (m *Cluster) EventCh() chan serf.Event {
	return m.serfCh
}

func (m *Cluster) Leave() error {
	return m.serf.Leave()
}

func (m *Cluster) Members() []serf.Member {
	return m.serf.Members()
}
