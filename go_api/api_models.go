package go_api

import "github.com/Suhaibinator/api/go_persistence"

type CollectionMeta struct {
	Language   string `json:"lang"`
	Title      string `json:"title"`
	ShortIntro string `json:"shortintro"`
}

type Collection struct {
	Name                 string           `json:"name"`
	HasBooks             bool             `json:"hasbooks"`
	HasChapters          bool             `json:"haschapters"`
	CollectionMeta       []CollectionMeta `json:"collection"`
	TotalHadith          int              `json:"totalhadith"`
	TotalAvailableHadith int              `json:"totalAvailableHadith"`
}

func ConvertDbCollectionToApiCollection(dbCollection go_persistence.HadithCollection) *Collection {
	collection := Collection{
		Name:        dbCollection.Name,
		HasBooks:    dbCollection.HasBooks == "yes",
		HasChapters: dbCollection.HasChapters == "yes",
		CollectionMeta: []CollectionMeta{
			{Language: "en", Title: dbCollection.EnglishTitle, ShortIntro: dbCollection.ShortIntro},
			{Language: "ar", Title: dbCollection.ArabicTitle, ShortIntro: dbCollection.ShortIntroArabic},
		},
		TotalHadith:          dbCollection.TotalHadith,
		TotalAvailableHadith: dbCollection.TotalAvailableHadith,
	}

	return &collection
}

type PaginatedCollections struct {
	Collections []Collection `json:"data"`
	Total       int          `json:"total"`
	Limit       int          `json:"limit"`
	PrevPage    int          `json:"previous"`
	NextPage    int          `json:"next"`
}

type BookMeta struct {
	Language string `json:"lang"`
	Name     string `json:"name"`
}

type Book struct {
	BookNumber        string     `json:"bookNumber"`
	BookMeta          []BookMeta `json:"book"`
	HadithStartNumber int        `json:"hadithStartNumber"`
	HadithEndNumber   int        `json:"hadithEndNumber"`
	NumberOfHadith    int        `json:"numberOfHadith"`
}

type PaginatedBooks struct {
	Books []Book `json:"data"`
	Total int    `json:"total"`
	Limit int    `json:"limit"`
	Prev  *int   `json:"previous"`
	Next  *int   `json:"next"`
}

type ChapterMeta struct {
	Language      string `json:"lang"`
	ChapterNumber string
	ChapterTitle  string  `json:"chapterTitle"`
	Intro         *string `json:"intro"`
	Ending        *string `json:"ending"`
}

type Chapter struct {
	BookNumber  string        `json:"bookNumber"`
	ChapterId   string        `json:"chapterId"`
	ChapterMeta []ChapterMeta `json:"chapter"`
}

type PaginatedChapters struct {
	Chapters []Chapter `json:"data"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Prev     *int      `json:"previous"`
	Next     *int      `json:"next"`
}

type HadithGradedBy struct {
	Grader string `json:"graded_by"`
	Grade  string `json:"grade"`
}

type HadithMeta struct {
	Language      string           `json:"lang"`
	ChapterNumber string           `json:"chapterNumber"`
	ChapterTitle  string           `json:"chapterTitle"`
	Urn           int              `json:"urn"`
	Body          string           `json:"body"`
	Grades        []HadithGradedBy `json:"grades"`
}

type Hadith struct {
	Collection   string       `json:"collection"`
	BookNumber   string       `json:"bookNumber"`
	ChapterId    string       `json:"chapterId"`
	HadithNumber string       `json:"hadithNumber"`
	HadithMeta   []HadithMeta `json:"hadith"`
}

type PaginatedHadiths struct {
	Hadiths []Hadith `json:"data"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Prev    *int     `json:"previous"`
	Next    *int     `json:"next"`
}
