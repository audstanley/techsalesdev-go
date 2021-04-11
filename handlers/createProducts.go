package handlers

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func PopulateDatabaseWithImages() {
	files, err := ioutil.ReadDir("./.products/photos")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// Get the extension for the photo file
		extension := filepath.Ext(f.Name())
		epoch := uint64(time.Now().Unix())

		// hash the file name (used as key in redis)
		md5Hash := getMD5Hash(f.Name())
		//fmt.Println(md5Hash)
		_, err := ioutil.ReadFile(".products/photos/" + f.Name())
		check(err)
		file, err := os.Open(".products/photos/" + f.Name())
		check(err)
		reader := bufio.NewReader(file)
		content, _ := ioutil.ReadAll(reader)
		encoded := "data:image/" + extension + ";base64, " + base64.StdEncoding.EncodeToString(content)

		p := &Product{
			Name:        RandStringBytes(8),
			Cost:        float64(rand.Intn(50-1) + 1),
			Description: RandStringBytes(8),
			Image:       encoded,
			InStock:     RandBool(),
			ProductId:   md5Hash,
			OnSale:      RandBool(),
			Iat:         epoch,
			Category:    RandStringBytes(8),
		}
		b, err := json.Marshal(p)
		err = ProductsClient.Set(RedisCtx, md5Hash, string(b), 0).Err()
		CheckRedisErrWithoutContext(err)
		file.Close()
	}

}
