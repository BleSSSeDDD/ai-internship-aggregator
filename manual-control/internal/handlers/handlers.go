package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/BleSSSeDDD/ai-internship-aggregator/manual-control/gen/go/vacancy"
	"github.com/BleSSSeDDD/ai-internship-aggregator/manual-control/internal/kafka"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	kafkaProducer kafka.Publisher
}

func NewHandlers(producer kafka.Publisher) *Handlers {
	return &Handlers{
		kafkaProducer: producer,
	}
}

func (h *Handlers) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handlers) Submit(c *gin.Context) {
	minSalary, _ := strconv.Atoi(c.PostForm("min_salary"))

	internship := &vacancy.CompanyInternship{
		CompanyName:            c.PostForm("company_name"),
		SourceUrl:              c.PostForm("source_url"),
		SourceSite:             c.PostForm("source_site"),
		PositionName:           c.PostForm("position_name"),
		TechStack:              splitTechStack(c.PostForm("tech_stack")),
		MinSalary:              int32(minSalary),
		Location:               c.PostForm("location"),
		InternshipDates:        c.PostForm("internship_dates"),
		SelectionProcess:       c.PostForm("selection_process"),
		Description:            c.PostForm("description"),
		ApplicationDeadline:    c.PostForm("application_deadline"),
		ContactInfo:            c.PostForm("contact_info"),
		ExperienceRequirements: c.PostForm("experience_requirements"),
	}

	if err := validate(internship); err != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	partition, offset, err := h.kafkaProducer.SendInternship(internship)
	if err != nil {
		slog.Error("failed to send internship", "company", internship.CompanyName, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	internshipJSON, _ := json.MarshalIndent(internship, "", "  ")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Success":        true,
		"Internship":     internship,
		"InternshipJSON": string(internshipJSON),
		"Partition":      partition,
		"Offset":         offset,
	})
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "manual-control",
	})
}

// validate возвращает текст ошибки для формы или "", если всё заполнено.
func validate(in *vacancy.CompanyInternship) string {
	switch {
	case in.CompanyName == "":
		return "Название компании обязательно"
	case in.PositionName == "":
		return "Название позиции обязательно"
	case len(in.TechStack) == 0:
		return "Технологии обязательны"
	}
	return ""
}

// splitTechStack разбирает строку "Go, Kafka, gRPC" в срез без пустых элементов.
func splitTechStack(raw string) []string {
	var result []string
	for _, item := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
