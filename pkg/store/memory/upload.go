package memory

import (
	"context"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

// IsUploadExists ...
func (m *Memory) IsUploadExists(ctx context.Context, user *api.User, filename, hash string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.uploads {
		if u.Filename == filename && u.Hash == hash {
			return true, nil
		}
	}

	return false, nil
}

// SaveNewUpload ...
func (m *Memory) SaveNewUpload(ctx context.Context, user *api.User, upload *api.Upload) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.uploads[upload.ID] = upload
	m.userUploads[upload.ID] = user.ID

	return nil
}

// LoadUpload ...
func (m *Memory) LoadUpload(ctx context.Context, user *api.User, id int64) (*api.Upload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, ok := m.uploads[id]
	if !ok {
		return nil, store.ErrNotFound
	}

	return u, nil
}

// LoadUploadByUUID ...
func (m *Memory) LoadUploadByUUID(ctx context.Context, user *api.User, uuid string) (*api.Upload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.uploads {
		if u.UUID == uuid {
			return u, nil
		}
	}

	return nil, store.ErrNotFound
}

// LoadUploads ...
func (m *Memory) LoadUploads(ctx context.Context, user *api.User) ([]*api.Upload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	uploads := []*api.Upload{}
	for uploadID, userID := range m.userUploads {
		if userID == user.ID {
			uploads = append(uploads, m.uploads[uploadID])
		}
	}

	return uploads, nil
}
