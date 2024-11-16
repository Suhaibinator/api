package go_persistence

import "gorm.io/gorm"

type ApplicationPersistence struct {
	dB *gorm.DB
}

func NewApplicationPersistence(dB *gorm.DB) *ApplicationPersistence {
	return &ApplicationPersistence{dB: dB}
}

func (ap *ApplicationPersistence) GetPaginatedHadithCollections(page int, limit int) ([]HadithCollection, error) {
	var collections []HadithCollection
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

func (ap *ApplicationPersistence) GetPaginatedBooksByCollection(collection string, page int, limit int) ([]Book, error) {
	var books []Book
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

func (ap *ApplicationPersistence) GetPaginatedChaptersByCollectionAndBookNumber(collection string, bookNumber string, page int, limit int) ([]Chapter, error) {
	var chapters []Chapter
	err := ap.dB.Where("collection = ? AND arabicBookID = ?", collection, bookNumber).Offset((page - 1) * limit).Limit(limit).Find(&chapters).Error
	if err != nil {
		return nil, err
	}
	return chapters, nil
}

func (ap *ApplicationPersistence) GetChapterByCollectionAndBookNumberAndChapterNumber(collection string, bookNumber string, chapterNumber string) (*Chapter, error) {
	var chapter Chapter
	err := ap.dB.Where("collection = ? AND arabicBookID = ? AND babID = ?", collection, bookNumber, chapterNumber).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (ap *ApplicationPersistence) GetPaginatedHadithsByCollectionAndBookNumberAndChapterNumber(collection string, bookNumber string, chapterNumber string, page int, limit int) ([]Hadith, error) {
	var hadiths []Hadith
	err := ap.dB.Where("collection = ? AND bookNumber = ? AND babID = ?", collection, bookNumber, chapterNumber).Offset((page - 1) * limit).Limit(limit).Find(&hadiths).Error
	if err != nil {
		return nil, err
	}
	return hadiths, nil
}

func (ap *ApplicationPersistence) GetHadithByCollectionAndHadithNumber(collection string, hadithNumber string) (*Hadith, error) {
	var hadith Hadith
	err := ap.dB.Where("collection = ? AND hadithNumber = ?", collection, hadithNumber).First(&hadith).Error
	if err != nil {
		return nil, err
	}
	return &hadith, nil
}

func (ap *ApplicationPersistence) GetPaginatedHadiths(page int, limit int) ([]Hadith, error) {
	var hadiths []Hadith
	err := ap.dB.Offset((page - 1) * limit).Limit(limit).Find(&hadiths).Error
	if err != nil {
		return nil, err
	}
	return hadiths, nil
}

func (ap *ApplicationPersistence) GetHadithByURN(urn int) (*Hadith, error) {
	var hadith Hadith
	err := ap.dB.Where("englishURN = ?", urn).First(&hadith).Error
	if err != nil {
		return nil, err
	}
	return &hadith, nil
}
