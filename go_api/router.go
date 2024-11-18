package go_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Suhaibinator/api/go_service"
	"github.com/gorilla/mux"
)

type ApplicationRouter struct {
	router             *mux.Router
	applicationService *go_service.ApplicationService
}

func NewApplicationRouter(applicationService *go_service.ApplicationService) *ApplicationRouter {
	router := mux.NewRouter()
	return &ApplicationRouter{router: router, applicationService: applicationService}
}

func (ar *ApplicationRouter) RegisterRoutes() {
	router := ar.router

	// Home route
	router.HandleFunc("/", ar.homeHandler).Methods("GET")

	// Register v1 routes
	ar.registerV1Routes()
}

func (ar *ApplicationRouter) registerV1Routes() {
	// Create a subrouter for v1 routes
	v1Router := ar.router.PathPrefix("/v1").Subrouter()

	// Create a subrouter for collections routes
	collectionsRouter := v1Router.PathPrefix("/collections").Subrouter()

	// Collections routes
	collectionsRouter.HandleFunc("", ar.apiGetAllCollections).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}", ar.apiCollectionHandler).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}/books", ar.apiGetBooksInCollectionHandler).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}/books/{bookNumber}", ar.apGetBookHandler).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}/books/{bookNumber}/chapters", ar.apiGetChaptersInBookInCollection).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}/books/{bookNumber}/chapters/{chapterId}", ar.apiGetChapterInBookInCollection).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}/books/{bookNumber}/hadiths", ar.apiGetHadithsInBook).Methods("GET")
	collectionsRouter.HandleFunc("/{collectionName}/hadiths/{hadithNumber}", ar.apiGetHadithInCollectionByHadithNumber).Methods("GET")

	// Other v1 routes
	v1Router.HandleFunc("/hadiths/{urn:[0-9]+}", ar.apiGetHadithByUrn).Methods("GET")
	// v1Router.HandleFunc("/hadiths", ar.apiGetHadithsByCollectionAndBookAndChapter).Methods("GET")
	v1Router.HandleFunc("/hadiths/random", ar.apiHadithsRandomHandler).Methods("GET")
}

func (ar *ApplicationRouter) Run(port int) {
	log.Printf("Starting server on port %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), ar.router)
}

func (ar *ApplicationRouter) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Welcome to sunnah.com API.</h1>"))
}

