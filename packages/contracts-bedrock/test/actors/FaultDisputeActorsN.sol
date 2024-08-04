// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { CommonBase } from "forge-std/Base.sol";

import { FaultDisputeGame } from "src/dispute/FaultDisputeGameN.sol";
import { IFaultDisputeGame } from "src/dispute/interfaces/IFaultDisputeGame.sol";

import "src/dispute/lib/Types.sol";
import "src/dispute/lib/LibDA.sol";

/// @title GameSolver
/// @notice The `GameSolver` contract is a contract that can produce an array of available
///         moves for a given `FaultDisputeGame` contract, from the eyes of an honest
///         actor. The `GameSolver` does not implement functionality for acting on the `Move`s
///         it suggests.
abstract contract GameSolver is CommonBase {
    /// @notice The `FaultDisputeGame` proxy that the `GameSolver` will be solving.
    FaultDisputeGame public immutable GAME;
    /// @notice The split depth of the game
    uint256 internal immutable SPLIT_DEPTH;
    /// @notice The max depth of the game
    uint256 internal immutable MAX_DEPTH;
    /// @notice The maximum L2 block number that the output bisection portion of the position tree
    ///         can handle.
    uint256 internal immutable MAX_L2_BLOCK_NUMBER;
    /// @notice 1<<Bits-1 of N-ary search
    uint256 internal immutable MAX_ATTACK_BRANCH;
    uint256 internal immutable N_BITS;
    uint64 internal immutable NOOP_ATTACK = type(uint64).max;

    /// @notice The L2 outputs that the `GameSolver` will be representing, keyed by L2 block number - 1.
    uint256[] public l2Outputs;
    uint256[] public counterL2Outputs;
    /// @notice The execution trace that the `GameSolver` will be representing.
    bytes public trace;
    bytes public counterTrace;
    /// @notice The raw absolute prestate data.
    bytes public absolutePrestateData;
    /// @notice The offset of previously processed claims in the `GAME` contract's `claimData` array.
    ///         Starts at 0 and increments by 1 for each claim processed.
    uint256 public processedBuf;
    /// @notice Signals whether or not the `GameSolver` agrees with the root claim of the
    ///         `GAME` contract.
    bool public agreeWithRoot;

    /// @notice The `MoveKind` enum represents a kind of interaction with the `FaultDisputeGame` contract.
    enum MoveKind {
        Attack,
        Step,
        AddLocalData
    }

    enum Actor {
        Self,
        Counter
    }

    /// @notice The `Move` struct represents a move in the game, and contains information
    ///         about the kind of move, the sender of the move, and the calldata to be sent
    ///         to the `FaultDisputeGame` contract by a consumer of this contract.
    struct Move {
        MoveKind kind;
        uint256 attackBranch;
        bytes data;
        uint256 value;
    }

    constructor(
        FaultDisputeGame _gameProxy,
        uint256[] memory _l2Outputs,
        uint256[] memory _counterL2Outputs,
        bytes memory _trace,
        bytes memory _counterTrace,
        bytes memory _preStateData
    ) {
        GAME = _gameProxy;
        SPLIT_DEPTH = GAME.splitDepth();
        MAX_DEPTH = GAME.maxGameDepth();
        MAX_L2_BLOCK_NUMBER = 2 ** (MAX_DEPTH - SPLIT_DEPTH);
        N_BITS = GAME.N_BITS();
        MAX_ATTACK_BRANCH = GAME.MAX_ATTACK_BRANCH();

        l2Outputs = _l2Outputs;
        counterL2Outputs = _counterL2Outputs;
        trace = _trace;
        counterTrace = _counterTrace;
        absolutePrestateData = _preStateData;
    }

    /// @notice Returns an array of `Move`s that can be taken from the perspective of an honest
    ///         actor in the `FaultDisputeGame` contract.
    function solveGame() external virtual returns (Move[] memory moves_);
}

