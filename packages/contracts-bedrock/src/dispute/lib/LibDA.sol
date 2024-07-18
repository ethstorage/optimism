// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

library LibDA {
    /// @notice The `ClaimData` struct represents the data associated with a Claim.

    /// @notice Represents a leaf DA item that can be verified against its root.
    /// @custom:field daType        Type of DA (either 4844 or calldata).
    /// @custom:field dataHash      Leaf hash (claim hash).
    /// @custom:field proof         Inclusion proof of the item.
    struct DAItem {
        uint256 daType;
        bytes32 dataHash;
        bytes proof;
    }

    uint256 constant DA_TYPE_CALLDATA = 0;
    uint256 constant DA_TYPE_EIP4844 = 1;

    function getClaimsHash(uint256 daType, uint256 nelemebts, bytes memory data) internal view returns (bytes32 root)  {
        if (daType == DA_TYPE_EIP4844) {
            // TODO: may specify which blob?
            // root = blobhash(0);
            // require(root != bytes32(0), "root must not zero");
            // return root;
            revert("EIP4844 not supported");
        }

        require(daType == DA_TYPE_CALLDATA, "unsupported DA type");
        require(nelemebts * 32 == data.length, "data must 32 * n");
        require(nelemebts > 0, "data must not empty");

        while (nelemebts != 1) {
            for (uint256 i = 0 ; i < nelemebts / 2; i++) {
                bytes32 hash;
                uint256 roff = i * 32 * 2;
                uint256 woff = i * 32;
                assembly {
                    hash := keccak256(add(add(data, 0x20), roff), 64)
                    mstore(add(add(data, 0x20), woff), hash)
                }
            }

            // directly copy the last item
            if (nelemebts % 2 == 1) {
                uint256 roff = (nelemebts - 1) * 32;
                uint256 woff = (nelemebts / 2) * 32;
                bytes32 hash;
                assembly {
                    hash := mload(add(add(data, 0x20), roff))
                    mstore(add(add(data, 0x20), woff), hash)
                }
            }

            nelemebts = (nelemebts + 1) / 2;
        }

        assembly {
            root := mload(add(data, 0x20))
        }
    }

    function verifyClaimHash(uint256 daType, bytes32 root, uint256 nelements, uint256 idx, bytes32 claimHash, bytes memory proof) internal pure {
        require(daType == 0, "unsupported DA type");
        bytes32 hash = claimHash;
        uint256 proofOff = 0;
        while (nelements != 1) {
            if (idx != nelements - 1 || nelements % 2 == 0) {
                bytes32 pHash;
                require(proofOff < proof.length, "no enough proof");
                assembly {
                    pHash := mload(add(add(proof, 0x20), proofOff))
                }
                proofOff += 32;
                if (idx % 2 == 0) {
                    hash = keccak256(abi.encode(hash, pHash));
                } else {
                    hash = keccak256(abi.encode(pHash, hash));
                }
            }
            nelements = (nelements + 1) / 2;
            idx = idx / 2;
        }
        require(root == hash, "proof failed");
    }
}

