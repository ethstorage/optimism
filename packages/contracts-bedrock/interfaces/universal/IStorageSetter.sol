// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IStorageSetter {
    function setBytes32(bytes32, bytes32) external;

    function __constructor__() external;
}
