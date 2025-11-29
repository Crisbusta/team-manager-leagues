package transporthttp

import (
	"net/http"

	"team-manager-leagues/internal/config"
	"team-manager-leagues/internal/middleware"
	"team-manager-leagues/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, svc *service.LeaguesService) *gin.Engine {
	r := gin.Default()

	// Middleware
	auth := middleware.AuthMiddleware(cfg)

	// Leagues
	leagues := r.Group("/leagues")
	leagues.Use(auth)
	{
		leagues.POST("", func(c *gin.Context) {
			var req struct {
				Name   string `json:"name"`
				Region string `json:"region"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			userID := c.GetString("userID")
			l, err := svc.CreateLeague(c.Request.Context(), userID, req.Name, req.Region)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"league": l})
		})

		leagues.GET("", func(c *gin.Context) {
			list, err := svc.ListLeagues(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"leagues": list})
		})

		leagues.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			l, err := svc.GetLeague(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if l == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"league": l})
		})

		// Series
		leagues.POST("/:id/series", func(c *gin.Context) {
			leagueID := c.Param("id")
			var req struct {
				Name   string `json:"name"`
				Format string `json:"format"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			s, err := svc.CreateSeries(c.Request.Context(), leagueID, req.Name, req.Format)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"series": s})
		})

		leagues.GET("/:id/series", func(c *gin.Context) {
			leagueID := c.Param("id")
			list, err := svc.ListSeries(c.Request.Context(), leagueID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"series": list})
		})
	}

	// Registrations
	regs := r.Group("/registrations")
	regs.Use(auth)
	{
		regs.POST("", func(c *gin.Context) {
			var req struct {
				TeamID   string `json:"teamId"`
				SeriesID string `json:"seriesId"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			userID := c.GetString("userID")
			reg, err := svc.RegisterTeam(c.Request.Context(), userID, req.TeamID, req.SeriesID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"registration": reg})
		})
	}
	return r
}
