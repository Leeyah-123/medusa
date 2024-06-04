package llm_assistant

import "fmt"

func GenerateFuzzHarnessPrompt(contractPath string, testContractPath string, contractContents string, testContractContents string, contractName string, testContractName string) string {
	prompt := "The following text in triple quotes is my Solidity file that resides at %s containing my main contracts: '''%s'''\n\n" +
		"I want to fuzz test the contract in this Solidity file. The following text in triple quotes denotes the fuzz harness that resides at %s: '''%s'''\n\n" +
		"Step 1 - If a contract named '%s' does not exist in the Solidity file, create it. This contract contains the tests for '%s'.\n" +
		"Step 2 - Go through the tests present in '%s', if any. If there are no tests present, create a test case for '%s'.If there are tests present, add another test case for '%s' to '%s'.\n" +
		"Note: The test case created should be a system invariant for '%s' and should be in the following form in triple quotes:" +
		"'''function testDeposit(uint256 _amount) public {\n" +
		"    // Let's bound the input to be _at most_ the ETH balance of this contract\n" +
		"    // The amount value will now in between [0, address(this).balance]\n" +
		"    uint256 amount = clampLte(_amount, address(this).balance);\n\n" +
		"    // Retrieve balance of user before deposit\n" +
		"    uint256 preBalance = depositContract.balances(address(this));\n\n" +
		"    // Call the deposit contract with a variable amount\n" +
		"    depositContract.deposit{value: _amount}();\n\n" +
		"    // Assert post-conditions\n" +
		"    assert(depositContract.balances(address(this)) == preBalance + amount);\n" +
		"}'''\n" +
		"Note: Return only the contents of the test file after the test case has been added to '%s'. If no solidity license specifier or version specifier was present in the test file, add it in based on the license and version in the main contract file.\n" +
		"Note: Make sure to properly document the code for better understanding.\n" +
		"Note: Make sure to include an import of the main contract in the test file. The imported path should be relative to the test file.\n" +
		"Note: You should not include any text in your response other than the generated test file, neither should your response be in markdown format as this response will be written directly to a solidity file.\n"

	return fmt.Sprintf(prompt, contractPath, contractContents, testContractPath, testContractContents, testContractName, contractName, testContractName, contractName, contractName, testContractName, contractName, testContractName)
}

func RegenerateFuzzHarnessPrompt(errEncountered string) string {
	prompt := "There is an error in the generated fuzzing harness. Here is the error in triple quotes: '''%s'''.\n" +
		"Please fix the error and re-generate the test file.\n." +
		"Note: Do not add any new invariants to the test file nor remove any invariants from the test file.\n" +
		"Note: Return only the contents of the re-generated test file.\n" +
		"Note: Make sure to properly document the code for better understanding, but do not add comments regarding the changes you have j made to the test file.\n" +
		"Note: You should not include any text in your response other than the generated test file, neither should your response be in markdown format as this response will be written directly to a solidity file.\n"

	return fmt.Sprintf(prompt, errEncountered)
}

