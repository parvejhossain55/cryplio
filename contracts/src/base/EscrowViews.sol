// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./EscrowOperations.sol";
import "../libraries/EscrowTypes.sol";

/**
 * @title EscrowViews
 * @dev Base contract containing view functions for escrow data
 */
abstract contract EscrowViews is EscrowOperations {
    
    // Constructor
    constructor(address _usdtToken) EscrowOperations(_usdtToken) {}
    
    /**
     * @dev Get escrow details
     * @param tradeId Trade identifier
     */
    function getEscrow(bytes32 tradeId) 
        external 
        view 
        returns (
            address buyer,
            address seller,
            address token,
            uint256 amount,
            uint256 createdAt,
            uint256 expiresAt,
            EscrowTypes.EscrowStatus status,
            address disputeRaiser,
            string memory disputeReason
        ) 
    {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return (
            escrow.buyer,
            escrow.seller,
            escrow.token,
            escrow.amount,
            escrow.createdAt,
            escrow.expiresAt,
            escrow.status,
            escrow.disputeRaiser,
            escrow.disputeReason
        );
    }
    
    /**
     * @dev Check if an escrow exists
     * @param tradeId Trade identifier
     */
    function escrowExists(bytes32 tradeId) external view returns (bool) {
        return _escrowExists(tradeId);
    }
    
    /**
     * @dev Get escrow status
     * @param tradeId Trade identifier
     */
    function getEscrowStatus(bytes32 tradeId) external view returns (EscrowTypes.EscrowStatus) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return escrow.status;
    }
    
    /**
     * @dev Check if escrow is expired
     * @param tradeId Trade identifier
     */
    function isEscrowExpired(bytes32 tradeId) external view returns (bool) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return EscrowTypes.isExpired(escrow.expiresAt);
    }
    
    /**
     * @dev Check if escrow can be refunded
     * @param tradeId Trade identifier
     */
    function canEscrowBeRefunded(bytes32 tradeId) external view returns (bool) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return EscrowTypes.canRefund(escrow.expiresAt, escrow.status);
    }
    
    /**
     * @dev Get time remaining until escrow expires
     * @param tradeId Trade identifier
     */
    function getTimeRemaining(bytes32 tradeId) external view returns (uint256) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        if (block.timestamp >= escrow.expiresAt) {
            return 0;
        }
        return escrow.expiresAt - block.timestamp;
    }
    
    /**
     * @dev Get time remaining until dispute window closes
     * @param tradeId Trade identifier
     */
    function getDisputeWindowRemaining(bytes32 tradeId) external view returns (uint256) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        uint256 disputeDeadline = escrow.expiresAt + EscrowTypes.DISPUTE_WINDOW;
        if (block.timestamp >= disputeDeadline) {
            return 0;
        }
        return disputeDeadline - block.timestamp;
    }
    
    /**
     * @dev Get contract balance for specific token
     */
    function getTokenBalance(address token) external view returns (uint256) {
        return IERC20(token).balanceOf(address(this));
    }
    
    /**
     * @dev Get USDT token balance
     */
    function getUSDTBalance() external view returns (uint256) {
        return IERC20(usdtToken).balanceOf(address(this));
    }
    
    /**
     * @dev Get total value locked in escrow
     */
    function getTotalValueLocked() external view returns (uint256) {
        return IERC20(usdtToken).balanceOf(address(this));
    }
    
    /**
     * @dev Get escrow count by status
     * @param status The status to count
     */
    function getEscrowCountByStatus(EscrowTypes.EscrowStatus status) external view returns (uint256) {
        // Note: This is a simplified implementation
        // In production, you'd want to maintain a separate mapping for counts
        // or use events to track this more efficiently
        uint256 count = 0;
        // This would require iterating through all escrows, which is gas-intensive
        // For now, returning 0 as placeholder
        return count;
    }
    
    /**
     * @dev Get user's escrows (buyer or seller)
     * @param user User address
     * @param isBuyer Whether to get escrows where user is buyer (true) or seller (false)
     * @param offset Starting offset
     * @param limit Maximum number of escrows to return
     */
    function getUserEscrows(
        address user, 
        bool isBuyer, 
        uint256 offset, 
        uint256 limit
    ) external view returns (bytes32[] memory tradeIds) {
        // Note: This is a simplified implementation
        // In production, you'd want to maintain separate mappings for user escrows
        // or use events to track this more efficiently
        tradeIds = new bytes32[](limit);
        // Placeholder implementation
        return tradeIds;
    }
}
