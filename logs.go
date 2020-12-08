/*
   Copyright (C) nerdctl authors.
   Copyright (C) containerd authors.

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
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

var logCommand = &cli.Command{
	Name:      "logs",
	Usage:     "Fetch the logs of a container",
	ArgsUsage: "CONTAINER-ID",
	Action:    logAction,
	Flags:     []cli.Flag{},
}

func logAction(clicontext *cli.Context) error {

	if clicontext.NArg() < 1 {
		return errors.New("context needs to be specified")
	}

	containerID := clicontext.Args().First()
	if containerID == "" {
		return errors.New("ID cannot be empty")
	}

	client, ctx, cancel, err := newClient(clicontext)
	if err != nil {
		return err
	}
	defer cancel()

	runtimeClient := runtimeapi.NewRuntimeServiceClient(client.Conn())
	if runtimeClient == nil {
		return errors.New("No RuntimeClient")
	}

	container, err := client.LoadContainer(ctx, containerID)
	if err != nil {
		return err
	}

	r, err := runtimeClient.Status(ctx, &runtimeapi.StatusRequest{Verbose: true})
	if err != nil {
		return err
	}
	fmt.Println(r)

	//FIXME: Now empty resp returned
	containerID = container.ID()
	statusResp, err := runtimeClient.ContainerStatus(ctx, &runtimeapi.ContainerStatusRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return err
	}
	fmt.Println(statusResp)

	//TODO: https://github.com/kubernetes-sigs/cri-tools/blob/7d3b28acc3ecdd6970787f592385f64dc8f7e2c7/cmd/crictl/container.go#L824
	//TODO: Marshal the statusResp and get the metadata from the marshaled statusResp
	//TODO: GetLogPath() from the metadata
	//TODO: call logs.ReadLogs()

	return nil
}
