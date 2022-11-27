package discovery

type Config struct {
	BindAddr  string
	Tags      map[string]string
	JoinAddrs []string
	NodeName  string
}
