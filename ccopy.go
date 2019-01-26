package inClusterServiceDiscovery

func ccopy(c *Config, s *Search) {
	s.namespace = c.Namespace
	if s.namespace == "" {
		s.namespace = "default"
	}
	s.service = c.Service
	s.name = c.Name
	s.typed = c.Type
}
