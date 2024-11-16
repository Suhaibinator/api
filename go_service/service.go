package go_service

import "github.com/Suhaibinator/api/go_persistence"

type ApplicationService struct {
	applicationPersistence *go_persistence.ApplicationPersistence
}

func NewApplicationService(applicationPersistence *go_persistence.ApplicationPersistence) *ApplicationService {
	return &ApplicationService{applicationPersistence: applicationPersistence}
}

func (as *ApplicationService) GetPaginatedHadithCollections(page int, limit int) ([]go_persistence.HadithCollection, error) {
	return as.applicationPersistence.GetPaginatedHadithCollections(page, limit)
}

func (as *ApplicationService) GetHadithCollectionByName(name string) (*go_persistence.HadithCollection, error) {
	return as.applicationPersistence.GetHadithCollectionByName(name)
}

func (as *ApplicationService) GetPaginatedBooksByCollection(collection string, page int, limit int) ([]go_persistence.Book, error) {
	return as.applicationPersistence.GetPaginatedBooksByCollection(collection, page, limit)
}

func (as *ApplicationService) GetBookByCollectionAndBookNumber(collection string, bookNumber string) (*go_persistence.Book, error) {
	return as.applicationPersistence.GetBookByCollectionAndBookNumber(collection, bookNumber)
}

func (as *ApplicationService) GetPaginatedChaptersByCollectionAndBookNumber(collection string, bookNumber string, page int, limit int) ([]go_persistence.Chapter, error) {
	return as.applicationPersistence.GetPaginatedChaptersByCollectionAndBookNumber(collection, bookNumber, page, limit)
}

func (as *ApplicationService) GetChapterByCollectionAndBookNumberAndChapterNumber(collection string, bookNumber string, chapterNumber string) (*go_persistence.Chapter, error) {
	return as.applicationPersistence.GetChapterByCollectionAndBookNumberAndChapterNumber(collection, bookNumber, chapterNumber)
}

func (as *ApplicationService) GetHadithsByCollectionAndBookNumberAndChapterNumber(collection string, bookNumber string, chapterNumber string, page int, limit int) ([]go_persistence.Hadith, error) {
	return as.applicationPersistence.GetPaginatedHadithsByCollectionAndBookNumberAndChapterNumber(collection, bookNumber, chapterNumber, page, limit)
}

func (as *ApplicationService) GetHadithByUrn(urn int) (*go_persistence.Hadith, error) {
	return as.applicationPersistence.GetHadithByURN(urn)
}
