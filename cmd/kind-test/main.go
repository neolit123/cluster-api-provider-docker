/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chuckha/cluster-api-provider-docker/kind/actions"
	"sigs.k8s.io/kind/pkg/cluster/constants"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Here we go!")
	clusterName := "my-cluster"
	version := "v1.14.2"
	for {
		// read input
		text, _ := reader.ReadString('\n')
		cleanText := strings.TrimSpace(text)
		inputs := strings.Split(cleanText, " ")
		switch inputs[0] {
		case "new-cluster":
			fmt.Println("Creating load balancer")
			lb, err := actions.SetUpLoadBalancer(clusterName)
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			ip, err := lb.IP()
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			if _, err := actions.CreateControlPlane(clusterName, inputs[1], ip, version); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		case "add-worker":
			if _, err := actions.AddWorker(clusterName, inputs[1], version); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		case "delete-node":
			if len(inputs) < 2 {
				fmt.Println("usage: delete-node my-cluster-worker1")
				continue
			}
			fmt.Println("Warning: If you are deleting a control plane node your cluster may break.")
			if err := actions.DeleteNode(clusterName, inputs[1]); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		case "add-control-plane":
			if _, err := actions.AddControlPlane(clusterName, inputs[1], version); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		case "set-cluster-name":
			fmt.Println("setting cluster name...")
			clusterName = inputs[1]
		case "set-version":
			fmt.Println("setting version")
			version = inputs[1]
		default:
			fmt.Println("Unknown command")
		}
		fmt.Println("Done!")
	}
}

func getName(clusterName, role string) string {
	ns, err := nodes.List(
		fmt.Sprintf("label=%s=%s", constants.ClusterLabelKey, clusterName),
		fmt.Sprintf("label=%s=%s", constants.NodeRoleKey, role))
	if err != nil {
		panic(err)
	}
	count := len(ns)
	suffix := fmt.Sprintf("%d", count)
	if count == 0 {
		suffix = ""
	}
	return fmt.Sprintf("%s-%s%s", clusterName, role, suffix)
}
