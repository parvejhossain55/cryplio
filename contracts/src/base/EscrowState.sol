// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import "../lib/EscrowTypes.sol";

/**
 * @title EscrowState
 * @dev Base contract containing all state variables and storage
 * @notice This contract is inherited by EscrowAuth and EscrowViews
 * @custom:security authorizedSigner mapping controls admin access
 */
abstract contract EscrowState {
    using EscrowTypes for EscrowTypes.Escrow;

    // State Variables
    /// @notice Mapping of tradeId to Escrow struct
    /// @dev SECURITY: tradeId must be unique - enforced in createEscrow
    mapping(bytes32 => EscrowTypes.Escrow) escrows;
    

    /// @notice Mapping of addresses authorized to call admin functions
    /// @dev SECURITY: Only owner can modify via addAuthorizedSigner/removeAuthorizedSigner
    mapping(address => bool) authorizedSigner;


    /// @notice Mapping of supported ERC20 tokens
    /// @dev SECURITY: Only tokens in this mapping can be used for escrow
    mapping(address => bool) public supportedTokens;


    /// @notice Address receiving platform fees (treasury wallet)
    /// @dev FINANCIAL: Must be set before fees can be collected
    address public feeRecipient;

    
    /// @notice Maximum fee cap to protect users (5% max)
    /// @dev SECURITY: Prevents exorbitant fees even if constants are modified
    uint256 public constant MAX_FEE_BPS = 500; // 5% maximum fee
    
    /// @notice Platform fee in basis points
    /// @dev MATH: 75 BPS = 0.75% = (amount * 75) / 10000
    ///      Example: 100 USDT * 75 / 10000 = 0.75 USDT fee
    ///      SECURITY: Constant cannot be changed after deployment
    ///      SAFETY: Must be <= MAX_FEE_BPS (enforced at compile time)
    uint256 public constant FEE_BPS = 75;
    
    /// @notice Refund fee in basis points (lower than release fee)
    /// @dev MATH: 25 BPS = 0.25% = (amount * 25) / 10000
    ///      Example: 100 USDT * 25 / 10000 = 0.25 USDT fee
    ///      RATIONALE: Lower fee on refunds to be fair to users while preventing abuse
    ///      SECURITY: Constant cannot be changed after deployment
    ///      SAFETY: Must be <= MAX_FEE_BPS (enforced at compile time)
    uint256 public constant REFUND_FEE_BPS = 25;
    

    /// @notice Contract constructor - initializes supported tokens
    /// @dev SECURITY: Requires at least one token to be supported
    /// @param _initialTokens Array of ERC20 token addresses to support initially
    constructor(address[] memory _initialTokens) {
        // Validate fee caps at deployment
        require(FEE_BPS <= MAX_FEE_BPS, "Fee exceeds maximum cap");
        require(REFUND_FEE_BPS <= MAX_FEE_BPS, "Refund fee exceeds maximum cap");
        
        uint256 len = _initialTokens.length;
        require(len > 0, "Initial tokens required");
        for (uint256 i = 0; i < len; ) {
            address token = _initialTokens[i];
            
            require(token != address(0), "Invalid token address");
            if (!supportedTokens[token]) {
                supportedTokens[token] = true;
            }
            unchecked { ++i; }
        }
    }

    /// @notice Add an address to authorized signers list
    /// @dev SECURITY: Internal function - only callable by owner via addAuthorizedSigner
    /// @param signer Address to authorize
    function _addAuthorizedSigner(address signer) internal {
        authorizedSigner[signer] = true;
    }

    /// @notice Remove an address from authorized signers list
    /// @dev SECURITY: Internal function - only callable by owner via removeAuthorizedSigner
    /// @param signer Address to remove authorization
    function _removeAuthorizedSigner(address signer) internal {
        authorizedSigner[signer] = false;
    }

    /// @notice Check if an address is authorized
    /// @dev SECURITY: Used by onlyAuthorized modifier
    /// @param signer Address to check
    /// @return bool True if authorized
    function _isAuthorized(address signer) internal view returns (bool) {
        return authorizedSigner[signer];
    }

    // Internal state management functions
    /// @notice Create a new escrow record in storage
    /// @dev Sets expiry time based on provided expiryTime parameter
    /// @param tradeId Unique trade identifier
    /// @param buyer Address receiving tokens
    /// @param seller Address providing tokens
    /// @param token ERC20 token address
    /// @param amount Token amount to lock
    function _createEscrow(
        bytes32 tradeId,
        address buyer,
        address seller,
        address token,
        uint256 amount,
        uint256 expiryTime
    ) internal {
        escrows[tradeId] = EscrowTypes.Escrow({
            tradeId: tradeId,
            buyer: buyer,
            seller: seller,
            token: token,
            amount: amount,
            expiresAt: uint48(block.timestamp + expiryTime),
            status: EscrowTypes.EscrowStatus.Locked
        });
    }

    /// @notice Update escrow status
    /// @dev SECURITY: Only callable internally after validation
    /// @param tradeId Trade identifier
    /// @param newStatus New status to set
    function _updateEscrowStatus(bytes32 tradeId, EscrowTypes.EscrowStatus newStatus) internal {
        escrows[tradeId].status = newStatus;
    }

    /// @notice Get escrow storage reference
    /// @dev SECURITY: Reverts if escrow not found
    /// @param tradeId Trade identifier
    /// @return Escrow storage reference
    function _getEscrow(bytes32 tradeId) internal view returns (EscrowTypes.Escrow storage) {
        EscrowTypes.Escrow storage escrow = escrows[tradeId];
        if (escrow.buyer == address(0)) {
            revert EscrowTypes.EscrowNotFound();
        }
        return escrow;
    }

    /// @notice Check if escrow exists
    /// @dev SECURITY: Uses buyer address as existence check (buyer != address(0))
    /// @param tradeId Trade identifier
    /// @return bool True if escrow exists
    function _escrowExists(bytes32 tradeId) internal view returns (bool) {
        return escrows[tradeId].buyer != address(0);
    }
}