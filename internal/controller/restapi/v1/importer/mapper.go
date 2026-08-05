package importer

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toImportExecuteResponse maps entity.ImportResult to openapi.ImportExecuteResponse DTO.
func toImportExecuteResponse(res entity.ImportResult) openapi.ImportExecuteResponse {
	return openapi.ImportExecuteResponse{
		Imported:    res.Imported,
		FailedCount: len(res.Failed),
	}
}
