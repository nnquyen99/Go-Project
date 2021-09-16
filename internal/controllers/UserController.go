package controllers

import (
	"demo3-gin/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type User struct {
	Id          string     `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"id" validate:"required,max=40"`
	Username    string     `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" validate:"required,username,max=100"`
	Email       string     `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" validate:"email,max=100"`
	Phone       string     `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" validate:"required,phone,max=18"`
	DateOfBirth *time.Time `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth"`
}

func GetAll(c *gin.Context) {
	db := models.ConnectDataBase()
	query := "select * from users"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users not found"})
		return
	}
	users := make([]models.Users, 0)
	for rows.Next() {
		var s models.Users
		if err := rows.Scan(&s.Id, &s.Username, &s.Email, &s.Phone, &s.DateOfBirth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rows not match"})
			return
		}
		users = append(users, s)
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
	defer db.Close()
}

func GetById(c *gin.Context) {
	db := models.ConnectDataBase()
	id := c.Param("id")
	declare := "declare @p1 nvarchar(40)  = '" + id + "' select * from users where id = @p1"
	rows, err := db.Query(declare)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messages": "Story not found"})
		return
	}
	user := User{}
	for rows.Next() {
		var id, username, email, phone string
		var dateOfBirth *time.Time
		err = rows.Scan(&id, &username, &email, &phone, &dateOfBirth)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Username = username
		user.Email = email
		user.Phone = phone
		user.DateOfBirth = dateOfBirth
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
	defer db.Close()
}

func Insert(c *gin.Context) {
	db := models.ConnectDataBase()
	var insert User
	err := c.ShouldBindJSON(&insert)
	if err == nil {
		createInput, err1 := db.Prepare("insert into users (id, username, email, phone, date_of_birth) values (@p1, @p2, @p3, @p4, @p5)")
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error db"})
		}
		_, err2 := createInput.Exec(insert.Id, insert.Username, insert.Email, insert.Phone, insert.DateOfBirth)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not match values"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "create success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": insert})
	defer db.Close()
}

func Update(c *gin.Context) {
	db := models.ConnectDataBase()
	var update User
	id := c.Param("id")
	err := c.ShouldBindJSON(&update)
	if err == nil {
		edit, err1 := db.Prepare("update users set username = @p1, email = @p2, phone = @p3, date_of_birth = @p4 where id = '" + id + "'")
		if err1 != nil {
			panic(err.Error())
		}
		_, err2 := edit.Exec(update.Username, update.Email, update.Phone, update.DateOfBirth)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "update fail: not match values"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Update success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": update})
	defer db.Close()
}

func Delete(c *gin.Context) {
	db := models.ConnectDataBase()
	delete, err := db.Prepare("delete from users where id = @p1")
	if err != nil {
		panic(err.Error())
	}
	delete.Exec(c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message ": " delete success"})
	c.JSON(http.StatusOK, gin.H{"data": true})
}
