package object_storage

type Config struct {
	AccessKey    string
	SecretKey    string
	Region       string
	BaseEndpoint string
	Bucket       string
	MetaFilename string
}
