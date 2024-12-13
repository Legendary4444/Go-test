package main
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func TestNotifyEndpoint(t *testing.T) {
	r := gin.Default()

	r.POST("/notify", func(c *gin.Context) {
		var payload struct {
			UserID string `json:"user_id"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		for _, user := range users {
			if user.ID == payload.UserID {
				go sendNotification(user)

				c.JSON(http.StatusOK,gin.H{
					"message": "Notification for user " + user.Name + " is being processed.",
				})
				return
			}
		}

		c.JSON(http.StatusNotFound,gin.H{"message": "User not found"})
	})

	tests := []struct {
		name string
		payload map[string]string
		expectedStatus int
		expectedBody string
	}{
		{
			name: "Valid User",
			payload: map[string]string{
				"user_id": "1",
			},
			expectedStatus: http.StatusOK,
			expectedBody: "Notification for user Alice is being processed.",
		},
		{
			name: "Invalid Payload",
			payload: map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody: "Invalid request payload",
		},
		{
			name: "User Not Found",
			payload: map[string]string{
				"user_id": "4",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name,func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)

			req, _ := http.NewRequest("POST", "/notify", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t,tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(),tt.expectedBody)
		})
	}
}