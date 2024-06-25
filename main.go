package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func main() {
	// Define flags
	configData := flag.String("config", "", "The config of secrets to retrieve")
	flag.Parse()

	// Validate flags
	if *configData == "" {
		log.Fatal("You must specify a config using --config")
	}

	data, err := base64.StdEncoding.DecodeString(*configData)
	if err != nil {
		log.Fatal("error:", err)
	}

	var secrets []map[string]string
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		log.Fatal("Failed to parse secret", err)
	}

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	// Create a Secrets Manager client
	client := secretsmanager.NewFromConfig(cfg)

	for _, secret := range secrets {
		// Retrieve the secret value
		secretValue, err := getSecretValue(client, secret["name"])
		if err != nil {
			log.Fatalf("failed to retrieve secret, %v", err)
		}

		// Write the secret value to the specified file
		err = os.WriteFile(secret["target"], []byte(fmt.Sprintf(secret["template"], secretValue)), 0644)
		if err != nil {
			log.Fatalf("failed to write secret to file, %v", err)
		}

		fmt.Printf("Secret value saved to %s\n", secret["target"])
	}
}

// getSecretValue retrieves the secret value from AWS Secrets Manager
func getSecretValue(client *secretsmanager.Client, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		var resourceNotFoundException *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundException) {
			return "", fmt.Errorf("the requested secret %s was not found", secretName)
		}
		return "", fmt.Errorf("failed to retrieve secret value: %v", err)
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	return "", fmt.Errorf("secret value is binary, which is not supported in this example")
}
