package user

import (
	vv "github.com/go-goyave/goyave-blog-example/http/validation"
	"goyave.dev/goyave/v5"
	v "goyave.dev/goyave/v5/validation"
)

func (ctrl *Controller) RegisterRequest(_ *goyave.Request) v.RuleSet {
	return v.RuleSet{
		{Path: v.CurrentElement, Rules: v.List{v.Required(), v.Object()}},
		{Path: "email", Rules: v.List{
			v.Required(), v.String(), v.Trim(), v.Email(), v.Max(320), v.Unique(ctrl.UserService.UniqueScope()),
		}},
		{Path: "username", Rules: v.List{v.Required(), v.String(), v.Trim(), v.Between(3, 100)}},
		{Path: "avatar", Rules: v.List{v.Nullable(), v.File(), v.Image(), v.Max(2048), v.FileCount(1)}},
		{Path: "password", Rules: v.List{v.Required(), v.String(), v.Between(6, 72), vv.Password()}},
	}
}

func (ctrl *Controller) UpdateRequest(_ *goyave.Request) v.RuleSet {
	return v.RuleSet{
		{Path: v.CurrentElement, Rules: v.List{v.Required(), v.Object()}},
		{Path: "email", Rules: v.List{
			v.String(), v.Trim(), v.Email(), v.Max(320), v.Unique(ctrl.UserService.UniqueScope()),
		}},
		{Path: "username", Rules: v.List{v.String(), v.Trim(), v.Between(3, 100)}},
		{Path: "avatar", Rules: v.List{v.Nullable(), v.File(), v.Image(), v.Max(2048), v.FileCount(1)}},
		{Path: "password", Rules: v.List{v.String(), v.Between(6, 72), vv.Password()}},
	}
}
