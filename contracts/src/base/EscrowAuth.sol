// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "./EscrowState.sol";
import "../lib/EscrowTypes.sol";
import "../interface/ICrpylioEscrow.sol";

/**
 * @title EscrowAuth
 * @dev Base contract containing authorization modifiers and management functions
 * @notice Handles all access control for the escrow system
 * @custom:security Owner controls authorized signers and fee recipient
 */
abstract contract EscrowAuth is ICryplioEscrow, EscrowState, Ownable, Pausable {
      // Modifiers
    /// @notice Restricts access to authorized admin signers only
    /// @dev Used for refund and force release functions
    modifier onlyAuthorized() {
        if (!_isAuthorized(msg.sender)) {
            revert EscrowTypes.NotAuthorizedSigner();
        }
        _;
    }
    
    /// @notice Validates that escrow exists before proceeding
    /// @param tradeId Trade identifier to check
    modifier validEscrow(bytes32 tradeId) {
        if (!_escrowExists(tradeId)) {
            revert EscrowTypes.EscrowNotFound();
        }
        _;
    }
    
    /// @notice Restricts access to escrow seller only
    /// @dev SECURITY: Prevents unauthorized release of funds
    /// @param tradeId Trade identifier
    modifier onlySeller(bytes32 tradeId) {
        EscrowTypes.Escrow storage escrow = _getEscrow(tradeId);
        if (escrow.seller != msg.sender) {
            revert EscrowTypes.Unauthorized();
        }
        _;
    }

    /// @param _initialTokens Array of supported ERC20 token addresses
    constructor(address[] memory _initialTokens) EscrowState(_initialTokens) Ownable(msg.sender) {}

    // Authorization management functions (only owner)
    /// @notice Add an address to authorized signers list
    /// @dev ACCESS: Only owner can call. Grants admin privileges for refund/force release.
    /// @param signer Address to authorize as admin
    function addAuthorizedSigner(address signer) external virtual override onlyOwner {
        require(signer != address(0), "Invalid signer");
        _addAuthorizedSigner(signer);
        emit EscrowTypes.AuthorizedSignerAdded(signer, msg.sender);
    }
    
    /// @notice Remove an address from authorized signers list
    /// @dev ACCESS: Only owner can call. Revokes admin privileges.
    /// @param signer Address to remove authorization from
    function removeAuthorizedSigner(address signer) external virtual override onlyOwner {
        _removeAuthorizedSigner(signer);
        emit EscrowTypes.AuthorizedSignerRemoved(signer, msg.sender);
    }

    /// @notice Set the fee recipient address (treasury wallet)
    /// @dev FINANCIAL: All platform fees will be sent to this address
    ///      ACCESS: Only owner can call
    ///      SECURITY: Cannot be set to zero address
    /// @param _feeRecipient New fee recipient address
    function setFeeRecipient(address _feeRecipient) external virtual override onlyOwner {
        require(_feeRecipient != address(0), "Invalid fee recipient");
        emit EscrowTypes.FeeRecipientUpdated(feeRecipient, _feeRecipient);
        feeRecipient = _feeRecipient;
    }

    // Token management functions
    /// @notice Add a token to the supported tokens list
    /// @dev ACCESS: Only owner can call. 
    ///      POLICY: Only verified stablecoins (e.g. USDT, USDC, DAI) should be added.
    /// @param token Token address to add
    function addSupportedToken(address token) external virtual override onlyOwner {
        require(token != address(0), "Invalid token address");
        require(!supportedTokens[token], "Token already supported");
        supportedTokens[token] = true;
        emit EscrowTypes.TokenSupportedAdded(token);
    }

    /// @notice Remove a token from the supported tokens list
    /// @dev ACCESS: Only owner can call
    /// @param token Token address to remove
    function removeSupportedToken(address token) external virtual override onlyOwner {
        require(supportedTokens[token], "Token not supported");
        supportedTokens[token] = false;
        emit EscrowTypes.TokenSupportedRemoved(token);
    }

    /// @notice Pause the contract (emergency)
    /// @dev ACCESS: Only owner can call
    function pause() external onlyOwner {
        _pause();
    }

    /// @notice Unpause the contract
    /// @dev ACCESS: Only owner can call
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @notice Check if an address is authorized as admin
    /// @dev Used by frontend to check admin status
    /// @param caller Address to check
    /// @return bool True if authorized, false otherwise
    function isAuthorized(address caller) external view virtual override returns (bool) {
        return _isAuthorized(caller);
    }
}