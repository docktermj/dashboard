// Support for the /api URL route.
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docktermj/dashboard/models/fileindex"
	log "github.com/docktermj/go-logger/logger"
	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

// ----------------------------------------------------------------------------
// WebHandler structures
// ----------------------------------------------------------------------------

type ApiHandler struct {
	datastore             fileindex.Interface
	EverythingCount       int
	DuplicatesSha256Count int
	UniqueSha256Count     int
}

// ----------------------------------------------------------------------------
// HTTP request handlers
// ----------------------------------------------------------------------------

type HttpResponseCount struct {
	Count int `json:"count"`
}

type HttpResponseAllColumns struct {
	Draw            int                    `json:"draw"`
	Data            []*fileindex.FileIndex `json:"data"`
	RecordsFiltered int                    `json:"recordsFiltered"`
	RecordsTotal    int                    `json:"recordsTotal"`
}

type HttpResponseDuplicatesSha56 struct {
	Draw            int                           `json:"draw"`
	Data            []*fileindex.DuplicatesSha256 `json:"data"`
	RecordsFiltered int                           `json:"recordsFiltered"`
	RecordsTotal    int                           `json:"recordsTotal"`
}

// ----------------------------------------------------------------------------
// ApiHandler methods
// ----------------------------------------------------------------------------

func (handler *ApiHandler) transformFileIndexToJson(rows []*fileindex.FileIndex) ([]byte, error) {

	response := HttpResponseAllColumns{
		Data: rows,
	}

	// Construct JSON response.

	response_json, err := json.Marshal(response)
	if err != nil {
		log.Errorf("handler.transformFileIndexToJson Error:  %v response: %v", err, response)
		return nil, err
	}

	return response_json, nil
}

func (handler *ApiHandler) getIntQueryParameters(request *http.Request, queryParameter string) (int, error) {

	valueString := request.URL.Query().Get(queryParameter)
	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		errMessage := fmt.Sprintf("Couldn't parse %s: %s Error:  %v", queryParameter, valueString, err)
		log.Errorf(errMessage)
		return 0, errors.New(errMessage)
	}

	return valueInt, nil
}

func (handler *ApiHandler) getDatatableQueryParameters(request *http.Request) (*fileindex.DatabaseQueryMetadata, error) {

	start, err := handler.getIntQueryParameters(request, "start")
	if err != nil {
		return nil, err
	}

	length, err := handler.getIntQueryParameters(request, "length")
	if err != nil {
		return nil, err
	}

	order0Column, err := handler.getIntQueryParameters(request, "order[0][column]")
	if err != nil {
		return nil, err
	}

	result := fileindex.DatabaseQueryMetadata{
		Start:          start,
		Limit:          length,
		OrderColumn:    order0Column + 1,
		OrderDirection: request.URL.Query().Get("order[0][dir]"),
		Search:         request.URL.Query().Get("search[value]"),
	}

	return &result, nil
}

func (handler *ApiHandler) getAllEverythingCount() int {

	if handler.EverythingCount == 0 {
		count, err := handler.datastore.AllEverythingCount()
		if err != nil {
			log.Errorf("handler.datastore.AllEverythingCount() Error:  %v", err)
			return 0
		}
		handler.EverythingCount = count
	}

	return handler.EverythingCount
}

func (handler *ApiHandler) getAllDuplicatesSha256Count() int {

	if handler.DuplicatesSha256Count == 0 {
		count, err := handler.datastore.AllDuplicatesSha256Count()
		if err != nil {
			log.Errorf("handler.datastore.AllDuplicatesSha256Count() Error:  %v", err)
			return 0
		}
		handler.DuplicatesSha256Count = count
	}

	return handler.DuplicatesSha256Count
}

func (handler *ApiHandler) getAllUniqueSha256Count() int {

	if handler.UniqueSha256Count == 0 {
		count, err := handler.datastore.AllUniqueSha256Count()
		if err != nil {
			log.Errorf("handler.datastore.getAllUniqueSha256Count() Error:  %v", err)
			return 0
		}
		handler.UniqueSha256Count = count
	}

	return handler.UniqueSha256Count
}

