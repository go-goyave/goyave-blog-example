package article

import (
	"goyave.dev/goyave/v5"
	v "goyave.dev/goyave/v5/validation"
)

func (ctrl *Controller) CreateRequest(_ *goyave.Request) v.RuleSet {
	return v.RuleSet{
		{Path: v.CurrentElement, Rules: v.List{v.Required(), v.Object()}},
		{Path: "title", Rules: v.List{v.Required(), v.String(), v.Trim(), v.Between(1, 200)}},
		{Path: "contents", Rules: v.List{v.Required(), v.String()}},
	}
}

func (ctrl *Controller) UpdateRequest(_ *goyave.Request) v.RuleSet {
	return v.RuleSet{
		{Path: v.CurrentElement, Rules: v.List{v.Required(), v.Object()}},
		{Path: "title", Rules: v.List{v.String(), v.Trim(), v.Between(1, 200)}},
		{Path: "contents", Rules: v.List{v.String()}},
	}
}
