package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	walletdomain "cryplio/internal/domain/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
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

	pkStr := pkHex
	if len(pkStr) > 2 && pkStr[:2] == "0x" {
		pkStr = pkStr[2:]
	}

	pk, err := crypto.HexToECDSA(pkStr)
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

func (c *EvmWalletClient) GenerateKeyPair() (string, string, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("generate key: %w", err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(crypto.FromECDSA(key))
	return address, privateKey, nil
}

func (c *EvmWalletClient) CreateDepositAddress(_ context.Context, _ int, userID string) (string, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}
	addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
	if !common.IsHexAddress(addr) {
		return "", fmt.Errorf("generated invalid address")
	}
	return addr, nil
}

func (c *EvmWalletClient) GetBalance(ctx context.Context, address string) (float64, error) {
	if !common.IsHexAddress(address) {
		return 0, fmt.Errorf("invalid address")
	}
	addr := common.HexToAddress(address)
	balance, err := c.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}

	// Convert wei to ether (float64)
	fbalance := new(big.Float).SetInt(balance)
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(params.Ether))
	result, _ := ethValue.Float64()

	return result, nil
}

func (c *EvmWalletClient) Send(ctx context.Context, tx *walletdomain.WalletTransaction, destination string) (string, error) {
	nonce, err := c.client.PendingNonceAt(ctx, c.fromAddr)
	if err != nil {
		return "", fmt.Errorf("get nonce: %w", err)
	}

	value := new(big.Int)
	// Convert float amount to wei
	amountWei := new(big.Float).Mul(big.NewFloat(tx.Amount), big.NewFloat(params.Ether))
	amountWei.Int(value)

	gasLimit := uint64(21000)
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("suggest gas price: %w", err)
	}

	toAddress := common.HexToAddress(destination)
	var data []byte

	ethtx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(ethtx, types.NewEIP155Signer(c.chainID), c.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign tx: %w", err)
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("send tx: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (c *EvmWalletClient) Watch(ctx context.Context, txHash string) error {
	// Simple polling mechanism for transaction confirmation
	hash := common.HexToHash(txHash)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			receipt, err := c.client.TransactionReceipt(ctx, hash)
			if err != nil {
				// Transaction might still be pending or not found yet
				time.Sleep(2 * time.Second)
				continue
			}
			if receipt.Status == types.ReceiptStatusFailed {
				return fmt.Errorf("transaction failed on-chain")
			}
			return nil // Transaction confirmed
		}
	}
}
