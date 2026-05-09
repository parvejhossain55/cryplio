// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title EscrowTypes
 * @dev Library containing shared types, constants, and events for escrow contracts
 */
library EscrowTypes {
    // Events
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

    // Structs
    struct Escrow {
        bytes32 tradeId;
        address buyer;
        address seller;
        address token;
        uint256 amount;
        uint256 createdAt;
        uint256 expiresAt;
        EscrowStatus status;
        address disputeRaiser;
        string disputeReason;
    }
    
    // Enums
    enum EscrowStatus {
        Pending,
        Locked,
        Released,
        Refunded,
        Disputed
    }
    
    // Constants
    uint256 public constant ESCROW_EXPIRY_TIME = 1 hours; // 1 hours for payment
    uint256 public constant DISPUTE_WINDOW = 2 days;      // 2 days to raise dispute
    
    // Errors
    error EscrowNotFound();
    error EscrowAlreadyExists();
    error InvalidAddresses();
    error InvalidAmount();
    error UnsupportedToken();
    error TransferFailed();
    error InvalidStatus();
    error Unauthorized();
    error EscrowExpired();
    error CannotRefund();
    error InvalidTradeId();
    
    // Validation functions
    function validateTradeId(bytes32 tradeId) internal pure {
        if (tradeId == bytes32(0)) {
            revert InvalidTradeId();
        }
    }
    
    function validateAddresses(address buyer, address seller) internal pure {
        if (buyer == address(0) || seller == address(0)) {
            revert InvalidAddresses();
        }
    }
    
    function validateAmount(uint256 amount) internal pure {
        if (amount == 0) {
            revert InvalidAmount();
        }
    }
    
    function validateToken(address token, address expectedToken) internal pure {
        if (token != expectedToken) {
            revert UnsupportedToken();
        }
    }
    
    function isExpired(uint256 expiresAt) internal view returns (bool) {
        return block.timestamp > expiresAt;
    }
    
    function canRefund(uint256 expiresAt, EscrowStatus status) internal view returns (bool) {
        return isExpired(expiresAt) || status == EscrowStatus.Disputed;
    }
}
