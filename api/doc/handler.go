package doc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocHandler struct {
	repo Repository
}

func NewDocHandler(db *gorm.DB) *DocHandler {
	return &DocHandler{repo: CreateRepository(db)}
}

func (d *DocHandler) GetDoc(c *gin.Context) {
	docs, err := d.repo.Get()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, docs)
}

func (d *DocHandler) CreateDoc(c *gin.Context) {
	var doc Doc

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := d.repo.Create(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (d *DocHandler) FindDoc(c *gin.Context) {
	doc, err := d.repo.Find(c.Param("id"))

	if err != nil {
		var code int
		switch err.(type) {
		case *NotFoundError:
			code = http.StatusNotFound
		default:
			code = http.StatusBadRequest
		}

		c.JSON(code, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, doc)
}

func (d *DocHandler) UpdateDoc(c *gin.Context) {
	var doc Doc

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := d.repo.Update(c.Param("id"), &doc); err != nil {
		var code int
		switch err.(type) {
		case *NotFoundError:
			code = http.StatusNotFound
		default:
			code = http.StatusBadRequest
		}

		c.JSON(code, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, doc)
}

func (d *DocHandler) DeleteDoc(c *gin.Context) {
	if err := d.repo.Delete(c.Param("id")); err != nil {
		var code int
		switch err.(type) {
		case *NotFoundError:
			code = http.StatusNotFound
		default:
			code = http.StatusBadRequest
		}

		c.JSON(code, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
