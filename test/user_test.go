package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-goyave/goyave-blog-example/database/model"
	userController "github.com/go-goyave/goyave-blog-example/http/controller/user"
	"github.com/go-goyave/goyave-blog-example/http/route"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/System-Glitch/goyave/v3/helper/filesystem"
	"github.com/System-Glitch/goyave/v3/validation"

	_ "github.com/System-Glitch/goyave/v3/database/dialect/mysql"
	_ "github.com/go-goyave/goyave-blog-example/http/validation"
)

type UserTestSuite struct {
	goyave.TestSuite
}

func (suite *UserTestSuite) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (suite *UserTestSuite) SetupTest() {
	suite.ClearDatabase()
}

func (suite *UserTestSuite) TestRegister() {
	suite.RunServer(route.Register, func() {
		headers := map[string]string{"Content-Type": "application/json"}
		request := map[string]interface{}{
			"username": "jack",
			"email":    "jack@example.org",
			"password": "super_Secret_password_2",
		}
		body, _ := json.Marshal(request)
		resp, err := suite.Post("/user", headers, bytes.NewReader(body))
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusCreated, resp.StatusCode)
			json := map[string]interface{}{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				suite.Contains(json, "id")
			}

			count := int64(0)
			if err := database.Conn().Model(&model.User{}).Where("email = ?", "jack@example.org").Count(&count).Error; err != nil {
				suite.Error(err)
			}
			suite.Equal(int64(1), count)
		}
	})
}

func (suite *UserTestSuite) TestRegisterValidationError() {
	suite.RunServer(route.Register, func() {
		resp, err := suite.Post("/user", nil, nil)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil { // Expect validation errors (field "username", "email" and "password" are required)
			defer resp.Body.Close()
			suite.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
			json := map[string]validation.Errors{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				suite.Contains(json["validationError"], "username")
				suite.Contains(json["validationError"], "email")
				suite.Contains(json["validationError"], "password")
			}
		}
	})
}

func (suite *UserTestSuite) TestRegisterNotUnique() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.UserGenerator)
		override := &model.User{
			Username: "jack",
			Email:    "jack@example.org",
		}
		factory.Override(override).Save(1)

		headers := map[string]string{"Content-Type": "application/json"}
		request := map[string]interface{}{
			"username": override.Username,
			"email":    override.Email,
			"password": "super_Secret_password_2",
		}
		body, _ := json.Marshal(request)
		resp, err := suite.Post("/user", headers, bytes.NewReader(body))
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
			json := map[string]validation.Errors{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				suite.Contains(json["validationError"], "username")
				suite.Contains(json["validationError"], "email")
			}
		}
	})
}

func (suite *UserTestSuite) TestRegisterWithImage() {
	suite.RunServer(route.Register, func() {
		const path = "resources/test/img/goyave_64.png"
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		suite.WriteField(writer, "email", "jack@example.org")
		suite.WriteField(writer, "username", "jack")
		suite.WriteField(writer, "password", "super_Secret_password_2")
		suite.WriteFile(writer, path, "image", filepath.Base(path))
		if err := writer.Close(); err != nil {
			suite.Error(err)
			return
		}
		headers := map[string]string{"Content-Type": writer.FormDataContentType()}

		resp, err := suite.Post("/user", headers, body)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusCreated, resp.StatusCode)
			json := map[string]interface{}{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				suite.Contains(json, "id")
			}

			u := &model.User{}
			if err := database.Conn().Where("email = ?", "jack@example.org").First(u).Error; err != nil {
				suite.Error(err)
			}

			actualPath := userController.StoragePath + u.Image.String
			if suite.FileExists(actualPath) {
				ref, err := suite.readFile(path)
				if err != nil {
					suite.Error(err)
					return
				}

				actual, err := suite.readFile(actualPath)
				if err != nil {
					suite.Error(err)
					return
				}

				suite.Equal(ref, actual)
				filesystem.Delete(actualPath)
			}
		}
	})
}

func (suite *UserTestSuite) TestBcryptPassword() {
	factory := database.NewFactory(model.UserGenerator)
	generatedUsers := factory.Generate(5).([]*model.User)
	passwords := make([]string, 0, len(generatedUsers))
	for _, u := range generatedUsers {
		passwords = append(passwords, u.Password)
	}

	db := database.Conn()
	if err := db.Create(generatedUsers).Error; err != nil {
		suite.Error(err)
		return
	}

	users := make([]*model.User, 0, len(generatedUsers))
	if err := db.Order("id asc").Find(&users).Error; err != nil {
		suite.Error(err)
	}

	for k, u := range users {
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwords[k])); err != nil {
			suite.Failf("Hash and password comparison failed", "%q %q: %w", u.Password, passwords[k], err)
		}
	}
}

