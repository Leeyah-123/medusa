package llm_assistant

import (
	"fmt"
	"github.com/crytic/medusa/fuzzing/contracts"
	"github.com/crytic/medusa/utils"
	"os"
	"os/exec"
)

// Stores all prompts and responses from OpenAI
var messages = TrainingPrompts()

func GenerateFuzzingHarness(contractDefinitions contracts.Contracts) error {
	// TODO: Accept main contracts instead, compile, then populate config with generated test contracts

	// Create a test file for each contract definition
	err := createTestFiles(contractDefinitions)
	if err != nil {
		return err
	}

	// Generate the fuzzing harness
	err = generateFuzzingHarness(contractDefinitions)
	if err != nil {
		return err
	}

	return nil
}

func generateFuzzingHarness(contractDefinitions contracts.Contracts) error {
	fmt.Println("Generating fuzzing harness...")

	// Obtain the source code for each contract definition
	for _, contractDefinition := range contractDefinitions {
		fmt.Println("Generating fuzzing harness for", contractDefinition.Name())
		testFilePath := generateTestFilePath(contractDefinition.SourcePath())

		// Read the contract source code
		contractSourceCode, err := os.ReadFile(contractDefinition.SourcePath())
		if err != nil {
			return err
		}

		// Read the test contract source code
		testContractSourceCode, err := os.ReadFile(testFilePath)
		if err != nil {
			return err
		}

		message := Message{
			Role:    "user",
			Content: GenerateFuzzHarnessPrompt(contractDefinition.SourcePath(), testFilePath, string(contractSourceCode), string(testContractSourceCode), contractDefinition.Name(), generateTestContractName(contractDefinition.Name())),
		}

		// Store prompt
		messages = append(messages, message)

		// Generate the fuzzing harness
		response, err := AskGPT4Turbo(messages)
		if err != nil {
			return err
		}
		fmt.Println("GPT 4 Response for", contractDefinition.Name(), "is", response)

		// Store response
		messages = append(messages, Message{
			Role:    "system",
			Content: response,
		})

		fmt.Println("Generated fuzzing harness for", contractDefinition.Name())
		// Write response to test file
		err = os.WriteFile(testFilePath, []byte(response), 0644)
		if err != nil {
			return err
		}

		for {
			// Validate generated test file
			stdErr, err := validateTestFile(testFilePath)
			if err == nil {
				break
			}

			fmt.Println("Regenerating test file due to error", string(stdErr))

			message := Message{
				Role:    "user",
				Content: RegenerateFuzzHarnessPrompt(string(stdErr)),
			}
			response, err := processMessageWithGPT4Turbo(message)
			if err != nil {
				return err
			}

			// Write response to test file
			err = os.WriteFile(generateTestFilePath(contractDefinition.SourcePath()), []byte(response), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processMessageWithGPT4Turbo(message Message) (string, error) {
	messages = append(messages, message)

	// Generate the fuzzing harness
	response, err := AskGPT4Turbo(messages)
	if err != nil {
		return "", err
	}

	// Store response
	messages = append(messages, Message{
		Role:    "system",
		Content: response,
	})

	return response, nil
}

func validateTestFile(testFilePath string) ([]byte, error) {
	command := exec.Command("crytic-compile", testFilePath, "--ignore-compile")

	_, stdErr, _, err := utils.RunCommandWithOutputAndError(command)

	return stdErr, err
}

func createTestFiles(contractDefinitions contracts.Contracts) error {
	// Create a test file for each contract definition if not exists
	for _, contractDefinition := range contractDefinitions {
		// Generate test file path using source path and timestamp of current time
		testFilePath := generateTestFilePath(contractDefinition.SourcePath())

		// Create the file if it does not exist
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			_, err = os.Create(testFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
