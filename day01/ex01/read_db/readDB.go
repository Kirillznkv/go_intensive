package read_db

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Recipes struct {
	XMLName  xml.Name `xml:"recipes" json:"-"`
	CakeList []Cake   `xml:"cake" json:"cake"`
}

type Cake struct {
	XMLName    xml.Name `xml:"cake" json:"-"`
	Name       string   `xml:"name" json:"name"`
	Time       string   `xml:"stovetime" json:"time"`
	Ingredient []Item   `xml:"ingredients>item" json:"ingredients"`
}

type Item struct {
	Name  string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}

func readStrDB(fName *string) string {
	var byteData []byte
	file, err := os.Open(*fName)
	if err != nil {
		log.Fatal(err)
	}
	if byteData, err = ioutil.ReadAll(file); err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(byteData)
}

func getFileFormat(fName *string) string {
	sp := strings.Split((*fName), ".")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}

func getDataFromJson(strData string) Recipes {
	var data Recipes
	if err := json.Unmarshal([]byte(strData), &data); err != nil {
		log.Fatal("ERROR: Json file")
	}
	return data
}

func getDataFromXml(strData string) Recipes {
	var data Recipes
	if err := xml.Unmarshal([]byte(strData), &data); err != nil {
		log.Fatal("ERROR: Xml file")
	}
	return data
}

func OutputJson(data *Recipes) {
	jsonData, _ := json.MarshalIndent(*data, "", "    ")
	fmt.Println(string(jsonData))
}

func OutputXml(data *Recipes) {
	xmlData, _ := xml.MarshalIndent(*data, "", "    ")
	fmt.Println(string(xmlData))
}

func ReadDB(fName *string) Recipes {
	var data Recipes
	strData := readStrDB(fName)
	format := getFileFormat(fName)
	if format == "json" {
		data = getDataFromJson(strData)
	} else if format == "xml" {
		data = getDataFromXml(strData)
	} else {
		fmt.Println("ERROR: file format")
	}
	return data
}