// ----------------------------------------------------------------------------
// HTTP handlers.
//    Signature: method(net/http.ResponseWriter, *net/http.Request)
// ----------------------------------------------------------------------------

func (handler *ApiHandler) everything(responseWriter http.ResponseWriter, request *http.Request) {

	// Get parameters from URL.

	datatableQueryParameters, err := handler.getDatatableQueryParameters(request)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.getDatatableQueryParameters() Error:  %v", err)
		return
	}

	draw, err := handler.getIntQueryParameters(request, "draw")
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.getIntQueryParameters() Error:  %v", err)
		return
	}

	// Get rows from database.

	rows, err := handler.datastore.Everything(datatableQueryParameters)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.Everything() Error:  %v", err)
		return
	}

	// Get total count of filtered

	count, err := handler.datastore.EverythingCount(datatableQueryParameters)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.EverythingCount() Error:  %v", err)
		return
	}

	response := HttpResponseAllColumns{
		Draw:            draw,
		Data:            rows,
		RecordsFiltered: count,
		RecordsTotal:    handler.getAllEverythingCount(),
	}

	// Construct JSON response.

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("json.Marshal() Error:  %v", err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) count(responseWriter http.ResponseWriter, request *http.Request) {

	// Get information from database.

	count, err := handler.datastore.AllEverythingCount()
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}

	// Update cached EverythingCount.

	handler.EverythingCount = count

	// Construct response.

	response := HttpResponseCount{
		count,
	}
	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) root(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte("index"))
}

