package controllers

import (
	"net/http"

	dbpkg "github.com/shimastripe/go-api-sokushukai/db"
	"github.com/shimastripe/go-api-sokushukai/helper"
	"github.com/shimastripe/go-api-sokushukai/models"
	"github.com/shimastripe/go-api-sokushukai/version"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	preloads := c.DefaultQuery("preloads", "")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

	pagination := dbpkg.Pagination{}
	db, err := pagination.Paginate(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = dbpkg.SetPreloads(preloads, db)

	var users []models.User
	if err := db.Select("*").Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// paging
	var index int
	if len(users) < 1 {
		index = 0
	} else {
		index = int(users[len(users)-1].ID)
	}
	pagination.SetHeaderLink(c, index)

	if version.Range(ver, "<", "1.0.0") {
		// conditional branch by version.
		// this version < 1.0.0 !!
		c.JSON(400, gin.H{"error": "this version (< 1.0.0) is not supported!"})
		return
	}

	fieldMap := []map[string]interface{}{}
	for key, _ := range users {
		fieldMap = append(fieldMap, helper.FieldToMap(users[key], fields))
	}
	_, ok := c.GetQuery("pretty")
	if ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func GetUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")
	preloads := c.DefaultQuery("preloads", "")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(preloads, db)

	var user models.User
	if err := db.Select("*").First(&user, id).Error; err != nil {
		content := gin.H{"error": "user with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if version.Range(ver, "<", "1.0.0") {
		// conditional branch by version.
		// this version < 1.0.0 !!
		c.JSON(400, gin.H{"error": "this version (< 1.0.0) is not supported!"})
		return
	}

	fieldMap := helper.FieldToMap(user, fields)
	_, ok := c.GetQuery("pretty")
	if ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func CreateUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	var user models.User
	c.Bind(&user)
	if db.Create(&user).Error != nil {
		content := gin.H{"error": err.Error()}
		c.JSON(500, content)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, user)
}

func UpdateUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var user models.User
	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "user with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}
	c.Bind(&user)
	db.Save(&user)

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var user models.User
	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "user with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}
	db.Delete(&user)

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
