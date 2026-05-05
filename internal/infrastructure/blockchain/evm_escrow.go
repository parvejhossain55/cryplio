package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	trading "cryplio/internal/domain/trading"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EvmEscrowClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	contract   common.Address
	chainID    *big.Int
}

func NewEvmEscrowClient(rawUrl string, pkHex string, contractAddr string) (*EvmEscrowClient, error) {
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

	return &EvmEscrowClient{
		client:     client,
		privateKey: pk,
		fromAddr:   fromAddr,
		contract:   common.HexToAddress(contractAddr),
		chainID:    chainID,
	}, nil
}

func (c *EvmEscrowClient) Lock(ctx context.Context, trade *trading.Trade) (string, string, error) {
	// In a real implementation, we would call the 'createEscrow' or 'lock' function on the contract
	// For now, this is a placeholder for the actual ABI call
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", "", err
	}
	auth.Context = ctx

	// Placeholder logic: assuming we have an Escrow contract instance
	// txn, err := c.contractInstance.Lock(auth, trade.BuyerID, trade.CryptoAmount, ...)

	fmt.Printf("Blockchain: Locking escrow for trade %s\n", trade.TradeID)

	return "0x..." + trade.TradeID.String()[:8], c.contract.Hex(), nil
}

func (c *EvmEscrowClient) Release(ctx context.Context, trade *trading.Trade) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", err
	}
	auth.Context = ctx

	// txn, err := c.contractInstance.Release(auth, trade.TradeID)

	fmt.Printf("Blockchain: Releasing escrow for trade %s\n", trade.TradeID)

	return "0x-release-" + trade.TradeID.String()[:8], nil
}

func (c *EvmEscrowClient) Refund(ctx context.Context, trade *trading.Trade) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", err
	}
	auth.Context = ctx

	// txn, err := c.contractInstance.Refund(auth, trade.TradeID)

	fmt.Printf("Blockchain: Refunding escrow for trade %s\n", trade.TradeID)

	return "0x-refund-" + trade.TradeID.String()[:8], nil
}
