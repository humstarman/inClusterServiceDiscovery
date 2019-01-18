package inClusterServiceDiscovery

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//v1beta1 "k8s.io/api/extensions/v1beta1"
)

const (
	Try   = 100
	Count = 3
)

type Search struct {
	ControllerName string
	ControllerType string
	Namespace      string
	Service        string
	Separator      string
	Total          int
	Client         *kubernetes.Clientset
	Ip             string
	Counts         []int
	Tmp            string
}

func CreateSearch(c *Config) (*Search, error) {
	s := Search{}
	ccopy(c, &s)
	s.Separator = ","
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

func (this *Search) Print() {
	log.Printf("Namespace: %v\n", this.Namespace)
	log.Printf("ControllerName: %v\n", this.ControllerName)
	log.Printf("ControllerType: %v\n", this.ControllerType)
	log.Printf("Service: %v\n", this.Service)
}

func (this *Search) Result() (string, error) {
	var ret int
	var err error
	switch this.ControllerType {
	case "daemonset", "ds":
		ret, err = this.Daemonset()
	case "deployment", "deploy":
		ret, err = this.Deployment()
	case "statefulset", "state", "s":
		ret, err = this.Statefulset()
	case "":
		ret, err = this.Endpoint()
	default:
		err = errors.New("err: wrong type of controller, as instance: deployment, statefulset or daemonset")
		ret = -1
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	this.Total = ret
	ip, err := this.GetEndpoints()
	return ip, err
}

func (this *Search) Daemonset() (int, error) {
	cli := this.Client
	namespace := this.Namespace
	name := this.ControllerName
	obj, err := cli.ExtensionsV1beta1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	total := int(obj.Status.DesiredNumberScheduled)
	return total, err
}

func (this *Search) Deployment() (int, error) {
	cli := this.Client
	namespace := this.Namespace
	name := this.ControllerName
	obj, err := cli.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	total := int(*(obj.Spec.Replicas)) // deployment, statefulset
	return total, err
}

func (this *Search) Statefulset() (int, error) {
	cli := this.Client
	namespace := this.Namespace
	name := this.ControllerName
	obj, err := cli.AppsV1beta1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	total := int(*(obj.Spec.Replicas)) // deployment, statefulset
	return total, err
}

func (this *Search) GetEndpoints() (string, error) {
	cli := this.Client
	namespace := this.Namespace
	svc := this.Service
	for try := 0; try < Try; try++ {
		eps, err := cli.CoreV1().Endpoints(namespace).Get(svc, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
			return "", err
		}
		n1 := len(eps.Subsets)
		for i := 0; i < n1; i++ {
			addrs := eps.Subsets[i].Addresses
			n2 := len(addrs)
			if n2 == this.Total {
				ips := ""
				sep := ""
				for j := 0; j < n2; j++ {
					ips += sep
					ips += fmt.Sprintf("%v", addrs[j].IP)
					sep = this.Separator
				}
				return ips, err
			}
			time.Sleep(3 * time.Second)
		}
	}
	msg := fmt.Sprintf("err: cannot find IP of %v.%v", this.Service, this.Namespace)
	err := errors.New(msg)
	log.Println(err)
	return "", err
}

func (this *Search) Endpoint() (string, error) {
	cli := this.Client
	namespace := this.Namespace
	svc := this.Service
	for try := 0; try < Try; try++ {
		for c := 0; i < Count; i++ {
			eps, err := cli.CoreV1().Endpoints(namespace).Get(svc, metav1.GetOptions{})
			if err != nil {
				log.Println(err)
				return "", err
			}
			addrs := eps.Subsets[0].Addresses
			n := len(addrs)
			this.Counts[c] = n
			ips := ""
			sep := ""
			for j := 0; j < n; j++ {
				ips += sep
				ips += fmt.Sprintf("%v", addrs[j].IP)
				sep = this.Separator
			}
			this.Tmp = ips
		}
		sum := 0
		max := -1
		for _, num := range this.Counts {
			sum += num
			if num > max {
				max = num
			}
		}
		if sum == max*Count {
			ret := this.Tmp
			return ret, nil
		}
		time.Sleep(3 * time.Second)
	}
	msg := fmt.Sprintf("err: cannot find IP of %v.%v", this.Service, this.Namespace)
	err := errors.New(msg)
	log.Println(err)
	return "", err
}
