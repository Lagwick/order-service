package section

type Client struct {
	Catalog Catalog `split_words:"true"`
}

type Catalog struct {
	Address string `env:"ADDRESS" envDefault:"localhost:50051"`
}
