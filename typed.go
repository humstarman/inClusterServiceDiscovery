package inClusterServiceDiscovery

import (
	//"fmt"
	"log"
	//"time"
	"strings"

	//"github.com/pkg/errors"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//v1beta1 "k8s.io/api/extensions/v1beta1"
)

const (
	Count     = 3
	Separator = ","
)

func CreateSearch(c *Config) (*Search, error) {
	s := Search{}
	ccopy(c, &s)
	s.Separator = Separator
	s.Counts = make([]int, Count, Count)
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	s.Client = cli
	return &s, err
}

func SimpleSearch(str string) (*Search, error) {
	s := Search{}
	s.Separator = Separator
	s.Counts = make([]int, Count, Count)
	vals := strings.Split(str, ".")
	s.Service = vals[0]
	if len(vals) == 1 {
		s.Namespace = "default"
	} else {
		s.Namespace = vals[1]
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	s.Client = cli
	return &s, err
}

func Create(x interface{}) (*Search, error) {
	switch x.(type) {
	case string:
		return SimpleSearch(x)
	case *Config:
		return CreateSearch(x)
	}
	err := fmt.Errorf("error: wrong type")
	return nil, err
}
