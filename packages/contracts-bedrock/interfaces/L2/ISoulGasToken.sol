// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title ISoulGasToken
/// @notice The interface for the SoulGasToken.
interface ISoulGasToken {
    function initialize(string memory _name, string memory _symbol, address _owner) external;

    function __constructor__() external;

    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function owner() external view returns (address);
    function admin() external view returns (address);
}
