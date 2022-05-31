package gotile

import (
	"GoTile/config"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var isCache bool = config.Configiure.Server.Cache
var cacheDir string = config.Configiure.Server.CacheDir

func MakeCache(x, y, z int64, pbf []byte, tablename string) error {
	if !isCache {
		return errors.New("cache is disable")
	}

	pbfCachePath := fmt.Sprintf("%s/%s/%d", cacheDir, tablename, z)
	err := os.MkdirAll(pbfCachePath, 777)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pdfCacheFile := fmt.Sprintf("%s/%d_%d.pbf", pbfCachePath, x, y)
	err = ioutil.WriteFile(pdfCacheFile, pbf, 777)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetCache(x, y, z int64, tablename string) ([]byte, error) {
	if !isCache {
		return nil, errors.New("cache is disable")
	}
	pbfCacheFile := fmt.Sprintf("%s/%s/%d/%d_%d.pbf", cacheDir, tablename, z, x, y)
	return ioutil.ReadFile(pbfCacheFile)
}
