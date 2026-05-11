package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	trading "cryplio/internal/domain/trading"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type EvmEscrowClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	contract   common.Address
	chainID    *big.Int
	auth       *bind.TransactOpts
	parsedABI  *abi.ABI
}

func NewEvmEscrowClient(rawUrl string, pkHex string, contractAddr string, abiPath string) (*EvmEscrowClient, error) {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("dial eth: %w", err)
	}

	pk, err := crypto.HexToECDSA(strings.TrimPrefix(pkHex, "0x"))
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

	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	// Set reasonable gas price and limit
	auth.GasLimit = 300000                       // Adjust based on actual gas usage
	auth.GasPrice = big.NewInt(params.GWei * 20) // 20 Gwei

	// Load and parse ABI from JSON file
	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("read ABI file: %w", err)
	}

	var abiWrapper struct {
		ABI json.RawMessage `json:"abi"`
	}
	if err := json.Unmarshal(abiBytes, &abiWrapper); err != nil {
		return nil, fmt.Errorf("unmarshal ABI JSON: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiWrapper.ABI)))
	if err != nil {
		return nil, fmt.Errorf("parse ABI: %w", err)
	}

	return &EvmEscrowClient{
		client:     client,
		privateKey: pk,
		fromAddr:   fromAddr,
		contract:   common.HexToAddress(contractAddr),
		chainID:    chainID,
		auth:       auth,
		parsedABI:  &parsedABI,
	}, nil
}

func (c *EvmEscrowClient) Lock(ctx context.Context, trade *trading.Trade) (string, string, error) {
	// Create contract binding
	boundContract := bind.NewBoundContract(c.contract, *c.parsedABI, c.client, c.client, c.client)

	// Convert trade ID to bytes32
	tradeIdBytes := common.HexToHash(trade.TradeID.String())

	// Convert addresses
	buyerAddr := common.HexToAddress(trade.BuyerID.String())
	sellerAddr := common.HexToAddress(trade.SellerID.String())

	// For MVP, we assume USDT token address - in production, get from config
	usdtAddr := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") // Mainnet USDT

	// Convert amount to wei (USDT has 6 decimals)
	amount := big.NewInt(int64(trade.CryptoAmount * 1000000)) // Convert to USDT smallest unit

	// Expiry time in seconds (e.g. 1 hour = 3600)
	expiryTime := big.NewInt(3600) // Default for MVP

	// Call createEscrow
	auth := *c.auth
	auth.Context = ctx

	tx, err := boundContract.Transact(&auth, "createEscrow", tradeIdBytes, buyerAddr, sellerAddr, usdtAddr, amount, expiryTime)
	if err != nil {
		return "", "", fmt.Errorf("call createEscrow: %w", err)
	}

	if tx == nil {
		return "", "", fmt.Errorf("transaction was nil")
	}

	return tx.Hash().Hex(), c.contract.Hex(), nil
}

func (c *EvmEscrowClient) Release(ctx context.Context, trade *trading.Trade) (string, error) {
	boundContract := bind.NewBoundContract(c.contract, *c.parsedABI, c.client, c.client, c.client)

	tradeIdBytes := common.HexToHash(trade.TradeID.String())

	auth := *c.auth
	auth.Context = ctx

	tx, err := boundContract.Transact(&auth, "releaseEscrow", tradeIdBytes)
	if err != nil {
		return "", fmt.Errorf("call releaseEscrow: %w", err)
	}

	if tx == nil {
		return "", fmt.Errorf("transaction was nil")
	}

	return tx.Hash().Hex(), nil
}

func (c *EvmEscrowClient) Refund(ctx context.Context, trade *trading.Trade) (string, error) {
	boundContract := bind.NewBoundContract(c.contract, *c.parsedABI, c.client, c.client, c.client)

	tradeIdBytes := common.HexToHash(trade.TradeID.String())

	auth := *c.auth
	auth.Context = ctx

	tx, err := boundContract.Transact(&auth, "refundEscrow", tradeIdBytes)
	if err != nil {
		return "", fmt.Errorf("call refundEscrow: %w", err)
	}

	if tx == nil {
		return "", fmt.Errorf("transaction was nil")
	}

	return tx.Hash().Hex(), nil
}

func (c *EvmEscrowClient) AdminRefund(ctx context.Context, trade *trading.Trade) (string, error) {
	return c.Refund(ctx, trade)
}

// GetEscrowStatus returns the current status of an escrow
func (c *EvmEscrowClient) GetEscrowStatus(ctx context.Context, tradeId string) (uint8, error) {
	boundContract := bind.NewBoundContract(c.contract, *c.parsedABI, c.client, c.client, c.client)

	tradeIdBytes := common.HexToHash(tradeId)

	var results []interface{}
	err := boundContract.Call(nil, &results, "getEscrow", tradeIdBytes)
	if err != nil {
		return 0, fmt.Errorf("call getEscrow: %w", err)
	}

	if len(results) < 6 {
		return 0, fmt.Errorf("unexpected result length")
	}

	status, ok := results[5].(uint8)
	if !ok {
		return 0, fmt.Errorf("invalid status type")
	}

	return status, nil
}
