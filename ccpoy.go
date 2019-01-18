package inClusterServiceDiscover

func ccopy(c *Config, s *Search) {
	s.Namespace = c.Namespace
        if s.Namespace == "" {
                s.Namespace = "default"
        }
        s.Service = c.Service
        s.ControllerName = c.ControllerName
        s.ControllerType = c.ControllerType
}
