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
	"github.com/tendermint/cosmos-rosetta-gateway/errors"
	crgtypes "github.com/tendermint/cosmos-rosetta-gateway/types"
)

// AccountBalance retrieves the account balance of an address
// rosetta requires us to fetch the block information too
func (on OnlineNetwork) AccountBalance(ctx context.Context, request *types.AccountBalanceRequest) (*types.AccountBalanceResponse, *types.Error) {
	var (
		height int64
		block  crgtypes.BlockResponse
		err    error
	)

	switch {
	case request.BlockIdentifier == nil:
		block, err = on.client.BlockByHeight(ctx, nil)
		if err != nil {
			return nil, errors.ToRosetta(err)
		}
	case request.BlockIdentifier.Hash != nil:
		block, err = on.client.BlockByHash(ctx, *request.BlockIdentifier.Hash)
		if err != nil {
			return nil, errors.ToRosetta(err)
		}
		height = block.Block.Index
	case request.BlockIdentifier.Index != nil:
		height = *request.BlockIdentifier.Index
		block, err = on.client.BlockByHeight(ctx, &height)
		if err != nil {
			return nil, errors.ToRosetta(err)
		}
	}

	accountCoins, err := on.client.Balances(ctx, request.AccountIdentifier.Address, &height)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	return &types.AccountBalanceResponse{
		BlockIdentifier: block.Block,
		Balances:        accountCoins,
		Metadata:        nil,
	}, nil
}

// Block gets the transactions in the given block
func (on OnlineNetwork) Block(ctx context.Context, request *types.BlockRequest) (*types.BlockResponse, *types.Error) {
	var (
		blockResponse crgtypes.BlockTransactionsResponse
		err           error
	)
	// block identifier is assumed not to be nil as rosetta will do this check for us
	// check if we have to query via hash or block number
	switch {
	case request.BlockIdentifier.Hash != nil:
		blockResponse, err = on.client.BlockTransactionsByHash(ctx, *request.BlockIdentifier.Hash)
		if err != nil {
			return nil, errors.ToRosetta(err)
		}
	case request.BlockIdentifier.Index != nil:
		blockResponse, err = on.client.BlockTransactionsByHeight(ctx, request.BlockIdentifier.Index)
		if err != nil {
			return nil, errors.ToRosetta(err)
		}
	default:
		err := errors.WrapError(errors.ErrBadArgument, "at least one of hash or index needs to be specified")
		return nil, errors.ToRosetta(err)
	}

	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier:       blockResponse.Block,
			ParentBlockIdentifier: blockResponse.ParentBlock,
			Timestamp:             blockResponse.MillisecondTimestamp,
			Transactions:          blockResponse.Transactions,
			Metadata:              nil,
		},
		OtherTransactions: nil,
	}, nil
}

// BlockTransaction gets the given transaction in the specified block, we do not need to check the block itself too
// due to the fact that tendermint achieves instant finality
func (on OnlineNetwork) BlockTransaction(ctx context.Context, request *types.BlockTransactionRequest) (*types.BlockTransactionResponse, *types.Error) {
	tx, err := on.client.GetTx(ctx, request.TransactionIdentifier.Hash)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	return &types.BlockTransactionResponse{
		Transaction: tx,
	}, nil
}

// Mempool fetches the transactions contained in the mempool
func (on OnlineNetwork) Mempool(ctx context.Context, _ *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	txs, err := on.client.Mempool(ctx)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	return &types.MempoolResponse{
		TransactionIdentifiers: txs,
	}, nil
}

// MempoolTransaction fetches a single transaction in the mempool
// NOTE: it is not implemented yet
func (on OnlineNetwork) MempoolTransaction(ctx context.Context, request *types.MempoolTransactionRequest) (*types.MempoolTransactionResponse, *types.Error) {
	tx, err := on.client.GetUnconfirmedTx(ctx, request.TransactionIdentifier.Hash)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	return &types.MempoolTransactionResponse{
		Transaction: tx,
	}, nil
}

func (on OnlineNetwork) NetworkList(_ context.Context, _ *types.MetadataRequest) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{NetworkIdentifiers: []*types.NetworkIdentifier{on.network}}, nil
}

func (on OnlineNetwork) NetworkOptions(_ context.Context, _ *types.NetworkRequest) (*types.NetworkOptionsResponse, *types.Error) {
	return on.networkOptions, nil
}

func (on OnlineNetwork) NetworkStatus(ctx context.Context, _ *types.NetworkRequest) (*types.NetworkStatusResponse, *types.Error) {
	block, err := on.client.BlockByHeight(ctx, nil)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	peers, err := on.client.Peers(ctx)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	syncStatus, err := on.client.Status(ctx)
	if err != nil {
		return nil, errors.ToRosetta(err)
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: block.Block,
		CurrentBlockTimestamp:  block.MillisecondTimestamp,
		GenesisBlockIdentifier: on.genesisBlockIdentifier,
		OldestBlockIdentifier:  nil,
		SyncStatus:             syncStatus,
		Peers:                  peers,
	}, nil
}
