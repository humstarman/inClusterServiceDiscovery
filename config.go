package inClusterServiceDiscovery

type Config struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Namespace string `"json:namespace"`
	Service   string `json:"service"`
}
