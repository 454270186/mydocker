package container

import "github.com/google/uuid"

// GetUUID generate a uuid for default container name
func GetUUID() string {
	return uuid.NewString()
}