func (handler *ApiHandler) duplicatesSha256(responseWriter http.ResponseWriter, request *http.Request) {

	// Get parameters from URL.

	datatableQueryParameters, err := handler.getDatatableQueryParameters(request)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.getDatatableQueryParameters() Error:  %v", err)
		return
	}

	draw, err := handler.getIntQueryParameters(request, "draw")
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.getIntQueryParameters() Error:  %v", err)
		return
	}

	// Get information from database.

	rows, err := handler.datastore.DuplicatesSha256(datatableQueryParameters)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.DuplicatesSha256() Error:  %v", err)
		return
	}

	count, err := handler.datastore.DuplicatesSha256Count(datatableQueryParameters)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.DuplicatesSha256Count() Error:  %v", err)
		return
	}

	response := HttpResponseDuplicatesSha56{
		Draw:            draw,
		Data:            rows,
		RecordsFiltered: count,
		RecordsTotal:    handler.getAllDuplicatesSha256Count(),
	}

	// Construct JSON response.

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("json.Marshal Error:  %v", err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) uniqueSha256(responseWriter http.ResponseWriter, request *http.Request) {

	// Get parameters from URL.

	datatableQueryParameters, err := handler.getDatatableQueryParameters(request)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.getDatatableQueryParameters() Error:  %v", err)
		return
	}

	draw, err := handler.getIntQueryParameters(request, "draw")
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.getIntQueryParameters() Error:  %v", err)
		return
	}

	// Get information from database.

	rows, err := handler.datastore.UniqueSha256(datatableQueryParameters)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.UniqueSha256() Error:  %v", err)
		return
	}

	count, err := handler.datastore.UniqueSha256Count(datatableQueryParameters)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.UniqueSha256Count() Error:  %v", err)
		return
	}

	response := HttpResponseAllColumns{
		Draw:            draw,
		Data:            rows,
		RecordsFiltered: count,
		RecordsTotal:    handler.getAllUniqueSha256Count(),
	}

	// Construct JSON response.

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("json.Marshal() Error:  %v", err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

// ----------------------------------------------------------------------------
// byGeneric helpers.
// ----------------------------------------------------------------------------

type byFunc func(string, *fileindex.DatabaseQueryMetadata) ([]*fileindex.FileIndex, error)
type byFuncCount func(string) (int, error)

func (handler *ApiHandler) byGeneric(request *http.Request, byFunction byFunc, byFunctionCount byFuncCount, arg string) ([]byte, error) {

	// Get parameters from URL.

	datatableQueryParameters, err := handler.getDatatableQueryParameters(request)
	if err != nil {
		log.Errorf("handler.getDatatableQueryParameters() Error:  %v", err)
		return nil, err
	}

	draw, err := handler.getIntQueryParameters(request, "draw")
	if err != nil {
		log.Errorf("handler.getIntQueryParameters() Error:  %v", err)
		return nil, err
	}

	// Get information from database.

	rows, err := byFunction(arg, datatableQueryParameters)
	if err != nil {
		log.Errorf("byFunction(%s) Error:  %v", arg, err)
		return nil, err
	}

	count, err := byFunctionCount(arg)
	if err != nil {
		log.Errorf("byFunctionCount(%s) Error:  %v", arg, err)
		return nil, err
	}

	response := HttpResponseAllColumns{
		Draw:            draw,
		Data:            rows,
		RecordsFiltered: count,
		RecordsTotal:    count,
	}

	// Construct JSON response.

	response_json, err := json.Marshal(response)
	if err != nil {
		log.Errorf("json.Marshal() Error:  %v", err)
		return nil, err
	}

	return response_json, nil
}

// ----------------------------------------------------------------------------
// HTTP handlers. byGeneric callers.
//    Signature: method(net/http.ResponseWriter, *net/http.Request)
// ----------------------------------------------------------------------------

func (handler *ApiHandler) id(responseWriter http.ResponseWriter, request *http.Request) {
	id_encoded := chi.URLParam(request, "id")
	id_decoded, err := url.QueryUnescape(id_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	response_json, err := handler.byGeneric(request, handler.datastore.ByID, handler.datastore.ByIDCount, id_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.ByID(%s) Error:  %v", id_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) modified(responseWriter http.ResponseWriter, request *http.Request) {

	// Pull information from HTTP request.

	modified_encoded := chi.URLParam(request, "modified")
	modified_decoded, err := url.QueryUnescape(modified_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Query data store.

	response_json, err := handler.byGeneric(request, handler.datastore.ByModified, handler.datastore.ByModifiedCount, modified_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.ByModified(%s) Error:  %v", modified_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) name(responseWriter http.ResponseWriter, request *http.Request) {

	// Pull information from HTTP request.

	name_encoded := chi.URLParam(request, "name")
	name_decoded, err := url.QueryUnescape(name_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	response_json, err := handler.byGeneric(request, handler.datastore.ByName, handler.datastore.ByNameCount, name_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.ByName(%s) Error:  %v", name_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) path(responseWriter http.ResponseWriter, request *http.Request) {

	// Pull information from HTTP request.

	path_encoded := chi.URLParam(request, "path")
	path_decoded, err := url.QueryUnescape(path_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	response_json, err := handler.byGeneric(request, handler.datastore.ByPath, handler.datastore.ByPathCount, path_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.ByPath(%s) Error:  %v", path_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) sha256(responseWriter http.ResponseWriter, request *http.Request) {

	// Pull information from HTTP request.

	sha256_encoded := chi.URLParam(request, "sha256")
	sha256_decoded, err := url.QueryUnescape(sha256_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	response_json, err := handler.byGeneric(request, handler.datastore.BySha256, handler.datastore.BySha256Count, sha256_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.BySha256(%s) Error:  %v", sha256_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) size(responseWriter http.ResponseWriter, request *http.Request) {

	// Pull information from HTTP request.

	size_encoded := chi.URLParam(request, "size")
	size_decoded, err := url.QueryUnescape(size_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	response_json, err := handler.byGeneric(request, handler.datastore.BySize, handler.datastore.BySizeCount, size_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.BySize(%s) Error:  %v", size_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}

func (handler *ApiHandler) volume(responseWriter http.ResponseWriter, request *http.Request) {

	// Pull information from HTTP request.

	volume_encoded := chi.URLParam(request, "volume")
	volume_decoded, err := url.QueryUnescape(volume_encoded)
	if err != nil {
		log.Fatal(err)
		return
	}

	response_json, err := handler.byGeneric(request, handler.datastore.ByVolume, handler.datastore.ByVolumeCount, volume_decoded)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.datastore.ByVolume(%s) Error:  %v", volume_decoded, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(response_json)
}
