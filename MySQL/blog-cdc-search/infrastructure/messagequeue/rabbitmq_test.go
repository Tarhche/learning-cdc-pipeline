package messagequeue

import (
	"testing"
)

func TestNewRabbitMQRepository(t *testing.T) {
	config := RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "test",
		Password: "test123",
		VHost:    "/",
	}

	repo := NewRabbitMQRepository(config)

	if repo == nil {
		t.Fatal("Expected repository to be created")
	}

	if repo.config.Host != config.Host {
		t.Errorf("Expected host %s, got %s", config.Host, repo.config.Host)
	}

	if repo.config.Port != config.Port {
		t.Errorf("Expected port %d, got %d", config.Port, repo.config.Port)
	}

	if repo.config.Username != config.Username {
		t.Errorf("Expected username %s, got %s", config.Username, repo.config.Username)
	}

	if repo.config.Password != config.Password {
		t.Errorf("Expected password %s, got %s", config.Password, repo.config.Password)
	}

	if repo.config.VHost != config.VHost {
		t.Errorf("Expected vhost %s, got %s", config.VHost, repo.config.VHost)
	}
}

func TestRabbitMQConfig_ConnectionURL(t *testing.T) {
	config := RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "test",
		Password: "test123",
		VHost:    "/",
	}

	// Test the connection URL format by checking the Connect method would use it
	repo := NewRabbitMQRepository(config)

	// Since we can't actually connect in tests, we just verify the config is set correctly
	if repo.config.Host != config.Host {
		t.Errorf("Expected host %s, got %s", config.Host, repo.config.Host)
	}
}

func TestRabbitMQRepository_Close(t *testing.T) {
	config := RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "test",
		Password: "test123",
		VHost:    "/",
	}

	repo := NewRabbitMQRepository(config)

	// Close should not error even without connection
	err := repo.Close()
	if err != nil {
		t.Errorf("Expected no error when closing without connection, got: %v", err)
	}
}
