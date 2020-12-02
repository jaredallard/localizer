// Copyright 2020 Jared Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package proxier

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type EventType string

var (
	EventAdded   EventType = "added"
	EventDeleted EventType = "deleted"
)

type ServiceEvent struct {
	EventType EventType
	Service   *corev1.Service
}

// CreateHandlers creates Kubernetes event handlers for all
// of our types. These then communicate with a port-forward
// worker to create Kubernetes port-forwards.
//nolint:gocritic // We're OK not naming these.
func CreateHandlers(ctx context.Context, requester chan<- PortForwardRequest, k kubernetes.Interface) (chan<- ServiceEvent, <-chan struct{}) {
	serviceChan := make(chan ServiceEvent)
	doneChan := make(chan struct{})

	go serviceProcessor(ctx, serviceChan, doneChan, requester, k)

	return serviceChan, doneChan
}

// Services
func serviceProcessor(ctx context.Context, event <-chan ServiceEvent, doneChan chan struct{}, requester chan<- PortForwardRequest, k kubernetes.Interface) {
	for {
		select {
		case <-ctx.Done():
			close(doneChan)
			return
		case s := <-event:
			info := ServiceInfo{
				Name:      s.Service.Name,
				Namespace: s.Service.Namespace,
				Type:      "",
			}

			// Skip this service for now.
			if info.Name == "kubernetes" {
				continue
			}

			if s.Service.Spec.ClusterIP == "None" {
				info.Type = ServiceTypeStatefulset
			}

			switch s.EventType {
			case EventAdded:
				ports := make([]int, len(s.Service.Spec.Ports))
				for i, p := range s.Service.Spec.Ports {
					ports[i] = int(p.Port)
				}

				switch info.Type {
				case ServiceTypeStandard:
					requester <- PortForwardRequest{
						CreatePortForwardRequest: &CreatePortForwardRequest{
							Service: info,
							Ports:   ports,
							Hostnames: []string{
								info.Name,
								fmt.Sprintf("%s.%s", info.Name, info.Namespace),
								fmt.Sprintf("%s.%s.svc", info.Name, info.Namespace),
								fmt.Sprintf("%s.%s.svc.cluster", info.Name, info.Namespace),
								fmt.Sprintf("%s.%s.svc.cluster.local", info.Name, info.Namespace),
							},
						},
					}
				case ServiceTypeStatefulset:
					endpoints, err := k.CoreV1().Endpoints(info.Namespace).Get(ctx, info.Name, metav1.GetOptions{})
					if err != nil {
						// TODO: expose error
						continue
					}

					for _, addresses := range endpoints.Subsets {
						for _, e := range addresses.Addresses {
							if e.TargetRef == nil {
								continue
							}

							if e.TargetRef.Kind != "Pod" {
								continue
							}

							name := fmt.Sprintf("%s.%s", e.TargetRef.Name, info.Name)
							requester <- PortForwardRequest{
								CreatePortForwardRequest: &CreatePortForwardRequest{
									Service:  info,
									Ports:    ports,
									Endpoint: &PodInfo{e.TargetRef.Name, e.TargetRef.Namespace},
									Hostnames: []string{
										info.Name,
										fmt.Sprintf("%s.%s", name, info.Namespace),
										fmt.Sprintf("%s.%s.svc", name, info.Namespace),
										fmt.Sprintf("%s.%s.svc.cluster", name, info.Namespace),
										fmt.Sprintf("%s.%s.svc.cluster.local", name, info.Namespace),
									},
								},
							}
						}
					}
				}
			case EventDeleted:
				requester <- PortForwardRequest{
					DeletePortForwardRequest: &DeletePortForwardRequest{
						Service: info,
					},
				}
			}
		}
	}
}
