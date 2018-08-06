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

package icode

import (
	"errors"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/icode"
	"github.com/urfave/cli"
)

func InvokeCmd() cli.Command {
	return cli.Command{
		Name:  "invoke",
		Usage: "it-chain icode invoke [icode id] [function name] [...args]",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return errors.New("not enough args")
			}

			icodeId := c.Args().Get(0)
			functionName := c.Args().Get(1)
			args := make([]string, 0)

			for i := 2; i < c.NArg(); i++ {
				args = append(args, c.Args().Get(i))
			}

			invoke(icodeId, functionName, args)

			return nil
		},
	}
}

func invoke(id string, functionName string, args []string) {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)

	defer client.Close()

	invokeCommand := command.ExecuteRequest{
		ICodeId:  id,
		Function: functionName,
		Args:     args,
		Method:   "invoke",
	}

	logger.Info(nil, "invoke icode ...")

	err := client.Call("icode.execute", invokeCommand, func(result icode.Result, err rpc.Error) {

		if !err.IsNil() {
			logger.Errorf(nil, "fail to invoke icode err: [%s]", err.Message)
			return
		}

		logger.Infof(nil, "%15s : [%s]", "icodeId", id)
		logger.Infof(nil, "%15s : [%s]", "function Name", functionName)

		for i, arg := range args {
			logger.Infof(nil, "%15s : [%s]", "arg"+string(i), arg)
		}

		logger.Infof(nil, "%15s : %t", "success", result.Result)

		logger.Infof(nil, "%s : %s", "[result]"+string(result.Data))
	})

	if err != nil {
		logger.Fatal(&logger.Fields{"err_msg": err.Error()}, "fatal err in query cmd")
	}
}
