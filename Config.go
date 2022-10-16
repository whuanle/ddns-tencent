package main

type Config struct {
	SecretId   string `json:"SecretId"`
	SecretKey  string `json:"SecretKey"`
	Domain     string `json:"Domain"`
	SubDomain  string `json:"SubDomain"`
	RecordType string `json:"RecordType"`
	RecordLine string `json:"RecordLine"`
	Value      string `json:"Value"`
	MX         int    `json:"MX"`
	TTL        int    `json:"TTL"`
	RecordId   int    `json:"RecordId"`
}
