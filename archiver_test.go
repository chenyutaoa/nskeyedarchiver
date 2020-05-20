package nskeyedarchiver_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	archiver "github.com/danielpaulus/nskeyedarchiver"
	"github.com/stretchr/testify/assert"
)

func TestArchiver3(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/payload_dump.json")
	if err != nil {
		log.Fatal(err)
	}

	var payloads []string
	json.Unmarshal([]byte(dat), &payloads)

	plistBytes, _ := hex.DecodeString(payloads[0])
	nska, err := archiver.Unarchive(plistBytes)
	value := nska[0]
	result, err := archiver.ArchiveBin(value)
	/*if assert.NoError(t, err) {
		output := convertToJSON(nska)
		print(output)
		assert.Equal(t, plistBytes, result)
		assert.NoError(t, err)
	}*/
	nska2, err2 := archiver.Unarchive(result)
	if assert.NoError(t, err2) {
		assert.Equal(t, nska2, nska)
	}

}

//TestDecoderJson tests if real DTX nsKeyedArchived plists can be decoded without error
func TestArchiver(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/payload_dump.json")
	if err != nil {
		log.Fatal(err)
	}

	var payloads []string
	json.Unmarshal([]byte(dat), &payloads)
	for _, plistHex := range payloads {
		plistBytes, _ := hex.DecodeString(plistHex)
		expectedNska, _ := archiver.Unarchive(plistBytes)

		archivedNskaBin, err := archiver.ArchiveBin(expectedNska[0])
		archivedNskaXml, err2 := archiver.ArchiveXML(expectedNska[0])

		if assert.NoError(t, err) && assert.NoError(t, err2) {
			actualNskaBin, unarchiverErrBin := archiver.Unarchive(archivedNskaBin)
			actualNskaXml, unarchiverErrXml := archiver.Unarchive([]byte(archivedNskaXml))
			if assert.NoError(t, unarchiverErrBin) && assert.NoError(t, unarchiverErrXml) {
				assert.Equal(t, expectedNska, actualNskaBin)
				assert.Equal(t, expectedNska, actualNskaXml)
			}
		}
	}
}

//TestDecoderJson tests if real DTX nsKeyedArchived plists can be decoded without error
func TestDecoderJson(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/payload_dump.json")
	if err != nil {
		log.Fatal(err)
	}

	var payloads []string
	json.Unmarshal([]byte(dat), &payloads)
	for _, plistHex := range payloads {
		plistBytes, _ := hex.DecodeString(plistHex)
		nska, err := archiver.Unarchive(plistBytes)
		output := convertToJSON(nska)
		print(output)
		assert.NoError(t, err)
	}
}

func TestDecoder(t *testing.T) {
	testCases := map[string]struct {
		filename string
		expected string
	}{
		"test one value":       {"onevalue", "[true]"},
		"test all primitives":  {"primitives", "[1,1,1,1.5,\"YXNkZmFzZGZhZHNmYWRzZg==\",true,\"Hello, World!\",\"Hello, World!\",\"Hello, World!\",false,false,42]"},
		"test arrays and sets": {"arrays", "[[1,1,1,1.5,\"YXNkZmFzZGZhZHNmYWRzZg==\",true,\"Hello, World!\",\"Hello, World!\",\"Hello, World!\",false,false,42],[true,\"Hello, World!\",42],[true],[42,true,\"Hello, World!\"]]"},
		"test nested arrays":   {"nestedarrays", "[[[true],[42,true,\"Hello, World!\"]]]"},
		"test dictionaries":    {"dict", "[{\"array\":[true,\"Hello, World!\",42],\"int\":1,\"string\":\"string\"}]"},
	}

	for _, tc := range testCases {
		dat, err := ioutil.ReadFile("fixtures/" + tc.filename + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		objects, err := archiver.Unarchive(dat)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, convertToJSON(objects))

		dat, err = ioutil.ReadFile("fixtures/" + tc.filename + ".bin")
		if err != nil {
			log.Fatal(err)
		}
		objects, err = archiver.Unarchive(dat)
		assert.Equal(t, tc.expected, convertToJSON(objects))
	}
}

func TestValidation(t *testing.T) {

	testCases := map[string]struct {
		filename string
	}{
		"$archiver key is missing":         {"missing_archiver"},
		"$archiver is not nskeyedarchiver": {"wrong_archiver"},
		"$top key is missing":              {"missing_top"},
		"$objects key is missing":          {"missing_objects"},
		"$version key is missing":          {"missing_version"},
		"$version is wrong":                {"wrong_version"},
		"plist is invalid":                 {"broken_plist"},
	}

	for _, tc := range testCases {
		dat, err := ioutil.ReadFile("fixtures/" + tc.filename + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		_, err = archiver.Unarchive(dat)
		assert.Error(t, err)
	}
}

func convertToJSON(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}
