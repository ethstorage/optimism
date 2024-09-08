// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IBigStepper, IPreimageOracle } from "src/dispute/interfaces/IBigStepper.sol";
import { PreimageOracle, PreimageKeyLib } from "src/cannon/PreimageOracle.sol";
import "src/dispute/lib/Types.sol";

/// @title AlphabetVM
/// @dev A mock VM for the purpose of testing the dispute game infrastructure. Note that this only works
///      for games with an execution trace subgame max depth of 3 (8 instructions per subgame).
contract AlphabetVM is IBigStepper {
    Claim internal immutable ABSOLUTE_PRESTATE;
    IPreimageOracle public oracle;
    uint256 internal immutable TRACE_DEPTH; // MaxGameDepth - SplitDepth

    constructor(Claim _absolutePrestate, PreimageOracle _oracle, uint256 _traceDepth) {
        ABSOLUTE_PRESTATE = _absolutePrestate;
        oracle = _oracle;
        // Add TRACE_DEPTH to get the starting trace index offset with `startingL2BlockNumber << TRACE_DEPTH`.
        TRACE_DEPTH = _traceDepth;
    }

    /// @inheritdoc IBigStepper
    function step(
        bytes calldata _stateData,
        bytes calldata,
        bytes32 _localContext
    )
        external
        view
        returns (bytes32 postState_)
    {
        uint256 traceIndex;
        uint256 claim;
        if ((keccak256(_stateData) << 8) == (Claim.unwrap(ABSOLUTE_PRESTATE) << 8)) {
            // If the state data is empty, then the absolute prestate is the claim.
            (bytes32 dat,) = oracle.readPreimage(
                PreimageKeyLib.localizeIdent(LocalPreimageKey.DISPUTED_L2_BLOCK_NUMBER, _localContext), 0
            );
            uint256 startingL2BlockNumber = ((uint256(dat) >> 128) & 0xFFFFFFFF) - 1;
            traceIndex = startingL2BlockNumber << TRACE_DEPTH;
            (uint256 absolutePrestateClaim) = abi.decode(_stateData, (uint256));
            claim = absolutePrestateClaim + traceIndex;
            // In actor test trace is byte1 type, so here the claim is truncated to uint8.
            claim = uint256(uint8(claim));
        } else {
            // Otherwise, decode the state data.
            (traceIndex, claim) = abi.decode(_stateData, (uint256, uint256));
            traceIndex++;
            claim++;
            claim = uint256(uint8(claim));
        }

        // STF: n -> n + 1
        postState_ = keccak256(abi.encode(traceIndex, claim));
        assembly {
            postState_ := or(and(postState_, not(shl(248, 0xFF))), shl(248, 1))
        }
    }
}
