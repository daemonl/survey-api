package surveys

import "fmt"

var AnimalList = []string{
	"dog",
	"cat",
	"bird",
	"horse",
	"snake",
}

var animalSet = map[string]struct{}{}

func init() {
	for _, animal := range AnimalList {
		animalSet[animal] = struct{}{}
	}
}

type AnimalResponse struct {
	Rating int `json:"rating" bson:"rating"`
	Owned  int `json:"owned" bson:"owned"`
}

type Response struct {
	Animals map[string]AnimalResponse `json:"animals" bson:"animals"`
	Age     int                       `json:"age" bson:"age"`
}

type StoredResponse struct {
	ID       string `json:"id" bson:"_id"`
	Response `bson:",inline"`
}

func rangeError(min, max int) string {
	return fmt.Sprintf("Must be between %d and %d", min, max)
}

func (r Response) Validate() map[string]string {
	issues := map[string]string{}
	for key, rating := range r.Animals {
		if _, ok := animalSet[key]; !ok {
			issues[fmt.Sprintf("animals.%s", key)] = "No Such Animal"
			continue
		}

		if rating.Rating > 10 || rating.Rating < 0 {
			issues[fmt.Sprintf("animals.%s.rating", key)] = rangeError(0, 10)
			continue
		}
		if rating.Owned > 1000 || rating.Owned < 0 {
			issues[fmt.Sprintf("animals.%s.owned", key)] = rangeError(0, 1000)
			continue
		}
	}

	if r.Age < 0 || r.Age > 150 {
		issues["age"] = rangeError(0, 150)
	}

	if len(issues) == 0 {
		return nil
	}
	return issues
}
