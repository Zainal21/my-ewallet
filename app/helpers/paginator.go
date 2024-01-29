package helpers

import (
	"log"

	"github.com/Zainal21/my-ewallet/app/appctx"
)

func MapPaginationResponseToApiResponse(paginatorResponse map[string]interface{}) *appctx.Response {
	apiResponse := appctx.NewResponse()
	//pagination data
	log.Println(paginatorResponse["current_page"])
	apiResponse.CurrentPage = paginatorResponse["current_page"].(int)
	apiResponse.LastPage = paginatorResponse["last_page"].(int)
	apiResponse.PerPage = paginatorResponse["per_page"].(int)
	apiResponse.Total = paginatorResponse["total"].(int)
	apiResponse.Status = "OK"

	apiResponse.WithData(paginatorResponse["data"]).WithMessage("SUCCESS")

	return apiResponse
}
