package section

type Client struct {
	Catalog Catalog `split_words:"true"`
}

type Catalog struct {
	GrpcAddress string `split_words:"true"`
}
