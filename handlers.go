package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

func GetPodsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		uc := GetOrCreateUserCache(userID, clientset)

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}

		uc.RLock()
		defer uc.RUnlock()

		start := (page - 1) * limit
		if start > len(uc.Data) {
			start = len(uc.Data)
		}
		end := start + limit
		if end > len(uc.Data) {
			end = len(uc.Data)
		}

		podsPage := uc.Data[start:end]

		c.JSON(http.StatusOK, gin.H{
			"page":  page,
			"limit": limit,
			"total": len(uc.Data),
			"pods":  podsPage,
		})
	}
}

func SearchPodsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		uc := GetOrCreateUserCache(userID, clientset)

		query := strings.ToLower(c.Query("q"))
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' required"})
			return
		}

		uc.RLock()
		defer uc.RUnlock()

		var results []Pod
		for _, pod := range uc.Data {
			if strings.Contains(strings.ToLower(pod.Name), query) ||
				strings.Contains(strings.ToLower(pod.Namespace), query) ||
				strings.Contains(strings.ToLower(pod.Status), query) {
				results = append(results, pod)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": results,
			"count":   len(results),
		})
	}
}
