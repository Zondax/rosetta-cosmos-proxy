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

	"github.com/coinbase/rosetta-sdk-go/types"
	crgerrs "github.com/tendermint/cosmos-rosetta-gateway/errors"
	crgtypes "github.com/tendermint/cosmos-rosetta-gateway/types"
)

// NewOffline instantiates the instance of an offline network
// whilst the offline network does not support the DataAPI,
// it supports a subset of the construction API.
func NewOffline(network *types.NetworkIdentifier, client crgtypes.Client) (crgtypes.API, error) {
	return OfflineNetwork{
		OnlineNetwork{
			client:         client,
			network:        network,
			networkOptions: networkOptionsFromClient(client),
		},
	}, nil
}

// OfflineNetwork implements an offline data API
// which is basically a data API that constantly
// returns errors, because it cannot be used if offline
type OfflineNetwork struct {
	OnlineNetwork
}

// Implement DataAPI in offline mode, which means no method is available
func (o OfflineNetwork) AccountBalance(_ context.Context, _ *types.AccountBalanceRequest) (*types.AccountBalanceResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) Block(_ context.Context, _ *types.BlockRequest) (*types.BlockResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) BlockTransaction(_ context.Context, _ *types.BlockTransactionRequest) (*types.BlockTransactionResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) Mempool(_ context.Context, _ *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) MempoolTransaction(_ context.Context, _ *types.MempoolTransactionRequest) (*types.MempoolTransactionResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) NetworkStatus(_ context.Context, _ *types.NetworkRequest) (*types.NetworkStatusResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) ConstructionSubmit(_ context.Context, _ *types.ConstructionSubmitRequest) (*types.TransactionIdentifierResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}

func (o OfflineNetwork) ConstructionMetadata(_ context.Context, _ *types.ConstructionMetadataRequest) (*types.ConstructionMetadataResponse, *types.Error) {
	return nil, crgerrs.ToRosetta(crgerrs.ErrOffline)
}
