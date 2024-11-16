package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database connection
func init() {
	var err error
	// Replace with your actual database connection string
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

// Models

func (hc *HadithCollection) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"name":        hc.Name,
		"hasBooks":    hc.HasBooks == "yes",
		"hasChapters": hc.HasChapters == "yes",
		"collection": []map[string]string{
			{"lang": "en", "title": hc.EnglishTitle, "shortIntro": hc.ShortIntro},
			{"lang": "ar", "title": hc.ArabicTitle, "shortIntro": hc.getArabicShortIntro()},
		},
		"totalHadith":          hc.TotalHadith,
		"totalAvailableHadith": hc.TotalAvailableHadith,
	}
}

func (hc *HadithCollection) getArabicShortIntro() string {
	if hc.ShortIntroArabic != "" {
		return hc.ShortIntroArabic
	}
	return hc.ShortIntro
}

var idNumberMap = map[int]string{
	-1:  "introduction",
	-35: "35b",
}

func GetNumberFromID(bookID int) string {
	if val, ok := idNumberMap[bookID]; ok {
		return val
	}
	return strconv.Itoa(bookID)
}

func GetIDFromNumber(bookNumber string) int {
	for k, v := range idNumberMap {
		if v == bookNumber {
			return k
		}
	}
	if num, err := strconv.Atoi(bookNumber); err == nil {
		return num
	}
	return 0 // default or error handling
}

func (b *Book) Serialize() map[string]interface{} {
	bookNumber := GetNumberFromID(b.OurBookID)
	return map[string]interface{}{
		"bookNumber": bookNumber,
		"book": []map[string]string{
			{"lang": "en", "name": b.EnglishName},
			{"lang": "ar", "name": b.ArabicName},
		},
		"hadithStartNumber": b.FirstNumber,
		"hadithEndNumber":   b.LastNumber,
		"numberOfHadith":    b.TotalNumber,
	}
}

func (c *Chapter) Serialize() map[string]interface{} {
	bookNumber := GetNumberFromID(c.ArabicBookID)
	return map[string]interface{}{
		"bookNumber": bookNumber,
		"chapterId":  fmt.Sprintf("%v", c.BabID),
		"chapter": []map[string]string{
			{
				"lang":          "en",
				"chapterNumber": c.EnglishBabNumber,
				"chapterTitle":  cleanupEnChapterTitle(c.EnglishBabName),
				"intro":         cleanupEnText(c.EnglishIntro),
				"ending":        cleanupEnText(c.EnglishEnding),
			},
			{
				"lang":          "ar",
				"chapterNumber": c.ArabicBabNumber,
				"chapterTitle":  cleanupChapterTitle(c.ArabicBabName),
				"intro":         cleanupText(c.ArabicIntro),
				"ending":        cleanupText(c.ArabicEnding),
			},
		},
	}
}

func (h *Hadith) GetGrade(fieldName string) []map[string]string {
	var grades []map[string]string
	gradeVal := ""

	// Use reflection or simple switch to get the field value
	switch fieldName {
	case "englishgrade1":
		gradeVal = h.EnglishGrade1
	case "arabicgrade1":
		gradeVal = h.ArabicGrade1
	}

	if gradeVal == "" {
		return grades
	}

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(gradeVal), &grades); err == nil {
		return grades
	}

	// If not JSON, construct from individual fields
	gradedBy := ""
	grade := gradeVal

	// Assume you have a relation to HadithCollection to get graded_by
	var collection HadithCollection
	if err := DB.Where("name = ?", h.Collection).First(&collection).Error; err == nil {
		if fieldName == "englishgrade1" {
			gradedBy = collection.EnglishTitle
		} else if fieldName == "arabicgrade1" {
			gradedBy = collection.ArabicTitle
		}
	}

	return []map[string]string{
		{"graded_by": gradedBy, "grade": grade},
	}
}

func (h *Hadith) Serialize() map[string]interface{} {
	gradesEn := h.GetGrade("englishgrade1")
	gradesAr := h.GetGrade("arabicgrade1")

	return map[string]interface{}{
		"collection":   h.Collection,
		"bookNumber":   h.BookNumber,
		"chapterId":    strconv.Itoa(h.BabID),
		"hadithNumber": h.HadithNumber,
		"hadith": []map[string]interface{}{
			{
				"lang":          "en",
				"chapterNumber": h.EnglishBabNumber,
				"chapterTitle":  cleanupEnChapterTitle(h.EnglishBabName),
				"urn":           h.EnglishURN,
				"body":          cleanupEnText(h.EnglishText),
				"grades":        gradesEn,
			},
			{
				"lang":          "ar",
				"chapterNumber": h.ArabicBabNumber,
				"chapterTitle":  cleanupChapterTitle(h.ArabicBabName),
				"urn":           h.ArabicURN,
				"body":          cleanupText(h.ArabicText),
				"grades":        gradesAr,
			},
		},
	}
}

// Helper functions for cleaning text
func cleanupText(text string) string {
	return strings.TrimSpace(text)
}

func cleanupEnText(text string) string {
	return strings.TrimSpace(text)
}

