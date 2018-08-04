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

import "errors"

type Version struct {
}

type ID = string
type MetaStatus = int

const (
	READY = iota
	UNDEPLOYED
	DEPLOYED
	DEPLOY_FAIL
)

type Meta struct {
	ICodeID        ID
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
	Status         MetaStatus
}

func NewMeta(id string, repositoryName string, gitUrl string, path string, commitHash string) *Meta {
	return &Meta{
		ICodeID:        id,
		RepositoryName: repositoryName,
		GitUrl:         gitUrl,
		Path:           path,
		CommitHash:     commitHash,
		Version:        Version{},
		Status:         0,
	}
}
func (m Meta) GetID() string {
	return m.ICodeID
}

type MetaRepository struct {
	metas map[ID]*Meta
}

func NewMetaRepository() *MetaRepository {
	return &MetaRepository{
		metas:        make(map[ID]*Meta),
	}
}

func (mr *MetaRepository) Save(meta *Meta) {
	mr.metas[meta.GetID()]=meta
}

func (mr *MetaRepository) Delete(id ID) {
	delete(mr.metas, id)
}

func (mr *MetaRepository) StatusChanged(id ID,status MetaStatus) error {
	meta := mr.metas[id]
	if meta==nil {
		return errors.New("no meta with ID : "+id)
	}
	meta.Status = status
	mr.Save(meta)
	return nil
}