/// @title HonestGameSolver
/// @notice The `HonestGameSolver` is an implementation of `GameSolver` which responds accordingly depending
///         on the state of the `FaultDisputeGame` contract in relation to their local opinion of the correct
///         order of output roots and the execution trace between each block `n` -> `n + 1` state transition.
contract HonestGameSolver is GameSolver {
    /// @notice The `Direction` enum represents the direction of a proposed move in the game,
    ///         or a lack thereof.
    enum Direction {
        Defend,
        Attack,
        Noop
    }

    constructor(
        FaultDisputeGame _gameProxy,
        uint256[] memory l2Outputs,
        uint256[] memory counterL2Outputs,
        bytes memory _honestTrace,
        bytes memory _dishonestTrace,
        bytes memory _preStateData
    )
        GameSolver(_gameProxy, l2Outputs, counterL2Outputs, _honestTrace, _dishonestTrace, _preStateData)
    {
        // Mark agreement with the root claim if the local opinion of the root claim is the same as the
        // observed root claim.
        agreeWithRoot = Claim.unwrap(outputAt(MAX_L2_BLOCK_NUMBER, Actor.Self)) == Claim.unwrap(_gameProxy.rootClaim());
    }

    ////////////////////////////////////////////////////////////////
    //                          EXTERNAL                          //
    ////////////////////////////////////////////////////////////////

    /// @notice Returns an array of `Move`s that can be taken from the perspective of an honest
    ///         actor in the `FaultDisputeGame` contract.
    function solveGame() external override returns (Move[] memory moves_) {
        uint256 numClaims = GAME.claimDataLen();

        // Pre-allocate the `moves_` array to the maximum possible length. Test environment, so
        // over-allocation is fine, and more easy to read than making a linked list in asm.
        moves_ = new Move[](numClaims - processedBuf);

        uint256 numMoves = 0;
        for (uint256 i = processedBuf; i < numClaims; i++) {
            // Grab the observed claim.
            IFaultDisputeGame.ClaimData memory observed = getClaimData(i);

            // Determine the direction of the next move to be taken.
            (uint64 moveDirection, Position movePos) = determineDirection(observed);

            // Continue if there is no move to be taken against the observed claim.
            if (moveDirection == NOOP_ATTACK) continue;

            if (movePos.depth() <= MAX_DEPTH) {
                // bisection
                moves_[numMoves++] = handleBisectionMove(moveDirection, movePos, i);
            } else {
                // instruction step
                moves_[numMoves++] = handleStepMove(moveDirection, observed.position, movePos, i);
            }
        }

        // Update the length of the `moves_` array to the number of moves that were added. This is
        // always a no-op or a truncation operation.
        assembly {
            mstore(moves_, numMoves)
        }

        // Increment `processedBuf` by the number of claims processed, so that next time around,
        // we don't attempt to process the same claims again.
        processedBuf += numClaims - processedBuf;
    }

    ////////////////////////////////////////////////////////////////
    //                          INTERNAL                          //
    ////////////////////////////////////////////////////////////////

    /// @dev Helper function to determine the direction of the next move to be taken.
    function determineDirection(IFaultDisputeGame.ClaimData memory _claimData)
        internal
        view
        returns (uint64 direction_, Position movePos_)
    {
        bool rightLevel = isRightLevel(_claimData.position);
        uint64 numClaims = _claimData.position.raw() == 1 || (_claimData.position.depth() == (SPLIT_DEPTH + N_BITS)) ? 1 : uint64(MAX_ATTACK_BRANCH - 1);
        Claim[] memory claims = new Claim[](numClaims);
        Claim[] memory counterClaims = new Claim[](numClaims);
        claims = subClaimsAt(_claimData.position, Actor.Self);
        if (_claimData.parentIndex == type(uint32).max || _claimData.position.depth() == SPLIT_DEPTH + N_BITS ) {
            // If we agree with the output/trace rootClaim and we agree with, ignore it.
            bool localAgree = claims[0].raw() == _claimData.claim.raw();
            if (localAgree) {
                return (NOOP_ATTACK, Position.wrap(0));
            }

            // The parent claim is the output/trace rootClaim. We must attack if we disagree per the game rules.
            direction_ = 0;
            movePos_ = _claimData.position.moveN(N_BITS, 0);
            return (direction_, movePos_);
        }

        // If the parent claim is not the root claim, first check if the observed claim is on a level that
        // agrees with the local view of the root claim. If it is, noop. If it is not, perform an attack or
        // defense depending on the local view of the observed claim.
        if (rightLevel) {
            // Never move against a claim on the right level. Even if it's wrong, if it's uncountered, it furthers
            // our goals.
            return (NOOP_ATTACK, Position.wrap(0));
        }
        counterClaims = subClaimsAt(_claimData.position, Actor.Counter);
        uint64 attackBranch_ = uint64(MAX_ATTACK_BRANCH);
        for (uint64 i = 0; i < uint64(MAX_ATTACK_BRANCH); i++) {
            if (claims[i].raw() != counterClaims[i].raw()) {
                attackBranch_ = i;
                break;
            }
        }
        return (attackBranch_, _claimData.position.moveN(N_BITS, attackBranch_));
    }

    /// @notice Returns a `Move` struct that represents an attack or defense move in the bisection portion
    ///         of the game.
    ///
    /// @dev Note: This function assumes that the `movePos` and `challengeIndex` are valid within the
    ///            output bisection context. This is enforced by the `solveGame` function.
    function handleBisectionMove(
        uint64 _attackBranch,
        Position _movePos,
        uint256 _challengeIndex
    )
        internal
        view
        returns (Move memory move_)
    {
        uint256 bond = GAME.getRequiredBond(_movePos);
        (,,,, Claim disputed,,) = GAME.claimData(_challengeIndex);

        move_ = Move({
            kind: MoveKind.Attack,
            attackBranch: _attackBranch,
            value: bond,
            data: abi.encodeCall(FaultDisputeGame.moveV2, (disputed, _challengeIndex, claimAt(_movePos), _attackBranch))
        });
    }

    /// @notice Returns a `Move` struct that represents a step move in the execution trace
    ///         bisection portion of the dispute game.
    /// @dev Note: This function assumes that the `movePos` and `challengeIndex` are valid within the
    ///            execution trace bisection context. This is enforced by the `solveGame` function.
    function handleStepMove(
        uint64 _direction,
        Position _parentPos,
        Position _movePos,
        uint256 _challengeIndex
    )
        internal
        view
        returns (Move memory move_)
    {
        bytes memory preStateTrace;
        LibDA.DAItem memory preStateItem;
        LibDA.DAItem memory postStateItem;

        // First, we need to find the pre/post state index depending on whether we
        // are making an attack step or a defense step. If the relative index at depth of the
        // move position is 0, the prestate is the absolute prestate and we need to
        // do nothing.
        if ((_movePos.indexAtDepth() % (2 ** (MAX_DEPTH - SPLIT_DEPTH))) != 0) {
            // Grab the trace up to the prestate's trace index.
            Position leafPos = Position.wrap(Position.unwrap(_parentPos) - 1 + _direction);
            preStateTrace = abi.encode(leafPos.traceIndex(MAX_DEPTH), stateAt(leafPos, Actor.Counter));
            preStateItem = daItemAtPos(_parentPos, _direction - 1, Actor.Counter);
        } else { // Left most branch
            preStateTrace = absolutePrestateData;
            preStateItem = LibDA.DAItem({
                daType: LibDA.DA_TYPE_CALLDATA,
                dataHash: GAME.absolutePrestate().raw(),
                proof: hex""
            });
        }

        if (_movePos.indexAtDepth() != (1<< (MAX_DEPTH - SPLIT_DEPTH)) - 1 ) {
             Claim postStateClaim;
            if (_direction < MAX_ATTACK_BRANCH) {
                postStateItem = daItemAtPos(_parentPos, _direction, Actor.Counter);
            } else {
                postStateItem = daItemAtPos(_parentPos, _direction, Actor.Self);
            }
        } else { // Right most branch
            postStateItem = LibDA.DAItem({
                daType: LibDA.DA_TYPE_CALLDATA,
                dataHash: GAME.rootClaim().raw(),
                proof: hex""
            });
        }

        FaultDisputeGame.StepProof memory stepProof = FaultDisputeGame.StepProof({
            preStateItem: preStateItem,
            postStateItem: postStateItem,
            vmProof: hex""
        });

        move_ = Move({
            kind: MoveKind.Step,
            attackBranch: _direction,
            value: 0,
            data: abi.encodeCall(FaultDisputeGame.stepV2, (_challengeIndex, _direction, preStateTrace, stepProof))
        });
    }

    ////////////////////////////////////////////////////////////////
    //                          HELPERS                           //
    ////////////////////////////////////////////////////////////////

    function daItemAtPos(Position _parentPos, uint64 branch, Actor actor) internal view returns (LibDA.DAItem memory daItem_) {
        Position leafPos = Position.wrap(Position.unwrap(_parentPos) + branch);
        Claim[] memory claims = new Claim[](MAX_ATTACK_BRANCH - 1);
        for (uint128 i = 0; i < branch; i++) {
            claims[i] = statehashAt(Position.wrap(Position.unwrap(_parentPos) + i), actor);
        }
        for (uint128 i = branch + 1; i < MAX_ATTACK_BRANCH; i++) {
            claims[i-1] = statehashAt(Position.wrap(Position.unwrap(_parentPos) + i), actor);
        }
        Claim claim_ = statehashAt(leafPos, actor);
        daItem_ = LibDA.DAItem({
            daType: LibDA.DA_TYPE_CALLDATA,
            dataHash: claim_.raw(),
            proof: abi.encodePacked(claims)
        });
    }

    /// @dev Helper function to get the `ClaimData` struct at a given index in the `GAME` contract's
    ///      `claimData` array.
    function getClaimData(uint256 _claimIndex) internal view returns (IFaultDisputeGame.ClaimData memory claimData_) {
        // thanks, solc
        (
            uint32 parentIndex,
            address countered,
            address claimant,
            uint128 bond,
            Claim claim,
            Position position,
            Clock clock
        ) = GAME.claimData(_claimIndex);
        claimData_ = IFaultDisputeGame.ClaimData({
            parentIndex: parentIndex,
            counteredBy: countered,
            claimant: claimant,
            bond: bond,
            claim: claim,
            position: position,
            clock: clock
        });
    }

    /// @notice Returns the player's subClaims that commits to a given position, swapping between
    ///         output bisection claims and execution trace bisection claims depending on the depth.
    /// @dev Prefer this function over `outputAt` or `statehashAt` directly.
    function subClaimsAt(Position _position, Actor actor) internal view returns (Claim[] memory claims_) {
        return _position.depth() > SPLIT_DEPTH ? statehashesAt(_position, actor) : outputsAt(_position, actor);
    }

    function claimAt(Position _position) internal view returns (Claim root_) {
        Claim[] memory claims_ = _position.depth() > SPLIT_DEPTH ? statehashesAt(_position, Actor.Self) : subClaimsAt(_position, Actor.Self);
        if (claims_.length == 1) { // It's the trace rootClaim
            root_ = claims_[0];
        } else {
            bytes memory input = abi.encodePacked(claims_); // bytes.concat(claims_[0].raw(), claims_[1].raw(), claims_[2].raw());
            root_ = Claim.wrap(LibDA.getClaimsHash(LibDA.DA_TYPE_CALLDATA, MAX_ATTACK_BRANCH, input));
        }
    }

    function outputAt(Position _position) internal view returns (Claim claim_) {
        return outputAt(_position.traceIndex(SPLIT_DEPTH) + 1, Actor.Self);
    }

    /// @notice Returns the mock output at the given position.
    function outputsAt(Position _position, Actor actor) internal view returns (Claim[] memory claims_) {
        // Don't allow for positions that are deeper than the split depth.
        if (_position.depth() > SPLIT_DEPTH) {
            revert("GameSolver: invalid position depth");
        }
        uint256 traceIndex = _position.traceIndex(SPLIT_DEPTH) + 1;
        uint8 depth = _position.depth();
        uint256 numClaims = _position.raw() == 1 ? 1 : MAX_ATTACK_BRANCH;
        claims_ = new Claim[](numClaims);
        uint256 offset = 1<< (SPLIT_DEPTH - depth);
        for (uint256 i = 0; i < numClaims; i++) {
            claims_[i] = outputAt(traceIndex + i * offset, actor);
        }
    }

    /// @notice Returns the mock output at the given L2 block number.
    function outputAt(uint256 _l2BlockNumber, Actor actor) internal view returns (Claim claim_) {
        uint256 output = actor == Actor.Self ? l2Outputs[_l2BlockNumber - 1] : counterL2Outputs[_l2BlockNumber - 1];
        return Claim.wrap(bytes32(output));
    }

    /// @notice Returns the player's claim that commits to a given trace index.
    function statehashAt(uint256 _traceIndex, Actor actor) internal view returns (Claim claim_) {
        bytes storage _trace = actor == Actor.Self ? trace: counterTrace;
        bytes32 hash =
            keccak256(abi.encode(_traceIndex >= _trace.length ? _trace.length - 1 : _traceIndex, stateAt(_traceIndex, actor)));
        assembly {
            claim_ := or(and(hash, not(shl(248, 0xFF))), shl(248, 1))
        }
    }

    function statehashAt(Position _position, Actor actor) internal view returns (Claim claim_) {
        uint256 traceIndex = _position.traceIndex(MAX_DEPTH);
        claim_ = statehashAt(traceIndex, actor);
    }

    /// @notice Returns the player's subClaims that commits to a given trace index.
    function statehashesAt(Position _position, Actor actor) internal view returns (Claim[] memory claims_) {
        uint256 depth = _position.depth();
        uint256 numClaims = depth == (SPLIT_DEPTH + N_BITS) ? 1 : MAX_ATTACK_BRANCH;
        uint256 traceIndex = _position.traceIndex(MAX_DEPTH);
        claims_ = new Claim[](numClaims);
        uint256 offset = 1<< (MAX_DEPTH - depth);
        for (uint256 i=0; i < numClaims; i++) {
            claims_[i] = statehashAt(traceIndex + i * offset, actor);
        }
    }

    /// @notice Returns the state at the trace index within the player's trace.
    function stateAt(Position _position, Actor actor) internal view returns (uint256 state_) {
        return stateAt(_position.traceIndex(MAX_DEPTH), actor);
    }

    /// @notice Returns the state at the trace index within the player's trace.
    function stateAt(uint256 _traceIndex, Actor actor) internal view returns (uint256 state_) {
        bytes storage trace = actor == Actor.Self ? trace: counterTrace;
        return uint256(uint8(_traceIndex >= trace.length ? trace[trace.length - 1] : trace[_traceIndex]));
    }

    /// @notice Returns whether or not the position is on a level which opposes the local opinion of the
    ///         root claim.
    function isRightLevel(Position _position) internal view returns (bool isRightLevel_) {
        isRightLevel_ = agreeWithRoot == ((_position.depth() / N_BITS) % 2 == 0);
    }
}

