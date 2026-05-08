// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/Escrow.sol";
import "../src/libraries/EscrowTypes.sol";

/**
 * @title DeployRefactored
 * @dev Deployment script for the refactored CryplioEscrow contract
 */
contract DeployRefactored is Script {
    CryplioEscrow public escrow;
    
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address usdtToken = vm.envAddress("USDT_TOKEN_ADDRESS");
        
        vm.startBroadcast(deployerPrivateKey);
        
        // Deploy the refactored escrow contract
        escrow = new CryplioEscrow(usdtToken);
        
        // Log deployment info
        console.log("CryplioEscrow deployed at:", address(escrow));
        console.log("USDT Token:", usdtToken);
        console.log("Deployer:", msg.sender);
        
        // Verify contract info
        (
            address usdtTokenAddress,
            uint256 escrowExpiryTime,
            uint256 disputeWindow
        ) = escrow.getContractInfo();
        
        console.log("USDT Token Address:", usdtTokenAddress);
        console.log("Escrow Expiry Time:", escrowExpiryTime);
        console.log("Dispute Window:", disputeWindow);
        
        vm.stopBroadcast();
    }
    
    /**
     * @dev Deploy with custom USDT token address
     */
    function deployCustom(address usdtToken) external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);
        escrow = new CryplioEscrowRefactored(usdtToken);
        vm.stopBroadcast();
    }
    
    /**
     * @dev Deploy for testing (uses mock USDT)
     */
    function deployTest() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);
        
        // Deploy with a mock USDT address for testing
        address mockUSDT = address(0x1234567890123456789012345678901234567890);
        escrow = new CryplioEscrowRefactored(mockUSDT);
        
        vm.stopBroadcast();
    }
}
