package main

import (
	"errors"
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
		reserve, err := v1.SandboxDailyReserve().Do()
		if err != nil {
			panic(err)
		}
		log.Printf("Submission Remaining Count: %d", reserve.SubmissionRemainingCount)
	}

	if false {
		log.Println("*** Submit URLs To Sandbox ***")
		urls := []string{"test0001", "test0002"}
		resp, err := v1.NewSubmitURLsToSandbox().AddURLs(urls).Do()
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
		submit, err := v1.SandboxSubmitFile().SetFileName(filePath)
		if err != nil {
			panic(err)
		}
		resp, err := submit.Do()
		if err != nil {
			panic(err)
		}
		log.Printf("Accepted: %v", filePath)
		log.Printf("ID: %s", resp.ID)
		log.Printf("MD5: %s", resp.Digest.MD5)
		log.Printf("SHA1: %s", resp.Digest.SHA1)
		log.Printf("SHA256: %s", resp.Digest.SHA256)
	}

	if false {
		log.Println("*** List Submissions ***")
		listSubmissions := v1.SandboxListSubmissions().OrderBy(vone.CreatedDateTime, vone.Asc)
		resp, err := listSubmissions.Do()
		//		log.Println("CCC", resp, err)
		if err != nil {
			panic(err)
		}
		for _, item := range resp.Items {
			log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
		}
		log.Printf("Next: %s", resp.NextLink)
		if resp.NextLink != "" {
			log.Println("*** List Submissions Next ***")
			resp2, err := listSubmissions.Next()
			if err != nil {
				panic(err)
			}
			for _, item := range resp2.Items {
				log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			}
		}
		log.Println("*** Iterate List Submissions ***")
		ls := v1.SandboxListSubmissions().OrderBy(vone.CreatedDateTime, vone.Asc)
		err = ls.IterateListSubmissions(func(item *vone.ListSubmissionsItem) error {
			log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			status := v1.SubmissionStatus(item.ID)
			result, err := status.Do()
			if err != nil {
				return err
			}
			log.Printf("ID: %v, ResourceLocation: %s", item.ID, result.ResourceLocation)
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
	if true {
		log.Println("*** Iterate List Submissions & Get Analysis Results ***")
		ls := v1.SandboxListSubmissions().OrderBy(vone.CreatedDateTime, vone.Asc)
		err := ls.IterateListSubmissions(func(item *vone.ListSubmissionsItem) error {
			//log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			results := v1.SandboxAnalysisResults(item.ID)
			result, err := results.Do()
			if err != nil {
				var perr *vone.VOneError
				if !errors.As(err, &perr) {
					return err
				}
				if perr.ErrorData.Code != "NotFound" {
					return err
				}
				log.Printf("ID: %v, NotFound", item.ID)
				return nil
			}
			log.Printf("ID: %v, RiskLevel: %s", item.ID, result.RiskLevel)
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
}