/// @title DisputeActor
/// @notice The `DisputeActor` contract is an abstract contract that represents an actor
///         that consumes the suggested moves from a `GameSolver` contract.
abstract contract DisputeActor {
    /// @notice The `GameSolver` contract used to determine the moves to be taken.
    GameSolver public solver;

    /// @notice Performs all available moves deemed by the attached solver.
    /// @return numMoves_ The number of moves that the actor took.
    /// @return success_ True if all moves were successful, false otherwise.
    function move() external virtual returns (uint256 numMoves_, bool success_);
}

/// @title HonestDisputeActor
/// @notice An actor that consumes the suggested moves from an `HonestGameSolver` contract. Note
///         that this actor *can* be dishonest if the trace is faulty, but it will always follow
///         the rules of the honest actor.
contract HonestDisputeActor is DisputeActor {
    FaultDisputeGame public immutable GAME;

    constructor(
        FaultDisputeGame _gameProxy,
        uint256[] memory _l2Outputs,
        uint256[] memory _counterL2Outputs,
        bytes memory _trace,
        bytes memory _counterTrace,
        bytes memory _preStateData
    ) {
        GAME = _gameProxy;
        solver = GameSolver(new HonestGameSolver(_gameProxy, _l2Outputs, _counterL2Outputs, _trace, _counterTrace, _preStateData));
    }

    /// @inheritdoc DisputeActor
    function move() external override returns (uint256 numMoves_, bool success_) {
        GameSolver.Move[] memory moves = solver.solveGame();
        numMoves_ = moves.length;

        // Optimistically assume success, will be set to false if any move fails.
        success_ = true;

        // Perform all available moves given to the actor by the solver.
        for (uint256 i = 0; i < moves.length; i++) {
            GameSolver.Move memory localMove = moves[i];

            // If the move is a step, we first need to add the starting L2 block number to the `PreimageOracle`
            // via the `FaultDisputeGame` contract.
            // TODO: This is leaky. Could be another move kind.
            if (localMove.kind == GameSolver.MoveKind.Step) {
                bytes memory moveData = localMove.data;
                uint256 challengeIndex;
                assembly {
                    challengeIndex := mload(add(moveData, 0x24))
                }
                LibDA.DAItem memory dummyItem = LibDA.DAItem({
                    daType: LibDA.DA_TYPE_CALLDATA,
                    dataHash: '00000000000000000000000000000000',
                    proof: hex""
                });
                GAME.addLocalData({
                    _ident: LocalPreimageKey.DISPUTED_L2_BLOCK_NUMBER,
                    _execLeafIdx: challengeIndex,
                    _partOffset: 0,
                    _daItem: dummyItem
                });
            }

            (bool innerSuccess,) = address(GAME).call{ value: localMove.value }(localMove.data);
            assembly {
                success_ := and(success_, innerSuccess)
            }
        }
    }

    fallback() external payable { }

    receive() external payable { }
}
