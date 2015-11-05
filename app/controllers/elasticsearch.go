package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mattbaird/elastigo/lib"
)

type ESDB struct {
	Conn  *elastigo.Conn
	Index string
	Type  string
}

func (k KeyDoc) RenderJSON() string {
	json, err := json.MarshalIndent(k, "", "    ")
	if err != nil {
		return ""
	}
	return string(json)
}

func NewIndex(domain string) *ESDB {
	es := new(ESDB)

	es.Conn = elastigo.NewConn()

	// Set the Elasticsearch Host to Connect to
	es.Conn.Domain = domain
	es.Index = "cipher_keys"
	es.Type = "key_doc"

	es.createIndex()

	return es
}

func (es *ESDB) IndexKey(cipher *Cipher) error {

	// index, _type, id, args, data
	_, err := es.Conn.Index(es.Index, es.Type, cipher.KeyID, nil, KeyDoc{
		Translation:     cipher.Translation,
		Key:             cipher.Key,
		KeyID:           cipher.KeyID,
		CipherName:      cipher.Name,
		Timestamp:       time.Now().Unix(),
		FoundWordsTotal: cipher.FoundWordsTotal,
		FoundWords:      cipher.FoundWords,
		KillCount:       cipher.KillCount,
	})

	return err

}

func (es *ESDB) IndexKeyDoc(keyDoc KeyDoc) error {

	// index, _type, id, args, data
	_, err := es.Conn.Index(es.Index, es.Type, keyDoc.KeyID, nil, keyDoc)

	return err

}

func (es *ESDB) GetKeyByDate(cipher *Cipher, offset int) (KeyDoc, error) {

	keyDoc := KeyDoc{}

	resp, err := es.GetKeys(cipher, 1, offset, "Timestamp")

	if err != nil || resp.Hits.Len() < 1 {
		return keyDoc, errors.New("No Keys Found!")
	}

	hit, err := resp.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return keyDoc, err
	}

	json.Unmarshal(hit, &keyDoc)

	return keyDoc, nil
}

func (es *ESDB) GetKeyByWordcount(cipher *Cipher, offset int) (KeyDoc, error) {

	keyDoc := KeyDoc{}

	resp, err := es.GetKeys(cipher, 1, offset, "FoundWordsTotal")

	if err != nil || resp.Hits.Len() < 1 {
		return keyDoc, errors.New("No Keys Found!")
	}

	hit, err := resp.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return keyDoc, err
	}

	json.Unmarshal(hit, &keyDoc)

	return keyDoc, nil
}

func (es *ESDB) GetKeyByHash(hash string) (KeyDoc, error) {

	keyDoc := KeyDoc{}

	resp, err := es.Conn.Get(es.Index, es.Type, hash, nil)
	if err != nil {
		return keyDoc, err
	}

	hit, err := resp.Source.MarshalJSON()

	if err != nil {
		return keyDoc, err
	}

	json.Unmarshal(hit, &keyDoc)

	return keyDoc, nil

}

func (es *ESDB) GetKeys(cipher *Cipher, size int, from int, sort string) (elastigo.SearchResult, error) {

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
	case "FoundWordsTotal":
		searchJson = searchJson + `"sort" : {
					"FoundWordsTotal" : {
						"order" : "desc"
					}
				}`
	}

	searchJson = searchJson + `}`

	return es.Conn.Search(es.Index, es.Type, nil, searchJson)

}

func (es *ESDB) GetKeyCount(cipher *Cipher) (int, error) {

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

	resp, err := es.Conn.Count(es.Index, es.Type, nil, searchJson)

	return resp.Count, err

}

func (es *ESDB) DeleteIndex(index string) (elastigo.BaseResponse, error) {

	return es.Conn.DeleteIndex(index)
}

func (es *ESDB) createIndex() error {

	mapping := `{
					"` + es.Type + `": {
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
						    "FoundWordsTotal": {
						        "type": "long"
						    }
						}
					}
				}`

	es.Conn.CreateIndex(es.Index)

	return es.Conn.PutMappingFromJSON(es.Index, es.Type, []byte(mapping))

}
