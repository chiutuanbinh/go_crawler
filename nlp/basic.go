package nlp

import (
	"bytes"
	"crawler/util"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var nlpDomain = util.Config.NLP.Domain
var nlpTokenizeURL = nlpDomain + "/tokenize"
var nlpNERUrl = nlpDomain + "/ner"
var timeout = 5 * time.Second

const (
	Noun              = "N"
	Verb              = "V"
	None              = "O"
	NounPhrase        = "Np"
	VerbPhrase        = "Vp"
	Punctuation       = "CH"
	BeginNounPhrase   = "B-NP"
	InNounPhrase      = "I-NP"
	BeginVerbPhrase   = "B-VP"
	InVerbPhrase      = "I-VP"
	BeginLocation     = "B-LOC"
	InLocation        = "I-LOC"
	BeginOrganization = "B-ORG"
	InOrganization    = "I-ORG"
	BeginPerson       = "B-PER"
	InPerson          = "I-PER"
)

type NamedEntitiesT struct {
	Text string
	Type string
}

type NLPResp struct {
	NounPhrases   []string
	NamedEntities []NamedEntitiesT
}

var empty = NLPResp{}

var nerTypeMapper = make(map[string]string)

const (
	LOCATION     = "LOCATION"
	PERSON       = "PERSON"
	ORGANIZATION = "ORGANIZATION"
	UNKNOWN      = "UNKNOWN"
)

func init() {
	nerTypeMapper["LOC"] = LOCATION
	nerTypeMapper["PER"] = PERSON
	nerTypeMapper["ORG"] = ORGANIZATION
}

func Tokenize(input string) []string {
	request, err := http.NewRequest("POST", nlpTokenizeURL, bytes.NewBuffer([]byte(input)))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println(string(body))
	var res []string
	json.Unmarshal([]byte(body), &res)
	log.Printf("%+v\n", res[1])
	return res

}

func NLPExtract(input string) NLPResp {
	request, err := http.NewRequest("POST", nlpNERUrl, bytes.NewBuffer([]byte(input)))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Fatal(err)
		return empty
	}

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return empty
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return empty
	}

	var nlpResp = NLPResp{NounPhrases: make([]string, 0), NamedEntities: make([]NamedEntitiesT, 0)}
	var nerBuffer bytes.Buffer
	var lastNeType = "0"
	var ok bool
	var res [][]string

	json.Unmarshal([]byte(body), &res)
	log.Printf("%v\n", nlpResp)
	for i := range res {
		entry := res[i]
		if entry[1] == NounPhrase {
			nlpResp.NounPhrases = append(nlpResp.NounPhrases, entry[0])
		}
		firstNERTagChar := entry[3][0]
		if nerBuffer.Len() == 0 {
			if firstNERTagChar == 'B' {
				nerBuffer.WriteString(entry[0])
				lastNeType, ok = nerTypeMapper[entry[3][2:]]
				log.Println(entry[3][2:])
				if !ok {
					lastNeType = UNKNOWN
				}
			}
		} else {
			if firstNERTagChar == 'B' {
				nlpResp.NamedEntities = append(nlpResp.NamedEntities, NamedEntitiesT{Text: nerBuffer.String(), Type: lastNeType})
				nerBuffer.Reset()
				nerBuffer.WriteString(entry[0])
				lastNeType, ok = nerTypeMapper[entry[3][2:]]
				if !ok {
					lastNeType = UNKNOWN
				}
			} else if firstNERTagChar == 'I' {
				nerBuffer.WriteString(" " + entry[0])
			}
		}

	}
	if nerBuffer.Len() != 0 {
		nlpResp.NamedEntities = append(nlpResp.NamedEntities, NamedEntitiesT{Text: nerBuffer.String(), Type: lastNeType})
	}
	return nlpResp
}
