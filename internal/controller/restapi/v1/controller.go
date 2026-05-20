package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	t        usecase.Translation
	u        usecase.User
	tk       usecase.Task
	catalog  usecase.Catalog
	media    usecase.Media
	homepage usecase.Homepage
	cart     usecase.Cart
	lead     usecase.Lead
	content  usecase.Content
	importer usecase.Importer
	l        logger.Interface
	v        *validator.Validate
}
