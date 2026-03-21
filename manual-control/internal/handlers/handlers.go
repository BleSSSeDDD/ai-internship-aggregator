package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BleSSSeDDD/ai-internship-aggregator/gen/go/vacancy"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/kafka"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	kafkaProducer *kafka.Producer
}

func NewHandlers(producer *kafka.Producer) *Handlers {
	return &Handlers{
		kafkaProducer: producer,
	}
}

// Index - главная страница с формой
func (h *Handlers) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// Submit - обработка отправки формы
func (h *Handlers) Submit(c *gin.Context) {
	// Парсим tech_stack (через запятую)
	techStackRaw := c.PostForm("tech_stack")
	techStack := strings.Split(techStackRaw, ",")
	var cleanTechStack []string
	for _, t := range techStack {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			cleanTechStack = append(cleanTechStack, trimmed)
		}
	}

	// Парсим зарплату
	minSalary, _ := strconv.Atoi(c.PostForm("min_salary"))

	// Создаем protobuf структуру
	internship := &vacancy.CompanyInternship{
		CompanyName:            c.PostForm("company_name"),
		SourceUrl:              c.PostForm("source_url"),
		SourceSite:             c.PostForm("source_site"),
		PositionName:           c.PostForm("position_name"),
		TechStack:              cleanTechStack,
		MinSalary:              int32(minSalary),
		Location:               c.PostForm("location"),
		InternshipDates:        c.PostForm("internship_dates"),
		SelectionProcess:       c.PostForm("selection_process"),
		Description:            c.PostForm("description"),
		ApplicationDeadline:    c.PostForm("application_deadline"),
		ContactInfo:            c.PostForm("contact_info"),
		ExperienceRequirements: c.PostForm("experience_requirements"),
	}

	// Валидация обязательных полей
	if internship.CompanyName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название компании обязательно"})
		return
	}
	if internship.PositionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название позиции обязательно"})
		return
	}
	if len(internship.TechStack) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Технологии обязательны"})
		return
	}

	// Отправляем в Kafka
	partition, offset, err := h.kafkaProducer.SendInternship(internship)
	if err != nil {
		log.Printf("❌ Ошибка отправки в Kafka: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Стажировка отправлена: %s, партиция: %d, оффсет: %d", internship.CompanyName, partition, offset)

	// Сохраняем в JSON для отображения
	internshipJSON, _ := json.MarshalIndent(internship, "", "  ")

	// Возвращаем успешный ответ с рендерингом страницы
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Success":        true,
		"Internship":     internship,
		"InternshipJSON": string(internshipJSON),
		"Partition":      partition,
		"Offset":         offset,
	})
}

// Health - проверка здоровья сервиса
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "admin-panel",
		"kafka":   "connected",
	})
}
