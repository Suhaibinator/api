package go_persistence

type HadithCollection struct {
	Name                         string  `gorm:"column:name" json:"name"`
	CollectionId                 int     `gorm:"column:collectionID" json:"collectionID"`
	Type                         string  `gorm:"column:type" json:"type"`
	EnglishTitle                 string  `gorm:"column:englishTitle" json:"englishTitle"`
	ArabicTitle                  string  `gorm:"column:arabicTitle" json:"arabicTitle"`
	HasVolumes                   string  `gorm:"column:hasvolumes" json:"hasvolumes"`
	HasBooks                     string  `gorm:"column:hasbooks" json:"hasbooks"`
	HasChapters                  string  `gorm:"column:haschapters" json:"haschapters"`
	NumHadith                    int     `gorm:"column:numhadith" json:"numhadith"`
	TotalHadith                  *int    `gorm:"column:totalhadith" json:"totalhadith"`
	EnglishGrade1                *string `gorm:"column:englishgrade1" json:"englishgrade1"`
	ArabicGrade1                 *string `gorm:"column:arabicgrade1" json:"arabicgrade1"`
	ShowEnglishTranslationNumber string  `gorm:"column:showEnglishTranslationNumber" json:"showEnglishTranslationNumber"`
	ShowInBookReference          string  `gorm:"column:showInBookReference" json:"showInBookReference"`
	ShowOnHome                   string  `gorm:"column:showOnHome" json:"showOnHome"`
	ShowTranslationProgress      string  `gorm:"column:showTranslationProgress" json:"showTranslationProgress"`
	ReferenceTemplate            *string `gorm:"column:reference_template" json:"reference_template"`
	Annotation                   *string `gorm:"column:annotation" json:"annotation"`
	ShortIntro                   string  `gorm:"column:shortintro" json:"shortintro"`
	About                        string  `gorm:"column:about" json:"about"`
	Status                       string  `gorm:"column:status" json:"status"`
	NumberingInfoDesc            string  `gorm:"column:numberinginfodesc" json:"numberinginfodesc"`
}

func (hc *HadithCollection) TableName() string {
	return "Collections"
}

type Book struct {
	OurBookID   int     `gorm:"column:ourBookID"`
	Collection  string  `gorm:"column:collection"`
	EnglishName *string `gorm:"column:englishBookName"`
	ArabicName  *string `gorm:"column:arabicBookName"`
	FirstNumber int     `gorm:"column:firstNumber"`
	LastNumber  int     `gorm:"column:lastNumber"`
	TotalNumber int     `gorm:"column:totalNumber"`
	Status      int     `gorm:"column:status"`
}

func (b *Book) TableName() string {
	return "BookData"
}

type Chapter struct {
	Collection       string  `gorm:"column:collection"`
	ArabicBookID     int     `gorm:"column:arabicBookID"`
	BabID            float64 `gorm:"column:babID"`
	EnglishBabNumber string  `gorm:"column:englishBabNumber"`
	EnglishBabName   string  `gorm:"column:englishBabName"`
	EnglishIntro     string  `gorm:"column:englishIntro"`
	EnglishEnding    string  `gorm:"column:englishEnding"`
	ArabicBabNumber  string  `gorm:"column:arabicBabNumber"`
	ArabicBabName    string  `gorm:"column:arabicBabName"`
	ArabicIntro      string  `gorm:"column:arabicIntro"`
	ArabicEnding     string  `gorm:"column:arabicEnding"`
}

func (c *Chapter) TableName() string {
	return "ChapterData"
}

type Hadith struct {
	Collection       string `gorm:"column:collection"`
	BookNumber       string `gorm:"column:bookNumber"`
	BabID            int    `gorm:"column:babID"`
	HadithNumber     string `gorm:"column:hadithNumber"`
	EnglishBabNumber string `gorm:"column:englishBabNumber"`
	EnglishBabName   string `gorm:"column:englishBabName"`
	EnglishURN       int    `gorm:"column:englishURN"`
	EnglishText      string `gorm:"column:englishText"`
	ArabicBabNumber  string `gorm:"column:arabicBabNumber"`
	ArabicBabName    string `gorm:"column:arabicBabName"`
	ArabicURN        int    `gorm:"column:arabicURN"`
	ArabicText       string `gorm:"column:arabicText"`
	EnglishGrade1    string `gorm:"column:englishgrade1"`
	ArabicGrade1     string `gorm:"column:arabicgrade1"`
}

func (h *Hadith) TableName() string {
	return "HadithTable"
}
