// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

// Forge
import {Script} from "forge-std/Script.sol";

// Scripts
import {DeployUtils} from "scripts/libraries/DeployUtils.sol";

// Contracts
import {StorageSetter} from "src/universal/StorageSetter.sol";

// Interfaces
import {IProxyAdmin} from "interfaces/universal/IProxyAdmin.sol";
import {IStorageSetter} from "interfaces/universal/IStorageSetter.sol";
import {IProxy} from "interfaces/universal/IProxy.sol";
import {ISoulGasToken} from "interfaces/L2/ISoulGasToken.sol";

/// @title UpgradeSoulGasToken
contract UpgradeSoulGasToken is Script {
    function run(address _storageSetter) public {
        vm.startBroadcast();
        upgradeSoulGasTokenImpl(_storageSetter);
        vm.stopBroadcast();
        checkOutput();
    }

    function upgradeSoulGasTokenImpl(address _storageSetter) internal {
        if (_storageSetter == address(0)) {
            _storageSetter = DeployUtils.create1({
                _name: "StorageSetter",
                _args: DeployUtils.encodeConstructor(
                    abi.encodeCall(IStorageSetter.__constructor__, ())
                )
            });
        }

        IProxyAdmin proxyAdmin = IProxyAdmin(
            0x4200000000000000000000000000000000000018
        );
        address payable soulGasToken = payable(
            0x4200000000000000000000000000000000000800
        );

        address impl = IProxy(soulGasToken).implementation();

        bytes memory data;
        data = encodeStorageSetterZeroOutInitializedSlot();
        upgradeAndCall(proxyAdmin, soulGasToken, _storageSetter, data);
        data = encodeSoulGasTokenInitializer();
        upgradeAndCall(proxyAdmin, soulGasToken, impl, data);
    }

    function encodeStorageSetterZeroOutInitializedSlot()
        internal
        pure
        returns (bytes memory)
    {
        return abi.encodeCall(IStorageSetter.setBytes32, (0, 0));
    }

    function encodeSoulGasTokenInitializer()
        internal
        view
        virtual
        returns (bytes memory)
    {
        return
            abi.encodeCall(
                ISoulGasToken.initialize,
                ("SoulQKC", "SoulQKC", tx.origin)
            );
    }

    /// @notice Makes an external call to the target to initialize the proxy with the specified data.
    /// First performs safety checks to ensure the target, implementation, and proxy admin are valid.
    function upgradeAndCall(
        IProxyAdmin _proxyAdmin,
        address _target,
        address _implementation,
        bytes memory _data
    ) internal {
        DeployUtils.assertValidContractAddress(address(_proxyAdmin));
        DeployUtils.assertValidContractAddress(_target);
        DeployUtils.assertValidContractAddress(_implementation);

        _proxyAdmin.upgradeAndCall(
            payable(address(_target)),
            _implementation,
            _data
        );
    }

    function checkOutput() public view {
        ISoulGasToken soulGasToken = ISoulGasToken(
            0x4200000000000000000000000000000000000800
        );

        require(
            keccak256(abi.encodePacked(soulGasToken.name())) ==
                keccak256("SoulQKC")
        );
        require(
            keccak256(abi.encodePacked(soulGasToken.symbol())) ==
                keccak256("SoulQKC")
        );
    }
}
