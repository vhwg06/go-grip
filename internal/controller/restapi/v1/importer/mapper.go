package importer

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	importermodule "github.com/evrone/go-clean-template/internal/module/importer"
)

// toImportExecuteResponse maps importermodule.ImportResult to openapi.ImportExecuteResponse DTO.
func toImportExecuteResponse(res importermodule.ImportResult) openapi.ImportExecuteResponse {
	return openapi.ImportExecuteResponse{
		Imported:    res.Imported,
		FailedCount: len(res.Failed),
	}
}
