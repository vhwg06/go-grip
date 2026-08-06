package content

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
)

// toStaticPageResponse maps contentmodule.StaticPage to openapi.StaticPageResponse DTO.
func toStaticPageResponse(p contentmodule.StaticPage) openapi.StaticPageResponse {
	content := p.Body
	templateKey := p.TemplateKey
	status := string(p.Status)
	return openapi.StaticPageResponse{
		Id:          p.ID,
		Slug:        p.Slug,
		Title:       p.Title,
		Content:     &content,
		Body:        &content,
		CreatedAt:   &p.UpdatedAt,
		Gallery:     &p.Gallery,
		TemplateKey: &templateKey,
		Status:      &status,
	}
}

// toHomepageConfigResponse maps contentmodule.HomepageBlock list to openapi.HomepageConfigResponse DTO.
func toHomepageConfigResponse(blocks []contentmodule.HomepageBlock) openapi.HomepageConfigResponse {
	bannerURL := "https://example.com/banner.png"
	metaTitle := "Go-Grip Store"
	featuredIDs := []string{}
	for _, b := range blocks {
		featuredIDs = append(featuredIDs, b.ID)
	}

	return openapi.HomepageConfigResponse{
		BannerUrl:          &bannerURL,
		FeaturedProductIds: &featuredIDs,
		MetaTitle:          &metaTitle,
	}
}
