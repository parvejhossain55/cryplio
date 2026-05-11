// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./EscrowState.sol";
import "../lib/EscrowTypes.sol";
import "../interface/ICrpylioEscrow.sol";

/**
 * @title EscrowViews
 * @dev Base contract containing view functions for escrow data
 * @notice All functions are read-only and do not modify state
 */
abstract contract EscrowViews is ICryplioEscrow, EscrowState {
    // EscrowState is initialized by EscrowOperations, no constructor needed here

    /**
     * @dev Check if a token is supported
     * @param token Address of the token to check
     */
    function isTokenSupported(address token) external view virtual override returns (bool) {
        return supportedTokens[token];
    }

    /**
     * @dev Get escrow details
     * @param tradeId Trade identifier
     */
    function getEscrow(bytes32 tradeId) 
        external 
        view 
        virtual
        override
        returns (
            address buyer,
            address seller,
            address token,
            uint256 amount,
            uint256 expiresAt,
            EscrowTypes.EscrowStatus status
        ) 
    {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return (
            escrow.buyer,
            escrow.seller,
            escrow.token,
            escrow.amount,
            escrow.expiresAt,
            escrow.status
        );
    }
    
    /**
     * @dev Check if an escrow exists
     * @param tradeId Trade identifier
     */
    function escrowExists(bytes32 tradeId) external view virtual override returns (bool) {
        return _escrowExists(tradeId);
    }
        
    /**
     * @dev Get escrow status
     * @param tradeId Trade identifier
     */
    function getEscrowStatus(bytes32 tradeId) external view virtual override returns (EscrowTypes.EscrowStatus) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return escrow.status;
    }

    /**
     * @dev Check if escrow is expired
     * @param tradeId Trade identifier
     */
    function isEscrowExpired(bytes32 tradeId) external view virtual override returns (bool) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return EscrowTypes.isExpired(escrow.expiresAt);
    }
    /**
     * @dev Check if escrow can be refunded
     * @param tradeId Trade identifier
     */
    function canEscrowBeRefunded(bytes32 tradeId) external view returns (bool) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        return EscrowTypes.canRefund(escrow.expiresAt);
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
}