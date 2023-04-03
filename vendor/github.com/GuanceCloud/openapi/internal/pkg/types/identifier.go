package types

import (
	"fmt"
	"strings"
)

// Identifier is a resource identifier.
type Identifier struct {
	ResourceType string
	PrimaryId    string
}

// String returns the identifier as a string.
func (i Identifier) String() string {
	return i.ResourceType + "/" + i.PrimaryId
}

// ParseIdentifier parses an identifier from a string.
func ParseIdentifier(id string) (Identifier, error) {
	tokens := strings.Split(id, "/")
	if len(tokens) != 2 {
		return Identifier{}, fmt.Errorf("invalid identifier: %s", id)
	}
	return Identifier{
		ResourceType: tokens[0],
		PrimaryId:    tokens[1],
	}, nil
}
