package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mattbaird/elastigo/lib"
	"github.com/revel/revel"
)

type KeyDoc struct {
	CipherName  string
	Translation string
	Key         []Character
	KeyID       string
	WordCount   int
	WordList    []string
	Timestamp   int64
}

func IndexKey(cipher *Cipher) {

	c := elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	c.Domain = revel.Config.StringDefault("elasticsearch.host", "localhost")

	index := "cipher_keys"
	_type := "key_doc"

	// index, _type, id, args, data
	_, err := c.Index(index, _type, cipher.KeyID, nil, KeyDoc{
		Translation: cipher.Translation,
		Key:         cipher.Key,
		KeyID:       cipher.KeyID,
		CipherName:  cipher.Name,
		Timestamp:   time.Now().Unix(),
	})
	if err != nil {
		revel.ERROR.Print(err)
	}

}

func GetNewestKey(cipher *Cipher) (KeyDoc, error) {
	resp, err := GetKeys(cipher, 1, 0, "Timestamp")

	if err != nil || resp.Hits.Len() < 1 {
		return KeyDoc{}, err
	}

	hit, err := resp.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return KeyDoc{}, err
	}

	keyDoc := KeyDoc{}

	json.Unmarshal(hit, &keyDoc)

	return keyDoc, nil
}

func GetBestKey(cipher *Cipher) (KeyDoc, error) {
	resp, err := GetKeys(cipher, 1, 0, "WordCount")

	if err != nil || resp.Hits.Len() < 1 {
		return KeyDoc{}, err
	}

	hit, err := resp.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return KeyDoc{}, err
	}

	keyDoc := KeyDoc{}

	json.Unmarshal(hit, &keyDoc)

	return keyDoc, nil
}

func GetKeyByHash(hash string) (KeyDoc, error) {
	c := elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	c.Domain = revel.Config.StringDefault("elasticsearch.host", "localhost")

	index := "cipher_keys"
	_type := "key_doc"

	resp, err := c.Get(index, _type, hash, nil)
	if err != nil {
		return KeyDoc{}, err
	}

	hit, err := resp.Source.MarshalJSON()

	if err != nil {
		return KeyDoc{}, err
	}

	keyDoc := KeyDoc{}

	json.Unmarshal(hit, &keyDoc)

	return keyDoc, nil

}

func GetKeys(cipher *Cipher, size int, from int, sort string) (elastigo.SearchResult, error) {

	c := elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	c.Domain = revel.Config.StringDefault("elasticsearch.host", "localhost")

	index := "cipher_keys"
	_type := "key_doc"

	searchJson := `{
			"from" : ` + fmt.Sprint(from) + `,
			"size" : ` + fmt.Sprint(size) + `,
		    "query": {
				"filtered": {
					"query": {
						"bool": {
							"must": [
							{
								"match": {
									"CipherName": "` + cipher.Name + `"
								}
							}
							]
						}
					}
				}
			},`

	switch sort {
	case "Timestamp":
		searchJson = searchJson + `"sort" : {
					"Timestamp" : {
						"order" : "desc"
					}
				}`
	case "WordCount":
		searchJson = searchJson + `"sort" : {
					"WordCount" : {
						"order" : "desc"
					}
				}`
	}

	searchJson = searchJson + `}`

	return c.Search(index, _type, nil, searchJson)

}

func GetKeyCount(cipher *Cipher) (int, error) {

	c := elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	c.Domain = revel.Config.StringDefault("elasticsearch.host", "localhost")

	index := "cipher_keys"
	_type := "key_doc"

	searchJson := `{
    	"query": {
			"filtered": {
				"query": {
					"bool": {
						"must": [
						{
							"match": {
								"CipherName": "` + cipher.Name + `"
							}
						}
						]
					}
				}
			}
		}
	}`

	resp, err := c.Count(index, _type, nil, searchJson)

	return resp.Count, err

}

func DeleteIndex(index string) (elastigo.BaseResponse, error) {

	c := elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	c.Domain = revel.Config.StringDefault("elasticsearch.host", "localhost")

	return c.DeleteIndex(index)
}

func CreateIndex() {

	c := elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	c.Domain = revel.Config.StringDefault("elasticsearch.host", "localhost")

	index := "cipher_keys"
	_type := "key_doc"

	mapping := `{
					"` + _type + `": {
						"_id": {
							"index": "analyzed",
							"path": "id"
						},
						"properties": {
						    "CipherName": {
						        "type": "string"
						    },
						    "Translation": {
						    	"type": "string"
						    },
						    "Key": {
						        "properties": {
						            "Letter": {
						                "type": "string"
						            },
						            "Symbol": {
						                "type": "string"
						            }
						        }
						    },
						    "Timestamp": {
						    	"type": "long"
						    },
						    "WordCount": {
						        "type": "long"
						    }
						}
					}
				}`

	c.CreateIndex(index)

	err := c.PutMappingFromJSON(index, _type, []byte(mapping))
	if err != nil {
		revel.ERROR.Print(err)
	}

}
