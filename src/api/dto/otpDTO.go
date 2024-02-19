package dto

type RequestOtpDTO struct {
	PhoneNumber string `json:"phoneNumber"`
}

type OtpDTO struct {
	Value string `json:"value"`
	Used  bool   `json:"used"`
}
