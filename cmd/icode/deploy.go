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
	"log"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/icode"
		"github.com/rs/xid"
	"github.com/urfave/cli"
)

func DeployCmd() cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "it-chain icode deploy [icode git url] [ssh path]",
		Action: func(c *cli.Context) error {

			gitUrl := c.Args().Get(0)
			sshPath := c.Args().Get(1)
			deploy(gitUrl, sshPath)

			return nil
		},
	}
}

func deploy(gitUrl string, sshPath string) {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)

	defer client.Close()

	deployCommand := command.Deploy{
		ICodeId: xid.New().String(),
		Url:     gitUrl,
		SshPath: sshPath,
	}

	log.Printf("deploying icode...")
	log.Printf("This may take a few minutes")

	err := client.Call("icode.deploy", deployCommand, func(meta icode.Meta, err rpc.Error) {

		if !err.IsNil() {
			log.Printf("fail to deploy icode err: [%s]", err.Message)
			return
		}

		log.Printf("[%s] icode has deployed", meta.ICodeID)
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
