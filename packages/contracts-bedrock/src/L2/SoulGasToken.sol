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
        // _minters are whitelist EOAs, only used when !IS_BACKED_BY_NATIVE
        mapping(address => bool) _minters;
        // _burners are whitelist EOAs to burn/withdraw SoulGasToken
        mapping(address => bool) _burners;
        // _allowSgtValue are whitelist contracts to consume sgt as msg.value
        // when IS_BACKED_BY_NATIVE
        mapping(address => bool) _allowSgtValue;
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

    constructor(bool isBackedByNative_) {
        IS_BACKED_BY_NATIVE = isBackedByNative_;
        initialize("", "", msg.sender);
    }

    /// @notice Initializer.
    function initialize(string memory name_, string memory symbol_, address owner_) public initializer {
        __Ownable_init();
        transferOwnership(owner_);

        // initialize the inherited ERC20Upgradeable
        __ERC20_init(name_, symbol_);
    }

    /// @notice deposit can be called by anyone to deposit native token for SoulGasToken when
    /// IS_BACKED_BY_NATIVE.
    function deposit() external payable {
        require(IS_BACKED_BY_NATIVE, "deposit should only be called when IS_BACKED_BY_NATIVE");

        _mint(_msgSender(), msg.value);
    }

    /// @notice batchDepositFor can be called by anyone to deposit native token for SoulGasToken in batch when
    /// IS_BACKED_BY_NATIVE.
    function batchDepositFor(address[] calldata accounts, uint256[] calldata values) external payable {
        require(accounts.length == values.length, "invalid arguments");

        require(IS_BACKED_BY_NATIVE, "batchDepositFor should only be called when IS_BACKED_BY_NATIVE");

        uint256 totalValue = 0;
        for (uint256 i = 0; i < accounts.length; i++) {
            _mint(accounts[i], values[i]);
            totalValue += values[i];
        }
        require(msg.value == totalValue, "unexpected msg.value");
    }

    /// @notice batchDepositForAll is similar to batchDepositFor, but the value is the same for all accounts.
    function batchDepositForAll(address[] calldata _accounts, uint256 _value) external payable {
        require(IS_BACKED_BY_NATIVE, "batchDepositForAll should only be called when IS_BACKED_BY_NATIVE");

        for (uint256 i = 0; i < _accounts.length; i++) {
            _mint(_accounts[i], _value);
        }
        require(msg.value == _value * _accounts.length, "unexpected msg.value");
    }

    /// @notice withdrawFrom is called by the burner to burn SoulGasToken and return the native token when
    /// IS_BACKED_BY_NATIVE.
    function withdrawFrom(address account, uint256 value) external {
        require(IS_BACKED_BY_NATIVE, "withdrawFrom should only be called when IS_BACKED_BY_NATIVE");

        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require($._burners[_msgSender()], "not the burner");

        _burn(account, value);
        payable(_msgSender()).transfer(value);
    }

    /// @notice batchWithdrawFrom is the batch version of withdrawFrom.
    function batchWithdrawFrom(address[] calldata accounts, uint256[] calldata values) external {
        require(accounts.length == values.length, "invalid arguments");

        require(IS_BACKED_BY_NATIVE, "batchWithdrawFrom should only be called when IS_BACKED_BY_NATIVE");

        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require($._burners[_msgSender()], "not the burner");

        uint256 totalValue = 0;
        for (uint256 i = 0; i < accounts.length; i++) {
            _burn(accounts[i], values[i]);
            totalValue += values[i];
        }

        payable(_msgSender()).transfer(totalValue);
    }

    /// @notice batchMint is called:
    ///                        1. by EOA minters to mint SoulGasToken in batch when !IS_BACKED_BY_NATIVE.
    ///                        2. by DEPOSITOR_ACCOUNT to refund SoulGasToken
    function batchMint(address[] calldata accounts, uint256[] calldata values) external {
        require(accounts.length == values.length, "invalid arguments");

        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require(_msgSender() == Constants.DEPOSITOR_ACCOUNT || $._minters[_msgSender()], "not a minter");

        for (uint256 i = 0; i < accounts.length; i++) {
            _mint(accounts[i], values[i]);
        }
    }

    /// @notice addMinters is called by the owner to add minters when !IS_BACKED_BY_NATIVE.
    function addMinters(address[] calldata minters_) external onlyOwner {
        require(!IS_BACKED_BY_NATIVE, "addMinters should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < minters_.length; i++) {
            $._minters[minters_[i]] = true;
        }
    }

    /// @notice delMinters is called by the owner to delete minters when !IS_BACKED_BY_NATIVE.
    function delMinters(address[] calldata minters_) external onlyOwner {
        require(!IS_BACKED_BY_NATIVE, "delMinters should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < minters_.length; i++) {
            delete $._minters[minters_[i]];
        }
    }

    /// @notice addBurners is called by the owner to add burners.
    function addBurners(address[] calldata burners_) external onlyOwner {
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < burners_.length; i++) {
            $._burners[burners_[i]] = true;
        }
    }

    /// @notice delBurners is called by the owner to delete burners.
    function delBurners(address[] calldata burners_) external onlyOwner {
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < burners_.length; i++) {
            delete $._burners[burners_[i]];
        }
    }

    /// @notice allowSgtValue is called by the owner to enable whitelist contracts to consume sgt as msg.value
    function allowSgtValue(address[] calldata contracts) external onlyOwner {
        require(IS_BACKED_BY_NATIVE, "allowSgtValue should only be called when IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < contracts.length; i++) {
            $._allowSgtValue[contracts[i]] = true;
            emit AllowSgtValue(contracts[i]);
        }
    }

    /// @notice allowSgtValue is called by the owner to disable whitelist contracts to consume sgt as msg.value
    function disallowSgtValue(address[] calldata contracts) external onlyOwner {
        require(IS_BACKED_BY_NATIVE, "disallowSgtValue should only be called when IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        uint256 i;
        for (i = 0; i < contracts.length; i++) {
            $._allowSgtValue[contracts[i]] = false;
            emit DisallowSgtValue(contracts[i]);
        }
    }

    /// @notice chargeFromOrigin is called when IS_BACKED_BY_NATIVE to charge for native balance
    ///         from tx.origin if caller is whitelisted.
    function chargeFromOrigin(uint256 amount) external returns (uint256 amountCharged) {
        require(IS_BACKED_BY_NATIVE, "chargeFromOrigin should only be called when IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require($._allowSgtValue[_msgSender()], "caller is not whitelisted");
        uint256 balance = balanceOf(tx.origin);
        if (balance == 0) {
            amountCharged = 0;
            return amountCharged;
        }
        if (balance >= amount) {
            amountCharged = amount;
        } else {
            amountCharged = balance;
        }
        _burn(tx.origin, amountCharged);
        payable(_msgSender()).transfer(amountCharged);
    }

    /// @notice burnFrom is called when !IS_BACKED_BY_NATIVE:
    ///                             1. by the burner to burn SoulGasToken.
    ///                             2. by DEPOSITOR_ACCOUNT to burn SoulGasToken.
    function burnFrom(address account, uint256 value) external {
        require(!IS_BACKED_BY_NATIVE, "burnFrom should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require(_msgSender() == Constants.DEPOSITOR_ACCOUNT || $._burners[_msgSender()], "not the burner");
        _burn(account, value);
    }

    /// @notice batchBurnFrom is the batch version of burnFrom.
    function batchBurnFrom(address[] calldata accounts, uint256[] calldata values) external {
        require(accounts.length == values.length, "invalid arguments");
        require(!IS_BACKED_BY_NATIVE, "batchBurnFrom should only be called when !IS_BACKED_BY_NATIVE");
        SoulGasTokenStorage storage $ = _getSoulGasTokenStorage();
        require(_msgSender() == Constants.DEPOSITOR_ACCOUNT || $._burners[_msgSender()], "not the burner");

        for (uint256 i = 0; i < accounts.length; i++) {
            _burn(accounts[i], values[i]);
        }
    }

    /// @notice transferFrom is disabled for SoulGasToken.
    function transfer(address, uint256) public virtual override returns (bool) {
        revert("transfer is disabled for SoulGasToken");
    }

    /// @notice transferFrom is disabled for SoulGasToken.
    function transferFrom(address, address, uint256) public virtual override returns (bool) {
        revert("transferFrom is disabled for SoulGasToken");
    }

    /// @notice approve is disabled for SoulGasToken.
    function approve(address, uint256) public virtual override returns (bool) {
        revert("approve is disabled for SoulGasToken");
    }

    /// @notice Returns whether SoulGasToken is backed by native token.
    function isBackedByNative() external view returns (bool) {
        return IS_BACKED_BY_NATIVE;
    }
}
