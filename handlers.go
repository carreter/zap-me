package main

import (
	"encoding/json"
	"fmt"
	"github.com/carreter/pavlok-go"
	"io"
	"net/http"
	"strconv"
)

type StimulusRequest struct {
	Code     string          `json:"authCode"`
	Stimulus pavlok.Stimulus `json:"stimulus"`
}

func parseFormRequest(r *http.Request) (StimulusRequest, error) {
	err := r.ParseForm()
	if err != nil {
		return StimulusRequest{}, err
	}

	code := r.Form.Get("authCode")
	value, err := strconv.Atoi(r.Form.Get("stimulus.value"))
	if err != nil {
		return StimulusRequest{}, fmt.Errorf("could not parse stimulus.value: %w", err)
	}

	return StimulusRequest{
		Code: code,
		Stimulus: pavlok.Stimulus{
			Type:  pavlok.StimulusType(r.Form.Get("stimulus.type")),
			Value: value,
		},
	}, nil
}

func parseJsonRequest(r *http.Request) (StimulusRequest, error) {
	body, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return StimulusRequest{}, err
	}

	var stimulusReq StimulusRequest
	err = json.Unmarshal(body, &stimulusReq)
	if err != nil {
		return StimulusRequest{}, err
	}

	return stimulusReq, nil
}

func StimulusHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	var stimulusReq StimulusRequest
	var err error
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		stimulusReq, err = parseFormRequest(r)
	} else if contentType == "application/json" {
		stimulusReq, err = parseJsonRequest(r)
	}
	if err != nil {
		log.Infof("invalid stimulus request: %v", err)
		w.WriteHeader(400)
		return
	}

	if stimulusReq.Code != authCode {
		log.Warningf("incorrect authCode (%v) on stimulus request from host %v", stimulusReq.Code, r.Host)
		w.WriteHeader(403)
		return
	}

	err = pavlokClient.SendStimulus(stimulusReq.Stimulus)
	if err != nil {
		log.Errorf("error sending stimulus: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
