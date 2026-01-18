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
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			userID := c.GetString("userID")
			l, err := svc.CreateLeague(c.Request.Context(), userID, req.Name, req.Region)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"league": l})
		})

		leagues.GET("", func(c *gin.Context) {
			list, err := svc.ListLeagues(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"leagues": list})
		})

		leagues.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			l, err := svc.GetLeague(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if l == nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"league": l})
		})

		leagues.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var req struct {
				Name   string `json:"name"`
				Region string `json:"region"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			l, err := svc.UpdateLeague(c.Request.Context(), id, req.Name, req.Region)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"league": l})
		})

		leagues.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := svc.DeleteLeague(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		// Series
		leagues.POST("/:id/series", func(c *gin.Context) {
			leagueID := c.Param("id")
			var req struct {
				Name   string `json:"name"`
				Format string `json:"format"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			s, err := svc.CreateSeries(c.Request.Context(), leagueID, req.Name, req.Format)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"series": s})
		})

		leagues.GET("/:id/series", func(c *gin.Context) {
			leagueID := c.Param("id")
			list, err := svc.ListSeries(c.Request.Context(), leagueID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"series": list})
		})

		leagues.PUT("/:id/series/:seriesId", func(c *gin.Context) {
			var req struct {
				Name   string `json:"name"`
				Format string `json:"format"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			seriesID := c.Param("seriesId")
			if err := svc.UpdateSeries(c.Request.Context(), seriesID, req.Name, req.Format); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			s, err := svc.GetSeries(c.Request.Context(), seriesID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"series": s})
		})

		leagues.DELETE("/:id/series/:seriesId", func(c *gin.Context) {
			seriesID := c.Param("seriesId")
			if err := svc.DeleteSeries(c.Request.Context(), seriesID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true})
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
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			userID := c.GetString("userID")
			reg, err := svc.RegisterTeam(c.Request.Context(), userID, req.TeamID, req.SeriesID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"registration": reg})
		})

		regs.GET("", func(c *gin.Context) {
			teamID := c.Query("teamId")
			seriesID := c.Query("seriesId")
			if teamID == "" && seriesID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "teamId or seriesId required"})
				return
			}
			if teamID != "" {
				list, err := svc.ListRegistrationsByTeam(c.Request.Context(), teamID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"registrations": list})
				return
			}
			list, err := svc.ListRegistrationsBySeries(c.Request.Context(), seriesID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"registrations": list})
		})

		regs.PUT("/:id", func(c *gin.Context) {
			var req struct {
				Status string `json:"status"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			regID := c.Param("id")
			if err := svc.UpdateRegistrationStatus(c.Request.Context(), regID, req.Status); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		regs.DELETE("/:id", func(c *gin.Context) {
			regID := c.Param("id")
			if err := svc.DeleteRegistration(c.Request.Context(), regID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true})
		})
	}
	return r
}
