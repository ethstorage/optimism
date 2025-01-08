// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC20Upgradeable } from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { Constants } from "src/libraries/Constants.sol";

/// @title SoulGasToken
/// @notice The SoulGasToken is a soul-bounded ERC20 contract which can be used to pay gas on L2.
///         It has 2 modes:
///             1. when IS_BACKED_BY_NATIVE(or in other words: SoulQKC mode), the token can be minted by
///                anyone depositing native token into the contract.
///             2. when !IS_BACKED_BY_NATIVE(or in other words: SoulETH mode), the token can only be
///                minted by whitelist minters specified by contract owner.
contract SoulGasToken is ERC20Upgradeable, OwnableUpgradeable {
    /// @custom:storage-location erc7201:openzeppelin.storage.SoulGasToken
    struct SoulGasTokenStorage {
        // minters are whitelist EOAs, only used when !IS_BACKED_BY_NATIVE
        mapping(address => bool) minters;
        // burners are whitelist EOAs to burn/withdraw SoulGasToken
        mapping(address => bool) burners;
        // allowSgtValue are whitelist contracts to consume sgt as msg.value
        // when IS_BACKED_BY_NATIVE
        mapping(address => bool) allowSgtValue;
    }

    /// @notice Emitted when sgt as msg.value is enabled for a contract.
    /// @param from     Address of the contract for which sgt as msg.value is enabled.
    event AllowSgtValue(address indexed from);
    /// @notice Emitted when sgt as msg.value is disabled for a contract.
    /// @param from     Address of the contract for which sgt as msg.value is disabled.
    event DisallowSgtValue(address indexed from);

    // keccak256(abi.encode(uint256(keccak256("openzeppelin.storage.SoulGasToken")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant _SOULGASTOKEN_STORAGE_LOCATION =
        0x135c38e215d95c59dcdd8fe622dccc30d04cacb8c88c332e4e7441bac172dd00;

    bool internal immutable IS_BACKED_BY_NATIVE;

    function _getSoulGasTokenStorage() private pure returns (SoulGasTokenStorage storage $) {
        assembly {
            $.slot := _SOULGASTOKEN_STORAGE_LOCATION
        }
    }

    constructor(bool _isBackedByNative) {
        IS_BACKED_BY_NATIVE = _isBackedByNative;
        initialize("", "", msg.sender);
    }

    /// @notice Initializer.
    function initialize(string memory _name, string memory _symbol, address _owner) public initializer {
        __Ownable_init();
        transferOwnership(_owner);

        // initialize the inherited ERC20Upgradeable
        __ERC20_init(_name, _symbol);
    }

    /// @notice deposit can be called by anyone to deposit native token for SoulGasToken when
    /// IS_BACKED_BY_NATIVE.
    function deposit() external payable {
        require(IS_BACKED_BY_NATIVE, "SGT: deposit should only be called when IS_BACKED_BY_NATIVE");

        _mint(_msgSender(), msg.value);
    }

    /// @notice batchDepositFor can be called by anyone to deposit native token for SoulGasToken in batch when
    /// IS_BACKED_BY_NATIVE.
    function batchDepositFor(address[] calldata _accounts, uint256[] calldata _values) external payable {
        require(_accounts.length == _values.length, "SGT: invalid arguments");

        require(IS_BACKED_BY_NATIVE, "SGT: batchDepositFor should only be called when IS_BACKED_BY_NATIVE");

        uint256 totalValue = 0;
        for (uint256 i = 0; i < _accounts.length; i++) {
            _mint(_accounts[i], _values[i]);
            totalValue += _values[i];
        }
        require(msg.value == totalValue, "SGT: unexpected msg.value");
    }

    /// @notice batchDepositForAll is similar to batchDepositFor, but the value is the same for all accounts.
    function batchDepositForAll(address[] calldata _accounts, uint256 _value) external payable {
        require(IS_BACKED_BY_NATIVE, "SGT: batchDepositForAll should only be called when IS_BACKED_BY_NATIVE");

        for (uint256 i = 0; i < _accounts.length; i++) {
            _mint(_accounts[i], _value);
        }
        require(msg.value == _value * _accounts.length, "SGT: unexpected msg.value");
    }

    /// @notice withdrawFrom is called by the burner to burn SoulGasToken and return the native token when
    /// IS_BACKED_BY_NATIVE.
    function withdrawFrom(address _account, uint256 _value) external {
        require(IS_BACKED_BY_NATIVE, "SGT: withdrawFrom should only be called when IS_BACKED_BY_NATIVE");

        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require($.burners[_msgSender()], "SGT: not the burner");

        _burn(_account, _value);
        payable(_msgSender()).transfer(_value);
    }

    /// @notice batchWithdrawFrom is the batch version of withdrawFrom.
    function batchWithdrawFrom(address[] calldata _accounts, uint256[] calldata _values) external {
        require(_accounts.length == _values.length, "SGT: invalid arguments");

        require(IS_BACKED_BY_NATIVE, "SGT: batchWithdrawFrom should only be called when IS_BACKED_BY_NATIVE");

        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require($.burners[_msgSender()], "SGT: not the burner");

        uint256 totalValue = 0;
        for (uint256 i = 0; i < _accounts.length; i++) {
            _burn(_accounts[i], _values[i]);
            totalValue += _values[i];
        }

        payable(_msgSender()).transfer(totalValue);
    }

    /// @notice batchMint is called:
    ///                        1. by EOA minters to mint SoulGasToken in batch when !IS_BACKED_BY_NATIVE.
    ///                        2. by DEPOSITOR_ACCOUNT to refund SoulGasToken
    function batchMint(address[] calldata _accounts, uint256[] calldata _values) external {
        require(_accounts.length == _values.length, "SGT: invalid arguments");

        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require(_msgSender() == Constants.DEPOSITOR_ACCOUNT || $.minters[_msgSender()], "SGT: not a minter");

        for (uint256 i = 0; i < _accounts.length; i++) {
            _mint(_accounts[i], _values[i]);
        }
    }

    /// @notice addMinters is called by the owner to add minters when !IS_BACKED_BY_NATIVE.
    function addMinters(address[] calldata _minters) external onlyOwner {
        require(!IS_BACKED_BY_NATIVE, "SGT: addMinters should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < _minters.length; i++) {
            $.minters[_minters[i]] = true;
        }
    }

    /// @notice delMinters is called by the owner to delete minters when !IS_BACKED_BY_NATIVE.
    function delMinters(address[] calldata _minters) external onlyOwner {
        require(!IS_BACKED_BY_NATIVE, "SGT: delMinters should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < _minters.length; i++) {
            delete $.minters[_minters[i]];
        }
    }

    /// @notice addBurners is called by the owner to add burners.
    function addBurners(address[] calldata _burners) external onlyOwner {
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < _burners.length; i++) {
            $.burners[_burners[i]] = true;
        }
    }

    /// @notice delBurners is called by the owner to delete burners.
    function delBurners(address[] calldata _burners) external onlyOwner {
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < _burners.length; i++) {
            delete $.burners[_burners[i]];
        }
    }

    /// @notice allowSgtValue is called by the owner to enable whitelist contracts to consume sgt as msg.value
    function allowSgtValue(address[] calldata _contracts) external onlyOwner {
        require(IS_BACKED_BY_NATIVE, "SGT: allowSgtValue should only be called when IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < _contracts.length; i++) {
            $.allowSgtValue[_contracts[i]] = true;
            emit AllowSgtValue(_contracts[i]);
        }
    }

    /// @notice allowSgtValue is called by the owner to disable whitelist contracts to consume sgt as msg.value
    function disallowSgtValue(address[] calldata _contracts) external onlyOwner {
        require(IS_BACKED_BY_NATIVE, "SGT: disallowSgtValue should only be called when IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < _contracts.length; i++) {
            $.allowSgtValue[_contracts[i]] = false;
            emit DisallowSgtValue(_contracts[i]);
        }
    }

    /// @notice chargeFromOrigin is called when IS_BACKED_BY_NATIVE to charge for native balance
    ///         from tx.origin if caller is whitelisted.
    function chargeFromOrigin(uint256 _amount) external returns (uint256 amountCharged_) {
        require(IS_BACKED_BY_NATIVE, "SGT: chargeFromOrigin should only be called when IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require($.allowSgtValue[_msgSender()], "SGT: caller is not whitelisted");
        uint256 balance = balanceOf(tx.origin);
        if (balance == 0) {
            amountCharged_ = 0;
            return amountCharged_;
        }
        if (balance >= _amount) {
            amountCharged_ = _amount;
        } else {
            amountCharged_ = balance;
        }
        _burn(tx.origin, amountCharged_);
        payable(_msgSender()).transfer(amountCharged_);
    }

    /// @notice burnFrom is called when !IS_BACKED_BY_NATIVE:
    ///                             1. by the burner to burn SoulGasToken.
    ///                             2. by DEPOSITOR_ACCOUNT to burn SoulGasToken.
    function burnFrom(address _account, uint256 _value) external {
        require(!IS_BACKED_BY_NATIVE, "SGT: burnFrom should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require(_msgSender() == Constants.DEPOSITOR_ACCOUNT || $.burners[_msgSender()], "SGT: not the burner");
        _burn(_account, _value);
    }

    /// @notice batchBurnFrom is the batch version of burnFrom.
    function batchBurnFrom(address[] calldata _accounts, uint256[] calldata _values) external {
        require(_accounts.length == _values.length, "SGT: invalid arguments");
        require(!IS_BACKED_BY_NATIVE, "SGT: batchBurnFrom should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require(_msgSender() == Constants.DEPOSITOR_ACCOUNT || $.burners[_msgSender()], "SGT: not the burner");

        for (uint256 i = 0; i < _accounts.length; i++) {
            _burn(_accounts[i], _values[i]);
        }
    }

    /// @notice transferFrom is disabled for SoulGasToken.
    function transfer(address, uint256) public virtual override returns (bool) {
        revert("SGT: transfer is disabled for SoulGasToken");
    }

    /// @notice transferFrom is disabled for SoulGasToken.
    function transferFrom(address, address, uint256) public virtual override returns (bool) {
        revert("SGT: transferFrom is disabled for SoulGasToken");
    }

    /// @notice approve is disabled for SoulGasToken.
    function approve(address, uint256) public virtual override returns (bool) {
        revert("SGT: approve is disabled for SoulGasToken");
    }

    /// @notice Returns whether SoulGasToken is backed by native token.
    function isBackedByNative() external view returns (bool) {
        return IS_BACKED_BY_NATIVE;
    }
}
