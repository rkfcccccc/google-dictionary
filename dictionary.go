package dictionary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var client = http.Client{}
var ErrNoDefinitionsFound = errors.New("no definition found")

func fetchData(ctx context.Context, language, word string) (*responseData, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com/async/callback:5493", nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	query := request.URL.Query()
	query.Add("fc", "ErUBCndBTlVfTnFUN29LdXdNSlQ2VlZoWUIwWE1HaElOclFNU29TOFF4ZGxGbV9zbzA3YmQ2NnJyQXlHNVlrb3l3OXgtREpRbXpNZ0M1NWZPeFo4NjQyVlA3S2ZQOHpYa292MFBMaDQweGRNQjR4eTlld1E4bDlCbXFJMBIWU2JzSllkLVpHc3J5OVFPb3Q2aVlDZxoiQU9NWVJ3QmU2cHRlbjZEZmw5U0lXT1lOR3hsM2xBWGFldw")
	query.Add("fcv", "3")
	query.Add("async", fmt.Sprintf("term:%q,corpus:%s,hhdr:true,hwdgt:true,wfp:true,ttl:,tsl:,ptl:", word, language))
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}

	defer response.Body.Close()
	response.Body.Read(make([]byte, 4)) // skipping )]}' at the beinning

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll for word %s: %v", string(word), err)
	}

	if response.StatusCode != 200 {
		return nil, ErrNoDefinitionsFound
	}

	var data responseData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &data, nil
}

func transform(data *responseData) ([]Entry, error) {
	singleResults := data.FeatureCallback.Payload.SingleResults
	entries := make([]Entry, len(singleResults))

	for i, singleResult := range singleResults {
		entry := singleResult.Entry

		entries[i].Word = entry.Headword
		entries[i].Meanings = make([]Meaning, len(entry.SenseFamilies))

		for j, senseFamily := range entry.SenseFamilies {
			if len(senseFamily.PartsOfSpeech) > 0 {
				entries[i].Meanings[j].PartOfSpeech = senseFamily.PartsOfSpeech[0].Value
			}

			entries[i].Meanings[j].Definitions = make([]Definition, len(senseFamily.Senses))

			for k, sense := range senseFamily.Senses {
				entries[i].Meanings[j].Definitions[k].Text = sense.Definition.Text

				if len(sense.ExampleGroups) == 0 {
					continue
				}

				entries[i].Meanings[j].Definitions[k].Examples = sense.ExampleGroups[0].Examples
			}
		}
	}

	return entries, nil
}

func GetWordEntry(ctx context.Context, language, word string) (*Entry, error) {
	data, err := fetchData(ctx, language, word)
	if err != nil {
		return nil, fmt.Errorf("fetchData: %w", err)
	}

	transformed, err := transform(data)
	if err != nil {
		return nil, fmt.Errorf("transform: %w", err)
	}

	if len(transformed) == 0 {
		return nil, ErrNoDefinitionsFound
	}

	return &transformed[0], nil
}
