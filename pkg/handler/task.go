package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mytasks"
	"net/http"
	"strconv"
)

func (h *Handler) getList(c *gin.Context) {
	lists, err := h.services.TaskList.GetAll()
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (h *Handler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.TaskList.Delete(id)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *Handler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	var input mytasks.UpdateTaskInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TaskList.Update(id, input)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