func (ar *ApplicationRouter) apiGetAllCollections(w http.ResponseWriter, r *http.Request) {
	const DEFAULT_LIMIT = 50
	const DEFAULT_PAGE = 1

	const MAX_LIMIT = 100

	limit := DEFAULT_LIMIT
	page := DEFAULT_PAGE

	// Get the page and limit from the query params
	limitFromQuery := r.URL.Query().Get("limit")
	if limitFromQuery != "" {
		// Convert the limit to an integer
		limitFromQueryAsInteger, strConvErr := strconv.Atoi(limitFromQuery)
		if strConvErr == nil && limitFromQueryAsInteger > 0 && limitFromQueryAsInteger < MAX_LIMIT {
			limit = limitFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	pageFromQuery := r.URL.Query().Get("page")
	if pageFromQuery != "" {
		// Convert the page to an integer
		pageFromQueryAsInteger, strConvErr := strconv.Atoi(pageFromQuery)
		if strConvErr == nil && pageFromQueryAsInteger > 0 {
			page = pageFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	// Get the collections from the service
	collections, err := ar.applicationService.GetPaginatedHadithCollections(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the collections to API collections
	apiCollections := make([]*Collection, len(collections))
	for i, collection := range collections {
		apiCollections[i] = ConvertDbCollectionToApiCollection(collection)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiCollections); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (ar *ApplicationRouter) apiCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]

	// Get the collection from the service
	collection, err := ar.applicationService.GetHadithCollectionByName(collectionName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the collection to an API collection
	apiCollection := ConvertDbCollectionToApiCollection(collection)
	result, jsonErr := json.Marshal(apiCollection)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiGetBooksInCollectionHandler(w http.ResponseWriter, r *http.Request) {
	const DEFAULT_LIMIT = 50
	const DEFAULT_PAGE = 1

	const MAX_LIMIT = 100

	vars := mux.Vars(r)
	collectionName := vars["collectionName"]

	limit := DEFAULT_LIMIT
	page := DEFAULT_PAGE

	// Get the page and limit from the query params
	limitFromQuery := r.URL.Query().Get("limit")
	if limitFromQuery != "" {
		// Convert the limit to an integer
		limitFromQueryAsInteger, strConvErr := strconv.Atoi(limitFromQuery)
		if strConvErr == nil && limitFromQueryAsInteger > 0 && limitFromQueryAsInteger < MAX_LIMIT {
			limit = limitFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	pageFromQuery := r.URL.Query().Get("page")
	if pageFromQuery != "" {
		// Convert the page to an integer
		pageFromQueryAsInteger, strConvErr := strconv.Atoi(pageFromQuery)
		if strConvErr == nil && pageFromQueryAsInteger > 0 {
			page = pageFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	// Get the books from the service
	books, err := ar.applicationService.GetPaginatedBooksByCollection(collectionName, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the books to API books
	apiBooks := make([]*Book, len(books))
	for i, book := range books {
		apiBooks[i] = ConvertDbBookToApiBook(book)
	}
	result, jsonErr := json.Marshal(apiBooks)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apGetBookHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := mux.Vars(r)["collectionName"]
	bookNumber := mux.Vars(r)["bookNumber"]

	// Get the book from the service
	book, err := ar.applicationService.GetBookByCollectionAndBookNumber(collectionName, bookNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the book to an API book
	apiBook := ConvertDbBookToApiBook(book)
	result, jsonErr := json.Marshal(apiBook)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiGetChaptersInBookInCollection(w http.ResponseWriter, r *http.Request) {
	const DEFAULT_LIMIT = 50
	const DEFAULT_PAGE = 1

	const MAX_LIMIT = 100

	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	bookNumber := vars["bookNumber"]

	limit := DEFAULT_LIMIT
	page := DEFAULT_PAGE

	// Get the page and limit from the query params
	limitFromQuery := r.URL.Query().Get("limit")
	if limitFromQuery != "" {
		// Convert the limit to an integer
		limitFromQueryAsInteger, strConvErr := strconv.Atoi(limitFromQuery)
		if strConvErr == nil && limitFromQueryAsInteger > 0 && limitFromQueryAsInteger < MAX_LIMIT {
			limit = limitFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	pageFromQuery := r.URL.Query().Get("page")
	if pageFromQuery != "" {
		// Convert the page to an integer
		pageFromQueryAsInteger, strConvErr := strconv.Atoi(pageFromQuery)
		if strConvErr == nil && pageFromQueryAsInteger > 0 {
			page = pageFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	// Get the chapters from the service
	chapters, err := ar.applicationService.GetPaginatedChaptersByCollectionAndBookNumber(collectionName, bookNumber, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the chapters to API chapters
	apiChapters := make([]*Chapter, len(chapters))
	for i, chapter := range chapters {
		apiChapters[i] = ConvertDbChapterToApiChapter(chapter)
	}
	result, jsonErr := json.Marshal(apiChapters)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiGetChapterInBookInCollection(w http.ResponseWriter, r *http.Request) {
	collectionName := mux.Vars(r)["collectionName"]
	bookNumber := mux.Vars(r)["bookNumber"]
	chapterId := mux.Vars(r)["chapterId"]

	// Get the chapter from the service
	chapter, err := ar.applicationService.GetChapterByCollectionAndBookNumberAndChapterNumber(collectionName, bookNumber, chapterId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the chapter to an API chapter
	apiChapter := ConvertDbChapterToApiChapter(chapter)
	result, jsonErr := json.Marshal(apiChapter)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiGetHadithsInBook(w http.ResponseWriter, r *http.Request) {
	const DEFAULT_LIMIT = 50
	const DEFAULT_PAGE = 1

	const MAX_LIMIT = 100

	limit := DEFAULT_LIMIT
	page := DEFAULT_PAGE

	// Get the page and limit from the query params
	limitFromQuery := r.URL.Query().Get("limit")
	if limitFromQuery != "" {
		// Convert the limit to an integer
		limitFromQueryAsInteger, strConvErr := strconv.Atoi(limitFromQuery)
		if strConvErr == nil && limitFromQueryAsInteger > 0 && limitFromQueryAsInteger < MAX_LIMIT {
			limit = limitFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	pageFromQuery := r.URL.Query().Get("page")
	if pageFromQuery != "" {
		// Convert the page to an integer
		pageFromQueryAsInteger, strConvErr := strconv.Atoi(pageFromQuery)
		if strConvErr == nil && pageFromQueryAsInteger > 0 {
			page = pageFromQueryAsInteger
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	bookNumber := vars["bookNumber"]

	// Get the hadiths from the service
	hadiths, err := ar.applicationService.GetPaginatedHadithsByCollectionAndBookNumber(collectionName, bookNumber, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the hadiths to API hadiths
	apiHadiths := make([]*Hadith, len(hadiths))
	for i, hadith := range hadiths {
		apiHadiths[i] = ConvertDbHadithToApiHadith(hadith)
	}
	result, jsonErr := json.Marshal(apiHadiths)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiGetHadithInCollectionByHadithNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	hadithNumber := vars["hadithNumber"]

	// Get the hadith from the service
	hadith, err := ar.applicationService.GetHadithByCollectionAndHadithNumber(collectionName, hadithNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the hadith to an API hadith
	apiHadith := ConvertDbHadithToApiHadith(hadith)
	result, jsonErr := json.Marshal(apiHadith)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiGetHadithsByCollectionAndBookAndChapter(w http.ResponseWriter, r *http.Request) {

	/*
		Need to get clarity on this. Not implemented in python code
		https://sunnah.stoplight.io/docs/api/hp9hzrfn7wia9-get-a-list-of-hadiths
		Why do we need hadithNumbnber as query paramter when we are returning a paginated list of hadiths?
	*/

}

func (ar *ApplicationRouter) apiGetHadithByUrn(w http.ResponseWriter, r *http.Request) {
	urn, err := strconv.Atoi(mux.Vars(r)["urn"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Get the hadith from the service
	hadith, err := ar.applicationService.GetHadithByUrn(urn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the hadith to an API hadith
	apiHadith := ConvertDbHadithToApiHadith(hadith)
	result, jsonErr := json.Marshal(apiHadith)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ar *ApplicationRouter) apiHadithsRandomHandler(w http.ResponseWriter, r *http.Request) {
	// Get the random hadith from the service
	hadith, err := ar.applicationService.GetRandomHadith()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the hadith to an API hadith
	apiHadith := ConvertDbHadithToApiHadith(hadith)
	result, jsonErr := json.Marshal(apiHadith)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
