package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/funkapparat/internal/app/exporter"
	"github.com/potibm/protokolapparat/pkg/common"
	"github.com/potibm/protokolapparat/pkg/news"
	"github.com/redis/go-redis/v9"
)

const streamName = "party:news:events"

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

func (h *EventHub) PublishCreate(ctx context.Context, entryID int64) {
	entry, err := h.getProtocolEntry(ctx, entryID)
	if err == nil {
		h.send(ctx, common.NewCreateEvent(entry))
	}
}

func (h *EventHub) PublishUpdate(ctx context.Context, entryID int64) {
	entry, err := h.getProtocolEntry(ctx, entryID)
	if err == nil {
		h.send(ctx, common.NewUpdateEvent(entry))
	}
}

func (h *EventHub) PublishDelete(ctx context.Context, entryID int64) {
	h.send(ctx, common.NewDeleteEvent(news.Entry{ID: entryID}))
}

func (h *EventHub) PublishFullSync(ctx context.Context) {
	if h.redis == nil {
		return
	}

	timetable, err := h.repo.GetAll(ctx)
	if err != nil {
		h.logger.Error("Failed to fetch timetable for sync", "error", err)

		return
	}

	mappedEntries := mapToEventListPayload(timetable)
	syncEvent := common.NewSyncEvent(mappedEntries)

	h.sendToStream(ctx, mappedEntries)

	h.logger.Info("Sent full state sync event to Redis", "count", len(syncEvent.Payload))
}

func (h *EventHub) getProtocolEntry(ctx context.Context, entryID int64) (news.Entry, error) {
	dbEntry, err := h.repo.GetByID(ctx, entryID)
	if err != nil {
		h.logger.Error("Failed to fetch schedule entry", "id", entryID, "error", err)

		return news.Entry{}, err
	}

	return mapToEventPayload(dbEntry), nil
}

func (h *EventHub) send(ctx context.Context, event common.Event[news.Entry]) {
	if h.redis == nil {
		return
	}

	if err := event.Validate(); err != nil {
		h.logger.Error("Tried to publish invalid event", "error", err, "action", event.Action)

		return
	}

	h.sendToStream(ctx, event)
	h.exporter.Ping()
}

func (h *EventHub) sendToStream(ctx context.Context, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal data for Redis", "error", err)

		return
	}

	err = h.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{"data": jsonData},
	}).Err()
	if err != nil {
		h.logger.Error("Redis XADD error", "error", err)
	}
}
