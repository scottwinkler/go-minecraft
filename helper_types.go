package minecraft

// Location represents an in game location
type Location struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Z     int    `json:"z"`
	World string `json:"world"`
}

// ResourceStatus represents a resource's status.
type ResourceStatus string

// List of available statuses.
const (
	ResourceStatusInitializing ResourceStatus = "initializing"
	ResourceStatusReady        ResourceStatus = "ready"
	ResourceStatusDeleting     ResourceStatus = "deleting"
	ResourceStatusUpdating     ResourceStatus = "updating"
)
