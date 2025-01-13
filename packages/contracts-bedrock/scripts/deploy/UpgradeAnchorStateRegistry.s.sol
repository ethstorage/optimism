// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

// Forge
import {Script} from "forge-std/Script.sol";

// Scripts
import {BaseDeployIO} from "scripts/deploy/BaseDeployIO.sol";
import {DeployUtils} from "scripts/libraries/DeployUtils.sol";

// Libraries
import {GameType, Hash} from "src/dispute/lib/Types.sol";
// Contracts
import {StorageSetter} from "src/universal/StorageSetter.sol";

// Interfaces
import {IAnchorStateRegistry} from "interfaces/dispute/IAnchorStateRegistry.sol";
import {IDisputeGameFactory} from "interfaces/dispute/IDisputeGameFactory.sol";
import {IAnchorStateRegistry} from "interfaces/dispute/IAnchorStateRegistry.sol";
import {IProxyAdmin} from "interfaces/universal/IProxyAdmin.sol";
import {ISuperchainConfig} from "interfaces/L1/ISuperchainConfig.sol";

/// @title UpgradeAnchorStateRegistryInput
contract UpgradeAnchorStateRegistryInput is BaseDeployIO {
    IDisputeGameFactory _disputeGameFactoryProxy;
    IProxyAdmin internal _opChainProxyAdmin;
    IAnchorStateRegistry internal _anchorStateRegistryProxy;
    ISuperchainConfig internal _superchainConfig;
    bytes internal _startingAnchorRoots;

    function set(bytes4 _sel, address _value) public {
        if (_sel == this.disputeGameFactoryProxy.selector) {
            require(
                _value != address(0),
                "UpgradeAnchorStateRegistryInput: disputeGameFactoryProxy cannot be zero address"
            );
            _disputeGameFactoryProxy = IDisputeGameFactory(_value);
        } else if (_sel == this.opChainProxyAdmin.selector) {
            require(
                _value != address(0),
                "UpgradeAnchorStateRegistryInput: opChainProxyAdmin cannot be zero address"
            );
            _opChainProxyAdmin = IProxyAdmin(_value);
        } else if (_sel == this.anchorStateRegistryProxy.selector) {
            require(
                _value != address(0),
                "UpgradeAnchorStateRegistryInput: anchorStateRegistryProxy cannot be zero address"
            );
            _anchorStateRegistryProxy = IAnchorStateRegistry(_value);
        } else if (_sel == this.superchainConfig.selector) {
            require(
                _value != address(0),
                "UpgradeAnchorStateRegistryInput: superchainConfig cannot be zero address"
            );
            _superchainConfig = ISuperchainConfig(_value);
        } else {
            revert(
                "UpgradeAnchorStateRegistryInput: unknown selector for address"
            );
        }
    }

    function set(bytes4 _sel, bytes memory _value) public {
        if (_sel == this.startingAnchorRoots.selector) {
            require(
                _value.length > 0,
                "UpgradeAnchorStateRegistryInput: startingAnchorRoots cannot be empty bytes"
            );
            _startingAnchorRoots = _value;
        } else {
            revert(
                "UpgradeAnchorStateRegistryInput: unknown selector for bytes"
            );
        }
    }

    function opChainProxyAdmin() public view returns (IProxyAdmin) {
        DeployUtils.assertValidContractAddress(address(_opChainProxyAdmin));
        return _opChainProxyAdmin;
    }

    function disputeGameFactoryProxy()
        public
        view
        returns (IDisputeGameFactory)
    {
        DeployUtils.assertValidContractAddress(
            address(_disputeGameFactoryProxy)
        );
        return _disputeGameFactoryProxy;
    }

    function anchorStateRegistryProxy() public view returns (address) {
        DeployUtils.assertValidContractAddress(
            address(_anchorStateRegistryProxy)
        );
        return address(_anchorStateRegistryProxy);
    }

    function superchainConfig() public view returns (ISuperchainConfig) {
        DeployUtils.assertValidContractAddress(address(_superchainConfig));
        return _superchainConfig;
    }

    function startingAnchorRoots() public view returns (bytes memory) {
        require(
            _startingAnchorRoots.length > 0,
            "UpgradeAnchorStateRegistryInput: startingAnchorRoots not set"
        );
        return _startingAnchorRoots;
    }
}

