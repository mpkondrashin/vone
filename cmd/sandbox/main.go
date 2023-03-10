/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	main.go - various SDK tests
*/

package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/mpkondrashin/vone"
)

func main() {
	v1 := vone.NewVOne(url, token)
	if true {
		log.Println("*** SandboxGetDailyReserve ***")
		reserve, err := v1.SandboxDailyReserve().Do(context.TODO())
		if err != nil {
			panic(err)
		}
		log.Printf("Submission Remaining Count: %d", reserve.SubmissionRemainingCount)
		log.Printf("Result: %v", reserve)
	}

	if false {
		log.Println("*** Submit URLs To Sandbox ***")
		urls := []string{"test0001", "test0002"}
		resp, _, err := v1.SandboxSubmitURLs().AddURLs(urls).Do(context.TODO())
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
		resp, _, err := submit.Do(context.TODO())
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
		resp, _, err = submit.Do(context.TODO())
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
			resp, _, err := submit.Do(context.TODO())
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
		resp, _, err := submit.Do(context.TODO())
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
		resp, err := listSubmissions.Do(context.TODO())
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
			resp2, err := listSubmissions.Next(context.TODO())
			if err != nil {
				panic(err)
			}
			for _, item := range resp2.Items {
				log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			}
		}
		log.Println("*** Iterate List Submissions ***")
		ls := v1.SandboxListSubmissions().OrderBy(vone.CreatedDateTime, vone.Asc)
		err = ls.IterateListSubmissions(context.TODO(), func(item *vone.ListSubmissionsItem) error {
			log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			status := v1.SandboxSubmissionStatus(item.ID)
			result, err := status.Do(context.TODO())
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
		err := ls.IterateListSubmissions(context.TODO(), func(item *vone.ListSubmissionsItem) error {
			//log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			results := v1.SandboxAnalysisResults(item.ID)
			result, err := results.Do(context.TODO())
			if err != nil {
				var perr *vone.Error
				if !errors.As(err, &perr) {
					return err
				}
				if perr.ErrorData.Code != vone.ErrorCodeNotFound {
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
		err := ls.IterateListSubmissions(context.TODO(), func(item *vone.ListSubmissionsItem) error {
			//log.Printf("ID: %v, action: %s, status: %s", item.ID, item.Action, item.Status)
			result, err := v1.SandboxAnalysisResults(item.ID).Do(context.TODO())
			if err != nil {
				var perr *vone.Error
				if !errors.As(err, &perr) {
					return err
				}
				if perr.ErrorData.Code != vone.ErrorCodeNotFound {
					return err
				}
				log.Printf("ID: %v, NotFound", item.ID)
				return nil
			}
			log.Printf("ID: %v, RiskLevel: %s", item.ID, result.RiskLevel)
			return v1.SandboxDownloadResults(item.ID).Store(context.TODO(), item.ID+".pdf")
		})
		if err != nil {
			log.Panic(err)
		}
	}
	if false {
		log.Println("*** Iterate List Submissions & Suspicious Objects ***")
		ls := v1.SandboxListSubmissions()
		err := ls.IterateListSubmissions(context.TODO(), func(item *vone.ListSubmissionsItem) error {
			result, err := v1.SandboxSuspiciousObjects(item.ID).Do(context.TODO())
			if err != nil {
				var perr *vone.Error
				if !errors.As(err, &perr) {
					return err
				}
				if perr.ErrorData.Code != vone.ErrorCodeNotFound {
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
