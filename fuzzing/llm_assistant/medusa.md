# Medusa

`medusa` is a cross-platform go-ethereum-based smart contract fuzzer inspired by Echidna. It provides parallelized fuzz
testing of smart contracts through CLI, or its Go API that allows custom user-extended testing methodology.

Traditional fuzz testing (e.g. with [`AFL`](https://lcamtuf.coredump.cx/afl/)) aims to generally explore a binary by providing
random inputs in an effort to identify new system states or crash the program (please note that this is a pretty crude generalization).
This model, however, does not translate to the smart contract ecosystem since you cannot cause a smart contract to "crash".
A transaction that reverts, for example, is not equivalent to a binary crashing or panicking.

Thus, with smart contracts, we have to change the fuzzing paradigm. When you hear of "fuzzing smart contracts", you are
not trying to crash the program but, instead, you are trying to validate the **invariants** of the program.

> **Definition**: An invariant is a property that remains unchanged after one or more operations are applied to it.

More generally, an invariant is a "truth" about some system. For smart contracts, this can take many faces.

1. **Mathematical invariants**: `a + b = b + a`. The commutative property is an invariant and any Solidity math library
   should uphold this property.
2. **ERC20 tokens**: The sum of all user balances should never exceed the total supply of the token.
3. **Automated market maker (e.g. Uniswap)**: `xy = k`. The constant-product formula is an invariant that maintains the
   economic guarantees of AMMs such as Uniswap.

> **Definition**: Smart contract fuzzing uses random sequences of transactions to test the invariants of the smart contract system.

Before we explore how to identify, write, and test invariants, it is beneficial to understand how smart contract fuzzing
works under-the-hood.

## Types of Invariants

> **Note**: In this context, property and invariant mean the same thing and are interchangeable

Defining and testing your invariants is critical to assessing the **expected system behavior**.

We like to break down invariants into two general categories: function-level invariants and system-level invariants.
Note that there are other ways of defining and scoping invariants, but this distinction is generally sufficient to
start fuzz testing even the most complex systems.

### Function-level invariants

A function-level invariant can be defined as follows:

> **Definition**: A function-level invariant is a property that arises from the execution of a specific function.

Let's take the following function from a smart contract:

```solidity
function deposit() public payable {
    // Make sure that the total deposited amount does not exceed the limit
    uint256 amount = msg.value;
    require(totalDeposited + amount <= MAX_DEPOSIT_AMOUNT);

    // Update the user balance and total deposited
    balances[msg.sender] += amount;
    totalDeposited += amount;

    emit Deposit(msg.sender, amount, totalDeposited);
}
```

The `deposit` function has the following function-level invariants:

1. The ETH balance of `msg.sender` must decrease by `amount`.
2. The ETH of `address(this)` must increase by `amount`.
3. `balances[msg.sender]` should increase by `amount`.
4. The `totalDeposited` value should increase by `amount`.

Note that there other properties that can also be tested for but the above should highlight what a function-level
invariant is. In general, function-level invariants can be identified by assessing what must be true _before_ the execution
of a function and what must be true _after_ the execution of that same function.

Let's now look at system-level invariants.

### System-level invariants

A system-level invariant can be defined as follows:

> **Definition**: A system-level invariant is a property that holds true across the _entire_ execution of a system

Thus, a system-level invariant is a lot more generalized than a function-level invariant. Here are two common examples
of a function-level invariant:

1. The `xy=k` constant product formula should always hold for Uniswap pools
2. No user's balance should ever exceed the total supply for an ERC20 token.

In the `deposit` function above, we also see the presence of a system-level invariant:

**The `totalDeposited` amount should always be less than or equal to the `MAX_DEPOSIT_AMOUNT`**.

Since the `totalDeposited` value can be affected by the presence of other functions in the system
(e.g. `withdraw` or `stake`), it is best tested at the system level instead of the function level.

### Writing Function-Level Invariants

Before we write the fuzz tests, let's look into how we would write a unit test for the `deposit` function:

```solidity
function testDeposit() public {
    // The amount of tokens to deposit
    uint256 amount = 10 ether;

    // Retrieve balance of user before deposit
    preBalance = depositContract.balances(address(this));

    // Call the deposit contract (let's assume this contract has 10 ether)
    depositContract.deposit{value: amount}();

    // Assert post-conditions
    assert(depositContract.balances(msg.sender) == preBalance + amount);
    // Add other assertions here
}
```

What we will notice about the test above is that it _fixes_ the value that is being sent. It is unable to test how the
`deposit` function behaves across a variety of input spaces. Thus, a function-level fuzz test can be thought of as a
"unit test on steroids". Instead of fixing the `amount`, we let the fuzzer control the `amount` value to any number between
`[0, type(uint256).max]` and see how the system behaves to that.

> **Note**: One of the core differences between a traditional unit test versus a fuzz test is that a fuzz test accepts input arguments that the fuzzer can control.

#### Writing a Fuzz Test for the `deposit` Function

Here is what a fuzz test for the `deposit` function would look like:

```solidity
function testDeposit(uint256 _amount) public {
    // Let's bound the input to be _at most_ the ETH balance of this contract
    // The amount value will now in between [0, address(this).balance]
    uint256 amount = clampLte(_amount, address(this).balance);

    // Retrieve balance of user before deposit
    uint256 preBalance = depositContract.balances(address(this));

    // Call the deposit contract with a variable amount
    depositContract.deposit{value: _amount}();

    // Assert post-conditions
    assert(depositContract.balances(address(this)) == preBalance + amount);
    // Add other assertions here
}
```

Notice that we bounded the `_amount` variable to be less than or equal to the test contract's ETH balance.
This type of bounding is very common when writing fuzz tests. Bounding allows you to only test values that are reasonable.
If `address(this)` doesn't have enough ETH, it does not make sense to try and call the `deposit` function. Additionally,
although we only tested one of the function-level invariants mentioned previously, writing the remaining
would follow a similar pattern as the one written above.

#### The contract to be tested and its corresponding test contract

```solidity
contract DepositContract {
    // @notice MAX_DEPOSIT_AMOUNT is the maximum amount that can be deposited into this contract
    uint256 public constant MAX_DEPOSIT_AMOUNT = 1_000_000e18;

    // @notice balances holds user balances
    mapping(address => uint256) public balances;

    // @notice totalDeposited represents the current deposited amount across all users
    uint256 public totalDeposited;

    // @notice Deposit event is emitted after a deposit occurs
    event Deposit(address depositor, uint256 amount, uint256 totalDeposited);

    // @notice deposit allows user to deposit into the system
    function deposit() public payable {
        // Make sure that the total deposited amount does not exceed the limit
        uint256 amount = msg.value;
        require(totalDeposited + amount <= MAX_DEPOSIT_AMOUNT);

        // Update the user balance and total deposited
        balances[msg.sender] += amount;
        totalDeposited += amount;

        emit Deposit(msg.sender, amount, totalDeposited);
    }
}

contract TestDepositContract {

    // @notice depositContract is an instance of DepositContract
    DepositContract depositContract;

    constructor() payable {
        // Deploy the deposit contract
        depositContract = new DepositContract();
    }

    // @notice testDeposit tests the DepositContract.deposit function
    function testDeposit(uint256 _amount) public {
        // Let's bound the input to be _at most_ the ETH balance of this contract
        // The amount value will now in between [0, address(this).balance]
        uint256 amount = clampLte(_amount, address(this).balance);

        // Retrieve balance of user before deposit
        uint256 preBalance = depositContract.balances(address(this));

        // Call the deposit contract with a variable amount
        depositContract.deposit{value: _amount}();

        // Assert post-conditions
        assert(depositContract.balances(address(this)) == preBalance + amount);
        // Add other assertions here
    }

    // @notice clampLte returns a value between [a, b]
    function clampLte(uint256 a, uint256 b) internal returns (uint256) {
        if (!(a <= b)) {
            uint256 value = a % (b + 1);
            return value;
        }
        return a;
    }

}
```

## Testing with `medusa`

`medusa` supports the following testing modes:

1. [Property Mode](https://secure-contracts.com/program-analysis/echidna/introduction/how-to-test-a-property.html)
2. [Assertion Mode](https://secure-contracts.com/program-analysis/echidna/basic/assertion-checking.html)

For more advanced information and documentation on how the various modes work and their pros/cons, check out [secure-contracts.com](https://secure-contracts.com/program-analysis/echidna/index.html)

### Writing property tests

Property tests are represented as functions within a Solidity contract whose names are prefixed with a prefix specified by the `testPrefixes` configuration option (`fuzz_` is the default test prefix). Additionally, they must take no arguments and return a `bool` indicating if the test succeeded.

```solidity
contract TestXY {
    uint x;
    uint y;

    function setX(uint value) public {
        x = value + 3;
    }

    function setY(uint value) public {
        y = value + 9;
    }

    function fuzz_never_specific_values() public returns (bool) {
        // ASSERTION: x should never be 10 at the same time y is 80
        return !(x == 10 && y == 80);
    }
}
```

`medusa` deploys your contract containing property tests and generates a sequence of calls to execute against all publicly accessible methods. After each function call, it calls upon your property tests to ensure they return a `true` (success) status.

#### Testing in property-mode

Invoking this fuzzing campaign, `medusa` will:

- Compile the given targets
- Start the configured number of worker threads, each with their own local Ethereum test chain.
- Deploy all contracts to each worker's test chain.
- Begin to generate and send call sequences to update contract state.
- Check property tests all succeed after each call executed.

Upon discovery of a failed property test, `medusa` will halt, reporting the call sequence used to violate any property test(s).
### Writing assertion tests

Although both property-mode and assertion-mode try to validate / invalidate invariants of the system, they do so in different ways. In property-mode, `medusa` will look for functions with a specific test prefix (e.g. `fuzz_`) and test those. In assertion-mode, `medusa` will test to see if a given call sequence can cause the Ethereum Virtual Machine (EVM) to "panic". The EVM has a variety of panic codes for different scenarios. For example, there is a unique panic code when an `assert(x)` statement returns `false` or when a division by zero is encountered. In assertion mode, which panics should or should not be treated as "failing test cases" can be toggled by updating the [Project Configuration](./Project-Configuration.md#fuzzing-configuration). By default, only `FailOnAssertion` is enabled. Check out the [Example Project Configuration File](https://github.com/crytic/medusa/wiki/Example-Project-Configuration-File) for a visualization of the various panic codes that can be enabled. An explanation of the various panic codes can be found in the [Solidity documentation](https://docs.soliditylang.org/en/latest/control-structures.html#panic-via-assert-and-error-via-require).

Please note that the behavior of assertion mode is different between `medusa` and Echidna. Echidna will only test for `assert(x)` statements while `medusa` provides additional flexibility.

```solidity
contract TestContract {
    uint x;
    uint y;

    function setX(uint value) public {
        x = value;

        // ASSERTION: x should be an even number
        assert(x % 2 == 0);
    }
}
```

During a call sequence, if `setX` is called with a `value` that breaks the assertion (e.g. `value = 3`), `medusa` will treat this as a failing property and report it back to the user.

#### Testing in assertion-mode

Invoking this fuzzing campaign, `medusa` will:

- Compile the given targets
- Start the configured number of worker threads, each with their own local Ethereum test chain.
- Deploy all contracts to each worker's test chain.
- Begin to generate and send call sequences to update contract state.
- Check to see if there any failing assertions after each call executed.

Upon discovery of a failed assertion, `medusa` will halt, reporting the call sequence used to violate any assertions.

### Testing with multiple modes

Note that we can run `medusa` with one, many, or no modes enabled.

```solidity
contract TestContract {
    int256 input;

    function set(int256 _input) public {
        input = _input;
    }

    function failing_assert_method(uint value) public {
        // ASSERTION: We always fail when you call this function.
        assert(false);
    }

    function fuzz_failing_property() public view returns (bool) {
        // ASSERTION: fail immediately.
        return false;
    }
}
```

## Cheatcodes Overview

Cheatcodes allow users to manipulate EVM state, blockchain behavior, provide easy ways to manipulate data, and much more.
The cheatcode contract is deployed at `0x7109709ECfa91a80626fF3989D68f67F5b1DD12D`.

### Cheatcode Interface

The following interface must be added to your Solidity project if you wish to use cheatcodes. Note that if you use Foundry
as your compilation platform that the cheatcode interface is already provided [here](https://book.getfoundry.sh/reference/forge-std/#forge-stds-test).
However, it is important to note that medusa does not support all the cheatcodes provided out-of-box
by Foundry (see below for supported cheatcodes).

```solidity
interface StdCheats {
    // Set block.timestamp
    function warp(uint256) external;

    // Set block.number
    function roll(uint256) external;

    // Set block.basefee
    function fee(uint256) external;

    // Set block.difficulty and block.prevrandao
    function difficulty(uint256) external;

    // Set block.chainid
    function chainId(uint256) external;

    // Sets the block.coinbase
    function coinbase(address) external;

    // Loads a storage slot from an address
    function load(address account, bytes32 slot) external returns (bytes32);

    // Stores a value to an address' storage slot
    function store(address account, bytes32 slot, bytes32 value) external;

    // Sets the *next* call's msg.sender to be the input address
    function prank(address) external;

    // Set msg.sender to the input address until the current call exits
    function prankHere(address) external;

    // Sets an address' balance
    function deal(address who, uint256 newBalance) external;

    // Sets an address' code
    function etch(address who, bytes calldata code) external;

    // Signs data
    function sign(uint256 privateKey, bytes32 digest)
        external
        returns (uint8 v, bytes32 r, bytes32 s);

    // Computes address for a given private key
    function addr(uint256 privateKey) external returns (address);

    // Gets the nonce of an account
    function getNonce(address account) external returns (uint64);

    // Sets the nonce of an account
    // The new nonce must be higher than the current nonce of the account
    function setNonce(address account, uint64 nonce) external;

    // Performs a foreign function call via terminal
    function ffi(string[] calldata) external returns (bytes memory);

    // Take a snapshot of the current state of the EVM
    function snapshot() external returns (uint256);

    // Revert state back to a snapshot
    function revertTo(uint256) external returns (bool);

    // Convert Solidity types to strings
    function toString(address) external returns(string memory);
    function toString(bytes calldata) external returns(string memory);
    function toString(bytes32) external returns(string memory);
    function toString(bool) external returns(string memory);
    function toString(uint256) external returns(string memory);
    function toString(int256) external returns(string memory);

    // Convert strings into Solidity types
    function parseBytes(string memory) external returns(bytes memory);
    function parseBytes32(string memory) external returns(bytes32);
    function parseAddress(string memory) external returns(address);
    function parseUint(string memory)external returns(uint256);
    function parseInt(string memory) external returns(int256);
    function parseBool(string memory) external returns(bool);
}
```

## Using cheatcodes

Below is an example snippet of how you would import the cheatcode interface into your project and use it.

```solidity
// Assuming cheatcode interface is in the same directory
import "./IStdCheats.sol";

// MyContract will utilize the cheatcode interface
contract MyContract {
    // Set up reference to cheatcode contract
    IStdCheats cheats = IStdCheats(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

    // This is a test function that will set the msg.sender's nonce to the provided input argument
    function testFunc(uint256 _x) public {
        // Ensure that the input argument is greater than msg.sender's current nonce
        require(_x > cheats.getNonce(msg.sender));

        // Set sender's nonce
        cheats.setNonce(msg.sender, x);

        // Assert that the nonce has been correctly updated
        assert(cheats.getNonce(msg.sender) == x);
    }
}
```

### Tips for Testing with Medusa

#### General

- **Use multiple testing modes:** Medusa supports property testing, assertion testing, and optimization testing. Use a combination of modes to thoroughly test your contracts.
- **Write clear and concise tests:** Your tests should be easy to read and understand. Avoid complex logic or unnecessary code.
- **Test edge cases:** Consider testing extreme values and unusual inputs to ensure your contracts handle them correctly.
- **Use a variety of test inputs:** Generate a diverse set of test inputs to cover a wide range of scenarios.
- **Monitor gas consumption:** Medusa can track gas consumption during testing. Use this information to identify areas where your contracts can be optimized.

#### Property Testing

- **Choose meaningful properties:** The properties you test should be important invariants of your contract.

#### Assertion Testing

- **Use assertions judiciously:** Assertions can be useful for catching errors, but they can also slow down testing. Use them only when necessary.
- **Test for both valid and invalid inputs:** Ensure your assertions check for both valid and invalid inputs to thoroughly test your contract's behavior.
- **Use pre-conditions and post-conditions to verify the state of the contract before and after a function call.:** Pre-conditions and post-conditions are assertions that can be used to verify the state of the contract before and after a function call. This can help to ensure that the function is called with the correct inputs, that it produces the expected outputs, and that the state of the contract is valid.


## More about medusa cheatcodes
### `addr`

#### Description

The `addr` cheatcode will compute the address for a given private key.

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Test with random private key
uint256 pkOne = 0x6df21769a2082e03f7e21f6395561279e9a7feb846b2bf740798c794ad196e00;
address addrOne = 0xdf8Ef652AdE0FA4790843a726164df8cf8649339;
address result = cheats.addr(pkOne);
assert(result == addrOne);
```

#### Function Signature

```solidity
function addr(uint256 privateKey) external returns (address);
```
### `chainId`

#### Description

The `chainId` cheatcode will set the `block.chainid`

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
cheats.chainId(777123);
assert(block.chainid == 777123);
```

#### Function Signature

```solidity
function chainId(uint256) external;
```

### `coinbase`

#### Description

The `coinbase` cheatcode will set the `block.coinbase`

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
cheats.coinbase(address(7));
assert(block.coinbase == address(7));
```

#### Function Signature

```solidity
function coinbase(address) external;
```

### `deal`

#### Description

The `deal` cheatcode will set the ETH balance of address `who` to `newBalance`

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
address acc = address(777);
cheats.deal(acc, x);
assert(acc.balance == x);
```

#### Function Signature

```solidity
function deal(address who, uint256 newBalance) external;
```

### `difficulty`

#### Description

The `difficulty` cheatcode will set the `block.difficulty` and the `block.prevrandao` value. At the moment, both values
are changed since the cheatcode does not check what EVM version is running.

Note that this behavior will change in the future.

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
cheats.difficulty(x);
assert(block.difficulty == x);
```

#### Function Signature

```solidity
function difficulty(uint256) external;
```

### `etch`

#### Description

The `etch` cheatcode will set the `who` address's bytecode to `code`.

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Obtain our original code hash for an account.
address acc = address(777);
bytes32 originalCodeHash;
assembly { originalCodeHash := extcodehash(acc) }

// Change value and verify.
cheats.etch(acc, address(someContract).code);
bytes32 updatedCodeHash;
assembly { updatedCodeHash := extcodehash(acc) }
assert(originalCodeHash != updatedCodeHash);
```

#### Function Signature

```solidity
function etch(address who, bytes calldata code) external;
```

### `fee`

#### Description

The `fee` cheatcode will set the `block.basefee`.

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
cheats.fee(7);
assert(block.basefee == 7);
```

#### Function Signature

```solidity
function fee(uint256) external;
```

### `ffi`

#### Description

The `ffi` cheatcode is used to call an arbitrary command on your host OS. Note that `ffi` must be enabled via the project
configuration file by setting `fuzzing.chainConfig.cheatCodes.enableFFI` to `true`.

Note that enabling `ffi` allows anyone to execute arbitrary commands on devices that run the fuzz tests which may
become a security risk.

#### Example with ABI-encoded hex

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Create command
string[] memory inputs = new string[](3);
inputs[0] = "echo";
inputs[1] = "-n";
// Encoded "hello"
inputs[2] = "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000568656C6C6F000000000000000000000000000000000000000000000000000000";

// Call cheats.ffi
bytes memory res = cheats.ffi(inputs);

// ABI decode
string memory output = abi.decode(res, (string));
assert(keccak256(abi.encodePacked(output)) == keccak256(abi.encodePacked("hello")));
```

#### Example with UTF8 encoding

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Create command
string[] memory inputs = new string[](3);
inputs[0] = "echo";
inputs[1] = "-n";
inputs[2] = "hello";

// Call cheats.ffi
bytes memory res = cheats.ffi(inputs);

// Convert to UTF-8 string
string memory output = string(res);
assert(keccak256(abi.encodePacked(output)) == keccak256(abi.encodePacked("hello")));
```

#### Function Signature

```solidity
function ffi(string[] calldata) external returns (bytes memory);
```

### `getNonce`

#### Description

The `getNonce` cheatcode will get the current nonce of `account`.

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Get nonce and verify that the sender has sent at least one transaction
address acc = address(msg.sender);
assert(cheats.getNonce(acc) > 0);
```

#### Function Signature

```solidity
function getNonce(address account) external returns (uint64);
```

### `load`

#### Description

The `load` cheatcode will load storage slot `slot` for `account`

#### Example

```solidity
contract TestContract {
    uint x = 123;
    function test() public {
        // Obtain our cheat code contract reference.
        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

        // Load and verify x
        bytes32 value = cheats.load(address(this), bytes32(uint(0)));
        assert(value == bytes32(uint(123)));
    }
}
```

#### Function Signature

```solidity
function load(address account, bytes32 slot) external returns (bytes32);
```

### `prank`

#### Description

The `prank` cheatcode will set the `msg.sender` for _only the next call_ to the specified input address. Note that,
contrary to [`prank` in Foundry](https://book.getfoundry.sh/cheatcodes/prank#description), calling the cheatcode contract will count as a
valid "next call"

#### Example

```solidity
contract TestContract {
    address owner = address(123);
    function transferOwnership(address _newOwner) public {
        require(msg.sender == owner);

        // Change ownership
        owner = _newOwner;
    }

    function test() public {
        // Obtain our cheat code contract reference.
        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

        // Prank, change ownership, and verify
        address newOwner = address(456);
        cheats.prank(owner);
        transferOwnership(newOwner);
        assert(owner == newOwner);
    }
    }
```

#### Function Signature

```solidity
function prank(address) external;
```

### `prankHere`

#### Description

The `prankHere` cheatcode will set the `msg.sender` to the specified input address until the current call exits. Compared
to `prank`, `prankHere` can persist for multiple calls.

#### Example

```solidity
contract TestContract {
    address owner = address(123);
    uint256 x = 0;
    uint256 y = 0;

    function updateX() public {
        require(msg.sender == owner);

        // Update x
        x = 1;
    }

    function updateY() public {
        require(msg.sender == owner);

        // Update y
        y = 1;
    }

    function test() public {
        // Obtain our cheat code contract reference.
        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

        // Prank, update variables, and verify
        cheats.prank(owner);
        updateX();
        updateY();
        assert((x == 1) && (y == 1));

        // Once this function returns, the `msg.sender` is reset
    }
}
```

#### Function Signature

```solidity
function prankHere(address) external;
```

### `roll`

#### Description

The `roll` cheatcode sets the `block.number`

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
cheats.roll(7);
assert(block.number == 7);
cheats.roll(9);
assert(block.number == 9);
```

#### Function Signature

```solidity
function roll(uint256) external;
```

### setNonce

#### Description

The `setNonce` cheatcode will set the nonce of `account` to `nonce`. Note that the `nonce` must be strictly greater than
the current nonce

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Set nonce and verify (assume nonce before `setNonce` was less than 7)
address acc = address(msg.sender);
cheats.setNonce(acc, 7);
assert(cheats.getNonce(acc) == 7);
```

#### Function Signature

```solidity
function setNonce(address account, uint64 nonce) external;
```

### `sign`

#### Description

The `sign` cheatcode will take in a private key `privateKey` and a hash digest `digest` to generate a `(v, r, s)`
signature

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

bytes32 digest = keccak256("Data To Sign");

// Call cheats.sign
(uint8 v, bytes32 r, bytes32 s) = cheats.sign(0x6df21769a2082e03f7e21f6395561279e9a7feb846b2bf740798c794ad196e00, digest);
address signer = ecrecover(digest, v, r, s);
assert(signer == 0xdf8Ef652AdE0FA4790843a726164df8cf8649339);
```

#### Function Signature

```solidity
function sign(uint256 privateKey, bytes32 digest)
external
returns (uint8 v, bytes32 r, bytes32 s);
```

### `snapshot` and `revertTo`

#### Description

The `snapshot` cheatcode will take a snapshot of the current state of the blockchain and return an identifier for the
snapshot.

On the flipside, the `revertTo` cheatcode will revert the EVM state back based on the provided identifier.

#### Example

```solidity
interface CheatCodes {
    function warp(uint256) external;

    function deal(address, uint256) external;

    function snapshot() external returns (uint256);

    function revertTo(uint256) external returns (bool);
}

struct Storage {
    uint slot0;
    uint slot1;
}

contract TestContract {
    Storage store;
    uint256 timestamp;

    function test() public {
        // Obtain our cheat code contract reference.
        CheatCodes cheats = CheatCodes(
            0x7109709ECfa91a80626fF3989D68f67F5b1DD12D
        );

        store.slot0 = 10;
        store.slot1 = 20;
        timestamp = block.timestamp;
        cheats.deal(address(this), 5 ether);

        // Save state
        uint256 snapshot = cheats.snapshot();

        // Change state
        store.slot0 = 300;
        store.slot1 = 400;
        cheats.deal(address(this), 500 ether);
        cheats.warp(12345);

        // Assert that state has been changed
        assert(store.slot0 == 300);
        assert(store.slot1 == 400);
        assert(address(this).balance == 500 ether);
        assert(block.timestamp == 12345);

        // Revert to snapshot
        cheats.revertTo(snapshot);

        // Ensure state has been reset
        assert(store.slot0 == 10);
        assert(store.slot1 == 20);
        assert(address(this).balance == 5 ether);
        assert(block.timestamp == timestamp);
    }
}
```

### `store`

#### Description

The `store` cheatcode will store `value` in storage slot `slot` for `account`

#### Example

```solidity
contract TestContract {
    uint x = 123;
    function test() public {
        // Obtain our cheat code contract reference.
        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

        // Store into x, verify it.
        cheats.store(address(this), bytes32(uint(0)), bytes32(uint(456)));
        assert(y == 456);
    }
}
```

#### Function Signature

```solidity
function store(address account, bytes32 slot, bytes32 value) external;
```

### warp

#### Description

The `warp` cheatcode sets the `block.timestamp`

#### Example

```solidity
// Obtain our cheat code contract reference.
IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

// Change value and verify.
cheats.warp(7);
assert(block.timestamp == 7);
cheats.warp(9);
assert(block.timestamp == 9);
```

#### Function Signature

```solidity
function warp(uint256) external;
```

