package testpostgres

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PgContainer struct {
	username  string
	password  string
	dbName    string
	container testcontainers.Container
}

// Создаем контейнер PostgreSQL
func RunContainer(ctx context.Context, image string, dbName string, username string, password string, strategy *wait.LogStrategy) (*PgContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: strategy,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	instance := &PgContainer{
		username:  username,
		password:  password,
		dbName:    dbName,
		container: container,
	}

	return instance, nil
}

func (c PgContainer) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}

func (c PgContainer) ConnectionDsn(ctx context.Context, params string) (string, error) {
	host, err := c.container.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("Could not get container host: %v", err)
	}

	port, err := c.container.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("Could not get container port: %v", err)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", c.username, c.password, host, port.Port(), c.dbName, params), nil
}
