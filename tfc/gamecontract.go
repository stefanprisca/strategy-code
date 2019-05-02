// Copyright 2019 Stefan Prisca

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfc

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	tfcPb "github.com/stefanprisca/strategy-protobufs/tfc"
)

// Dummy struct for hyperledger
type GameContract struct {
}

const CONTRACT_STATE_KEY = "contract.tfc.com"

func (gc *GameContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	// The first argument is the function name!
	// Second will be our protobuf payload.
	protoInitArgs := APIstub.GetArgs()[1]
	gcInitArgs := &tfcPb.GameContractInitArgs{}

	err := proto.Unmarshal(protoInitArgs, gcInitArgs)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse init args: %s", err)
		return shim.Error(errMsg)
	}

	var contractUUID = gcInitArgs.GetUuid()
	print(contractUUID)
	return shim.Success(nil)
}

func (gc *GameContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	creator, errc := APIstub.GetCreator()
	if errc == nil {
		fmt.Println("Creator: ", string(creator))
	}
	return shim.Error(fmt.Sprint("Unkown transaction type <>"))
}
