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

package adapter

import (
	"os"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/api"
)

type DeployCommandHandler struct {
	icodeApi api.ICodeApi
}

func NewDeployCommandHandler(icodeApi api.ICodeApi) *DeployCommandHandler {
	return &DeployCommandHandler{
		icodeApi: icodeApi,
	}
}

func (d *DeployCommandHandler) HandleDeployCommand(deployCommand command.Deploy) (icode.Meta, rpc.Error) {

	savePath := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/.tmp/"

	meta, err := d.icodeApi.Deploy(deployCommand.ICodeId, savePath, deployCommand.Url, deployCommand.SshPath)

	if err != nil {
		return icode.Meta{}, rpc.Error{Message: err.Error()}
	}

	return *meta, rpc.Error{}
}
