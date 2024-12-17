// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ISoulGasToken {
    function initialize(string memory _name, string memory _symbol, address _owner) external;

    function deposit() external payable;

    function batchDepositFor(address[] calldata _accounts, uint256[] calldata _values) external payable;

    function batchDepositForAll(address[] calldata _accounts, uint256 _value) external payable;

    function withdrawFrom(address _account, uint256 _value) external;

    function batchWithdrawFrom(address[] calldata _accounts, uint256[] calldata _values) external;

    function batchMint(address[] calldata _accounts, uint256[] calldata _values) external;

    function addMinters(address[] calldata _minters) external;

    function delMinters(address[] calldata _minters) external;

    function addBurners(address[] calldata _burners) external;

    function delBurners(address[] calldata _burners) external;

    function allowSgtValue(address[] calldata _contracts) external;

    function disallowSgtValue(address[] calldata _contracts) external;

    function chargeFromOrigin(uint256 _amount) external returns (uint256 amountCharged_);

    function burnFrom(address _account, uint256 _value) external;

    function batchBurnFrom(address[] calldata _accounts, uint256[] calldata _values) external;

    function transfer(address, uint256) external returns (bool);

    function transferFrom(address, address, uint256) external returns (bool);

    function approve(address, uint256) external returns (bool);

    function isBackedByNative() external view returns (bool);

    event AllowSgtValue(address indexed from);
    event DisallowSgtValue(address indexed from);
}
