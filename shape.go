package minecraft

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Shapes = (*shapes)(nil)

// Shapes describes all the Shape related methods that the Minecraft
// API supports.
type Shapes interface {
	// List all the Shapes.
	List(ctx context.Context, options ShapeListOptions) (*ShapeList, error)

	// Create a new Shape with the given options.
	Create(ctx context.Context, options ShapeCreateOptions) (*Shape, error)

	// Read an Shape by its ID.
	Read(ctx context.Context, ShapeID string) (*Shape, error)

	// Update an Shape by its ID.
	Update(ctx context.Context, ShapeID string, options ShapeUpdateOptions) (*Shape, error)

	// Delete an Shape by its ID.
	Delete(ctx context.Context, ShapeID string) error
}

// shapes implements Shapes.
type shapes struct {
	client *Client
}

// ShapeList represents a list of Shapes.
type ShapeList struct {
	Items []*Shape
}

// Shape represents a Minecraft Shape.
type Shape struct {
	ID           string `json:"id"`
	*Location    `json:"location"`
	ShapeType    `json:"shapeType"`
	Material     string         `json:"material"`
	PreviousData []string       `json:"previousData"`
	Dimensions   interface{}    `json:dimensions`
	Status       ResourceStatus `json:"status"`
}

// ShapeType represents a shape type.
type ShapeType string

// List of available shape types.
const (
	ShapeTypeCube     ShapeType = "cube"
	ShapeTypeCylinder ShapeType = "cylinder"
)

// Dimensions is not being cast properly. This is a ghetto fix
func castDimensions(dimensions interface{}, shapeType ShapeType) interface{} {
	d := dimensions.(map[string]interface{})
	m = make(map[string]int)
	for key, value := range d {
		m[key] = value.(string)
	}
	switch shapeType {
	case ShapeTypeCube:
		return NewCubeDimensions(m["lengthX"], m["heightY"], m["widthZ"])
	case ShapeTypeCylinder:
		return NewCylinderDimensions(m["height"], m["radius"])
	}
}

// CubeDimensions is a Dimensions implementation
type CubeDimensions struct {
	LengthX int `json:"lengthX"`
	HeightY int `json:"heightY"`
	WidthZ  int `json:"widthZ"`
}

// NewCubeDimensions is a constructor for CubeDimensions
func NewCubeDimensions(lengthX int, heightY int, widthZ int) *CubeDimensions {
	return &CubeDimensions{LengthX: lengthX, HeightY: heightY, WidthZ: widthZ}
}

// CylinderDimensions is a Dimensions implementation
type CylinderDimensions struct {
	Height int `json:"height"`
	Radius int `json:"radius"`
}

// NewCylinderDimensions is a constructor for CubeDimensions
func NewCylinderDimensions(height int, radius int) *CylinderDimensions {
	return &CylinderDimensions{Height: height, Radius: radius}
}

// ShapeListOptions represents the options for listing Shapes.
type ShapeListOptions struct {
	Limit int
}

// List all Shapes.
func (s *shapes) List(ctx context.Context, options ShapeListOptions) (*ShapeList, error) {
	req, err := s.client.newRequest("GET", "shapes", &options)
	if err != nil {
		return nil, err
	}

	shapel := &ShapeList{}
	err = s.client.do(ctx, req, shapel)
	if err != nil {
		return nil, err
	}
	for _, shp := range shapel.Items {
		shp.Dimensions = castDimensions(shp.Dimensions)
	}
	return shapel, nil
}

// ShapeCreateOptions represents the options for creating a Shape.
type ShapeCreateOptions struct {
	*Location  `json:"location"`
	ShapeType  `json:"shapeType"`
	Material   string      `json:"material"`
	Dimensions interface{} `json:"dimensions"`
}

func (e ShapeCreateOptions) valid() error {
	if !notNil(&e.ShapeType) {
		return errors.New("shape type is required")
	}
	if !validString(&e.Material) {
		return errors.New("material is required")
	}
	if !notNil(&e.Dimensions) {
		return errors.New("dimensions is required")
	}
	if !validLocation(e.Location) {
		return errors.New("location is required")
	}
	return nil
}

// Create a new Shape with the given options.
func (s *shapes) Create(ctx context.Context, options ShapeCreateOptions) (*Shape, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("POST", "shapes", &options)
	if err != nil {
		return nil, err
	}

	shp := &Shape{}
	err = s.client.do(ctx, req, shp)
	if err != nil {
		return nil, err
	}
	shp.Dimensions = castDimensions(shp.Dimensions)
	return shp, nil
}

// Read an Shape by id.
func (s *shapes) Read(ctx context.Context, ShapeID string) (*Shape, error) {
	if !validStringID(&ShapeID) {
		return nil, errors.New("invalid value for ShapeID")
	}

	u := fmt.Sprintf("shapes/%s", url.QueryEscape(ShapeID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	shp := &Shape{}
	err = s.client.do(ctx, req, shp)
	if err != nil {
		return nil, err
	}
	shp.Dimensions = castDimensions(shp.Dimensions)
	return shp, nil
}

// ShapeUpdateOptions represents the options for updating a Shape.
type ShapeUpdateOptions struct {
	*Location  `json:"location"`
	ShapeType  `json:"shapeType"`
	Material   string      `json:"material"`
	Dimensions interface{} `json:dimensions`
}

// Update attributes of an existing Shape.
func (s *shapes) Update(ctx context.Context, ShapeID string, options ShapeUpdateOptions) (*Shape, error) {
	if !validStringID(&ShapeID) {
		return nil, errors.New("invalid value for ShapeID")
	}

	u := fmt.Sprintf("shapes/%s", url.QueryEscape(ShapeID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	shp := &Shape{}
	err = s.client.do(ctx, req, shp)
	if err != nil {
		return nil, err
	}
	shp.Dimensions = castDimensions(shp.Dimensions)
	return shp, nil
}

// Delete a Shape by its ID.
func (s *shapes) Delete(ctx context.Context, ShapeID string) error {
	if !validStringID(&ShapeID) {
		return errors.New("invalid value for ShapeID")
	}

	u := fmt.Sprintf("shapes/%s", url.QueryEscape(ShapeID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
