package memory

import (
	"context"
	"sync"
	"time"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// Memory implements DatastoreInterface using in-memory storage
type Memory struct {
	applications map[uint]*model.Application
	nameIndex    map[string]uint
	nextID       uint
	mutex        sync.RWMutex
}

// New creates a new Memory datastore instance
func New() (datastore.DatastoreInterface, error) {
	logger.Info("Initialized in-memory datastore")

	return &Memory{
		applications: make(map[uint]*model.Application),
		nameIndex:    make(map[string]uint),
		nextID:       1,
	}, nil
}

// CreateApplication creates a new application
func (m *Memory) CreateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if name already exists
	if _, exists := m.nameIndex[app.Name]; exists {
		return nil, datastore.ErrDuplicateKey
	}

	// Set ID and timestamps
	app.ID = m.nextID
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	m.nextID++

	// Store application
	m.applications[app.ID] = app
	m.nameIndex[app.Name] = app.ID

	return app, nil
}

// GetApplicationByID retrieves an application by ID
func (m *Memory) GetApplicationByID(ctx context.Context, id uint) (*model.Application, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	app, exists := m.applications[id]
	if !exists {
		return nil, datastore.ErrNotFound
	}

	return app, nil
}

// GetApplicationByName retrieves an application by name
func (m *Memory) GetApplicationByName(ctx context.Context, name string) (*model.Application, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	id, exists := m.nameIndex[name]
	if !exists {
		return nil, datastore.ErrNotFound
	}

	app := m.applications[id]
	return app, nil
}

// ListApplications retrieves a paginated list of applications
func (m *Memory) ListApplications(ctx context.Context, page, pageSize int) ([]*model.Application, int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := int64(len(m.applications))

	// Convert map to slice
	apps := make([]*model.Application, 0, len(m.applications))
	for _, app := range m.applications {
		apps = append(apps, app)
	}

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(apps) {
		return []*model.Application{}, total, nil
	}

	if end > len(apps) {
		end = len(apps)
	}

	paginatedApps := apps[start:end]
	return paginatedApps, total, nil
}

// UpdateApplication updates an existing application
func (m *Memory) UpdateApplication(ctx context.Context, app *model.Application) (*model.Application, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if application exists
	existing, exists := m.applications[app.ID]
	if !exists {
		return nil, datastore.ErrNotFound
	}

	// Check if name changed and new name already exists
	if existing.Name != app.Name {
		if _, nameExists := m.nameIndex[app.Name]; nameExists {
			return nil, datastore.ErrDuplicateKey
		}
		// Update name index
		delete(m.nameIndex, existing.Name)
		m.nameIndex[app.Name] = app.ID
	}

	// Update timestamps
	app.CreatedAt = existing.CreatedAt
	app.UpdatedAt = time.Now()

	// Store updated application
	m.applications[app.ID] = app

	return app, nil
}

// DeleteApplication deletes an application by ID
func (m *Memory) DeleteApplication(ctx context.Context, id uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	app, exists := m.applications[id]
	if !exists {
		return datastore.ErrNotFound
	}

	// Remove from both maps
	delete(m.applications, id)
	delete(m.nameIndex, app.Name)

	return nil
}

// Migrate runs database migrations (no-op for memory)
func (m *Memory) Migrate() error {
	logger.Info("Memory datastore migration completed (no-op)")
	return nil
}

// Close closes the datastore (no-op for memory)
func (m *Memory) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear all data
	m.applications = make(map[uint]*model.Application)
	m.nameIndex = make(map[string]uint)
	m.nextID = 1

	logger.Info("Memory datastore closed")
	return nil
}

// HealthCheck checks the datastore health (always healthy for memory)
func (m *Memory) HealthCheck() error {
	return nil
}
