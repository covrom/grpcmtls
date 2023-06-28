package grpcmtls

type CertFiles struct {
	ServerCertPemFilePath,
	ServerKeyPemFilePath,
	CAFilePath,
	ClientCertPemFilePath,
	ClientKeyPemFilePath string
}

func (cf *CertFiles) ServerCertFiles() ServerCertFiles {
	return ServerCertFiles{
		ServerCertPemFilePath: cf.ServerCertPemFilePath,
		ServerKeyPemFilePath:  cf.ServerKeyPemFilePath,
		CAFilePath:            cf.CAFilePath,
	}
}

func (cf *CertFiles) ClientCertFiles() ClientCertFiles {
	return ClientCertFiles{
		CAFilePath:            cf.CAFilePath,
		ClientCertPemFilePath: cf.ClientCertPemFilePath,
		ClientKeyPemFilePath:  cf.ClientKeyPemFilePath,
	}
}

type ServerCertFiles struct {
	ServerCertPemFilePath,
	ServerKeyPemFilePath,
	CAFilePath string
}

type ClientCertFiles struct {
	CAFilePath,
	ClientCertPemFilePath,
	ClientKeyPemFilePath string
}
