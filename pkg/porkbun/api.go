package porkbun

import (
	"bacon/pkg/client"
	"bacon/pkg/helpers"
	"encoding/json"
	"fmt"
)

// https://porkbun.com/api/json/v3/documentation
const (
	PING   string = "https://porkbun.com/api/json/v3/ping"
	GET    string = "https://porkbun.com/api/json/v3/dns/retrieve"
	CREATE string = "https://porkbun.com/api/json/v3/dns/create"
	DELETE string = "https://porkbun.com/api/json/v3/dns/delete"
)

type checkable interface {
	checkStatus() bool
	messageAsError() error
}

type baseRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (res *baseRes) checkStatus() bool {
	switch res.Status {
	case "SUCCESS":
		return true
	default:
		return false
	}
}

func (res *baseRes) messageAsError() error {
	return fmt.Errorf("%s", res.Message)
}

func unmarshalAndCheckStatus(data []byte, body checkable) error {
	err := json.Unmarshal(data, &body)
	if err != nil {
		return err
	}

	if !body.checkStatus() {
		return body.messageAsError()
	}

	return nil
}

func ping(auth PorkAuth) error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	body, err := helpers.PostJsonAndRead(PING, auth)
	if err != nil {
		return err
	}

	ping := pingRes{}
	err = unmarshalAndCheckStatus(body, &ping)
	if err != nil {
		return err
	}

	return nil
}

func getRecordsBytes(auth PorkAuth, domain string) ([]byte, error) {
	body, err := helpers.PostJsonAndRead(GET+"/"+domain, auth)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getRecords(auth PorkAuth, domain string) ([]PorkbunRecord, error) {
	type recordsRes struct {
		baseRes
		Records []PorkbunRecord `json:"records"`
	}

	bytes, err := getRecordsBytes(auth, domain)
	if err != nil {
		return nil, err
	}

	allRecords := recordsRes{}
	err = unmarshalAndCheckStatus(bytes, &allRecords)
	if err != nil {
		return nil, fmt.Errorf("error parsing records response: %w", err)
	}

	return allRecords.Records, nil
}

func createFromPorkbunRecord(auth PorkAuth, domain string, porkRecord PorkbunRecord) (string, error) {
	type createBody struct {
		PorkAuth
		PorkbunRecord
	}

	// Porkbun says that id is a string, but it's actually an int
	type createRes struct {
		baseRes
		Id int `json:"id"`
	}

	toCreate := createBody{
		PorkAuth:      auth,
		PorkbunRecord: porkRecord,
	}

	// Don't include the domain in create requests
	toCreate.Name = trimDomain(toCreate.Name, domain)

	body, err := helpers.PostJsonAndRead(CREATE+"/"+domain, toCreate)
	if err != nil {
		return "", err
	}

	created := createRes{}
	err = unmarshalAndCheckStatus(body, &created)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(created.Id), nil
}

func delete(auth PorkAuth, domain string, id string) error {
	body, err := helpers.PostJsonAndRead(DELETE+"/"+domain+"/"+id, auth)
	if err != nil {
		return err
	}

	deleted := baseRes{}
	err = unmarshalAndCheckStatus(body, &deleted)
	if err != nil {
		return err
	}

	return nil
}

func deploy(auth PorkAuth, domain string, records []client.Record, shouldCreate bool, shouldDelete bool) (err error) {
	old, err := getRecords(auth, domain)
	if err != nil {
		return fmt.Errorf("error getting records: %w", err)
	}

	new := make([]PorkbunRecord, len(records))
	for i, clientRecord := range records {
		new[i] = ConvertToPorkbunRecord(clientRecord)
	}

	toDelete := helpers.DifferenceByHasher(old, new, PorkbunRecord.HashFuzzy)
	toCreate := helpers.DifferenceByHasher(new, old, PorkbunRecord.HashFuzzy)

	if shouldDelete {
		fmt.Printf("Deleting %d records...\n", len(toDelete))
		for _, target := range toDelete {
			// These always have an ID from Porkbun
			err = delete(auth, domain, target.Id)
			if err != nil {
				return err
			}
			fmt.Println("-", target)
		}
	} else {
		fmt.Printf("Would delete %d records:\n", len(toDelete))
		for _, target := range toDelete {
			fmt.Println("-", target)
		}
	}
	if shouldCreate {
		fmt.Printf("Creating %d records...\n", len(toCreate))
		for _, target := range toCreate {
			id, err := createFromPorkbunRecord(auth, domain, target)
			if err != nil {
				return err
			}
			fmt.Println("-", id)
		}
	} else {
		fmt.Printf("Would create %d records:\n", len(toCreate))
		for _, target := range toCreate {
			fmt.Println("-", target)
		}
	}

	if !shouldCreate && !shouldDelete {
		fmt.Println("Mock deployment complete")
	} else if shouldCreate && shouldDelete {
		fmt.Println("Deployment complete!")
	} else {
		fmt.Println("Partial deployment complete!")
	}
	return nil
}
