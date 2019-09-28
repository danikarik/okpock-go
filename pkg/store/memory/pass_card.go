package memory

import (
	"context"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

// SaveNewPassCard ...
func (m *Memory) SaveNewPassCard(ctx context.Context, project *api.Project, passcard *api.PassCardInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.passCards[passcard.ID] = passcard
	m.projectPassCards[passcard.ID] = project.ID

	return nil
}

// LoadPassCard ...
func (m *Memory) LoadPassCard(ctx context.Context, project *api.Project, id int64) (*api.PassCardInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.passCards[id]
	if !ok {
		return nil, store.ErrNotFound
	}

	return p, nil
}

// LoadPassCardBySerialNumber ...
func (m *Memory) LoadPassCardBySerialNumber(ctx context.Context, project *api.Project, serialNumber string) (*api.PassCardInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.passCards {
		if p.Data != nil {
			if p.Data.SerialNumber == serialNumber {
				return p, nil
			}
		}
	}

	return nil, store.ErrNotFound
}

// LoadPassCards ...
func (m *Memory) LoadPassCards(ctx context.Context, project *api.Project, opts *api.PagingOptions) (*api.PassCardInfoList, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := []*api.PassCardInfo{}
	for passCardID, projectID := range m.projectPassCards {
		if projectID == project.ID {
			data = append(data, m.passCards[passCardID])
		}
	}

	return &api.PassCardInfoList{Opts: opts, Data: data}, nil
}

// LoadPassCardsByBarcodeMessage ...
func (m *Memory) LoadPassCardsByBarcodeMessage(ctx context.Context, project *api.Project, message string, opts *api.PagingOptions) (*api.PassCardInfoList, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := []*api.PassCardInfo{}
	for _, p := range m.passCards {
		if p.Data != nil {
			for _, barcode := range p.Data.Barcodes {
				if barcode.Message == message {
					data = append(data, p)
				}
			}
		}
	}

	return &api.PassCardInfoList{Opts: opts, Data: data}, nil
}

// UpdatePassCard ...
func (m *Memory) UpdatePassCard(ctx context.Context, data *api.PassCard, passcard *api.PassCardInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data.CopyFrom(passcard.Data)
	passcard.Data = data
	passcard.UpdatedAt = time.Now()
	m.passCards[passcard.ID] = passcard

	return nil
}
