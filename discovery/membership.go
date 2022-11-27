package discovery

import (
	"github.com/hashicorp/serf/serf"
	"net"
)

type Membership struct {
	cluster *serf.Serf
	serfCh  chan serf.Event
}

func NewMembership(c Config) (*Membership, error) {
	m := &Membership{}
	if err := m.setupSerf(c); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Membership) setupSerf(c Config) error {
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

	m.cluster, err = serf.Create(config)
	if err != nil {
		return err
	}

	if c.JoinAddrs != nil {
		_, err = m.cluster.Join(c.JoinAddrs, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Membership) IsLocal(member serf.Member) bool {
	return m.cluster.LocalMember().Name == member.Name
}

func (m *Membership) EventCh() chan serf.Event {
	return m.serfCh
}

func (m *Membership) Leave() error {
	return m.cluster.Leave()
}
