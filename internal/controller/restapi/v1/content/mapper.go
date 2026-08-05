package content

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toStaticPageResponse maps entity.StaticPage to openapi.StaticPageResponse DTO.
func toStaticPageResponse(p entity.StaticPage) openapi.StaticPageResponse {
	content := p.Body
	return openapi.StaticPageResponse{
		Id:        p.ID,
		Slug:      p.Slug,
		Title:     p.Title,
		Content:   &content,
		CreatedAt: &p.UpdatedAt,
	}
}

// toHomepageConfigResponse maps entity.HomepageBlock list to openapi.HomepageConfigResponse DTO.
func toHomepageConfigResponse(blocks []entity.HomepageBlock) openapi.HomepageConfigResponse {
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
