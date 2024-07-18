// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { LibDA } from "src/dispute/lib/LibDA.sol";

/// @notice Tests for `LibDA`
contract LibDA_Test is Test {
    function test_calldata_one() public view {
        bytes32 root;
        bytes memory input = "00000000000000000000000000000000";
        root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 1, input);
        assertEq(root, bytes32("00000000000000000000000000000000"));
        input = "10000000000000000000000000000001";
        root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 1, input);
        assertEq(root, bytes32("10000000000000000000000000000001"));
    }

    function test_calldata_two() public view {
        bytes32 root;
        bytes memory input = "0000000000000000000000000000000010000000000000000000000000000001";
        root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 2, input);
        assertEq(root, keccak256(abi.encode(bytes32("00000000000000000000000000000000"), bytes32("10000000000000000000000000000001"))));
    }

    function test_calldata_three() public view {
        bytes32 root;
        bytes memory input = "000000000000000000000000000000001000000000000000000000000000000120000000000000000000000000000002";
        root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 3, input);
        assertEq(root, keccak256(abi.encode(
            keccak256(abi.encode(bytes32("00000000000000000000000000000000"), bytes32("10000000000000000000000000000001"))),
            bytes32("20000000000000000000000000000002")
        )));
    }

    function test_calldata_three1() public view {
        bytes memory claimData1 = abi.encode(1, 1);
        bytes memory claimData2 = abi.encode(2, 2);
        bytes memory claimData3 = abi.encode(3, 3);
        bytes32 claim1 = keccak256(claimData1);
        bytes32 claim2 = keccak256(claimData2);
        bytes32 claim3 = keccak256(claimData3);

        bytes memory input = abi.encodePacked(claim1, claim2, claim3);

        bytes32 root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 3, input);
        assertEq(root, keccak256(abi.encode(
            keccak256(abi.encode(claim1, claim2)),
            claim3
        )));
    }

    function test_calldata_seven() public view {
        bytes32 root;
        bytes memory input = "00000000000000000000000000000000100000000000000000000000000000012000000000000000000000000000000230000000000000000000000000000003400000000000000000000000000000045000000000000000000000000000000560000000000000000000000000000006";
        root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 7, input);
        assertEq(root,
            keccak256(abi.encode(
                keccak256(abi.encode(
                    keccak256(abi.encode(
                        bytes32("00000000000000000000000000000000"),
                        bytes32("10000000000000000000000000000001"))),
                    keccak256(abi.encode(
                        bytes32("20000000000000000000000000000002"),
                        bytes32("30000000000000000000000000000003"))))),
                keccak256(abi.encode(
                    keccak256(abi.encode(
                        bytes32("40000000000000000000000000000004"),
                        bytes32("50000000000000000000000000000005"))),
                    bytes32("60000000000000000000000000000006")))
        )));
    }

    function test_calldata_prove_three() public view {
        bytes32 root;
        bytes memory input = "000000000000000000000000000000001000000000000000000000000000000120000000000000000000000000000002";
        root = LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, 3, input);
        LibDA.verifyClaimHash(LibDA.DA_TYPE_CALLDATA, root, 3, 0, "00000000000000000000000000000000", "1000000000000000000000000000000120000000000000000000000000000002");
        LibDA.verifyClaimHash(LibDA.DA_TYPE_CALLDATA, root, 3, 1, "10000000000000000000000000000001", "0000000000000000000000000000000020000000000000000000000000000002");
        LibDA.verifyClaimHash(LibDA.DA_TYPE_CALLDATA, root, 3, 2, "20000000000000000000000000000002", bytes.concat(keccak256(abi.encode(bytes32("00000000000000000000000000000000"), bytes32("10000000000000000000000000000001")))));
    }
}
