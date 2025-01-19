// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract EthStorageContract {
    uint256 internal immutable COST;
    uint256 public kvEntryCount;

    /// @notice Emitted when a BLOB is appended.
    /// @param kvIdx    The index of the KV pair
    /// @param kvSize   The size of the KV pair
    /// @param dataHash The hash of the data
    event PutBlob(uint256 indexed kvIdx, uint256 indexed kvSize, bytes32 indexed dataHash);

    constructor(uint256 _cost) {
        COST = _cost;
    }

    /// @notice Write a large value to KV store.  If the KV pair exists, overrides it.
    ///         Otherwise, will append the KV to the KV array.
    /// @param _key The key of the KV pair
    /// @param _blobIdx The index of the blob
    /// @param _length The length of the blob
    function putBlob(bytes32 _key, uint256 _blobIdx, uint256 _length) public payable virtual {
        bytes32 dataHash = blobhash(_blobIdx);
        require(dataHash != 0, "EthStorageContract: failed to get blob hash");
        require(_key != 0, "EthStorageContract: failed to get blob key");
        require(msg.value >= upfrontPayment(), "DecentralizedKV: not enough batch payment");

        uint256 kvIndex = kvEntryCount;
        kvEntryCount = kvEntryCount + 1;

        emit PutBlob(kvIndex, _length, dataHash);
    }

    /// @notice Evaluate the storage cost of a single put().
    function upfrontPayment() public view virtual returns (uint256) {
        return COST ;
    }
}
