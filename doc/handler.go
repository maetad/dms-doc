package doc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocHandler struct {
	db *gorm.DB
}

func NewDocHandler(db *gorm.DB) *DocHandler {
	return &DocHandler{db: db}
}

func (u *DocHandler) GetDoc(c *gin.Context) {
	var docs []Doc

	r := u.db.Find(&docs)

	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, docs)
}

func (u *DocHandler) CreateDoc(c *gin.Context) {
	var doc Doc

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	r := u.db.Create(&doc)

	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (u *DocHandler) FindDoc(c *gin.Context) {
	var doc Doc

	r := u.db.First(&doc, c.Param("id"))

	if err := r.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, doc)
}

func (u *DocHandler) UpdateDoc(c *gin.Context) {
	var doc Doc

	r := u.db.First(&doc, c.Param("id"))

	if err := r.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	u.db.Save(&doc)

	c.JSON(http.StatusOK, doc)
}

func (u *DocHandler) DeleteDoc(c *gin.Context) {
	var doc Doc

	r := u.db.First(&doc, c.Param("id"))

	if err := r.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	u.db.Delete(&doc)

	c.JSON(http.StatusNoContent, gin.H{})
}