func TrainingPrompts() []Message {
	return []Message{
		{
			Role:    "system",
			Content: "# Medusa\n\n`medusa` is a cross-platform go-ethereum-based smart contract fuzzer inspired by Echidna. It provides parallelized fuzz\ntesting of smart contracts through CLI, or its Go API that allows custom user-extended testing methodology.\n\nTraditional fuzz testing (e.g. with [`AFL`](https://lcamtuf.coredump.cx/afl/)) aims to generally explore a binary by providing\nrandom inputs in an effort to identify new system states or crash the program (please note that this is a pretty crude generalization).\nThis model, however, does not translate to the smart contract ecosystem since you cannot cause a smart contract to \"crash\".\nA transaction that reverts, for example, is not equivalent to a binary crashing or panicking.\n\nThus, with smart contracts, we have to change the fuzzing paradigm. When you hear of \"fuzzing smart contracts\", you are\nnot trying to crash the program but, instead, you are trying to validate the **invariants** of the program.\n\n> **Definition**: An invariant is a property that remains unchanged after one or more operations are applied to it.\n\nMore generally, an invariant is a \"truth\" about some system. For smart contracts, this can take many faces.\n\n1. **Mathematical invariants**: `a + b = b + a`. The commutative property is an invariant and any Solidity math library\n   should uphold this property.\n2. **ERC20 tokens**: The sum of all user balances should never exceed the total supply of the token.\n3. **Automated market maker (e.g. Uniswap)**: `xy = k`. The constant-product formula is an invariant that maintains the\n   economic guarantees of AMMs such as Uniswap.\n\n> **Definition**: Smart contract fuzzing uses random sequences of transactions to test the invariants of the smart contract system.\n\nBefore we explore how to identify, write, and test invariants, it is beneficial to understand how smart contract fuzzing\nworks under-the-hood.\n\n## Types of Invariants\n\n> **Note**: In this context, property and invariant mean the same thing and are interchangeable\n\nDefining and testing your invariants is critical to assessing the **expected system behavior**.\n\nWe like to break down invariants into two general categories: function-level invariants and system-level invariants.\nNote that there are other ways of defining and scoping invariants, but this distinction is generally sufficient to\nstart fuzz testing even the most complex systems.\n\n### Function-level invariants\n\nA function-level invariant can be defined as follows:\n\n> **Definition**: A function-level invariant is a property that arises from the execution of a specific function.\n\nLet's take the following function from a smart contract:\n\n```solidity\nfunction deposit() public payable {\n    // Make sure that the total deposited amount does not exceed the limit\n    uint256 amount = msg.value;\n    require(totalDeposited + amount <= MAX_DEPOSIT_AMOUNT);\n\n    // Update the user balance and total deposited\n    balances[msg.sender] += amount;\n    totalDeposited += amount;\n\n    emit Deposit(msg.sender, amount, totalDeposited);\n}\n```\n\nThe `deposit` function has the following function-level invariants:\n\n1. The ETH balance of `msg.sender` must decrease by `amount`.\n2. The ETH of `address(this)` must increase by `amount`.\n3. `balances[msg.sender]` should increase by `amount`.\n4. The `totalDeposited` value should increase by `amount`.\n\nNote that there other properties that can also be tested for but the above should highlight what a function-level\ninvariant is. In general, function-level invariants can be identified by assessing what must be true _before_ the execution\nof a function and what must be true _after_ the execution of that same function.\n\nLet's now look at system-level invariants.\n\n### System-level invariants\n\nA system-level invariant can be defined as follows:\n\n> **Definition**: A system-level invariant is a property that holds true across the _entire_ execution of a system\n\nThus, a system-level invariant is a lot more generalized than a function-level invariant. Here are two common examples\nof a function-level invariant:\n\n1. The `xy=k` constant product formula should always hold for Uniswap pools\n2. No user's balance should ever exceed the total supply for an ERC20 token.\n\nIn the `deposit` function above, we also see the presence of a system-level invariant:\n\n**The `totalDeposited` amount should always be less than or equal to the `MAX_DEPOSIT_AMOUNT`**.\n\nSince the `totalDeposited` value can be affected by the presence of other functions in the system\n(e.g. `withdraw` or `stake`), it is best tested at the system level instead of the function level.\n\n### Writing Function-Level Invariants\n\nBefore we write the fuzz tests, let's look into how we would write a unit test for the `deposit` function:\n\n```solidity\nfunction testDeposit() public {\n    // The amount of tokens to deposit\n    uint256 amount = 10 ether;\n\n    // Retrieve balance of user before deposit\n    preBalance = depositContract.balances(address(this));\n\n    // Call the deposit contract (let's assume this contract has 10 ether)\n    depositContract.deposit{value: amount}();\n\n    // Assert post-conditions\n    assert(depositContract.balances(msg.sender) == preBalance + amount);\n    // Add other assertions here\n}\n```\n\nWhat we will notice about the test above is that it _fixes_ the value that is being sent. It is unable to test how the\n`deposit` function behaves across a variety of input spaces. Thus, a function-level fuzz test can be thought of as a\n\"unit test on steroids\". Instead of fixing the `amount`, we let the fuzzer control the `amount` value to any number between\n`[0, type(uint256).max]` and see how the system behaves to that.\n\n> **Note**: One of the core differences between a traditional unit test versus a fuzz test is that a fuzz test accepts input arguments that the fuzzer can control.\n\n#### Writing a Fuzz Test for the `deposit` Function\n\nHere is what a fuzz test for the `deposit` function would look like:\n\n```solidity\nfunction testDeposit(uint256 _amount) public {\n    // Let's bound the input to be _at most_ the ETH balance of this contract\n    // The amount value will now in between [0, address(this).balance]\n    uint256 amount = clampLte(_amount, address(this).balance);\n\n    // Retrieve balance of user before deposit\n    uint256 preBalance = depositContract.balances(address(this));\n\n    // Call the deposit contract with a variable amount\n    depositContract.deposit{value: _amount}();\n\n    // Assert post-conditions\n    assert(depositContract.balances(address(this)) == preBalance + amount);\n    // Add other assertions here\n}\n```\n\nNotice that we bounded the `_amount` variable to be less than or equal to the test contract's ETH balance.\nThis type of bounding is very common when writing fuzz tests. Bounding allows you to only test values that are reasonable.\nIf `address(this)` doesn't have enough ETH, it does not make sense to try and call the `deposit` function. Additionally,\nalthough we only tested one of the function-level invariants mentioned previously, writing the remaining\nwould follow a similar pattern as the one written above.\n\n#### The contract to be tested and its corresponding test contract\n\n```solidity\ncontract DepositContract {\n    // @notice MAX_DEPOSIT_AMOUNT is the maximum amount that can be deposited into this contract\n    uint256 public constant MAX_DEPOSIT_AMOUNT = 1_000_000e18;\n\n    // @notice balances holds user balances\n    mapping(address => uint256) public balances;\n\n    // @notice totalDeposited represents the current deposited amount across all users\n    uint256 public totalDeposited;\n\n    // @notice Deposit event is emitted after a deposit occurs\n    event Deposit(address depositor, uint256 amount, uint256 totalDeposited);\n\n    // @notice deposit allows user to deposit into the system\n    function deposit() public payable {\n        // Make sure that the total deposited amount does not exceed the limit\n        uint256 amount = msg.value;\n        require(totalDeposited + amount <= MAX_DEPOSIT_AMOUNT);\n\n        // Update the user balance and total deposited\n        balances[msg.sender] += amount;\n        totalDeposited += amount;\n\n        emit Deposit(msg.sender, amount, totalDeposited);\n    }\n}\n\ncontract TestDepositContract {\n\n    // @notice depositContract is an instance of DepositContract\n    DepositContract depositContract;\n\n    constructor() payable {\n        // Deploy the deposit contract\n        depositContract = new DepositContract();\n    }\n\n    // @notice testDeposit tests the DepositContract.deposit function\n    function testDeposit(uint256 _amount) public {\n        // Let's bound the input to be _at most_ the ETH balance of this contract\n        // The amount value will now in between [0, address(this).balance]\n        uint256 amount = clampLte(_amount, address(this).balance);\n\n        // Retrieve balance of user before deposit\n        uint256 preBalance = depositContract.balances(address(this));\n\n        // Call the deposit contract with a variable amount\n        depositContract.deposit{value: _amount}();\n\n        // Assert post-conditions\n        assert(depositContract.balances(address(this)) == preBalance + amount);\n        // Add other assertions here\n    }\n\n    // @notice clampLte returns a value between [a, b]\n    function clampLte(uint256 a, uint256 b) internal returns (uint256) {\n        if (!(a <= b)) {\n            uint256 value = a % (b + 1);\n            return value;\n        }\n        return a;\n    }\n\n}\n```\n\n## Testing with `medusa`\n\n`medusa` supports the following testing modes:\n\n1. [Property Mode](https://secure-contracts.com/program-analysis/echidna/introduction/how-to-test-a-property.html)\n2. [Assertion Mode](https://secure-contracts.com/program-analysis/echidna/basic/assertion-checking.html)\n\nFor more advanced information and documentation on how the various modes work and their pros/cons, check out [secure-contracts.com](https://secure-contracts.com/program-analysis/echidna/index.html)\n\n### Writing property tests\n\nProperty tests are represented as functions within a Solidity contract whose names are prefixed with a prefix specified by the `testPrefixes` configuration option (`fuzz_` is the default test prefix). Additionally, they must take no arguments and return a `bool` indicating if the test succeeded.\n\n```solidity\ncontract TestXY {\n    uint x;\n    uint y;\n\n    function setX(uint value) public {\n        x = value + 3;\n    }\n\n    function setY(uint value) public {\n        y = value + 9;\n    }\n\n    function fuzz_never_specific_values() public returns (bool) {\n        // ASSERTION: x should never be 10 at the same time y is 80\n        return !(x == 10 && y == 80);\n    }\n}\n```\n\n`medusa` deploys your contract containing property tests and generates a sequence of calls to execute against all publicly accessible methods. After each function call, it calls upon your property tests to ensure they return a `true` (success) status.\n\n#### Testing in property-mode\n\nInvoking this fuzzing campaign, `medusa` will:\n\n- Compile the given targets\n- Start the configured number of worker threads, each with their own local Ethereum test chain.\n- Deploy all contracts to each worker's test chain.\n- Begin to generate and send call sequences to update contract state.\n- Check property tests all succeed after each call executed.\n\nUpon discovery of a failed property test, `medusa` will halt, reporting the call sequence used to violate any property test(s).\n### Writing assertion tests\n\nAlthough both property-mode and assertion-mode try to validate / invalidate invariants of the system, they do so in different ways. In property-mode, `medusa` will look for functions with a specific test prefix (e.g. `fuzz_`) and test those. In assertion-mode, `medusa` will test to see if a given call sequence can cause the Ethereum Virtual Machine (EVM) to \"panic\". The EVM has a variety of panic codes for different scenarios. For example, there is a unique panic code when an `assert(x)` statement returns `false` or when a division by zero is encountered. In assertion mode, which panics should or should not be treated as \"failing test cases\" can be toggled by updating the [Project Configuration](./Project-Configuration.md#fuzzing-configuration). By default, only `FailOnAssertion` is enabled. Check out the [Example Project Configuration File](https://github.com/crytic/medusa/wiki/Example-Project-Configuration-File) for a visualization of the various panic codes that can be enabled. An explanation of the various panic codes can be found in the [Solidity documentation](https://docs.soliditylang.org/en/latest/control-structures.html#panic-via-assert-and-error-via-require).\n\nPlease note that the behavior of assertion mode is different between `medusa` and Echidna. Echidna will only test for `assert(x)` statements while `medusa` provides additional flexibility.\n\n```solidity\ncontract TestContract {\n    uint x;\n    uint y;\n\n    function setX(uint value) public {\n        x = value;\n\n        // ASSERTION: x should be an even number\n        assert(x % 2 == 0);\n    }\n}\n```\n\nDuring a call sequence, if `setX` is called with a `value` that breaks the assertion (e.g. `value = 3`), `medusa` will treat this as a failing property and report it back to the user.\n\n#### Testing in assertion-mode\n\nInvoking this fuzzing campaign, `medusa` will:\n\n- Compile the given targets\n- Start the configured number of worker threads, each with their own local Ethereum test chain.\n- Deploy all contracts to each worker's test chain.\n- Begin to generate and send call sequences to update contract state.\n- Check to see if there any failing assertions after each call executed.\n\nUpon discovery of a failed assertion, `medusa` will halt, reporting the call sequence used to violate any assertions.\n\n### Testing with multiple modes\n\nNote that we can run `medusa` with one, many, or no modes enabled.\n\n```solidity\ncontract TestContract {\n    int256 input;\n\n    function set(int256 _input) public {\n        input = _input;\n    }\n\n    function failing_assert_method(uint value) public {\n        // ASSERTION: We always fail when you call this function.\n        assert(false);\n    }\n\n    function fuzz_failing_property() public view returns (bool) {\n        // ASSERTION: fail immediately.\n        return false;\n    }\n}\n```\n\n## Cheatcodes Overview\n\nCheatcodes allow users to manipulate EVM state, blockchain behavior, provide easy ways to manipulate data, and much more.\nThe cheatcode contract is deployed at `0x7109709ECfa91a80626fF3989D68f67F5b1DD12D`.\n\n### Cheatcode Interface\n\nThe following interface must be added to your Solidity project if you wish to use cheatcodes. Note that if you use Foundry\nas your compilation platform that the cheatcode interface is already provided [here](https://book.getfoundry.sh/reference/forge-std/#forge-stds-test).\nHowever, it is important to note that medusa does not support all the cheatcodes provided out-of-box\nby Foundry (see below for supported cheatcodes).\n\n```solidity\ninterface StdCheats {\n    // Set block.timestamp\n    function warp(uint256) external;\n\n    // Set block.number\n    function roll(uint256) external;\n\n    // Set block.basefee\n    function fee(uint256) external;\n\n    // Set block.difficulty and block.prevrandao\n    function difficulty(uint256) external;\n\n    // Set block.chainid\n    function chainId(uint256) external;\n\n    // Sets the block.coinbase\n    function coinbase(address) external;\n\n    // Loads a storage slot from an address\n    function load(address account, bytes32 slot) external returns (bytes32);\n\n    // Stores a value to an address' storage slot\n    function store(address account, bytes32 slot, bytes32 value) external;\n\n    // Sets the *next* call's msg.sender to be the input address\n    function prank(address) external;\n\n    // Set msg.sender to the input address until the current call exits\n    function prankHere(address) external;\n\n    // Sets an address' balance\n    function deal(address who, uint256 newBalance) external;\n\n    // Sets an address' code\n    function etch(address who, bytes calldata code) external;\n\n    // Signs data\n    function sign(uint256 privateKey, bytes32 digest)\n        external\n        returns (uint8 v, bytes32 r, bytes32 s);\n\n    // Computes address for a given private key\n    function addr(uint256 privateKey) external returns (address);\n\n    // Gets the nonce of an account\n    function getNonce(address account) external returns (uint64);\n\n    // Sets the nonce of an account\n    // The new nonce must be higher than the current nonce of the account\n    function setNonce(address account, uint64 nonce) external;\n\n    // Performs a foreign function call via terminal\n    function ffi(string[] calldata) external returns (bytes memory);\n\n    // Take a snapshot of the current state of the EVM\n    function snapshot() external returns (uint256);\n\n    // Revert state back to a snapshot\n    function revertTo(uint256) external returns (bool);\n\n    // Convert Solidity types to strings\n    function toString(address) external returns(string memory);\n    function toString(bytes calldata) external returns(string memory);\n    function toString(bytes32) external returns(string memory);\n    function toString(bool) external returns(string memory);\n    function toString(uint256) external returns(string memory);\n    function toString(int256) external returns(string memory);\n\n    // Convert strings into Solidity types\n    function parseBytes(string memory) external returns(bytes memory);\n    function parseBytes32(string memory) external returns(bytes32);\n    function parseAddress(string memory) external returns(address);\n    function parseUint(string memory)external returns(uint256);\n    function parseInt(string memory) external returns(int256);\n    function parseBool(string memory) external returns(bool);\n}\n```\n\n## Using cheatcodes\n\nBelow is an example snippet of how you would import the cheatcode interface into your project and use it.\n\n```solidity\n// Assuming cheatcode interface is in the same directory\nimport \"./IStdCheats.sol\";\n\n// MyContract will utilize the cheatcode interface\ncontract MyContract {\n    // Set up reference to cheatcode contract\n    IStdCheats cheats = IStdCheats(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n    // This is a test function that will set the msg.sender's nonce to the provided input argument\n    function testFunc(uint256 _x) public {\n        // Ensure that the input argument is greater than msg.sender's current nonce\n        require(_x > cheats.getNonce(msg.sender));\n\n        // Set sender's nonce\n        cheats.setNonce(msg.sender, x);\n\n        // Assert that the nonce has been correctly updated\n        assert(cheats.getNonce(msg.sender) == x);\n    }\n}\n```\n\n### Tips for Testing with Medusa\n\n#### General\n\n- **Use multiple testing modes:** Medusa supports property testing, assertion testing, and optimization testing. Use a combination of modes to thoroughly test your contracts.\n- **Write clear and concise tests:** Your tests should be easy to read and understand. Avoid complex logic or unnecessary code.\n- **Test edge cases:** Consider testing extreme values and unusual inputs to ensure your contracts handle them correctly.\n- **Use a variety of test inputs:** Generate a diverse set of test inputs to cover a wide range of scenarios.\n- **Monitor gas consumption:** Medusa can track gas consumption during testing. Use this information to identify areas where your contracts can be optimized.\n\n#### Property Testing\n\n- **Choose meaningful properties:** The properties you test should be important invariants of your contract.\n\n#### Assertion Testing\n\n- **Use assertions judiciously:** Assertions can be useful for catching errors, but they can also slow down testing. Use them only when necessary.\n- **Test for both valid and invalid inputs:** Ensure your assertions check for both valid and invalid inputs to thoroughly test your contract's behavior.\n- **Use pre-conditions and post-conditions to verify the state of the contract before and after a function call.:** Pre-conditions and post-conditions are assertions that can be used to verify the state of the contract before and after a function call. This can help to ensure that the function is called with the correct inputs, that it produces the expected outputs, and that the state of the contract is valid.",
		},
		{
			Role:    "system",
			Content: "## More about medusa cheatcodes\n### `addr`\n\n#### Description\n\nThe `addr` cheatcode will compute the address for a given private key.\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Test with random private key\nuint256 pkOne = 0x6df21769a2082e03f7e21f6395561279e9a7feb846b2bf740798c794ad196e00;\naddress addrOne = 0xdf8Ef652AdE0FA4790843a726164df8cf8649339;\naddress result = cheats.addr(pkOne);\nassert(result == addrOne);\n```\n\n#### Function Signature\n\n```solidity\nfunction addr(uint256 privateKey) external returns (address);\n```\n### `chainId`\n\n#### Description\n\nThe `chainId` cheatcode will set the `block.chainid`\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\ncheats.chainId(777123);\nassert(block.chainid == 777123);\n```\n\n#### Function Signature\n\n```solidity\nfunction chainId(uint256) external;\n```\n\n### `coinbase`\n\n#### Description\n\nThe `coinbase` cheatcode will set the `block.coinbase`\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\ncheats.coinbase(address(7));\nassert(block.coinbase == address(7));\n```\n\n#### Function Signature\n\n```solidity\nfunction coinbase(address) external;\n```\n\n### `deal`\n\n#### Description\n\nThe `deal` cheatcode will set the ETH balance of address `who` to `newBalance`\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\naddress acc = address(777);\ncheats.deal(acc, x);\nassert(acc.balance == x);\n```\n\n#### Function Signature\n\n```solidity\nfunction deal(address who, uint256 newBalance) external;\n```\n\n### `difficulty`\n\n#### Description\n\nThe `difficulty` cheatcode will set the `block.difficulty` and the `block.prevrandao` value. At the moment, both values\nare changed since the cheatcode does not check what EVM version is running.\n\nNote that this behavior will change in the future.\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\ncheats.difficulty(x);\nassert(block.difficulty == x);\n```\n\n#### Function Signature\n\n```solidity\nfunction difficulty(uint256) external;\n```\n\n### `etch`\n\n#### Description\n\nThe `etch` cheatcode will set the `who` address's bytecode to `code`.\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Obtain our original code hash for an account.\naddress acc = address(777);\nbytes32 originalCodeHash;\nassembly { originalCodeHash := extcodehash(acc) }\n\n// Change value and verify.\ncheats.etch(acc, address(someContract).code);\nbytes32 updatedCodeHash;\nassembly { updatedCodeHash := extcodehash(acc) }\nassert(originalCodeHash != updatedCodeHash);\n```\n\n#### Function Signature\n\n```solidity\nfunction etch(address who, bytes calldata code) external;\n```\n\n### `fee`\n\n#### Description\n\nThe `fee` cheatcode will set the `block.basefee`.\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\ncheats.fee(7);\nassert(block.basefee == 7);\n```\n\n#### Function Signature\n\n```solidity\nfunction fee(uint256) external;\n```\n\n### `ffi`\n\n#### Description\n\nThe `ffi` cheatcode is used to call an arbitrary command on your host OS. Note that `ffi` must be enabled via the project\nconfiguration file by setting `fuzzing.chainConfig.cheatCodes.enableFFI` to `true`.\n\nNote that enabling `ffi` allows anyone to execute arbitrary commands on devices that run the fuzz tests which may\nbecome a security risk.\n\n#### Example with ABI-encoded hex\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Create command\nstring[] memory inputs = new string[](3);\ninputs[0] = \"echo\";\ninputs[1] = \"-n\";\n// Encoded \"hello\"\ninputs[2] = \"0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000568656C6C6F000000000000000000000000000000000000000000000000000000\";\n\n// Call cheats.ffi\nbytes memory res = cheats.ffi(inputs);\n\n// ABI decode\nstring memory output = abi.decode(res, (string));\nassert(keccak256(abi.encodePacked(output)) == keccak256(abi.encodePacked(\"hello\")));\n```\n\n#### Example with UTF8 encoding\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Create command\nstring[] memory inputs = new string[](3);\ninputs[0] = \"echo\";\ninputs[1] = \"-n\";\ninputs[2] = \"hello\";\n\n// Call cheats.ffi\nbytes memory res = cheats.ffi(inputs);\n\n// Convert to UTF-8 string\nstring memory output = string(res);\nassert(keccak256(abi.encodePacked(output)) == keccak256(abi.encodePacked(\"hello\")));\n```\n\n#### Function Signature\n\n```solidity\nfunction ffi(string[] calldata) external returns (bytes memory);\n```\n\n### `getNonce`\n\n#### Description\n\nThe `getNonce` cheatcode will get the current nonce of `account`.\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Get nonce and verify that the sender has sent at least one transaction\naddress acc = address(msg.sender);\nassert(cheats.getNonce(acc) > 0);\n```\n\n#### Function Signature\n\n```solidity\nfunction getNonce(address account) external returns (uint64);\n```\n\n### `load`\n\n#### Description\n\nThe `load` cheatcode will load storage slot `slot` for `account`\n\n#### Example\n\n```solidity\ncontract TestContract {\n    uint x = 123;\n    function test() public {\n        // Obtain our cheat code contract reference.\n        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n        // Load and verify x\n        bytes32 value = cheats.load(address(this), bytes32(uint(0)));\n        assert(value == bytes32(uint(123)));\n    }\n}\n```\n\n#### Function Signature\n\n```solidity\nfunction load(address account, bytes32 slot) external returns (bytes32);\n```\n\n### `prank`\n\n#### Description\n\nThe `prank` cheatcode will set the `msg.sender` for _only the next call_ to the specified input address. Note that,\ncontrary to [`prank` in Foundry](https://book.getfoundry.sh/cheatcodes/prank#description), calling the cheatcode contract will count as a\nvalid \"next call\"\n\n#### Example\n\n```solidity\ncontract TestContract {\n    address owner = address(123);\n    function transferOwnership(address _newOwner) public {\n        require(msg.sender == owner);\n\n        // Change ownership\n        owner = _newOwner;\n    }\n\n    function test() public {\n        // Obtain our cheat code contract reference.\n        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n        // Prank, change ownership, and verify\n        address newOwner = address(456);\n        cheats.prank(owner);\n        transferOwnership(newOwner);\n        assert(owner == newOwner);\n    }\n    }\n```\n\n#### Function Signature\n\n```solidity\nfunction prank(address) external;\n```\n\n### `prankHere`\n\n#### Description\n\nThe `prankHere` cheatcode will set the `msg.sender` to the specified input address until the current call exits. Compared\nto `prank`, `prankHere` can persist for multiple calls.\n\n#### Example\n\n```solidity\ncontract TestContract {\n    address owner = address(123);\n    uint256 x = 0;\n    uint256 y = 0;\n\n    function updateX() public {\n        require(msg.sender == owner);\n\n        // Update x\n        x = 1;\n    }\n\n    function updateY() public {\n        require(msg.sender == owner);\n\n        // Update y\n        y = 1;\n    }\n\n    function test() public {\n        // Obtain our cheat code contract reference.\n        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n        // Prank, update variables, and verify\n        cheats.prank(owner);\n        updateX();\n        updateY();\n        assert((x == 1) && (y == 1));\n\n        // Once this function returns, the `msg.sender` is reset\n    }\n}\n```\n\n#### Function Signature\n\n```solidity\nfunction prankHere(address) external;\n```\n\n### `roll`\n\n#### Description\n\nThe `roll` cheatcode sets the `block.number`\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\ncheats.roll(7);\nassert(block.number == 7);\ncheats.roll(9);\nassert(block.number == 9);\n```\n\n#### Function Signature\n\n```solidity\nfunction roll(uint256) external;\n```\n\n### setNonce\n\n#### Description\n\nThe `setNonce` cheatcode will set the nonce of `account` to `nonce`. Note that the `nonce` must be strictly greater than\nthe current nonce\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Set nonce and verify (assume nonce before `setNonce` was less than 7)\naddress acc = address(msg.sender);\ncheats.setNonce(acc, 7);\nassert(cheats.getNonce(acc) == 7);\n```\n\n#### Function Signature\n\n```solidity\nfunction setNonce(address account, uint64 nonce) external;\n```\n\n### `sign`\n\n#### Description\n\nThe `sign` cheatcode will take in a private key `privateKey` and a hash digest `digest` to generate a `(v, r, s)`\nsignature\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\nbytes32 digest = keccak256(\"Data To Sign\");\n\n// Call cheats.sign\n(uint8 v, bytes32 r, bytes32 s) = cheats.sign(0x6df21769a2082e03f7e21f6395561279e9a7feb846b2bf740798c794ad196e00, digest);\naddress signer = ecrecover(digest, v, r, s);\nassert(signer == 0xdf8Ef652AdE0FA4790843a726164df8cf8649339);\n```\n\n#### Function Signature\n\n```solidity\nfunction sign(uint256 privateKey, bytes32 digest)\nexternal\nreturns (uint8 v, bytes32 r, bytes32 s);\n```\n\n### `snapshot` and `revertTo`\n\n#### Description\n\nThe `snapshot` cheatcode will take a snapshot of the current state of the blockchain and return an identifier for the\nsnapshot.\n\nOn the flipside, the `revertTo` cheatcode will revert the EVM state back based on the provided identifier.\n\n#### Example\n\n```solidity\ninterface CheatCodes {\n    function warp(uint256) external;\n\n    function deal(address, uint256) external;\n\n    function snapshot() external returns (uint256);\n\n    function revertTo(uint256) external returns (bool);\n}\n\nstruct Storage {\n    uint slot0;\n    uint slot1;\n}\n\ncontract TestContract {\n    Storage store;\n    uint256 timestamp;\n\n    function test() public {\n        // Obtain our cheat code contract reference.\n        CheatCodes cheats = CheatCodes(\n            0x7109709ECfa91a80626fF3989D68f67F5b1DD12D\n        );\n\n        store.slot0 = 10;\n        store.slot1 = 20;\n        timestamp = block.timestamp;\n        cheats.deal(address(this), 5 ether);\n\n        // Save state\n        uint256 snapshot = cheats.snapshot();\n\n        // Change state\n        store.slot0 = 300;\n        store.slot1 = 400;\n        cheats.deal(address(this), 500 ether);\n        cheats.warp(12345);\n\n        // Assert that state has been changed\n        assert(store.slot0 == 300);\n        assert(store.slot1 == 400);\n        assert(address(this).balance == 500 ether);\n        assert(block.timestamp == 12345);\n\n        // Revert to snapshot\n        cheats.revertTo(snapshot);\n\n        // Ensure state has been reset\n        assert(store.slot0 == 10);\n        assert(store.slot1 == 20);\n        assert(address(this).balance == 5 ether);\n        assert(block.timestamp == timestamp);\n    }\n}\n```\n\n### `store`\n\n#### Description\n\nThe `store` cheatcode will store `value` in storage slot `slot` for `account`\n\n#### Example\n\n```solidity\ncontract TestContract {\n    uint x = 123;\n    function test() public {\n        // Obtain our cheat code contract reference.\n        IStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n        // Store into x, verify it.\n        cheats.store(address(this), bytes32(uint(0)), bytes32(uint(456)));\n        assert(y == 456);\n    }\n}\n```\n\n#### Function Signature\n\n```solidity\nfunction store(address account, bytes32 slot, bytes32 value) external;\n```\n\n### warp\n\n#### Description\n\nThe `warp` cheatcode sets the `block.timestamp`\n\n#### Example\n\n```solidity\n// Obtain our cheat code contract reference.\nIStdCheats cheats = CheatCodes(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);\n\n// Change value and verify.\ncheats.warp(7);\nassert(block.timestamp == 7);\ncheats.warp(9);\nassert(block.timestamp == 9);\n```\n\n#### Function Signature\n\n```solidity\nfunction warp(uint256) external;\n```",
		},
		{
			Role:    "system",
			Content: "You are a smart contract auditing assistant. You'll generate fuzz tests to be run by `medusa`. Do not generate tests for `Foundry`, only generate tests to be run by `medusa`. Make use of medusa's cheatcodes where necessary",
		},
		{
			Role:    "system",
			Content: "You will be provided with main contracts which you should carefully examine to find possible vulnerabilities/invariants that should be tested for.",
		},
		{
			Role:    "system",
			Content: "Instructions given to you following preceded by 'Note:' should be taken into utmost importance.",
		},
	}
}
