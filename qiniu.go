package bogger

import (
	"io"
	"time"

	"qiniupkg.com/api.v7/kodo"
)

type Qiniu struct {
	conf   Config
	client *kodo.Client
	bucket kodo.Bucket
}

func NewQiniu(conf Config) *Qiniu {
	kodoConfig := &kodo.Config{
		AccessKey: conf.Ak,
		SecretKey: conf.Sk,
		Scheme:    "https",
		RSHost:    "https://rs.qbox.me",
		RSFHost:   "https://rsf.qbox.me",
	}
	client := kodo.New(conf.Zone, kodoConfig)
	return &Qiniu{
		conf:   conf,
		client: client,
		bucket: client.Bucket(conf.Bucket),
	}
}

func (q *Qiniu) Uptoken(key string, lifeMinute uint32, secure bool) string {
	if key != "" {
		key = ":" + key
	}
	uphost := q.conf.UpHost
	if secure {
		uphost = q.conf.UpHostSecure
	}
	if lifeMinute > q.conf.MaxUpLifeMinute || lifeMinute < 1 {
		lifeMinute = q.conf.UpLifeMinute
	}
	putPolicy := &kodo.PutPolicy{
		Scope:   q.conf.Bucket + key,
		UpHosts: []string{uphost},
		Expires: uint32(time.Now().Unix()) + lifeMinute*60,
	}
	return q.client.MakeUptoken(putPolicy)
}

func (q *Qiniu) List(prefix string) (items []kodo.ListItem, err error) {
	items, _, _, err = q.bucket.List(nil, prefix, "", "", 0)
	if err == io.EOF {
		err = nil
	}
	return
}

func (q *Qiniu) Delete(key string) error {
	return q.bucket.Delete(nil, key)
}
