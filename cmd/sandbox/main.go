package main

import (
	"errors"
	"log"
	"os"
	"strings"

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
	if true {
		log.Println("*** SandboxGetDailyReserve ***")
		reserve, err := v1.SandboxDailyReserve().Do()
		if err != nil {
			panic(err)
		}
		log.Printf("Submission Remaining Count: %d", reserve.SubmissionRemainingCount)
		log.Printf("Result: %v", reserve)
	}

	if false {
		log.Println("*** Submit URLs To Sandbox ***")
		urls := []string{"test0001", "test0002"}
		resp, _, err := v1.SandboxSubmitURLs().AddURLs(urls).Do()
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
		submit, err := v1.SandboxSubmitFile().SetFilePath(filePath)
		if err != nil {
			panic(err)
		}
		resp, _, err := submit.Do()
		if err != nil {
			panic(err)
		}
		log.Printf("Accepted: %v", filePath)
		log.Printf("ID: %s", resp.ID)
		log.Printf("MD5: %s", resp.Digest.MD5)
		log.Printf("SHA1: %s", resp.Digest.SHA1)
		log.Printf("SHA256: %s", resp.Digest.SHA256)

		virus := strings.NewReader(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
		fileName := "eicar.com"
		submit = v1.SandboxSubmitFile()
		if err := submit.SetReader(virus, fileName); err != nil {
			panic(err)
		}
		resp, _, err = submit.Do()
		if err != nil {
			panic(err)
		}
		log.Printf("Accepted: %v", fileName)
		log.Printf("ID: %s", resp.ID)
		log.Printf("MD5: %s", resp.Digest.MD5)
		log.Printf("SHA1: %s", resp.Digest.SHA1)
		log.Printf("SHA256: %s", resp.Digest.SHA256)

	}
	if false {
		log.Println("*** Submit File To Sandbox ***")
		dir, err := os.ReadDir(".")
		if err != nil {
			panic(err)
		}
		for _, each := range dir {
			if each.IsDir() {
				continue
			}
			if !strings.HasSuffix(each.Name(), ".exe") {
				continue
			}
			filePath := each.Name()
			submit, err := v1.SandboxSubmitFile().SetFilePath(filePath)
			if err != nil {
				panic(err)
			}
			resp, _, err := submit.Do()
			if err != nil {
				panic(err)
			}
			log.Printf("Accepted: %v", filePath)
			log.Printf("ID: %s", resp.ID)
			log.Printf("MD5: %s", resp.Digest.MD5)
			log.Printf("SHA1: %s", resp.Digest.SHA1)
			log.Printf("SHA256: %s", resp.Digest.SHA256)

		}
		filePath := "main.go"
		submit, err := v1.SandboxSubmitFile().SetFilePath(filePath)
		if err != nil {
			panic(err)
		}
		resp, _, err := submit.Do()
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
			status := v1.SandboxSubmissionStatus(item.ID)
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
	if false {
		log.Println("*** Iterate List Submissions & Get Analysis Results ***")
		ls := v1.SandboxListSubmissions().OrderBy(vone.CreatedDateTime, vone.Asc)
		err := ls.IterateListSubmissions(func(item *vone.ListSubmissionsItem) error {
			//log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			results := v1.SandboxAnalysisResults(item.ID)
			result, err := results.Do()
			if err != nil {
				var perr *vone.Error
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
	if false {
		log.Println("*** Iterate List Submissions & Download Result ***")
		ls := v1.SandboxListSubmissions()
		err := ls.IterateListSubmissions(func(item *vone.ListSubmissionsItem) error {
			//log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			result, err := v1.SandboxAnalysisResults(item.ID).Do()
			if err != nil {
				var perr *vone.Error
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
			return v1.SandboxDownloadResults(item.ID).Store(item.ID + ".pdf")
		})
		if err != nil {
			log.Panic(err)
		}
	}
	if false {
		log.Println("*** Iterate List Submissions & Suspicious Objects ***")
		ls := v1.SandboxListSubmissions()
		err := ls.IterateListSubmissions(func(item *vone.ListSubmissionsItem) error {
			result, err := v1.SandboxSuspiciousObjects(item.ID).Do()
			if err != nil {
				var perr *vone.Error
				if !errors.As(err, &perr) {
					return err
				}
				if perr.ErrorData.Code != "NotFound" {
					return err
				}
				log.Printf("ID: %v, NotFound", item.ID)
				return nil
			}
			for _, each := range result.Items {
				log.Printf("ID: %v, RiskLevel: %s, SHA1: %s, IP: %s", item.ID, each.RiskLevel, each.RootSHA1, each.IP)
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
}
