package config

type Config struct {
	Database struct {
		Mongo struct {
			DataSource string
			DB         string

			NFTCollection string
			NFT           string
			Tx            string
		}
	}
}
