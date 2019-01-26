package inClusterServiceDiscovery

import (
	"fmt"
	"log"
	//"time"
	"strings"

	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//v1beta1 "k8s.io/api/extensions/v1beta1"
)

func createSearch(c *Config) (*Search, error) {
	s := Search{}
	ccopy(c, &s)
	s.separator = separator
	s.counts = make([]int, count, count)
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
	s.client = cli
	return &s, err
}

func simpleSearch(str string) (*Search, error) {
	s := Search{}
	s.separator = separator
	s.counts = make([]int, count, count)
	vals := strings.Split(str, ".")
	s.service = vals[0]
	if len(vals) == 1 {
		s.namespace = "default"
	} else {
		s.namespace = vals[1]
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
	s.client = cli
	return &s, err
}

func Create(x interface{}) (*Search, error) {
	switch x := x.(type) {
	case string:
		return simpleSearch(x)
	case *Config:
		return createSearch(x)
	}
	err := fmt.Errorf("error: wrong type")
	return nil, err
}
