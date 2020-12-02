package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const serviceAccount = `{ "type": "service_account", "project_id": "elementalserver-8c6d0", "private_key_id": "78ba8b0bb00e5233e4ac4cb5e640ec6d6c56eb6b", "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDEw7lgdH3/H5F2\nMUkY6ImrSkS9sqvKmgAlW1CLxn8p/V4yeSbFvFiUv9mBTlKYHAnE/wmez8F1OyOT\ncnvG0AFicmkV/2QJSjkW8yvkFB7S2Br7rWkPsmacA1RlWHlctybSUo9/HkooDl5j\nkqtbJsT2Vt/nWlsDMmZTKTCJoRQuP8jATbeA16XcQ/Gl8mFap1lC89OE7wxrx5Eq\nTKqt4AkhKhBQhoNr9VTFRCBkhaBMq830/MKaZpbLDVa6qTk4c9QOs9hohtpUAup6\n2nvxX6k9l/K9nPJI9Gp7lOhQ951/yws8c+DGbuSdjHnj1sKQlWmOY6hMkS2i4YHX\nWWcEheYnAgMBAAECggEAGb+DApw74KbA4jaQ2jGT0lZlqG05DcoZOso4QBI5kcUW\nDoTMDhQXg1+XltQo+r6wiJbXK3EEX9LdVO4mRF3z0G4oUjiZXp3X2qj3lWEMp4qf\n/U8z8FnoE4JcCOcK+pb8/YjQPlI4YgV/VIhc5BCutY2ovx2Ty1dNDJTXRStO+L4k\n1rnkaL/zwbxUUBLEb4rYuAjw7cOyxk081yXdogs34IQd/9aLVZnIri5A6rSDYFFw\n1A2aJzNBolur+/5m57dT6OqVjgD99zjWEJ9lq61oNMtYVetGY9bFis59JnUdOUHg\n11B2RMMOAu4lqDaCFEvmqiw09OEVgynlA72XUhaAAQKBgQD9+FMOjM9dai9RoDBZ\nrgEQQfqwPXS1z24Fqlx1zBgI6wUCvVU1qfnMBCQsrAwzfev+nt44yOBctuPr7ZUC\nvBvSGO6Dd50Wwnd24wHrhd+45FX6pHbdSItbkpMyuA6OHHigXJNKaHib6G7EqHs5\nIK88GrawYtYD2abWfQJwtJhTgQKBgQDGVlhqyOrGfBub1NSFsey5oltCH7KuMV1F\nozNoOlpU4EtOaYE1g71+9uNSf3n1vO8PLCLj8ZoF2fk6VV2ttME2fXUjfZz9HUgw\ncfOZpD5j3XM/s8mJFncbR+5eC/KAoXA0oUVvhkffdEaGtrInrobOOcyETGqdPcjA\n6oAvKFntpwKBgQCWJDA18djFiPjgcKsk2VGXounpNuvAcBjDEKwIl9e9rfMQY430\nY8BhdDFOl4e/CTpzFMibGWZKaXTlDVeCfmKUGlknL5eW1PB7QEjqTAKu845A1unO\neAyq3kRXP6ibKwnFA/Wvj4N96DNT36a5ZzExfzlxnXyYWhvfwZenuZw0AQKBgH0c\nqKer2BWe4leZmPpBM4AiP4jlr/QcJacw/NOpw6O43Sg4e45DbTzzBpDa4xc1uGOM\nxvGdTTiVuJaolPBnjl4OI99gdLBiUVBmAXGQ3t5mKjYr9lyotDecV2wyAyZLMBmz\nBbcFML9vfLGr+5P2jwj2AuINxk8sU0AGbRfST3APAoGBAMOJsmiBoXBiOxi52lyM\nZN8jxyTHd9LwnTgPHQd2JedBi7EIJ3j3T+QP3Z3SENMMImQr6MOda8otrTyqpMTp\nDS+pTomSwTCCEir7bSVpi7QMejchURVYM/PmMwhso1vocZBM3YHvxLtGAnFOu7BM\nQ39vHDC9jyj00STzo/+fD6X3\n-----END PRIVATE KEY-----\n", "client_email": "firebase-adminsdk-7nmm6@elementalserver-8c6d0.iam.gserviceaccount.com", "client_id": "113854670633531537114", "auth_uri": "https://accounts.google.com/o/oauth2/auth", "token_uri": "https://oauth2.googleapis.com/token", "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs", "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-7nmm6%40elementalserver-8c6d0.iam.gserviceaccount.com" }`

var colors map[string]map[string]string

// Element has the data for a created element
type Element struct {
	Color     string   `json:"color"`
	Comment   string   `json:"comment"`
	CreatedOn int      `json:"createdOn"`
	Creator   string   `json:"creator"`
	Name      string   `json:"name"`
	Parents   []string `json:"parents"`
	Pioneer   string   `json:"pioneer"`
}

type element struct {
	Name    string
	Parents []element
	Color   string
}

const elemName string = "Life"

var cache map[string]Element = make(map[string]Element, 0)

func main() {
	opt := option.WithCredentialsJSON([]byte(serviceAccount))
	config := &firebase.Config{
		DatabaseURL:   "https://elementalserver-8c6d0.firebaseio.com",
		ProjectID:     "elementalserver-8c6d0",
		StorageBucket: "elementalserver-8c6d0.appspot.com",
	}
	var err error
	fireapp, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		panic(err)
	}
	store, err := fireapp.Firestore(context.Background())
	if err != nil {
		panic(err)
	}

	colorData, _ := ioutil.ReadFile("colors.json")
	json.Unmarshal(colorData, &colors)

	thing := calcElem(elemName, store)
	out, err := json.Marshal(thing)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

func calcElem(name string, store *firestore.Client) element {
	_, exists := cache[name]
	var elem Element
	if exists {
		elem = cache[name]
	} else {
		dat, err := store.Collection("elements").Doc(name).Get(context.Background())
		if err != nil {
			panic(err)
		}
		dat.DataTo(&elem)
		cache[name] = elem
	}
	thing := element{
		Name:    elem.Name,
		Color:   colors[strings.Split(elem.Color, "_")[0]]["color"],
		Parents: make([]element, 0),
	}
	for _, parent := range elem.Parents {
		thing.Parents = append(thing.Parents, calcElem(parent, store))
	}
	if (len(thing.Parents) > 0) && (thing.Parents[0].Name == thing.Parents[1].Name) {
		thing.Parents = []element{thing.Parents[0]}
	}
	return thing
}