/// @title UpgradeAnchorStateRegistryOutput
contract UpgradeAnchorStateRegistryOutput is BaseDeployIO {
    IAnchorStateRegistry internal _anchorStateRegistryImpl;
    StorageSetter internal _storageSetter;

    function set(bytes4 _sel, address _value) public {
        if (_sel == this.anchorStateRegistryImpl.selector) {
            require(
                _value != address(0),
                "UpgradeAnchorStateRegistryOutput: anchorStateRegistryImpl cannot be zero address"
            );
            _anchorStateRegistryImpl = IAnchorStateRegistry(_value);
        } else if (_sel == this.storageSetter.selector) {
            require(
                _value != address(0),
                "UpgradeAnchorStateRegistryOutput: storageSetter cannot be zero address"
            );
            _storageSetter = StorageSetter(_value);
        } else {
            revert("UpgradeAnchorStateRegistryOutput: unknown selector");
        }
    }

    function anchorStateRegistryImpl() public view returns (address) {
        DeployUtils.assertValidContractAddress(
            address(_anchorStateRegistryImpl)
        );
        return address(_anchorStateRegistryImpl);
    }

    function storageSetter() public view returns (address) {
        DeployUtils.assertValidContractAddress(address(_storageSetter));
        return address(_storageSetter);
    }
    function checkOutput(UpgradeAnchorStateRegistryInput _input) public view {
        IAnchorStateRegistry.StartingAnchorRoot[]
            memory startingAnchorRoots = abi.decode(
                _input.startingAnchorRoots(),
                (IAnchorStateRegistry.StartingAnchorRoot[])
            );
        (Hash root, uint256 l2BlockNumber) = IAnchorStateRegistry(
            _input.anchorStateRegistryProxy()
        ).anchors(startingAnchorRoots[0].gameType);
        require(
            Hash.unwrap(root) ==
                Hash.unwrap(startingAnchorRoots[0].outputRoot.root),
            "UpgradeAnchorStateRegistryOutput: root mismatch"
        );
        require(
            l2BlockNumber == startingAnchorRoots[0].outputRoot.l2BlockNumber,
            "UpgradeAnchorStateRegistryOutput: l2BlockNumber mismatch"
        );
    }
}

/// @title UpgradeAnchorStateRegistry
contract UpgradeAnchorStateRegistry is Script {
    function run(
        UpgradeAnchorStateRegistryInput _input,
        UpgradeAnchorStateRegistryOutput _output
    ) public {
        upgradeAnchorStateRegistryImpl(_input, _output);
        _output.checkOutput(_input);
    }

    function upgradeAnchorStateRegistryImpl(
        UpgradeAnchorStateRegistryInput _input,
        UpgradeAnchorStateRegistryOutput _output
    ) internal {
        _output.set(
            _output.anchorStateRegistryImpl.selector,
            DeployUtils.create1({
                _name: "AnchorStateRegistry",
                _args: abi.encode(_input.disputeGameFactoryProxy())
            })
        );

        _output.set(
            _output.storageSetter.selector,
            DeployUtils.create1({_name: "StorageSetter", _args: ""})
        );

        bytes memory data;
        data = encodeStorageSetterZeroOutInitializedSlot();
        upgradeAndCall(
            _input.opChainProxyAdmin(),
            _input.anchorStateRegistryProxy(),
            _output.storageSetter(),
            data
        );
        data = encodeAnchorStateRegistryInitializer(_input);
        upgradeAndCall(
            _input.opChainProxyAdmin(),
            _input.anchorStateRegistryProxy(),
            _output.anchorStateRegistryImpl(),
            data
        );
    }

    function encodeStorageSetterZeroOutInitializedSlot()
        internal
        pure
        returns (bytes memory)
    {
        return
            abi.encodeWithSelector(
                bytes4(keccak256("setBytes32(bytes32,bytes32)")),
                0,
                0
            );
    }

    function encodeAnchorStateRegistryInitializer(
        UpgradeAnchorStateRegistryInput _input
    ) internal view virtual returns (bytes memory) {
        // this line fails in the op-deployer tests because it is not passing in any data
        IAnchorStateRegistry.StartingAnchorRoot[]
            memory startingAnchorRoots = abi.decode(
                _input.startingAnchorRoots(),
                (IAnchorStateRegistry.StartingAnchorRoot[])
            );
        return
            abi.encodeWithSelector(
                IAnchorStateRegistry.initialize.selector,
                startingAnchorRoots,
                _input.superchainConfig()
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
}
