package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"fmt"
)

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	DOB string `json:"dob"`
}

var users = []User{
	{ID: "1", Name: "Alice", DOB: "1998-04-25"},
	{ID: "2", Name: "Bob", DOB: "1990-07-15"},
	{ID: "3", Name: "Charlie", DOB: "1985-12-10"},
}

func sendNotification(user User) {
	fmt.Printf("Starting notification for %s... \n", user.Name)
	time.Sleep(1 * time.Second)
	fmt.Printf("Notification sent to %s!\n", user.Name)
}

func calculateAge(dob string, refDate time.Time) (int,error) {
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0, err
	}
	years := refDate.Year() - parsedDOB.Year()
	if refDate.YearDay() < parsedDOB.YearDay() {
		years--
	}

	return years, nil
}

func main() {
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

				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Notification for user %s is being processed.", user.Name),
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	})

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		for _, user := range users {
			if user.ID == id {
				nextYear := time.Now().AddDate(1, 0, 0)
				ageNextYear, err := calculateAge(user.DOB, nextYear)

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date of birth format"})
					return
				}

				c.JSON(http.StatusOK,gin.H{
					"id": user.ID,
					"name": user.Name,
					"dob": user.DOB,
					"age_next_year": ageNextYear,
				})
				return
			}

		}

		c.JSON(http.StatusNotFound,gin.H{"message": "User not found"})
	})

	r.Run(":8080")
}