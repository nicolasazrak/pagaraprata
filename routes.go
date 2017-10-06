package main

import (
	"crypto/subtle"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

/* GET / renders index page */

func (s *Server) index(c *gin.Context) {
	content, err := s.loadFile("static/index.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sorry, algo salio mal"})
		return
	}

	c.Header("Content-type", "text/html; charset=utf-8")
	c.Writer.Write(content)
}

/* GET /version */

func (s *Server) version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": version})
}

/* POST /api/debts creates a new debt */

func (s *Server) createDebt(c *gin.Context) {
	debt := s.store.NewDebt(&DebtData{})
	err := s.store.SaveDebt(debt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("X-Created-Id", debt.ID)
	c.JSON(http.StatusCreated, debt)
}

/* GET /api/debts/:id returns a debt */

func (s *Server) getDebt(c *gin.Context) {
	debt, err := s.store.GetDebt(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if debt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, debt.Data)
}

/* PUT /api/debts/:id/key updates a debt */

func (s *Server) updateDebt(c *gin.Context) {
	debt, err := s.store.GetDebt(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if debt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	if subtle.ConstantTimeCompare([]byte(debt.Secret), []byte(c.Param("secret"))) != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid secret"})
		return
	}

	newData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read data"})
	}

	debt.Data = &DebtData{}
	err = json.Unmarshal(newData, debt.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.store.SaveDebt(debt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.getDebt(c)
}
