package bogger

type Config struct {
	Zone         int
	Ak           string
	Sk           string
	Bucket       string
	UpLifeMinute uint32 `default:"30"`
	UpHost       string `default:"https://up.qbox.me"`
}