func (suite *UserTestSuite) TestShow() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.UserGenerator)
		override := &model.User{
			Email: "jack@example.org",
		}
		generatedUser := factory.Override(override).Save(1).([]*model.User)[0]
		token, err := auth.GenerateToken("jack@example.org")
		if err != nil {
			suite.Error(err)
			return
		}

		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}
		resp, err := suite.Get("/user", headers)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusOK, resp.StatusCode)
			user := &model.User{}
			err := suite.GetJSONBody(resp, user)
			suite.Nil(err)
			if err == nil {
				suite.Equal(generatedUser.ID, user.ID)
				suite.Equal(generatedUser.Email, user.Email)
				suite.Equal(generatedUser.Username, user.Username)
				suite.Equal("", user.Password) // Password is hidden
			}
		}
	})
}

func (suite *UserTestSuite) TestImage() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.UserGenerator)
		str := null.String{}
		str.String = "test_profile_picture.png"
		str.Valid = true
		override := &model.User{
			Image: str,
		}
		u := factory.Override(override).Save(1).([]*model.User)[0]

		// Create temp profile picture
		destPath := userController.StoragePath + "test_profile_picture.png"
		refPath := "resources/test/img/goyave_64.png"
		input, err := ioutil.ReadFile(refPath)
		if err != nil {
			suite.Error(err)
			return
		}

		err = ioutil.WriteFile(destPath, input, 0660)
		if err != nil {
			suite.Error(err)
			return
		}
		defer filesystem.Delete(destPath)

		resp, err := suite.Get(fmt.Sprintf("/user/%d/image", u.ID), nil)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusOK, resp.StatusCode)
			suite.Equal("image/png", resp.Header.Get("Content-Type"))
			body := suite.GetBody(resp)

			ref, err := suite.readFile(refPath)
			if err != nil {
				suite.Error(err)
				return
			}

			suite.Equal(ref, body)
		}
	})
}

func (suite *UserTestSuite) TestImageDefault() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.UserGenerator) // UserGenerator doesn't set the Image field
		user := factory.Save(1).([]*model.User)[0]

		resp, err := suite.Get(fmt.Sprintf("/user/%d/image", user.ID), nil)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusOK, resp.StatusCode)
			suite.Equal("image/png", resp.Header.Get("Content-Type"))
			body := suite.GetBody(resp)

			ref, err := suite.readFile("resources/img/default_profile_picture.png")
			if err != nil {
				suite.Error(err)
				return
			}

			suite.Equal(ref, body)
		}
	})
}

func (suite *UserTestSuite) TestUpdate() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.UserGenerator)
		user := factory.Save(1).([]*model.User)[0]
		token, err := auth.GenerateToken(user.Email)
		if err != nil {
			suite.Error(err)
			return
		}

		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}
		request := map[string]interface{}{
			"username": user.Username + "_edited",
		}
		body, _ := json.Marshal(request)
		resp, err := suite.Patch("/user", headers, bytes.NewReader(body))
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusNoContent, resp.StatusCode)

			updatedUser := &model.User{}
			if err := database.Conn().Model(&model.User{}).Where("email = ?", user.Email).First(updatedUser).Error; err != nil {
				suite.Error(err)
			}
			suite.Equal(request["username"], updatedUser.Username)
		}
	})
}

func (suite *UserTestSuite) TestUpdateImage() {
	suite.RunServer(route.Register, func() {
		const path = "resources/test/img/goyave_64.png"

		factory := database.NewFactory(model.UserGenerator)
		user := factory.Save(1).([]*model.User)[0]

		// Set image for user (to check deletion on update)
		destPath := userController.StoragePath + "test_profile_picture.png"
		input, err := ioutil.ReadFile(path)
		if err != nil {
			suite.Error(err)
			return
		}

		err = ioutil.WriteFile(destPath, input, 0660)
		if err != nil {
			suite.Error(err)
			return
		}
		user.Image.String = "test_profile_picture.png"
		user.Image.Valid = true
		database.Conn().Save(user)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		suite.WriteFile(writer, path, "image", filepath.Base(path))
		if err := writer.Close(); err != nil {
			suite.Error(err)
			return
		}

		token, err := auth.GenerateToken(user.Email)
		if err != nil {
			suite.Error(err)
			return
		}

		headers := map[string]string{
			"Content-Type":  writer.FormDataContentType(),
			"Authorization": "Bearer " + token,
		}

		resp, err := suite.Patch("/user", headers, body)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusNoContent, resp.StatusCode)

			u := &model.User{}
			if err := database.Conn().Where("email = ?", user.Email).First(u).Error; err != nil {
				suite.Error(err)
			}

			// Previous image has been deleted
			suite.NoFileExists(userController.StoragePath + user.Image.String)

			actualPath := userController.StoragePath + u.Image.String
			if suite.FileExists(actualPath) {

				ref, err := suite.readFile(path)
				if err != nil {
					suite.Error(err)
					return
				}

				actual, err := suite.readFile(actualPath)
				if err != nil {
					suite.Error(err)
					return
				}
				suite.Equal(ref, actual)

				filesystem.Delete(actualPath)
			}
		}
	})
}

func TestUserSuite(t *testing.T) {
	goyave.RunTest(t, new(UserTestSuite))
}
