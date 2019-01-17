package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//v1beta1 "k8s.io/api/extensions/v1beta1"
)

var name = flag.String("m", "", "The name of Controller object")
var namespace = flag.String("n", "default", "Namespace")
var svc = flag.String("s", "", "The name of Service object")
var total_try = flag.Int("y", 100, "todo")
var separator = flag.String("e", ",", "todo")

func init() {
	flag.Parse()
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	// creates the clientset
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Daemonset
	//obj, err := cli.ExtensionsV1beta1().DaemonSets(*namespace).Get(*name, metav1.GetOptions{})
	// Deployment
	obj, err := cli.ExtensionsV1beta1().Deployments(*namespace).Get(*name, metav1.GetOptions{})
	// Statfulset
	//obj, err := cli.AppsV1beta1().StatefulSets(*namespace).Get(*name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	total := int(*(obj.Spec.Replicas)) // deployment, statfulset
	//total := int(obj.Spec.Replicas) // daemonset
	//log.Printf("total: %v\n", total)
	for try := 0; try < *total_try; try++ {
		eps, err := cli.CoreV1().Endpoints(*namespace).Get(*svc, metav1.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		n1 := len(eps.Subsets)
		for i := 0; i < n1; i++ {
			addrs := eps.Subsets[i].Addresses
			n2 := len(addrs)
			if n2 == total {
				ips := ""
				sep := ""
				for j := 0; j < n2; j++ {
					ips += sep
					ips += fmt.Sprintf("%v", addrs[j].IP)
					sep = *separator
				}
				fmt.Println(ips)
				return
			}
			time.Sleep(3 * time.Second)
		}
	}
	log.Fatal("err")
}
