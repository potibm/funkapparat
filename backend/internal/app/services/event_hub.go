package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/funkapparat/internal/app/exporter"
	"github.com/redis/go-redis/v9"
)

type ScheduleSource interface {
	GetAll(ctx context.Context) (domain.AnnouncementList, error)
	GetByID(ctx context.Context, id int64) (*domain.Announcement, error)
}

type EventHub struct {
	exporter *exporter.Manager
	redis    *redis.Client
	repo     ScheduleSource
	logger   *slog.Logger
}

func NewEventHub(e *exporter.Manager, redisClient *redis.Client, repo ScheduleSource) *EventHub {
	logger := slog.Default().With("component", "EventHub")

	return &EventHub{
		exporter: e,
		redis:    redisClient,
		repo:     repo,
		logger:   logger,
	}
}

func (h *EventHub) Publish(ctx context.Context, entryID int64, action ActionType) {
	if h.redis != nil {
		var eventDTO AnnouncementEventDTO

		if action == ActionDelete {
			eventDTO = AnnouncementEventDTO{
				Action:    action,
				Timestamp: time.Now().Unix(),
				Payload:   AnnouncementDTO{ID: entryID},
			}
		} else {
			entry, err := h.repo.GetByID(ctx, entryID)
			if err != nil {
				h.logger.Error("Failed to fetch announcement", "id", entryID, "error", err)

				return
			}

			eventDTO = mapToEventDTO(entry, action)
		}

		h.sendToStream(ctx, eventDTO)
	}

	h.exporter.Ping()
}

func (h *EventHub) PublishFullSync(ctx context.Context) {
	if h.redis == nil {
		return
	}

	timetable, err := h.repo.GetAll(ctx)
	if err != nil {
		h.logger.Error("Failed to fetch all announcements for sync", "error", err)

		return
	}

	syncEvent := AnnouncementSyncEventDTO{
		Action:    ActionSync,
		Timestamp: time.Now().Unix(),
		Payload:   mapToAnnouncementListDTO(timetable),
	}

	h.sendToStream(ctx, syncEvent)

	h.logger.Info("Sent full state sync event to Redis", "count", len(syncEvent.Payload))
}

func (h *EventHub) sendToStream(ctx context.Context, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal data for Redis", "error", err)

		return
	}

	err = h.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "party:news:events",
		Values: map[string]interface{}{"data": jsonData},
	}).Err()
	if err != nil {
		h.logger.Error("Redis XADD error", "error", err)
	}
}
