package miniopkg

import "github.com/rysmaadit/go-template/config"

type ClientConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
}

func MinioInit() *ClientConfig {
	mysqlConfig := &ClientConfig{
		Endpoint:   config.Init().MinioEndpoint,
		AccessKey:  config.Init().MinioAccessKey,
		SecretKey:  config.Init().MinioSecretKey,
		Region:     config.Init().MinioRegion,
		BucketName: config.Init().MinioBucket,
	}
	return mysqlConfig
}
