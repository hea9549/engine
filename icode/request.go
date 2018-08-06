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

import "github.com/rs/xid"

type Request struct {
	Uuid     string
	ICodeID  string
	Function string
	Args     []string
	Type     string
}

func NewInvoke(txId string,icodeId string, function string,args []string) Request {
	return Request{
		Uuid:     txId,
		ICodeID:  icodeId,
		Function: function,
		Args:     args,
		Type:     "invoke",
	}
}

func NewQuery(icodeId string, function string,args []string) Request {
	return Request{
		Uuid:     xid.New().String(),
		ICodeID:  icodeId,
		Function: function,
		Args:     args,
		Type:     "query",
	}
}
