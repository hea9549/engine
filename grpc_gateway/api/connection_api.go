/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/grpc_gateway"
)

type ConnectionApi struct {
	grpcService  grpc_gateway.GrpcService
	eventService common.EventService
}

func NewConnectionApi(grpcService grpc_gateway.GrpcService, eventService common.EventService) *ConnectionApi {
	return &ConnectionApi{
		grpcService:  grpcService,
		eventService: eventService,
	}
}

func (c ConnectionApi) CreateConnection(address string) (grpc_gateway.Connection, error) {

	logger.Infof(nil, "[gRPC-Gateway] Dialing - Address: [%s]", address)

	connection, err := c.grpcService.Dial(address)
	if err != nil {
		logger.Fatalf(nil, "[gRPC-Gateway] Fail to dial - Err: [%s]", err)
		return grpc_gateway.Connection{}, err
	}

	err = c.eventService.Publish("connection.created", createConnectionCreatedEvent(connection))
	if err != nil {
		return connection, err
	}

	logger.Infof(nil, "[gRPC-Gateway] Connection created - Address [%s], ConnectionID [%s]", connection.Address, connection.ConnectionID)

	return connection, nil
}

func createConnectionCreatedEvent(connection grpc_gateway.Connection) event.ConnectionCreated {
	return event.ConnectionCreated{
		Address:      connection.Address,
		ConnectionID: connection.ConnectionID,
	}
}

func (c ConnectionApi) CloseConnection(connectionID string) error {
	logger.Infof(nil, "[gRPC-Gateway] Close connection - ConnectionID [%s]", connectionID)

	c.grpcService.CloseConnection(connectionID)

	return c.eventService.Publish("connection.closed", createConnectionClosedEvent(connectionID))
}

func createConnectionClosedEvent(connectionID string) event.ConnectionClosed {
	return event.ConnectionClosed{
		ConnectionID: connectionID,
	}
}

func (c ConnectionApi) OnConnection(connection grpc_gateway.Connection) {
	logger.Infof(nil, "[gRPC-Gateway] Connection created - Address [%s], ConnectionID [%s]", connection.Address, connection.ConnectionID)

	if err := c.eventService.Publish("connection.created", createConnectionCreatedEvent(connection)); err != nil {
		logger.Infof(nil, "[gRPC-Gateway] Fail to publish connection createdEvent - ConnectionID: [%s]", connection.ConnectionID)
	}
}

func (c ConnectionApi) OnDisconnection(connection grpc_gateway.Connection) {
	logger.Infof(nil, "[gRPC-Gateway] Connection closed - ConnectionID [%s]", connection.ConnectionID)

	if err := c.eventService.Publish("connection.closed", connection); err != nil {
		logger.Infof(nil, "[gRPC-Gateway] Fail to publish connection createdEvent - ConnectionID: [%s]", connection.ConnectionID)
	}
}

func (c ConnectionApi) GetAllConnections() ([]grpc_gateway.Connection, error) {
	return c.grpcService.GetAllConnections()
}
