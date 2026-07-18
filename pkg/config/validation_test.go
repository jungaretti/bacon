package config

import (
	"testing"
)

func TestValidateConfigValid(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "bacondemo.com",
				Type: "ALIAS",
				Ttl:  600,
				Data: "pixie.porkbun.com",
			},
			{
				Name: "*.bacondemo.com",
				Type: "CNAME",
				Ttl:  600,
				Data: "pixie.porkbun.com",
			},
			{
				Name:     "bacondemo.com",
				Type:     "MX",
				Ttl:      600,
				Data:     "in1-smtp.messagingengine.com",
				Priority: 10,
			},
		},
	}

	if err := ValidateConfiguration(config); err != nil {
		t.Fatal("expected valid config, got", err)
	}
}

func TestValidateConfigMissingHost(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Type: "ALIAS",
				Ttl:  600,
				Data: "pixie.porkbun.com",
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record is missing host field")
	}
}

func TestValidateConfigAcmeChallengeHost(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "_acme-challenge.bacondemo.com",
				Type: "TXT",
				Ttl:  600,
				Data: "c_V4WaKPWlisAvnvTZ62BOuLiQDpkC2cOtahW8TDthw",
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record host begins with _acme-challenge")
	}
}

func TestValidateConfigMissingType(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "bacondemo.com",
				Ttl:  600,
				Data: "pixie.porkbun.com",
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record is missing type field")
	}
}

func TestValidateConfigInvalidType(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "bacondemo.com",
				Type: "FAKE",
				Ttl:  600,
				Data: "pixie.porkbun.com",
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record has an invalid type")
	}
}

func TestValidateConfigMissingTtl(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "bacondemo.com",
				Type: "ALIAS",
				Data: "pixie.porkbun.com",
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record is missing ttl field")
	}
}

func TestValidateConfigInvalidTtl(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "bacondemo.com",
				Type: "ALIAS",
				Ttl:  300,
				Data: "pixie.porkbun.com",
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record has an invalid ttl")
	}
}

func TestValidateConfigMissingContent(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "bacondemo.com",
				Type: "ALIAS",
				Ttl:  600,
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a record is missing content field")
	}
}

func TestValidateConfigInvalidPriority(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name:     "bacondemo.com",
				Type:     "ALIAS",
				Ttl:      600,
				Data:     "pixie.porkbun.com",
				Priority: 20,
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when priority is set on an ALIAS record")
	}
}

func TestValidateConfigCnameSameHost(t *testing.T) {
	config := Config{
		Domain: "bacondemo.com",
		Records: []Record{
			{
				Name: "*.bacondemo.com",
				Type: "CNAME",
				Ttl:  600,
				Data: "pixie.porkbun.com",
			},
			{
				Name:     "*.bacondemo.com",
				Type:     "MX",
				Ttl:      600,
				Data:     "in1-smtp.messagingengine.com",
				Priority: 10,
			},
		},
	}

	if err := ValidateConfiguration(config); err == nil {
		t.Fatal("expected error when a CNAME and another record share a host")
	}
}
