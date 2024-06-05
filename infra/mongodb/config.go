package mongodb

type Config struct {
	URI      string `validate:"required"`
	Database string `validate:"required"`
}
