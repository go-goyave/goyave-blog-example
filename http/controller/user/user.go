package user

import (
	"net/http"

	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/mitchellh/go-homedir"
	"goyave.dev/goyave/v4"
	"goyave.dev/goyave/v4/database"
	"goyave.dev/goyave/v4/util/fsutil"
)

var (
	// StoragePath the path used to store user profile pictures
	StoragePath string = ""
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		goyave.ErrLogger.Fatal(err)
	}
	StoragePath = home + "/storage/"
}

// Register insert a new user in the database
func Register(response *goyave.Response, request *goyave.Request) {
	user := &model.User{
		Email:    request.String("email"),
		Username: request.String("username"),
		Password: request.String("password"),
	}
	if request.Has("image") { // image is nullable
		image := request.File("image")[0]
		user.Image.String = image.Save(StoragePath, user.Username+"-"+image.Header.Filename)
		user.Image.Valid = true
	}

	if err := database.Conn().Create(user).Error; err != nil {
		if user.Image.Valid {
			fsutil.Delete(StoragePath + user.Image.String)
		}
		response.Error(err)
	} else {
		response.JSON(http.StatusCreated, map[string]uint{"id": user.ID})
	}
}

// Show returns the authenticated user
func Show(response *goyave.Response, request *goyave.Request) {
	response.JSON(http.StatusOK, request.User)
}

// Image returns the profile picture of the authenticated user.
// A default profile picture is sent if the user doesn't have a profile picture.
func Image(response *goyave.Response, request *goyave.Request) {
	user := model.User{}
	result := database.Conn().First(&user, request.Params["id"])
	if response.HandleDatabaseError(result) {
		path := ""
		if user.Image.Valid {
			path = StoragePath + user.Image.String
		} else {
			path = "resources/img/default_profile_picture.png"
		}
		if err := response.File(path); err != nil {
			response.Error(err)
		}
	}
}

// Update replaces the record of the authenticated user
// If the profile picture is modified, the previous one is deleted.
func Update(response *goyave.Response, request *goyave.Request) {
	db := database.Conn()
	user := request.User.(*model.User)
	if request.Has("image") {
		path := StoragePath + user.Image.String
		if user.Image.Valid && fsutil.FileExists(path) {
			fsutil.Delete(path)
		}

		if request.Data["image"] != nil {
			image := request.File("image")[0]
			actualName := image.Save(StoragePath, user.Username+"-"+image.Header.Filename)
			request.Data["image"] = actualName
		}
	}

	updates := map[string]interface{}{}
	for c := range UpdateRequest {
		if request.Has(c) {
			updates[c] = request.Data[c]
		}
	}

	if len(updates) <= 0 {
		return
	}

	if err := db.Model(request.User).Updates(updates).Error; err != nil {
		response.Error(err)
	}
}
