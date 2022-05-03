package dictionary

// what an outrage...
type responseData struct {
	FeatureCallback struct {
		Payload struct {
			SingleResults []*struct {
				Entry struct {
					Headword string `json:"headword"`

					SenseFamilies []*struct {
						PartsOfSpeech []struct {
							Value string `json:"value"`
						} `json:"parts_of_speech"`

						Senses []struct {
							ExampleGroups []struct {
								Examples []string `json:"examples"`
							} `json:"example_groups"`

							Definition struct {
								Text string `json:"text"`
							} `json:"definition"`
						} `json:"senses"`
					} `json:"sense_families"`
				} `json:"entry"`
			} `json:"single_results"`
		} `json:"payload"`
	} `json:"feature-callback"`
}

type Definition struct {
	Examples []string
	Text     string
}

type Meaning struct {
	PartOfSpeech string
	Definitions  []Definition
}

type Entry struct {
	Word     string
	Meanings []Meaning
}
