// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../libraries/EscrowTypes.sol";

/**
 * @title IEscrowCore
 * @dev Interface for core escrow functionality
 */
interface IEscrowCore {
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
    
    // Core functions
    function createEscrow(
        bytes32 tradeId,
        address buyer,
        address seller,
        address token,
        uint256 amount
    ) external returns (bool);
    
    function releaseEscrow(bytes32 tradeId) external returns (bool);
    
    function refundEscrow(bytes32 tradeId) external returns (bool);
    
    function raiseDispute(bytes32 tradeId, string calldata reason) external returns (bool);
    
    // View functions
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
        );
    
    function escrowExists(bytes32 tradeId) external view returns (bool);
    
    function getEscrowStatus(bytes32 tradeId) external view returns (EscrowTypes.EscrowStatus);
    
    function isEscrowExpired(bytes32 tradeId) external view returns (bool);
    
    function canEscrowBeRefunded(bytes32 tradeId) external view returns (bool);
    
    function getTimeRemaining(bytes32 tradeId) external view returns (uint256);
    
    function getDisputeWindowRemaining(bytes32 tradeId) external view returns (uint256);
    
    // State functions
    function usdtToken() external view returns (address);
    
    function isAuthorized(address caller) external view returns (bool);
    
    function isContractAuthorized(address contractAddr) external view returns (bool);
    
    function isSignerAuthorized(address signer) external view returns (bool);
}
