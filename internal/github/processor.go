package github

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/varavell/mcard/pkg/mcardhttp"

	"github.com/pkg/errors"
)

// AggregatorConfig is used for encapsulate config needed for gathering user and repo results
type Config struct {
	Gclient mcardhttp.Client
}

func (a *Config) RepoDataByUser(user string) (map[string]interface{}, error) {
	userData := make(map[string]interface{})
	repoData := make(map[string][]string)

	repos, err := a.Gclient.List(user)
	if err != nil {
		return userData, errors.New("unable to get the list of repos for user")
	}

	for _, repo := range repos {
		var langlist []string
		langObj, err := a.Gclient.ListLanguages(*repo.LanguagesURL)
		if err != nil {
			return userData, errors.New("unable to get the list of languages for repos")
		}
		for lang, _ := range langObj {
			langlist = append(langlist, lang)
		}
		repoData[*repo.Name] = langlist
	}

	userData[user] = repoData

	return userData, nil
}

func (a *Config) Produce(jobs <-chan string, producer chan map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		userData, err := a.RepoDataByUser(j)
		if err != nil {
			panic(err)
		}
		producer <- userData
	}
}

func (a *Config) Consume(producer chan map[string]interface{}, finish chan bool, outputPath string) {
	// Create the threadsafe map.
	finalMap := make(map[string]interface{})
	var mutex = &sync.RWMutex{}
	for repoData := range producer {
		for key, value := range repoData {
			mutex.Lock()
			finalMap[key] = value
			mutex.Unlock()
		}
	}

	userData, err := json.Marshal(finalMap)
	if err != nil {
		finish <- false
		return
	}

	_ = ioutil.WriteFile(outputPath, userData, 0644)

	finish <- true
}
