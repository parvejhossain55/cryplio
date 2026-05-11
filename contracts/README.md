# Cryplio Smart Contract

A secure, non-custodial P2P escrow system specialized for stablecoins (USDT, USDC, DAI).

## 1. Overview
Cryplio Escrow facilitates trustless trades between buyers and sellers using stablecoins. It locks seller funds in a smart contract and releases them only after payment confirmation or through an authorized dispute resolution process. By focusing exclusively on stablecoins, the platform ensures that trade values remain stable during the escrow period.

## 2. Key Features
- **Non-Custodial**: Users maintain control over their funds through the smart contract.
- **Meta-Transactions (EIP-712)**: Supports gasless escrow creation for a better user experience.
- **Dispute Resolution**: Authorized admin signers can resolve conflicts by refunding or force-releasing funds.
- **Emergency Pause**: Built-in circuit breaker for emergency situations.
- **Fee Mechanism**: 
    - **Standard Release Fee**: 0.75% (75 BPS)
    - **Refund Fee**: 0.25% (25 BPS)
- **Gas Optimized**: Efficient storage packing and caching to minimize transaction costs.

## 3. Architecture
The system is built using a modular inheritance structure:
- **[CryplioEscrow.sol](file:///home/parvej/Project/Cryplio_Smart_Contract/src/CryplioEscrow.sol)**: The main entry point that combines all functionalities.
- **[EscrowOperations.sol](file:///home/parvej/Project/Cryplio_Smart_Contract/src/base/EscrowOperations.sol)**: Handles state-changing logic (creation, release, refund).
- **[EscrowAuth.sol](file:///home/parvej/Project/Cryplio_Smart_Contract/src/base/EscrowAuth.sol)**: Manages access control, admin signers, and emergency pause.
- **[EscrowState.sol](file:///home/parvej/Project/Cryplio_Smart_Contract/src/base/EscrowState.sol)**: Defines the storage layout and internal state management.
- **[EscrowTypes.sol](file:///home/parvej/Project/Cryplio_Smart_Contract/src/lib/EscrowTypes.sol)**: Central library for events, errors, structs, and validation logic.
- **[ICryplioEscrow.sol](file:///home/parvej/Project/Cryplio_Smart_Contract/src/interface/ICrpylioEscrow.sol)**: External interface defining the public API.

## 4. Development & Usage
This project uses [Foundry](https://book.getfoundry.sh/).

### Installation
```shell
forge install
```

### Local Development

1. **Start Anvil**:
   ```shell
   anvil
   ```
2. **Deploy Mock USDT** (Optional, for local testing):
   ```shell
   source .env && forge script script/DeployMockUSDT.s.sol:DeployMockUSDT --rpc-url http://127.0.0.1:8545 --broadcast
   ```
3. **Deploy CryplioEscrow**:
   ```shell
   source .env && forge script script/Deploy.s.sol:DeployCryplioEscrow --rpc-url http://127.0.0.1:8545 --broadcast
   ```

### Build
```shell
forge build
```

### Test
```shell
forge test --match-path test/CryplioEscrow.t.sol
```

### Deploy
Deployment scripts are located in `script/`.

1. Create a `.env` file based on [.env.example](./.env.example):
   ```shell
   cp .env.example .env
   ```
2. Fill in your variables in `.env`.
3. Run the deployment script:
   ```shell
   source .env && forge script script/Deploy.s.sol:DeployCryplioEscrow --rpc-url $SEPOLIA_RPC_URL --broadcast --verify
   ```

**Required Environment Variables:**
- `PRIVATE_KEY`: Deployer's private key.
- `FEE_RECIPIENT`: Address to receive platform fees.
- `INITIAL_TOKEN_ADDRESS`: First supported ERC20 token.
- `INITIAL_SIGNER`: First authorized admin signer.

## 5. Integration Guide

### Creating an Escrow (On-chain)
```solidity
escrow.createEscrow(tradeId, buyer, seller, token, amount, expiryTime);
```

### Creating an Escrow (Meta-transaction)
1. Sign the `CreateEscrow` struct hash following EIP-712.
2. Submit via:
```solidity
escrow.createEscrowMeta(tradeId, buyer, seller, token, amount, expiryTime, nonce, signature);
```

## 6. Security & Safety
- **Reentrancy Protection**: All state-changing functions use OpenZeppelin's `ReentrancyGuard`.
- **Emergency Pause**: The owner can pause new escrow creations in case of emergency.
- **Token Whitelisting**: Only owner-approved ERC20 tokens can be used.
- **Signatures**: Uses EIP-712 standard for all meta-transactions to prevent replay attacks.

## 7. Interaction Guide (using Cast)

### Read Operations
```shell
# Check if token is supported
cast call <ESCROW_ADDR> "isTokenSupported(address)(bool)" <TOKEN_ADDR> --rpc-url $RPC_URL
```

### Write Operations
```shell
# 1. Approve tokens
cast send <TOKEN_ADDR> "approve(address,uint256)" <ESCROW_ADDR> <AMOUNT> --private-key $PRIV_KEY --rpc-url $RPC_URL

# 2. Create Escrow
cast send <ESCROW_ADDR> "createEscrow(bytes32,address,address,address,uint256,uint256)" <TRADE_ID> <BUYER> <SELLER> <TOKEN> <AMOUNT> <EXPIRY> --private-key $PRIV_KEY --rpc-url $RPC_URL
```

## License
This project is licensed under the MIT License.
