package main

import (
	"fmt"
	"os"
	
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/cloudfoundry-incubator/runtime-schema/cc_messages"
	"github.com/linsun/pod-builder/transformer"
	"github.com/pivotal-golang/lager"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {
	logger := lager.NewLogger("my-app")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))
	address := "http://9.37.192.140:8080"

	desireRequest := cc_messages.DesireAppRequestFromCC{
		ProcessGuid:  "myapp",
		DropletUri:   "source-url-1",
		Stack:        "stack-1",
		StartCommand: "npm start",
		Environment: []*models.EnvironmentVariable{
			{Name: "HOST", Value: "0.0.0.0"},
			{Name: "env-key-2", Value: "env-value-2"},
		},
		MemoryMB:        256,
		DiskMB:          1024,
		FileDescriptors: 16,
		NumInstances:    2,
		//RoutingInfo:     routeInfo1,
		LogGuid: "log-guid-1",
		ETag:    "1234567.1890",
		Ports:   []uint32{8080},
	}

	// // create data container dockerfile
	// err := WriteRuntimeDockerfile(desireRequest.DropletUri, desireRequest.StartCommand)
	// if err != nil {
	// 	logger.Fatal("write-file-failed", err)
	// }

	k8sClient, err := unversioned.New(&restclient.Config{
		Host: address,
	})

	if err != nil {
		logger.Fatal("Can't create Kubernetes Client", err, lager.Data{"address": address})
	}

	_, err = k8sClient.ServerVersion()
	if err != nil {
		logger.Fatal("Can't connect to Kubernetes API", err, lager.Data{"address": address})
	}

	fmt.Printf("connected to Kube API %v\n", lager.Data{"address": address})
	logger.Debug("Connected to Kubernetes API %s", lager.Data{"address": address})

	newPod, err := transformer.DesiredAppToPod(desireRequest)
	if err != nil {
		logger.Fatal("desired-app-to-pod", err)
	}

	_, err = k8sClient.Pods("linsun").Create(newPod)
	if err != nil {
		logger.Fatal("create-pod", err)
	}
}
