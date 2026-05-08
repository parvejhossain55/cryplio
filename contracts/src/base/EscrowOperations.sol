// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./EscrowAuth.sol";
import "../libraries/EscrowTypes.sol";

/**
 * @title EscrowOperations
 * @dev Base contract containing core escrow operations
 */
abstract contract EscrowOperations is EscrowAuth, ReentrancyGuard {
    
    // Events (re-exported for convenience)
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
    
    // Constructor
    constructor(address _usdtToken) EscrowAuth(_usdtToken) {}
    
    /**
     * @dev Create a new escrow for a trade
     * @param tradeId Unique identifier for the trade
     * @param buyer Address of the buyer
     * @param seller Address of the seller
     * @param token Address of the ERC20 token (USDT)
     * @param amount Amount of tokens to lock
     */
    function createEscrow(
        bytes32 tradeId,
        address buyer,
        address seller,
        address token,
        uint256 amount
    ) external onlyAuthorized nonReentrant returns (bool) {
        // Validate inputs
        EscrowTypes.validateTradeId(tradeId);
        EscrowTypes.validateAddresses(buyer, seller);
        EscrowTypes.validateAmount(amount);
        EscrowTypes.validateToken(token, usdtToken);
        
        // Check if escrow already exists
        if (_escrowExists(tradeId)) {
            revert EscrowTypes.EscrowAlreadyExists();
        }
        
        // Transfer USDT from buyer to this contract
        IERC20 usdt = IERC20(usdtToken);
        bool success = usdt.transferFrom(buyer, address(this), amount);
        if (!success) {
            revert EscrowTypes.TransferFailed();
        }
        
        // Create escrow record
        _createEscrow(tradeId, buyer, seller, token, amount);
        
        emit EscrowCreated(tradeId, buyer, seller, token, amount, block.timestamp);
        return true;
    }
    
    /**
     * @dev Release escrow funds to seller
     * @param tradeId Trade identifier
     */
    function releaseEscrow(bytes32 tradeId) 
        external 
        onlyAuthorized 
        validEscrow(tradeId) 
        validStatus(tradeId, EscrowTypes.EscrowStatus.Locked)
        nonReentrant 
        returns (bool) 
    {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        
        // Check if still within dispute window
        if (block.timestamp > escrow.expiresAt + EscrowTypes.DISPUTE_WINDOW) {
            revert EscrowTypes.EscrowExpired();
        }
        
        // Update status
        _updateEscrowStatus(tradeId, EscrowTypes.EscrowStatus.Released);
        
        // Transfer USDT to seller
        IERC20 usdt = IERC20(usdtToken);
        bool success = usdt.transfer(escrow.seller, escrow.amount);
        if (!success) {
            revert EscrowTypes.TransferFailed();
        }
        
        emit EscrowReleased(tradeId, escrow.seller, escrow.amount, block.timestamp);
        return true;
    }
    
    /**
     * @dev Refund escrow funds to buyer
     * @param tradeId Trade identifier
     */
    function refundEscrow(bytes32 tradeId) 
        external 
        onlyAuthorized 
        validEscrow(tradeId) 
        canRefund(tradeId)
        nonReentrant 
        returns (bool) 
    {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        
        // Update status
        _updateEscrowStatus(tradeId, EscrowTypes.EscrowStatus.Refunded);
        
        // Transfer USDT back to buyer
        IERC20 usdt = IERC20(usdtToken);
        bool success = usdt.transfer(escrow.buyer, escrow.amount);
        if (!success) {
            revert EscrowTypes.TransferFailed();
        }
        
        emit EscrowRefunded(tradeId, escrow.buyer, escrow.amount, block.timestamp);
        return true;
    }
    
    /**
     * @dev Raise a dispute for an escrow
     * @param tradeId Trade identifier
     * @param reason Dispute reason
     */
    function raiseDispute(bytes32 tradeId, string calldata reason) 
        external 
        onlyAuthorized 
        validEscrow(tradeId) 
        validStatus(tradeId, EscrowTypes.EscrowStatus.Locked)
        nonReentrant 
        returns (bool) 
    {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        
        // Check if still within dispute window
        if (block.timestamp > escrow.expiresAt + EscrowTypes.DISPUTE_WINDOW) {
            revert EscrowTypes.EscrowExpired();
        }
        
        // Update escrow with dispute info
        _updateEscrowStatus(tradeId, EscrowTypes.EscrowStatus.Disputed);
        _setDisputeInfo(tradeId, msg.sender, reason);
        
        emit DisputeRaised(tradeId, msg.sender, reason, block.timestamp);
        return true;
    }
}
