package archiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	plist "howett.net/plist"
)

//toUidList type asserts a []interface{} to a []plist.UID by iterating through the list.
func toUidList(list []interface{}) []plist.UID {
	l := len(list)
	result := make([]plist.UID, l)
	for i := 0; i < l; i++ {
		result[i] = list[i].(plist.UID)
	}
	return result
}

//plistFromBytes decodes a binary or XML based PLIST using the amazing github.com/DHowett/go-plist library and returns an interface{} or propagates the error raised by the library.
func plistFromBytes(plistBytes []byte) (interface{}, error) {
	var test interface{}
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))

	err := decoder.Decode(&test)
	if err != nil {
		return test, err
	}
	return test, nil
}

//ToPlist converts a given struct to a Plist using the
//github.com/DHowett/go-plist library. Make sure your struct is exported.
//It returns a string containing the plist.
func ToPlist(data interface{}) string {
	buf := &bytes.Buffer{}
	encoder := plist.NewEncoder(buf)
	encoder.Encode(data)
	return buf.String()
}

//Print an object as JSON for debugging purposes, careful log.Fatals on error
func printAsJSON(obj interface{}) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf("Error while marshalling Json:%s", err)
	}
	fmt.Print(string(b))
}

//verifyCorrectArchiver makes sure the nsKeyedArchived plist has all the right keys and values and returns an error otherwise
func verifyCorrectArchiver(nsKeyedArchiverData map[string]interface{}) error {
	if val, ok := nsKeyedArchiverData[archiverKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", archiverKey)
	} else {
		if stringValue := val.(string); stringValue != nsKeyedArchiver {
			return fmt.Errorf("Invalid value: %s for key '%s', expected: '%s'", stringValue, archiverKey, nsKeyedArchiver)
		}
	}
	if _, ok := nsKeyedArchiverData[topKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", topKey)
	}

	if _, ok := nsKeyedArchiverData[objectsKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", objectsKey)
	}

	if val, ok := nsKeyedArchiverData[versionKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", versionKey)
	} else {
		if stringValue := val.(uint64); stringValue != versionValue {
			return fmt.Errorf("Invalid value: %d for key '%s', expected: '%d'", stringValue, versionKey, versionValue)
		}
	}

	return nil
}
