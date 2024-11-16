package go_api

import (
	"net/http"

	"github.com/Suhaibinator/api/go_service"
	"github.com/gorilla/mux"
)

type ApplicationRouter struct {
	router             *mux.Router
	applicationService *go_service.ApplicationService
}

func NewApplicationRouter(router *mux.Router, applicationService *go_service.ApplicationService) *ApplicationRouter {
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
	collectionsRouter.HandleFunc("/{name}", ar.apiCollectionHandler).Methods("GET")
	collectionsRouter.HandleFunc("/{name}/books", ar.apiGetBooksInCollectionHandler).Methods("GET")
	collectionsRouter.HandleFunc("/{name}/books/{bookNumber}", ar.apGetBookHandler).Methods("GET")
	collectionsRouter.HandleFunc("/{collection_name}/books/{bookNumber}/hadiths", ar.apiGetHadithsInBookInCollection).Methods("GET")
	collectionsRouter.HandleFunc("/{collection_name}/hadiths/{hadithNumber}", ar.apiGetHadithInCollection).Methods("GET")
	collectionsRouter.HandleFunc("/{collection_name}/books/{bookNumber}/chapters", ar.apiGetChaptersInBook).Methods("GET")
	collectionsRouter.HandleFunc("/{collection_name}/books/{bookNumber}/chapters/{chapterId}", ar.apiGetHadithsInChapter).Methods("GET")

	// Other v1 routes
	v1Router.HandleFunc("/hadiths/{urn:[0-9]+}", ar.apiGetHadithByUrn).Methods("GET")
	v1Router.HandleFunc("/hadiths", ar.apiHadithsHandler).Methods("GET")
	v1Router.HandleFunc("/hadiths/random", ar.apiHadithsRandomHandler).Methods("GET")
}

func (ar *ApplicationRouter) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to sunnah.com API.</h1>"))
}

func (ar *ApplicationRouter) apiGetAllCollections(w http.ResponseWriter, r *http.Request) {

}

func (ar *ApplicationRouter) apiCollectionHandler(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiGetBooksInCollectionHandler(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apGetBookHandler(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiGetHadithsInBookInCollection(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiGetHadithInCollection(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiGetChaptersInBook(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiGetHadithsInChapter(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiGetHadithByUrn(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiHadithsHandler(w http.ResponseWriter, r *http.Request) {
}

func (ar *ApplicationRouter) apiHadithsRandomHandler(w http.ResponseWriter, r *http.Request) {
}
