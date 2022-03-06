package pork

func Ping(auth Auth) (string, error) {
	body, err := postAndDecode(auth, PK_PING)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func RetrieveRecords(auth Auth, domain string) (string, error) {
	body, err := postAndDecode(auth, PK_DNS_RETRIEVE+domain)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
