package exporter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/potibm/funkapparat/internal/app/domain"
)

const defaultTimeout = 30 * time.Second

type Exporter interface {
	Name() string
	Export(ctx context.Context, announcements domain.AnnouncementList) error
}

type Manager struct {
	exporters    []Exporter
	db           AllLoader
	debounceTime time.Duration
	timer        *time.Timer
	mu           sync.Mutex
	logger       *slog.Logger
}

type AllLoader interface {
	GetAll(ctx context.Context) (domain.AnnouncementList, error)
}

func NewManager(source AllLoader, debounce time.Duration) *Manager {
	logger := slog.Default()

	return &Manager{
		exporters:    []Exporter{},
		db:           source,
		debounceTime: debounce,
		logger:       logger.With("component", "Exporter"),
	}
}

func (m *Manager) Register(e Exporter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.exporters = append(m.exporters, e)
}

func (m *Manager) Ping() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.timer != nil {
		m.timer.Stop()
	}

	m.timer = time.AfterFunc(m.debounceTime, func() {
		m.RunAll()
	})
}

func (m *Manager) RunAll() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	announcements, err := m.db.GetAll(ctx)
	if err != nil {
		m.logger.Error("Error fetching announcements", "error", err)

		return
	}

	var wg sync.WaitGroup
	for _, e := range m.exporters {
		wg.Add(1)

		go func(exp Exporter) {
			defer wg.Done()

			m.logger.Info("Starting", "exporter", exp.Name())

			if err := exp.Export(ctx, announcements); err != nil {
				m.logger.Error("Failed", "exporter", exp.Name(), "error", err)
			} else {
				m.logger.Info("Finished successfully", "exporter", exp.Name())
			}
		}(e)
	}

	wg.Wait()
}
