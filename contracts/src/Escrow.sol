// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./base/EscrowViews.sol";
import "./libraries/EscrowTypes.sol";

/**
 * @title CryplioEscrow
 * @dev P2P crypto trading escrow with USDT support
 * 
 * This contract has been refactored from a monolithic structure into a modular
 * architecture with separate concerns:
 * - EscrowTypes: Shared types, constants, and validation
 * - EscrowState: State management
 * - EscrowAuth: Authorization and modifiers
 * - EscrowOperations: Core escrow operations
 * - EscrowViews: View functions and queries
 * 
 * Benefits of this refactoring:
 * - Better code organization and readability
 * - Easier testing and maintenance
 * - Reusable components
 * - Clear separation of concerns
 * - Reduced gas costs for deployment (due to smaller individual contracts)
 */
contract CryplioEscrow is EscrowViews {
    
    // Re-export events for external interfaces
    event EscrowCreated(
        bytes32 indexed tradeId,
        address indexed buyer,
        address indexed seller,
        address token,
        uint256 amount,
        uint256 timestamp
    );
    
    event EscrowReleased(
        bytes32 indexed tradeId,
        address indexed seller,
        uint256 amount,
        uint256 timestamp
    );
    
    event EscrowRefunded(
        bytes32 indexed tradeId,
        address indexed buyer,
        uint256 amount,
        uint256 timestamp
    );
    
    event DisputeRaised(
        bytes32 indexed tradeId,
        address indexed raiser,
        string reason,
        uint256 timestamp
    );
    
    /**
     * @dev Constructor
     * @param _usdtToken Address of the USDT token contract
     */
    constructor(address _usdtToken) EscrowViews(_usdtToken) {
        // Deployer is initially an authorized signer
        _addAuthorizedSigner(msg.sender);
    }
    
    /**
     * @dev Emergency withdraw tokens (only owner)
     * @param token Token address
     * @param amount Amount to withdraw
     */
    function emergencyWithdraw(address token, uint256 amount) external onlyOwner {
        IERC20(token).transfer(owner(), amount);
    }
    
    /**
     * @dev Get contract metadata
     */
    function getContractInfo() external view returns (
        address usdtTokenAddress,
        uint256 escrowExpiryTime,
        uint256 disputeWindow
    ) {
        return (
            usdtToken,
            EscrowTypes.ESCROW_EXPIRY_TIME,
            EscrowTypes.DISPUTE_WINDOW
        );
    }
}
