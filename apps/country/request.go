package country

type Params struct {
	Province string `query:"province"`
	City     string `query:"city"`
	District string `query:"district"`
	Zip      string `query:"zip"`
	Children int    `query:"children"`
	Search   string `query:"search"`
}
