package country

type Country struct {
	Country string `json:"country"`
}

type Province struct {
	Country  string `json:"country"`
	Province string `json:"province"`
}

type Cities struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	Cities   []City `json:"cities"`
}
type City struct {
	City string `json:"city"`
}

type Districts struct {
	Province  string     `json:"province"`
	City      string     `json:"city"`
	Districts []District `json:"districts"`
}
type District struct {
	District string `json:"district"`
}

type Zips struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Zips     []Zip  `json:"zips"`
}
type Zip struct {
	Zip string `json:"zip"`
}

type Search struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
}
