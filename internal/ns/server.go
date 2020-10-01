// Copyright (c) 2020 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !windows

package ns

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/vishvananda/netns"

	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
)

type nsServer struct {
	ns netns.NsHandle
}

// NewServer - simple server that will change the nsNet of the client to ns before calling the next chain element
// and return it to its original netNS before returning
func NewServer(ns netns.NsHandle) networkservice.NetworkServiceServer {
	return &nsServer{ns: ns}
}

func (n *nsServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	curNetns, err := netns.Get()
	if err != nil {
		return nil, err
	}
	err = netns.Set(n.ns)
	if err != nil {
		return nil, err
	}
	conn, err := next.Server(ctx).Request(ctx, request)
	if err != nil {
		return nil, err
	}
	err = netns.Set(curNetns)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (n *nsServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	curNetns, err := netns.Get()
	if err != nil {
		return nil, err
	}
	err = netns.Set(n.ns)
	if err != nil {
		return nil, err
	}
	_, err = next.Server(ctx).Close(ctx, conn)
	if err != nil {
		return nil, err
	}
	err = netns.Set(curNetns)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
