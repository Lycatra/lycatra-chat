package config

import (
    "os"
)

type MatrixConfig struct {
    Homeserver string
    AccessToken string
    RoomID string
    Approvers []string
}

type GithubConfig struct {
    Repos []string
    PollIntervalSeconds int
}

type DeployConfig struct {
    ComposeFile string
    Services map[string]string
    HealthcheckURL string
    RollbackOnFailure bool
}

type AppConfig struct {
    Matrix MatrixConfig
    Github GithubConfig
    Deploy DeployConfig
}

func FromEnv() AppConfig {
    // Skeleton: read minimal values from env; expand later
    return AppConfig{
        Matrix: MatrixConfig{
            Homeserver: os.Getenv("MATRIX_HOMESERVER"),
            AccessToken: os.Getenv("MATRIX_TOKEN"),
            RoomID: os.Getenv("MATRIX_ROOM_ID"),
        },
    }
}


