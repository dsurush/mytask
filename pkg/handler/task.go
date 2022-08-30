package handler

import (
	"fmt"
	"mytasks"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getList(c *gin.Context) {
	fmt.Println(len(h.cache.GetItems()))
	if len(h.cache.GetItems()) != 0 {
		fmt.Println("from cache")
		c.JSON(http.StatusOK, h.cache.GetItems())
		return
	}

	lists, err := h.services.TaskList.GetAll()
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("from not cache")
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

	if err = h.cache.Delete(int64(id)); err != nil {
		logrus.Error(err.Error())
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
	fmt.Println(len(h.cache.GetItems()))
	h.cache.Set(int64(id), input, 0)
	fmt.Println(len(h.cache.GetItems()))
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
