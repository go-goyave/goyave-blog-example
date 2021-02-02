package middleware

import (
	"errors"
	"net/http"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/database"
	dbModel "github.com/go-goyave/goyave-blog-example/database/model"
	"gorm.io/gorm"
)

// Owner checks if the authenticated user is the owner of the requested resource.
//  - param: the name of the URL parameter to use
//  - column: the name of the foreign key column
//  - model: the model of the requested resource
func Owner(param, column string, model interface{}) goyave.Middleware {
	return func(next goyave.Handler) goyave.Handler {
		return func(response *goyave.Response, request *goyave.Request) {

			if p, ok := request.Params[param]; ok {
				id := request.User.(*dbModel.User).ID
				if err := database.Conn().Model(model).Select("1").Where(column+" = ?", id).First(map[string]interface{}{}, p).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						response.Status(http.StatusForbidden)
					} else {
						response.Error(err)
					}
					return
				}
			}
			next(response, request)
		}
	}
}
