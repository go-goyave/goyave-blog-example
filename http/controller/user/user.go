package user

import (
	"net/http"

	"github.com/System-Glitch/goyave-blog-example/database/model"
	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/database"
)

// Register insert a new user in the database
func Register(response *goyave.Response, request *goyave.Request) {
	user := &model.User{
		Email:    request.String("email"),
		Username: request.String("username"),
		Password: request.String("password"),
	}
	if request.Has("image") { // image is nullable
		user.Image.String = request.String("image")
		user.Image.Valid = true
	}

	if err := database.Conn().Create(user).Error; err != nil {
		response.Error(err)
	} else {
		response.JSON(http.StatusCreated, map[string]uint{"id": user.ID})
	}
}

// Show returns the authenticated user
func Show(response *goyave.Response, request *goyave.Request) {
	response.JSON(http.StatusOK, request.User)
}

// Update replaces the record of the authenticated user
func Update(response *goyave.Response, request *goyave.Request) {
	db := database.Conn()
	if err := db.Model(request.User).Updates(request.Data).Error; err != nil {
		response.Error(err)
	}
}
