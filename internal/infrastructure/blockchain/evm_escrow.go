package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	tradingdomain "cryplio/internal/domain/trading"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EvmEscrowClient struct {
	client          *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	fromAddr        common.Address
	contractAddress common.Address
	chainID         *big.Int
	parsedABI       abi.ABI
}

func NewEvmEscrowClient(rawUrl, pkHex, contractAddr, abiPath string) (*EvmEscrowClient, error) {
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

	// Load ABI
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("read abi file: %w", err)
	}

	// The ABI might be a raw JSON array or a full Foundry/Hardhat artifact object
	var abiArray []byte
	var artifact struct {
		Abi interface{} `json:"abi"`
	}

	if err := json.Unmarshal(abiData, &artifact); err == nil && artifact.Abi != nil {
		// It's an artifact, re-marshal the 'abi' field
		abiArray, _ = json.Marshal(artifact.Abi)
	} else {
		// Assume it's a raw ABI array
		abiArray = abiData
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiArray)))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}

	return &EvmEscrowClient{
		client:          client,
		privateKey:      pk,
		fromAddr:        fromAddr,
		contractAddress: common.HexToAddress(contractAddr),
		chainID:         chainID,
		parsedABI:       parsedABI,
	}, nil
}

func (c *EvmEscrowClient) Lock(ctx context.Context, trade *tradingdomain.Trade) (string, string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", "", fmt.Errorf("create transactor: %w", err)
	}

	if trade.BuyerAddress == nil || !common.IsHexAddress(*trade.BuyerAddress) {
		return "", "", fmt.Errorf("invalid buyer address")
	}
	if trade.SellerAddress == nil || !common.IsHexAddress(*trade.SellerAddress) {
		return "", "", fmt.Errorf("invalid seller address")
	}
	if trade.TokenAddress != nil && !common.IsHexAddress(*trade.TokenAddress) {
		return "", "", fmt.Errorf("invalid token address")
	}

	tradeID := common.HexToHash(trade.TradeID.String())

	contract := bind.NewBoundContract(c.contractAddress, c.parsedABI, c.client, c.client, c.client)

	expiry := big.NewInt(time.Now().Add(time.Hour).Unix())

	tokenAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	if trade.TokenAddress != nil {
		tokenAddr = common.HexToAddress(*trade.TokenAddress)
	}

	buyerAddr := common.HexToAddress(*trade.BuyerAddress)
	sellerAddr := common.HexToAddress(*trade.SellerAddress)

	tx, err := contract.Transact(auth, "createEscrow",
		tradeID,
		buyerAddr,
		sellerAddr,
		tokenAddr,
		big.NewInt(int64(trade.CryptoAmount*1e18)),
		expiry,
	)
	if err != nil {
		return "", "", fmt.Errorf("execute createEscrow tx: %w", err)
	}

	return tx.Hash().Hex(), c.contractAddress.Hex(), nil
}

func (c *EvmEscrowClient) Release(ctx context.Context, trade *tradingdomain.Trade) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", fmt.Errorf("create transactor: %w", err)
	}

	tradeID := common.HexToHash(trade.TradeID.String())
	contract := bind.NewBoundContract(c.contractAddress, c.parsedABI, c.client, c.client, c.client)
	tx, err := contract.Transact(auth, "releaseEscrow", tradeID)
	if err != nil {
		return "", fmt.Errorf("execute releaseEscrow tx: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (c *EvmEscrowClient) Refund(ctx context.Context, trade *tradingdomain.Trade) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", fmt.Errorf("create transactor: %w", err)
	}

	tradeID := common.HexToHash(trade.TradeID.String())
	contract := bind.NewBoundContract(c.contractAddress, c.parsedABI, c.client, c.client, c.client)
	tx, err := contract.Transact(auth, "refundEscrow", tradeID)
	if err != nil {
		return "", fmt.Errorf("execute refundEscrow tx: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (c *EvmEscrowClient) AdminRefund(ctx context.Context, trade *tradingdomain.Trade) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return "", fmt.Errorf("create transactor: %w", err)
	}

	tradeID := common.HexToHash(trade.TradeID.String())
	contract := bind.NewBoundContract(c.contractAddress, c.parsedABI, c.client, c.client, c.client)
	// Using forceReleaseEscrow or similar for admin intervention if needed,
	// but here we'll stick to refundEscrow for now or match ABI if there's a specific admin one.
	tx, err := contract.Transact(auth, "refundEscrow", tradeID)
	if err != nil {
		return "", fmt.Errorf("execute admin refundEscrow tx: %w", err)
	}

	return tx.Hash().Hex(), nil
}
