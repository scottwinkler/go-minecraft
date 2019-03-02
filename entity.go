package minecraft

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Entities = (*entities)(nil)

// Entities describes all the entity related methods that the Minecraft
// API supports.
type Entities interface {
	// List all the entities.
	List(ctx context.Context, options EntityListOptions) (*EntityList, error)

	// Create a new entity with the given options.
	Create(ctx context.Context, options EntityCreateOptions) (*Entity, error)

	// Read an entity by its ID.
	Read(ctx context.Context, entityID string) (*Entity, error)

	// Update an entity by its ID.
	Update(ctx context.Context, entityID string, options EntityUpdateOptions) (*Entity, error)

	// Delete an entity by its ID.
	Delete(ctx context.Context, entityID string) error
}

// entities implements Entities.
type entities struct {
	client *Client
}

// EntityList represents a list of entities.
type EntityList struct {
	Items []*Entity
}

// Entity represents a Minecraft entity.
type Entity struct {
	ID         string `json:"id"`
	*Location  `json:"location"`
	EntityType string         `json:"entityType"`
	CustomName string         `json:"customName"`
	Status     ResourceStatus `json:"status"`
}

// EntityListOptions represents the options for listing entities.
type EntityListOptions struct {
	Limit int
}

// List all entities.
func (s *entities) List(ctx context.Context, options EntityListOptions) (*EntityList, error) {
	req, err := s.client.newRequest("GET", "entities", &options)
	if err != nil {
		return nil, err
	}

	entityl := &EntityList{}
	err = s.client.do(ctx, req, entityl)
	if err != nil {
		return nil, err
	}

	return entityl, nil
}

// EntityCreateOptions represents the options for creating an entity.
type EntityCreateOptions struct {
	*Location  `json:"location"`
	EntityType string `json:"entityType"`
	CustomName string `json:"customName"`
}

func (e EntityCreateOptions) valid() error {
	if !validString(&e.EntityType) {
		return errors.New("entity type is required")
	}
	if !validString(&e.CustomName) {
		return errors.New("custom name is required")
	}
	if !validLocation(e.Location) {
		return errors.New("location is required")
	}
	return nil
}

// Create a new entity with the given options.
func (s *entities) Create(ctx context.Context, options EntityCreateOptions) (*Entity, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("POST", "entities", &options)
	if err != nil {
		return nil, err
	}

	ent := &Entity{}
	err = s.client.do(ctx, req, ent)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

// Read an entity by id.
func (s *entities) Read(ctx context.Context, entityID string) (*Entity, error) {
	if !validStringID(&entityID) {
		return nil, errors.New("invalid value for entityID")
	}

	u := fmt.Sprintf("entities/%s", url.QueryEscape(entityID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ent := &Entity{}
	err = s.client.do(ctx, req, ent)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

// EntityUpdateOptions represents the options for updating an entity.
type EntityUpdateOptions struct {
	*Location  `json:"location"`
	EntityType string `json:"entityType"`
	CustomName string `json:"customName"`
}

// Update attributes of an existing entity.
func (s *entities) Update(ctx context.Context, entityID string, options EntityUpdateOptions) (*Entity, error) {
	if !validStringID(&entityID) {
		return nil, errors.New("invalid value for entityID")
	}

	u := fmt.Sprintf("entities/%s", url.QueryEscape(entityID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	ent := &Entity{}
	err = s.client.do(ctx, req, ent)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

// Delete an entity by its ID.
func (s *entities) Delete(ctx context.Context, entityID string) error {
	if !validStringID(&entityID) {
		return errors.New("invalid value for entityID")
	}

	u := fmt.Sprintf("entities/%s", url.QueryEscape(entityID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
