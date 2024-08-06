package mock

import (
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsManagerClient interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

type MockSecretsManager struct {
	secrets map[string]string
	lock    sync.RWMutex
}

func NewMockSecretsManager() *MockSecretsManager {
	return &MockSecretsManager{
		secrets: make(map[string]string),
	}
}

func (m *MockSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	secretID := *input.SecretId
	if secret, ok := m.secrets[secretID]; ok {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: &secret,
		}, nil
	}
	return nil, errors.New("secret not found")
}

func (m *MockSecretsManager) AddSecret(secretID, secretValue string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.secrets[secretID] = secretValue
}
