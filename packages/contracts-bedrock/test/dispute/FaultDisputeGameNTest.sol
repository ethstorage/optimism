// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { GameType, Claim, Duration } from "src/dispute/lib/LibUDT.sol";
import { FaultDisputeGame } from "src/dispute/FaultDisputeGameN.sol";
import { IAnchorStateRegistry } from "src/dispute/interfaces/IAnchorStateRegistry.sol";
import { IDelayedWETH } from "src/dispute/interfaces/IDelayedWETH.sol";
import { IBigStepper } from "src/dispute/interfaces/IBigStepper.sol";

contract FaultDisputeGameTest is FaultDisputeGame {
    constructor(
        GameType _gameType,
        Claim _absolutePrestate,
        uint256 _maxGameDepth,
        uint256 _splitDepth,
        Duration _clockExtension,
        Duration _maxClockDuration,
        IBigStepper _vm,
        IDelayedWETH _weth,
        IAnchorStateRegistry _anchorStateRegistry,
        uint256 _l2ChainId
    )
        FaultDisputeGame(
            _gameType,
            _absolutePrestate,
            _maxGameDepth,
            _splitDepth,
            _clockExtension,
            _maxClockDuration,
            _vm,
            _weth,
            _anchorStateRegistry,
            _l2ChainId
        )
    { }

    // For testing convenience and to minimize changes in the testing code, the submission of the "claims" value is
    // omitted during the attack. In contract testing, the value of "claims" is already known and does not need to be
    // submitted via calldata or EIP-4844 during the attack.
    function attackV2(Claim _disputed, uint256 _parentIndex, Claim _claim, uint64 _attackBranch) public payable {
        moveV2(_disputed, _parentIndex, _claim, _attackBranch);
    }
}
