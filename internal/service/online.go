/********************************************************************************
 	Apache License 2.0
 	Copyright (c) 2020-2021 Tendermint
 	Copyright (c) 2022 Zondax AG

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 *********************************************************************************/

package service

import (
	"context"
	"time"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tendermint/cosmos-rosetta-gateway/errors"
	crgerrs "github.com/tendermint/cosmos-rosetta-gateway/errors"
	crgtypes "github.com/tendermint/cosmos-rosetta-gateway/types"
)

// genesisBlockFetchTimeout defines a timeout to fetch the genesis block
const genesisBlockFetchTimeout = 15 * time.Second

// NewOnlineNetwork builds a single network adapter.
// It will get the Genesis block on the beginning to avoid calling it everytime.
func NewOnlineNetwork(network *types.NetworkIdentifier, client crgtypes.Client) (crgtypes.API, error) {
	ctx, cancel := context.WithTimeout(context.Background(), genesisBlockFetchTimeout)
	defer cancel()

	var genesisHeight int64 = 1
	block, err := client.BlockByHeight(ctx, &genesisHeight)
	if err != nil {
		return OnlineNetwork{}, err
	}

	return OnlineNetwork{
		client:                 client,
		network:                network,
		networkOptions:         networkOptionsFromClient(client),
		genesisBlockIdentifier: block.Block,
	}, nil
}

// OnlineNetwork groups together all the components required for the full rosetta implementation
type OnlineNetwork struct {
	client crgtypes.Client // used to query cosmos app + tendermint

	network        *types.NetworkIdentifier      // identifies the network, it's static
	networkOptions *types.NetworkOptionsResponse // identifies the network options, it's static

	genesisBlockIdentifier *types.BlockIdentifier // identifies genesis block, it's static
}

// AccountsCoins - relevant only for UTXO based chain
// see https://www.rosetta-api.org/docs/AccountApi.html#accountcoins
func (o OnlineNetwork) AccountCoins(_ context.Context, _ *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

// networkOptionsFromClient builds network options given the client
func networkOptionsFromClient(client crgtypes.Client) *types.NetworkOptionsResponse {
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: crgtypes.SpecVersion,
			NodeVersion:    client.Version(),
		},
		Allow: &types.Allow{
			OperationStatuses:       client.OperationStatuses(),
			OperationTypes:          client.SupportedOperations(),
			Errors:                  errors.SealAndListErrors(),
			HistoricalBalanceLookup: true,
		},
	}
}
