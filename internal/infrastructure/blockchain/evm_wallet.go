package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	walletdomain "cryplio/internal/domain/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EvmWalletClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	chainID    *big.Int
}

func NewEvmWalletClient(rawUrl string, pkHex string) (*EvmWalletClient, error) {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("dial eth: %w", err)
	}

	pk, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &EvmWalletClient{
		client:     client,
		privateKey: pk,
		fromAddr:   fromAddr,
		chainID:    chainID,
	}, nil
}

func (c *EvmWalletClient) CreateDepositAddress(ctx context.Context, cryptoID int, userID string) (string, error) {
	// In a real custodial setup, we might generate a new HD wallet address or use a proxy contract
	// For MVP, we use the main fromAddr or generate a predictable sub-address
	return c.fromAddr.Hex(), nil
}

func (c *EvmWalletClient) Send(ctx context.Context, tx *walletdomain.WalletTransaction, destination string) (string, error) {
	nonce, err := c.client.PendingNonceAt(ctx, c.fromAddr)
	if err != nil {
		return "", err
	}

	value := big.NewInt(0) // Assuming token transfer, value is 0
	gasLimit := uint64(21000)
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(destination)
	var data []byte // Placeholder for ERC20 transfer data

	ethtx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(ethtx, types.NewEIP155Signer(c.chainID), c.privateKey)
	if err != nil {
		return "", err
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func (c *EvmWalletClient) Watch(ctx context.Context, txHash string) error {
	hash := common.HexToHash(txHash)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
			_, isPending, err := c.client.TransactionByHash(ctx, hash)
			if err != nil {
				return err
			}
			if !isPending {
				return nil
			}
		}
	}
}

// time.After used in Watch. Need to import time.
