package main

import (
	"log"

	"github.com/mpkondrashin/vone"
)

/*
SubmitURLsToSandboxDataResponse [][]struct {
	Status  int `json:"status"`
	Headers []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"headers,omitempty"`
	Body struct {
		URL    string `json:"url"`
		ID     string `json:"id"`
		Digest struct {
			Md5    string `json:"md5"`
			Sha1   string `json:"sha1"`
			Sha256 string `json:"sha256"`
		} `json:"digest"`
	} `json:"body,omitempty"`
	BodyError struct {
		URL   string `json:"url"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	} `json:"body,omitempty"`
}*/

func main() {
	v1 := vone.NewVOne(url, token)
	if false {
		log.Println("*** SandboxGetDailyReserve ***")
		reserve, err := v1.SandboxGetDailyReserve()
		if err != nil {
			panic(err)
		}
		log.Printf("Submission Remaining Count: %d", reserve.SubmissionRemainingCount)
	}
	if false {
		log.Println("*** Submit URLs To Sandbox ***")
		urls := []string{"test0001", "test0002"}
		resp, err := v1.SubmitURLsToSandbox(urls)
		if err != nil {
			panic(err)
		}
		//log.Println("RESP", resp)
		for _, each := range resp {
			if each.Status == 207 {
				log.Printf("%v: %v", each.Body.URL, each.Body.ID)
			} else {
				log.Printf("%v: %v", each.Body.URL, each.Body.Error.Message)
			}
		}
	}
	if false {
		log.Println("*** Submit File To Sandbox ***")
		filePath := "main.go"
		resp, err := v1.SubmitFileToSandbox(filePath)
		if err != nil {
			panic(err)
		}
		log.Printf("Accepted: %v", filePath)
		log.Printf("ID: %s", resp.ID)
		log.Printf("MD5: %s", resp.Digest.MD5)
		log.Printf("SHA1: %s", resp.Digest.SHA1)
		log.Printf("SHA256: %s", resp.Digest.SHA256)
	}
	if true {
		log.Println("*** List Submissions ***")
		resp, err := v1.ListSubmissions()
		if err != nil {
			panic(err)
		}
		for _, item := range resp.Items {
			log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
		}
		log.Println("*** List Submissions Next ***")
		resp2, err := v1.ListSubmissionsNext(resp)
		if err != nil {
			panic(err)
		}
		for _, item := range resp2.Items {
			log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
		}

	}
}
