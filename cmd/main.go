package main

import (
	"flag"
	"fmt"
	"log"

	ifces "github.com/Carlos-Descalzi/MessageTriggeredJobRunner/client/v1alpha1"
	mtjobrunner "github.com/Carlos-Descalzi/MessageTriggeredJobRunner/pkg/controller/mtjobrunner"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"knative.dev/pkg/signals"
)

var (
	masterURL   string
	kubeconfig  string
	logLevel    string
	workerCount int
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kube config")
	flag.StringVar(&masterURL, "masterURL", "", "Kubernetes API Server")
	flag.StringVar(&logLevel, "log-level", "info", "Log level")
	flag.IntVar(&workerCount, "worker-count", 3, "Number of workers")
}

func main() {
	fmt.Printf("Starting\n")
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()

	logger, err := setupLogger(logLevel)

	defer logger.Sync()

	if err != nil {
		log.Fatalf("Unable to setup logger: %v", err)
	}

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)

	if err != nil {
		logger.Panic(err)
	}

	httpClient, err := rest.HTTPClientFor(cfg)
	if err != nil {
		logger.Panic(err)
	}

	if err != nil {
		logger.Panic(err)
	}

	ifces.AddToScheme(scheme.Scheme)

	client, err := kubernetes.NewForConfigAndClient(cfg, httpClient)

	if err != nil {
		logger.Panic(err)
	}

	jIfce, err := ifces.MessageTriggeredJobInterfaceNew(cfg, httpClient)

	if err != nil {
		logger.Panic(err)
	}

	lIfce, err := ifces.MessageListenerInterfaceNew(cfg, httpClient)

	if err != nil {
		logger.Panic(err)
	}

	controller, err := mtjobrunner.MessageTriggeredJobControllerNew(client, lIfce, jIfce, workerCount, logger)

	if err != nil {
		logger.Panic(err)
	}

	stopCh := signals.SetupSignalHandler()

	controller.Start()
	<-stopCh
	controller.Close()
	fmt.Print("Exit\n")
}
