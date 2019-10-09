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

package main

import (
	. "github.com/hyperledger/sawtooth-seth/common"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/hyperledger/sawtooth-seth/seth-cli/client"
)

type ContractList struct {
	Positional struct {
		Alias string `positional-arg-name:"alias" required:"true" description:"Alias of the imported key associated with the contract to be created"`
	} `positional-args:"true"`
}

func (args *ContractList) Name() string {
	return "list"
}

func (args *ContractList) Register(parent *flags.Command) error {
	help := `List the addresses of all contracts that could have been created based on the current nonce of the account owned by the private key with the given alias. Note that not all addresses may be valid, since the nonce increments whenever a transaction is sent from an account`
	_, err := parent.AddCommand(
		args.Name(), help,
		"", args,
	)
	if err != nil {
		return err
	}
	return nil
}

func (args *ContractList) Run(config *Config) error {
	client := client.New(config.Url)

	priv, err := LoadKey(args.Positional.Alias)
	if err != nil {
		return err
	}

	nonce, err := client.LookupAccountNonce(priv)
	if err != nil {
		return err
	}

	addr, err := PrivToEvmAddr(priv)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Address of contracts by nonce for account with alias `%v`\n",
		args.Positional.Alias,
	)
	var i uint64
	for i = 1; i < nonce; i++ {
		var derived = addr.Derive(i)
		entry, err := client.Get(derived.Bytes())

		// Abort on error
		if err != nil {
			fmt.Printf("Encountered error while listing contracts: %s", err)
			break
		}

		// Skip non-existent contracts
		if entry == nil {
			continue
		}

		fmt.Printf("%v: %v\n", i, derived)
	}

	return nil
}
