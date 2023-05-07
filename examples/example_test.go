package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformExample(t *testing.T) {
	// list all example modules
	var moduleDirs []string
	entries, err := os.ReadDir(".")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, entry := range entries {
		if entry.IsDir() {
			moduleDirs = append(moduleDirs, entry.Name())
		}
	}

	// get env vars
	envVars := map[string]string{"TF_LOG_PROVIDER": "WARN"}
	for _, envVar := range []string{
		"GUANCE_ACCESS_TOKEN",
	} {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Environment variable %s must be set for acceptance tests", envVar)
		}
		envVars[envVar] = os.Getenv(envVar)
	}

	for _, name := range moduleDirs {
		t.Run(name, func(t *testing.T) {
			// retryable errors in terraform testing.
			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: name,
				EnvVars:      envVars,
			})

			// Clean up resources with "terraform destroy" at the end of the test.
			defer terraform.Destroy(t, terraformOptions)

			// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
			terraform.Apply(t, terraformOptions)
		})
	}
}
