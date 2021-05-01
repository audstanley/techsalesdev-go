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

func populateCategories(m *map[string]string) {
	(*m)["gameboy.jpeg"] = "PCB"
	(*m)["joystick.jpg"] = "PCB"
	(*m)["pcb_01.png"] = "PCB"
	(*m)["pcb_02.png"] = "PCB"
	(*m)["pcb_03.png"] = "PCB"
	(*m)["pcb_04.png"] = "PCB"
	(*m)["pcb_05.png"] = "PCB"
	(*m)["pcb_06.png"] = "PCB"
	(*m)["pcb_07.png"] = "PCB"
	(*m)["pcb_08.png"] = "PCB"
	(*m)["pcb_09.png"] = "PCB"
	(*m)["pcb_10.png"] = "PCB"
	(*m)["wires_01.jpg"] = "Wires"
	(*m)["wires_02.jpg"] = "Wires"
	(*m)["wires_03.png"] = "Wires"
	(*m)["wires_04.png"] = "Wires"
	(*m)["wires_05.png"] = "Wires"
	(*m)["wires_06.png"] = "Wires"
	(*m)["wires_07.png"] = "Wires"
	(*m)["wires_08.png"] = "Wires"
	(*m)["wires_09.png"] = "Wires"
	(*m)["wires_10.png"] = "Wires"
	(*m)["wires_11.png"] = "Wires"
	(*m)["wires_12.png"] = "Wires"
	(*m)["wires_13.png"] = "Wires"
	(*m)["wires_14.png"] = "Wires"
	(*m)["wires_15.png"] = "Wires"
	(*m)["wires_16.png"] = "Wires"
	(*m)["wires_17.png"] = "Wires"
	(*m)["diodes_01.jpg"] = "Diodes"
	(*m)["diodes_02.jpg"] = "Diodes"
	(*m)["diodes_03.jpg"] = "Diodes"
	(*m)["diodes_04.jpg"] = "Diodes"
	(*m)["diodes_05.jpg"] = "Diodes"
	(*m)["diodes_06.jpg"] = "Diodes"
	(*m)["diodes_07.jpg"] = "Diodes"
	(*m)["diodes_08.jpg"] = "Diodes"
	(*m)["diodes_09.jpg"] = "Diodes"
	(*m)["diodes_10.jpg"] = "Diodes"
	(*m)["diodes_11.jpg"] = "Diodes"
	(*m)["diodes_12.jpg"] = "Diodes"
	(*m)["diodes_13.jpg"] = "Diodes"
	(*m)["diodes_14.jpg"] = "Diodes"
	(*m)["diodes_15.jpg"] = "Diodes"
	(*m)["diodes_16.jpg"] = "Diodes"
	(*m)["diodes_17.jpg"] = "Diodes"
	(*m)["capacitors_01.jpg"] = "Caps"
	(*m)["capacitors_02.png"] = "Caps"
	(*m)["capacitors_03.jpg"] = "Caps"
	(*m)["capacitors_04.jpg"] = "Caps"
	(*m)["capacitors_05.png"] = "Caps"
	(*m)["capacitors_06.png"] = "Caps"
	(*m)["capacitors_07.jpg"] = "Caps"
	(*m)["capacitors_08.jpg"] = "Caps"
	(*m)["capacitors_09.jpg"] = "Caps"
	(*m)["capacitors_10.jpg"] = "Caps"
	(*m)["capacitors_11.jpg"] = "Caps"
	(*m)["capacitors_12.png"] = "Caps"
	(*m)["capacitors_13.png"] = "Caps"
	(*m)["capacitors_14.png"] = "Caps"
	(*m)["capacitors_15.png"] = "Caps"
	(*m)["capacitors_16.jpg"] = "Caps"
	(*m)["capacitors_17.png"] = "Caps"

}

func putInRandomCategory() string {
	var categories = []string{"PCB", "Soldering", "Wires", "Caps", "Diodes"}
	return categories[rand.Intn(4-0)]
}

func PopulateDatabaseWithImages() {
	m := map[string]string{}
	populateCategories(&m)
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
		c := m[f.Name()]
		if c == "" {
			c = putInRandomCategory()
		}

		p := &Product{
			Name:        RandStringBytes(8),
			Cost:        float64(rand.Intn(50-1) + 1),
			Description: RandStringBytes(8),
			Image:       encoded,
			InStock:     RandBool(),
			ProductId:   md5Hash,
			OnSale:      RandBool(),
			Iat:         epoch,
			Category:    c,
		}

		b, err := json.Marshal(p)
		err = ProductsClient.Set(RedisCtx, md5Hash, string(b), 0).Err()
		CheckRedisErrWithoutContext(err)
		if p.Category == "PCB" {
			err = PCBClient.Set(RedisCtx, md5Hash, string(b), 0).Err()
			CheckRedisErrWithoutContext(err)
		} else if p.Category == "Wires" {
			err = WiresClient.Set(RedisCtx, md5Hash, string(b), 0).Err()
			CheckRedisErrWithoutContext(err)
		} else if p.Category == "Diodes" {
			err = DiodesClient.Set(RedisCtx, md5Hash, string(b), 0).Err()
			CheckRedisErrWithoutContext(err)
		} else if p.Category == "Caps" {
			err = CapsClient.Set(RedisCtx, md5Hash, string(b), 0).Err()
			CheckRedisErrWithoutContext(err)
		}

		file.Close()
	}

}
