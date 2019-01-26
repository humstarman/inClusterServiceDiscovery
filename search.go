package inClusterServiceDiscovery

import (
	"fmt"
	"log"
	"time"

	//"github.com/pkg/errors"
	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
	//v1beta1 "k8s.io/api/extensions/v1beta1"
)

const (
	trys = 100
)

type Search struct {
	name      string
	typed     string
	namespace string
	service   string
	separator string
	total     int
	client    *kubernetes.Clientset
	ip        string
	counts    []int
	tmp       string
}

func (this *Search) Print() {
	log.Printf("Namespace: %v\n", this.namespace)
	log.Printf("ControllerName: %v\n", this.name)
	log.Printf("ControllerType: %v\n", this.typed)
	log.Printf("Service: %v\n", this.service)
}

func (this *Search) Result() (string, error) {
	var ret int
	var err error
	switch this.typed {
	case "daemonset", "ds":
		ret, err = this.daemonset()
	case "deployment", "deploy":
		ret, err = this.deployment()
	case "statefulset", "state", "s":
		ret, err = this.statefulset()
	case "":
		ips, err := this.endpoint()
		if err != nil {
			log.Println(err)
			return "", err
		}
		return ips, nil
	default:
		err = fmt.Errorf("err: wrong type of controller, as instance: deployment, statefulset or daemonset")
		ret = -1
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	this.total = ret
	ip, err := this.getEndpoints()
	return ip, nil
}

func (this *Search) daemonset() (int, error) {
	cli := this.client
	namespace := this.namespace
	name := this.name
	obj, err := cli.ExtensionsV1beta1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	total := int(obj.Status.DesiredNumberScheduled)
	return total, nil
}

func (this *Search) deployment() (int, error) {
	cli := this.client
	namespace := this.namespace
	name := this.name
	obj, err := cli.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	total := int(*(obj.Spec.Replicas)) // deployment, statefulset
	return total, nil
}

func (this *Search) statefulset() (int, error) {
	cli := this.client
	namespace := this.namespace
	name := this.name
	obj, err := cli.AppsV1beta1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	total := int(*(obj.Spec.Replicas)) // deployment, statefulset
	return total, nil
}

func (this *Search) getEndpoints() (string, error) {
	cli := this.client
	namespace := this.namespace
	svc := this.service
	for try := 0; try < trys; try++ {
		eps, err := cli.CoreV1().Endpoints(namespace).Get(svc, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
			return "", err
		}
		n1 := len(eps.Subsets)
		for i := 0; i < n1; i++ {
			addrs := eps.Subsets[i].Addresses
			n2 := len(addrs)
			if n2 == this.total {
				ips := ""
				sep := ""
				for j := 0; j < n2; j++ {
					ips += sep
					ips += fmt.Sprintf("%v", addrs[j].IP)
					sep = this.separator
				}
				return ips, nil
			}
			time.Sleep(3 * time.Second)
		}
	}
	msg := fmt.Sprintf("err: cannot find IP of %v.%v", this.service, this.namespace)
	err := fmt.Errorf(msg)
	log.Println(err)
	return "", err
}

func (this *Search) endpoint() (string, error) {
	cli := this.client
	namespace := this.namespace
	svc := this.service
	for try := 0; try < trys; try++ {
		for c := 0; c < Count; c++ {
			eps, err := cli.CoreV1().Endpoints(namespace).Get(svc, metav1.GetOptions{})
			if err != nil {
				log.Println(err)
				return "", err
			}
			addrs := eps.Subsets[0].Addresses
			n := len(addrs)
			this.counts[c] = n
			ips := ""
			sep := ""
			for j := 0; j < n; j++ {
				ips += sep
				ips += fmt.Sprintf("%v", addrs[j].IP)
				sep = this.separator
			}
			this.tmp = ips
		}
		sum := 0
		max := -1
		for _, num := range this.counts {
			sum += num
			if num > max {
				max = num
			}
		}
		if sum == max*Count {
			ret := this.tmp
			return ret, nil
		}
		time.Sleep(3 * time.Second)
	}
	msg := fmt.Sprintf("err: cannot find IP of %v.%v", this.service, this.namespace)
	err := fmt.Errorf(msg)
	log.Println(err)
	return "", err
}
