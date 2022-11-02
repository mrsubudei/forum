package usecase

type Repository interface {
	Get(jsonBlob []byte) ([]byte, error)
	Create(jsonBlob []byte) ([]byte, error)
	Change(jsonBlob []byte) ([]byte, error)
	Delete(jsonBlob []byte) error
}
