package utils

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TestCCStub struct {
	State    map[string][]byte
	function string
	args     []string
}

func (tStub TestCCStub) GetArgs() [][]byte       { return nil }
func (tStub TestCCStub) GetStringArgs() []string { return nil }

func (tStub *TestCCStub) SetFunctionAndParameters(function string, args ...string) {
	tStub.function = function
	tStub.args = args
}
func (tStub TestCCStub) GetFunctionAndParameters() (string, []string) {
	return tStub.function, tStub.args
}

func (tStub TestCCStub) GetArgsSlice() ([]byte, error) { return nil, nil }
func (tStub TestCCStub) GetTxID() string               { return "" }
func (tStub TestCCStub) GetChannelID() string          { return "" }
func (tStub TestCCStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	return pb.Response{}
}
func (tStub TestCCStub) GetState(key string) ([]byte, error) {
	return tStub.State[key], nil
}
func (tStub TestCCStub) PutState(key string, value []byte) error {
	tStub.State[key] = value
	return nil
}
func (tStub TestCCStub) DelState(key string) error                               { return nil }
func (tStub TestCCStub) SetStateValidationParameter(key string, ep []byte) error { return nil }
func (tStub TestCCStub) GetStateValidationParameter(key string) ([]byte, error)  { return nil, nil }
func (tStub TestCCStub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}
func (tStub TestCCStub) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}
func (tStub TestCCStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return "", nil
}
func (tStub TestCCStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return "", nil, nil
}
func (tStub TestCCStub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}
func (tStub TestCCStub) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetPrivateData(collection, key string) ([]byte, error) { return nil, nil }
func (tStub TestCCStub) PutPrivateData(collection string, key string, value []byte) error {
	return nil
}
func (tStub TestCCStub) DelPrivateData(collection, key string) error { return nil }
func (tStub TestCCStub) SetPrivateDataValidationParameter(collection, key string, ep []byte) error {
	return nil
}
func (tStub TestCCStub) GetPrivateDataValidationParameter(collection, key string) ([]byte, error) {
	return nil, nil
}
func (tStub TestCCStub) GetPrivateDataByRange(collection, startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (tStub TestCCStub) GetCreator() ([]byte, error)                    { return nil, nil }
func (tStub TestCCStub) GetTransient() (map[string][]byte, error)       { return nil, nil }
func (tStub TestCCStub) GetBinding() ([]byte, error)                    { return nil, nil }
func (tStub TestCCStub) GetDecorations() map[string][]byte              { return nil }
func (tStub TestCCStub) GetSignedProposal() (*pb.SignedProposal, error) { return nil, nil }
func (tStub TestCCStub) GetTxTimestamp() (*timestamp.Timestamp, error)  { return nil, nil }
func (tStub TestCCStub) SetEvent(name string, payload []byte) error     { return nil }
