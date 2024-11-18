package go_persistence

import "gorm.io/gorm"

type ApplicationPersistence struct {
	dB *gorm.DB
}

func NewApplicationPersistence(dB *gorm.DB) *ApplicationPersistence {
	return &ApplicationPersistence{dB: dB}
}

// getRandomFunction returns the appropriate random function based on the database dialect
func getRandomFunction(db *gorm.DB) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "RAND()"
	case "postgres":
		return "RANDOM()"
	case "sqlite", "sqlite3":
		return "RANDOM()"
	default:
		// Default to RANDOM(), which is common in many SQL dialects
		return "RANDOM()"
	}
}

func (ap *ApplicationPersistence) GetPaginatedHadithCollections(page int, limit int) ([]*HadithCollection, error) {
	var collections []*HadithCollection
	err := ap.dB.Offset((page - 1) * limit).Limit(limit).Find(&collections).Error
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (ap *ApplicationPersistence) GetHadithCollectionByName(name string) (*HadithCollection, error) {
	var collection HadithCollection
	err := ap.dB.Where("name = ?", name).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (ap *ApplicationPersistence) GetPaginatedBooksByCollection(collection string, page int, limit int) ([]*Book, error) {
	var books []*Book
	err := ap.dB.Where("collection = ?", collection).Offset((page - 1) * limit).Limit(limit).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (ap *ApplicationPersistence) GetBookByCollectionAndBookNumber(collection string, bookNumber string) (*Book, error) {
	var book Book
	err := ap.dB.Where("collection = ? AND ourBookID = ?", collection, bookNumber).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (ap *ApplicationPersistence) GetPaginatedChaptersByCollectionAndBookNumber(collection string, bookNumber string, page int, limit int) ([]*Chapter, error) {
	var chapters []*Chapter
	err := ap.dB.Where("collection = ? AND ourBookID = ?", collection, bookNumber).Offset((page - 1) * limit).Limit(limit).Find(&chapters).Error
	if err != nil {
		return nil, err
	}
	return chapters, nil
}

func (ap *ApplicationPersistence) GetChapterByCollectionAndBookNumberAndChapterNumber(collection string, bookNumber string, chapterNumber string) (*Chapter, error) {
	var chapter *Chapter
	err := ap.dB.Where("collection = ? AND ourBookID = ? AND babID = ?", collection, bookNumber, chapterNumber).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

func (ap *ApplicationPersistence) GetPaginatedHadithsByCollectionAndBookNumber(collection string, bookNumber string, page int, limit int) ([]*Hadith, error) {
	var hadiths []*Hadith
	err := ap.dB.Where("collection = ? AND bookNumber = ?", collection, bookNumber).Offset((page - 1) * limit).Limit(limit).Find(&hadiths).Error
	if err != nil {
		return nil, err
	}
	return hadiths, nil
}

func (ap *ApplicationPersistence) GetHadithByCollectionAndBookNumberAndChapterNumberAndHadithNumber(collection string, bookNumber string, chapterNumber string, hadithNumber string) (*Hadith, error) {
	var hadith *Hadith
	err := ap.dB.Where("collection = ? AND bookNumber = ? AND babID = ? AND hadithNumber = ?", collection, bookNumber, chapterNumber, hadithNumber).First(&hadith).Error
	if err != nil {
		return nil, err
	}
	return hadith, nil
}

func (ap *ApplicationPersistence) GetHadithByCollectionAndHadithNumber(collection, hadithNumber string) (*Hadith, error) {
	var hadith *Hadith
	err := ap.dB.Where("collection = ? AND hadithNumber = ?", collection, hadithNumber).First(&hadith).Error
	if err != nil {
		return nil, err
	}
	return hadith, nil
}

func (ap *ApplicationPersistence) GetHadithByUrn(urn int) (*Hadith, error) {
	var hadith *Hadith
	err := ap.dB.Where("englishURN = ?", urn).First(&hadith).Error
	if err != nil {
		return nil, err
	}
	return hadith, nil
}

func (ap *ApplicationPersistence) GetRandomHadithInCollection(collection string) (*Hadith, error) {
	var hadith *Hadith
	err := ap.dB.Where("collection = ?", collection).Order(getRandomFunction(ap.dB)).First(&hadith).Error
	if err != nil {
		return nil, err
	}
	return hadith, nil
}
