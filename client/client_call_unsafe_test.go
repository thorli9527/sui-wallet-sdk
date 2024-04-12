package client

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/thorli9527/sui-wallet-sdk/sui_types"

	"github.com/stretchr/testify/require"
	"github.com/thorli9527/sui-wallet-sdk/account"
	"github.com/thorli9527/sui-wallet-sdk/types"
)

func TestClient_TransferObject(t *testing.T) {
	cli := MainnetClient(t)
	signer := SuiAddressNoErr("0x57188743983628b3474648d8aa4a9ee8abebe8f6816243773d7e8ed4fd833a28")
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 2)
	coin := coins.Data[0]

	txn, err := cli.TransferObject(
		context.Background(), *signer, *recipient,
		coin.CoinObjectId, nil, types.NewSafeSuiBigInt(SUI(0.01).Uint64()),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_TransferSui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.0001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 1, 0)
	require.Nil(t, err)

	txn, err := cli.TransferSui(
		context.Background(), *signer, *recipient,
		pickedCoins.Coins[0].CoinObjectId,
		types.NewSafeSuiBigInt(amount),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_PayAllSui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.Nil(t, err)

	txn, err := cli.PayAllSui(
		context.Background(), *signer, *recipient,
		pickedCoins.CoinIds(),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_Pay(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := Address
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas

	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := cli.Pay(
		context.Background(), *signer,
		pickedCoins.CoinIds(),
		[]suiAddress{*recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(amount),
		},
		nil,
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_PaySui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := Address
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := cli.PaySui(
		context.Background(), *signer,
		pickedCoins.CoinIds(),
		[]suiAddress{*recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(amount),
		},
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_SplitCoin(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.01).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)
	splitCoins := []types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amount / 2)}

	txn, err := cli.SplitCoin(
		context.Background(), *signer,
		pickedCoins.Coins[0].CoinObjectId,
		splitCoins,
		nil, types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, false)
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.01).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)

	txn, err := cli.SplitCoinEqual(
		context.Background(), *signer,
		pickedCoins.Coins[0].CoinObjectId,
		types.NewSafeSuiBigInt(uint64(2)),
		nil, types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_MergeCoins(t *testing.T) {
	cli := ChainClient(t)
	signer := Address
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= 3)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txn, err := cli.MergeCoins(
		context.Background(), *signer,
		coin1.CoinObjectId, coin2.CoinObjectId,
		&coin3.CoinObjectId, coin3.Balance,
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}

func TestClient_Publish(t *testing.T) {
	t.Log("TestClient_Publish TODO")
	// cli := ChainClient(t)

	// txnBytes, err := cli.Publish(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_MoveCall(t *testing.T) {
	t.Log("TestClient_MoveCall TODO")
	// cli := ChainClient(t)

	// txnBytes, err := cli.MoveCall(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_BatchTransaction(t *testing.T) {
	t.Log("TestClient_BatchTransaction TODO")
	// cli := ChainClient(t)

	// txnBytes, err := cli.BatchTransaction(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

// @return types.DryRunTransactionBlockResponse
func simulateCheck(
	t *testing.T,
	cli *Client,
	txBytes suiBase64Data,
	showJson bool,
) *types.DryRunTransactionBlockResponse {
	simulate, err := cli.DryRunTransaction(context.Background(), txBytes)
	require.Nil(t, err)
	require.Equal(t, simulate.Effects.Data.V1.Status.Error, "")
	require.True(t, simulate.Effects.Data.IsSuccess())
	if showJson {
		data, err := json.Marshal(simulate)
		require.Nil(t, err)
		t.Log(string(data))
		t.Log("gasFee = ", simulate.Effects.Data.GasFee())
	}
	return simulate
}

func executeTxn(
	t *testing.T,
	cli *Client,
	txBytes suiBase64Data,
	acc *account.Account,
) *types.SuiTransactionBlockResponse {
	// First of all, make sure that there are no problems with simulated trading.
	simulate, err := cli.DryRunTransaction(context.Background(), txBytes)
	require.Nil(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// sign and send
	signature, err := acc.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	require.NoError(t, err)
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	resp, err := cli.ExecuteTransactionBlock(
		context.TODO(), txBytes, []any{signature}, &options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	require.NoError(t, err)
	t.Log(resp)
	return resp
}
