package bogger

type Config struct {
	Zone            int
	Ak              string `validate:"required"`
	Sk              string `validate:"required"`
	Bucket          string `validate:"required"`
	UpLifeMinute    uint32 `default:"30ns"`
	MaxUpLifeMinute uint32 `default:"60ns"`
	UpHost          string `default:"http://upload.qiniu.com"`
	UpHostSecure    string `default:"https://up.qbox.me"`
}
