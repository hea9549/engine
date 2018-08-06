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

package service

import (
	"errors"
	"fmt"

	"github.com/it-chain/engine/icode"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/pb"
)

type TesseractContainerService struct {
	tesseract      *tesseract.Tesseract
	containerIdMap map[icode.ID]string // key : iCodeId, value : containerId
}

func NewTesseractContainerService() *TesseractContainerService {
	tesseractObj := &TesseractContainerService{
		tesseract:      tesseract.New(),
		containerIdMap: make(map[icode.ID]string, 0),
	}
	return tesseractObj
}

func (cs TesseractContainerService) StartContainer(meta icode.Meta) error {

	tesseractIcodeInfo := tesseract.ICodeInfo{
		Name:      meta.RepositoryName,
		Directory: meta.Path,
	}

	containerId, err := cs.tesseract.SetupContainer(tesseractIcodeInfo)

	if err != nil {
		return err
	}

	cs.containerIdMap[meta.ICodeID] = containerId

	return nil
}

func (cs TesseractContainerService) ExecuteRequest(request icode.Request, callback func(result *icode.Result, err error)) error {

	containerId, found := cs.containerIdMap[request.ICodeID]

	if !found {
		return errors.New(fmt.Sprintf("no container for iCode : %s", request.ICodeID))
	}

	tesseractReq := tesseract.Request{
		Uuid:     request.Uuid,
		TypeName: request.Type,
		FuncName: request.Function,
		Args:     request.Args,
	}

	tesseractCallBack := func(tesRes *pb.Response, err error) {
		icodeResult := icode.Result{
			Uuid:   tesRes.Uuid,
			Type:   tesRes.Type,
			Result: tesRes.Result,
			Data:   tesRes.Data,
			Error:  tesRes.Error,
		}
		callback(&icodeResult, err)
	}

	err := cs.tesseract.Request(containerId, tesseractReq, tesseractCallBack)

	if err != nil {
		return err
	}

	return nil
}

func (cs TesseractContainerService) StopContainer(id icode.ID) error {
	containerId := cs.containerIdMap[id]
	if containerId == "" {
		return errors.New(fmt.Sprintf("no container with icode id %s:", id))
	}
	err := cs.tesseract.StopContainerById(containerId)
	if err != nil {
		return err
	}
	delete(cs.containerIdMap, id)
	return nil
}
