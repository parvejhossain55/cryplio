// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../libraries/EscrowTypes.sol";

/**
 * @title EscrowState
 * @dev Base contract containing all state variables and storage for escrow functionality
 */
abstract contract EscrowState {
    using EscrowTypes for EscrowTypes.Escrow;
    
    // State variables
    mapping(bytes32 => EscrowTypes.Escrow) public escrows;
    mapping(address => bool) public authorizedContracts;
    mapping(address => bool) public authorizedSigners;
    
    // Token interface
    address public usdtToken;
    
    // Constructor
    constructor(address _usdtToken) {
        require(_usdtToken != address(0), "Invalid USDT token address");
        usdtToken = _usdtToken;
    }
    
    // Internal state management functions
    function _createEscrow(
        bytes32 tradeId,
        address buyer,
        address seller,
        address token,
        uint256 amount
    ) internal {
        escrows[tradeId] = EscrowTypes.Escrow({
            tradeId: tradeId,
            buyer: buyer,
            seller: seller,
            token: token,
            amount: amount,
            createdAt: block.timestamp,
            expiresAt: block.timestamp + EscrowTypes.ESCROW_EXPIRY_TIME,
            status: EscrowTypes.EscrowStatus.Locked,
            disputeRaiser: address(0),
            disputeReason: ""
        });
    }
    
    function _updateEscrowStatus(bytes32 tradeId, EscrowTypes.EscrowStatus newStatus) internal {
        escrows[tradeId].status = newStatus;
    }
    
    function _setDisputeInfo(bytes32 tradeId, address raiser, string memory reason) internal {
        escrows[tradeId].disputeRaiser = raiser;
        escrows[tradeId].disputeReason = reason;
    }
    
    function _getEscrow(bytes32 tradeId) internal view returns (EscrowTypes.Escrow storage) {
        if (escrows[tradeId].buyer == address(0)) {
            revert EscrowTypes.EscrowNotFound();
        }
        return escrows[tradeId];
    }
    
    function _escrowExists(bytes32 tradeId) internal view returns (bool) {
        return escrows[tradeId].buyer != address(0);
    }
    
    // Authorization state management
    function _addAuthorizedContract(address contractAddr) internal {
        authorizedContracts[contractAddr] = true;
    }
    
    function _removeAuthorizedContract(address contractAddr) internal {
        authorizedContracts[contractAddr] = false;
    }
    
    function _addAuthorizedSigner(address signer) internal {
        authorizedSigners[signer] = true;
    }
    
    function _removeAuthorizedSigner(address signer) internal {
        authorizedSigners[signer] = false;
    }
    
    function _isAuthorized(address caller) internal view returns (bool) {
        return authorizedContracts[caller] || authorizedSigners[caller];
    }
}
