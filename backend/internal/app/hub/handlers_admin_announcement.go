package hub

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/funkapparat/internal/app/repository"
)

func (s *Server) listAnnouncements(c *gin.Context) {
	params := parseAnnouncementListParams(c)

	filters := parseAnnouncementListFilters(c)

	annoucenments, total, err := s.announcementRepo.List(c.Request.Context(), params, filters)
	if err != nil {
		respondWithInternalServerProblem(c, "Failed to list announcements: "+err.Error())

		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	c.JSON(http.StatusOK, annoucenments)
}

func (s *Server) getAnnouncement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	announcement, err := s.announcementRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondWithNotFoundProblem(c, "Announcement with ID "+strconv.FormatInt(id, 10)+" not found")

		return
	}

	c.JSON(http.StatusOK, announcement)
}

func (s *Server) createAnnouncement(c *gin.Context) {
	var announcement domain.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	if err := s.announcementRepo.Save(c.Request.Context(), &announcement); err != nil {
		slog.Error("Create Announcement: Failed to create announcement", "error", err)
		respondWithInternalServerProblem(c, "Failed to create announcement: "+err.Error())

		return
	} else {
		s.eventHub.PublishCreate(c, announcement.ID)

		slog.Info("Create Announcement: Successfully created announcement", "id", announcement.ID)
	}

	c.JSON(http.StatusCreated, announcement)
}

func (s *Server) updateAnnouncement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	var announcement domain.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	announcement.ID = id

	if err := s.announcementRepo.Save(c.Request.Context(), &announcement); err != nil {
		slog.Error("Update Announcement: Failed to update announcement", "id", id, "error", err)
		respondWithInternalServerProblem(c, "Failed to update announcement: "+err.Error())

		return
	} else {
		s.eventHub.PublishUpdate(c, announcement.ID)

		slog.Info("Update Announcement: Successfully updated announcement", "id", id)
	}

	c.JSON(http.StatusOK, announcement)
}

func (s *Server) deleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	if err := s.announcementRepo.Delete(c.Request.Context(), id); err != nil {
		respondWithInternalServerProblem(c, "Failed to delete announcement: "+err.Error())

		return
	}

	s.eventHub.PublishDelete(c, id)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func parseAnnouncementListParams(c *gin.Context) repository.AnnouncementListParams {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "20"))

	return repository.AnnouncementListParams{
		Offset: start,
		Limit:  end - start,
		Sort:   c.DefaultQuery("_sort", "id"),
		Order:  c.DefaultQuery("_order", "DESC"),
	}
}

func parseAnnouncementListFilters(c *gin.Context) repository.AnnouncementListFilters {
	filters := repository.AnnouncementListFilters{}

	if q := c.Query("q"); q != "" {
		filters.Query = &q
	}

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			filters.ID = &id
		}
	}

	return filters
}
