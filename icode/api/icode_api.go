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
	"github.com/it-chain/engine/icode"
	"context"
	"time"
	"github.com/it-chain/engine/common/logger"
)

type ICodeApi struct {
	ContainerService icode.ContainerService
	GitService       icode.GitService
	Repository       *icode.MetaRepository
	EventService     icode.EventService
}

func NewIcodeApi(containerService icode.ContainerService, gitService icode.GitService, repository *icode.MetaRepository, eventService icode.EventService) *ICodeApi {

	return &ICodeApi{
		ContainerService: containerService,
		GitService:       gitService,
		Repository:       repository,
		EventService:     eventService,
	}
}

func (i ICodeApi) Deploy(id string, baseSaveUrl string, gitUrl string, sshPath string) (*icode.Meta, error) {

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err := i.GitService.Clone(id, baseSaveUrl, gitUrl, sshPath)

	if err != nil {
		return nil, err
	}

	//start ICode with container
	if err = i.ContainerService.StartContainer(*meta); err != nil {
		return nil, err
	}

	i.Repository.Save(meta)
	i.EventService.MetaCreated(*meta)

	return meta, nil
}

func (i ICodeApi) UnDeploy(id icode.ID) error {
	// stop iCode container
	err := i.ContainerService.StopContainer(id)
	if err != nil {
		return err
	}

	i.Repository.Delete(id)
	i.EventService.MetaDeleted(id)

	return nil
}

func (i ICodeApi) ExecuteTransactionList(RequestList []icode.Request) []*icode.Result {

	resultList := make([]*icode.Result, 0)

	for _, transaction := range RequestList {
		result,err := i.ExecuteRequest(transaction)
		if err != nil{
			resultList = append(resultList, result)
		}
	}

	return resultList
}

func (i ICodeApi) ExecuteRequest(request icode.Request) (*icode.Result, error) {
	//todo request timeout setting by outside
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	c := make(chan *icode.Result, 1)
	err := i.ContainerService.ExecuteRequest(request, func(result *icode.Result, err error) {
		logger.Error(nil, "error in execute request response : "+err.Error())
		c <- result
	})

	if err != nil {
		logger.Error(nil, "error in execute request : "+err.Error())
		return nil, err
	}

	select {
	case <-ctx.Done():
		logger.Error(nil, "error in execute request, timeout : "+ctx.Err().Error())
		return nil, ctx.Err()
	case result := <-c:
		return result, nil
	}
}