func cleanupChapterTitle(title string) string {
	return strings.TrimSpace(title)
}

func cleanupEnChapterTitle(title string) string {
	return strings.TrimSpace(title)
}

func main() {
	router := mux.NewRouter()

	// Home route
	router.HandleFunc("/", homeHandler).Methods("GET")

	// API routes
	router.HandleFunc("/v1/collections", apiCollectionsHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{name}", apiCollectionHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{name}/books", apiCollectionBooksHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{name}/books/{bookNumber}", apiCollectionBookHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{collection_name}/books/{bookNumber}/hadiths", apiCollectionBookHadithsHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{collection_name}/hadiths/{hadithNumber}", apiCollectionHadithHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{collection_name}/books/{bookNumber}/chapters", apiCollectionBookChaptersHandler).Methods("GET")
	router.HandleFunc("/v1/collections/{collection_name}/books/{bookNumber}/chapters/{chapterId}", apiCollectionBookChapterHandler).Methods("GET")
	router.HandleFunc("/v1/hadiths/{urn:[0-9]+}", apiHadithHandler).Methods("GET")
	router.HandleFunc("/v1/hadiths/random", apiHadithsRandomHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to sunnah.com API.</h1>")
}

// Helper functions for pagination and error handling (same as before)
func getPaginationParams(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	return page, limit
}

func handleDBError(w http.ResponseWriter, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler functions
func apiCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	var collections []HadithCollection
	err := DB.Offset(offset).Limit(limit).Find(&collections).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	serializedCollections := make([]map[string]interface{}, len(collections))
	for i, collection := range collections {
		serializedCollections[i] = collection.Serialize()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializedCollections)
}

func apiCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var collection HadithCollection
	err := DB.Where("name = ?", name).First(&collection).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection.Serialize())
}

func apiCollectionBooksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	var books []Book
	err := DB.Where("collection = ? AND status = ?", name, 4).
		Order("ABS(ourBookID)").
		Offset(offset).
		Limit(limit).
		Find(&books).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	serializedBooks := make([]map[string]interface{}, len(books))
	for i, book := range books {
		serializedBooks[i] = book.Serialize()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializedBooks)
}

func apiCollectionBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	bookNumber := vars["bookNumber"]

	bookID := GetIDFromNumber(bookNumber)

	var book Book
	err := DB.Where("collection = ? AND status = ? AND ourBookID = ?", name, 4, bookID).
		First(&book).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book.Serialize())
}

func apiCollectionBookHadithsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionName := vars["collection_name"]
	bookNumber := vars["bookNumber"]

	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	var hadiths []Hadith
	err := DB.Where("collection = ? AND bookNumber = ?", collectionName, bookNumber).
		Order("englishURN").
		Offset(offset).
		Limit(limit).
		Find(&hadiths).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	serializedHadiths := make([]map[string]interface{}, len(hadiths))
	for i, hadith := range hadiths {
		serializedHadiths[i] = hadith.Serialize()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializedHadiths)
}

func apiCollectionHadithHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionName := vars["collection_name"]
	hadithNumber := vars["hadithNumber"]

	var hadith Hadith
	err := DB.Where("collection = ? AND hadithNumber = ?", collectionName, hadithNumber).
		First(&hadith).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hadith.Serialize())
}

func apiCollectionBookChaptersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionName := vars["collection_name"]
	bookNumber := vars["bookNumber"]

	bookID := GetIDFromNumber(bookNumber)

	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	var chapters []Chapter
	err := DB.Where("collection = ? AND arabicBookID = ?", collectionName, bookID).
		Order("babID").
		Offset(offset).
		Limit(limit).
		Find(&chapters).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	serializedChapters := make([]map[string]interface{}, len(chapters))
	for i, chapter := range chapters {
		serializedChapters[i] = chapter.Serialize()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializedChapters)
}

func apiCollectionBookChapterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionName := vars["collection_name"]
	bookNumber := vars["bookNumber"]
	chapterIdStr := vars["chapterId"]

	chapterID, err := strconv.ParseFloat(chapterIdStr, 64)
	if err != nil {
		http.Error(w, "Invalid chapterId", http.StatusBadRequest)
		return
	}

	bookID := GetIDFromNumber(bookNumber)

	var chapter Chapter
	err = DB.Where("collection = ? AND arabicBookID = ? AND babID = ?", collectionName, bookID, chapterID).
		First(&chapter).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapter.Serialize())
}

func apiHadithHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urnStr := vars["urn"]

	urn, err := strconv.Atoi(urnStr)
	if err != nil {
		http.Error(w, "Invalid URN", http.StatusBadRequest)
		return
	}

	var hadith Hadith
	err = DB.Where("arabicURN = ? OR englishURN = ?", urn, urn).
		First(&hadith).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hadith.Serialize())
}

func apiHadithsRandomHandler(w http.ResponseWriter, r *http.Request) {
	var hadith Hadith
	err := DB.Where("collection = ?", "riyadussalihin").
		Order(gorm.Expr("RAND()")).
		First(&hadith).Error
	if err != nil {
		handleDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hadith.Serialize())
}
