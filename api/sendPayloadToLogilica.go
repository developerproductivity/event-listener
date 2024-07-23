package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	types "github.com/event-listener/types"
)

func UploadPlanningData(repoId string, payload []types.CiBuildPayload) {
	postBody, _ := json.Marshal(payload)
	logilicaUrl := fmt.Sprintf("https://logilica.io/api/import/v1/ci_build/%v/create", repoId)
	contentType := "application/json"

	client := &http.Client{}
	req, err := http.NewRequest("POST", logilicaUrl, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-lgca-domain", "redhat")
	req.Header.Add("X-lgca-token", os.Getenv("LOGILICA_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
}
