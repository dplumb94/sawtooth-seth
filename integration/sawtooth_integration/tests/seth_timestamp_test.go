/**
 * Copyright 2017 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

 package tests

import (
  "github.com/hyperledger/burrow/binary"
  "encoding/base64"
  "encoding/hex"
  "fmt"
  "net/http"
  "github.com/golang/protobuf/proto"
  c "github.com/hyperledger/sawtooth-seth/seth-cli/client"
  . "github.com/hyperledger/sawtooth-seth/common"
  . "github.com/hyperledger/sawtooth-seth/protobuf/block_info_pb2"
  "testing"
)

const (
  WAIT = 300
  PRIV = "274AAF0BFAF68C5F1CB510D2851A71CF720C7218A2C1637D635F6850FF009774"
  INIT_TIMESTAMP = "6060604052341561000f57600080fd5b5b60a68061001e6000396000f300606060405263ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663ffcfdd768114603c575b600080fd5b3415604657600080fd5b605160043515156063565b60405190815260200160405180910390f35b6000811560705750426074565b5060005b5b9190505600a165627a7a723058207f054877f30897a342148fdfdf717bdeef570ffd8d1f918cf0feca42f50b3e1b0029"
  TIMESTAMP = "ffcfdd760000000000000000000000000000000000000000000000000000000000000001"
)

// Test getting the timestamp of the current block
func TestBlockNumber(t *testing.T) {
  client := c.New("http://rest-api:8080")
  priv, _ := hex.DecodeString(PRIV)
  init, _ := hex.DecodeString(INIT_TIMESTAMP)
  nonce := uint64(0)

  // Create the EOA
  _, err := client.CreateExternalAccount(priv, nil, nil, 0, WAIT)
  if err != nil {
    t.Fatal(err.Error())
  }
  nonce += 1

  // Create the Contract
  contractCreateResult, err := client.CreateContractAccount(priv, init, nil, nonce, 1000, WAIT)
  if err != nil {
   t.Fatal(err.Error())
  }
  nonce += 1

  cmd, _ := hex.DecodeString(TIMESTAMP)
  contractCallResult, err := client.MessageCall(priv, contractCreateResult.Address, cmd, nonce, 1000, WAIT, false)
  timestamp := binary.Uint64FromWord256(binary.RightPadWord256(contractCallResult.ReturnValue))

  // Get number of current block from BlockInfo
  blockInfoAddr, err := NewBlockInfoAddr(2);
  if err != nil {
    t.Fatal(err.Error())
  }
  resp, err := http.Get(fmt.Sprintf("%s/state/%s?wait=%v", client.Url, blockInfoAddr, WAIT))
  if err != nil {
    t.Fatal(err.Error())
  }
  body, err := c.ParseRespBody(resp)
  if err != nil {
    t.Fatal(err.Error())
  }
  buf, err := base64.StdEncoding.DecodeString(body.Data.(string))
  if err != nil {
    t.Fatal(err.Error())
  }
  blockInfo := &BlockInfo{}
  err = proto.Unmarshal(buf, blockInfo)
  expectedTimestamp := blockInfo.GetTimestamp()

  if expectedTimestamp != timestamp {
    t.Fatalf("Incorrect timestamp: %v, expected: %v", timestamp, expectedTimestamp)
  }
}
