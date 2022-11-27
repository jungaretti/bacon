package dns

type Config struct {
	Domain  string         `yaml:"domain"`
	Records []ConfigRecord `yaml:"records"`
}

type ConfigRecord struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Ttl  string `yaml:"ttl"`
	Data string `yaml:"data"`
}

func (r ConfigRecord) GetName() string {
	return r.Name
}

func (r ConfigRecord) GetType() string {
	return r.Type
}

func (r ConfigRecord) GetTtl() string {
	return r.Ttl
}

func (r ConfigRecord) GetData() string {
	return r.Data
}

var _ Record = ConfigRecord{}

func ConfigFromRecord(r Record) ConfigRecord {
	return ConfigRecord{
		Name: r.GetName(),
		Type: r.GetType(),
		Ttl:  r.GetTtl(),
		Data: r.GetData(),
	}
}
