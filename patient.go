package main

import (
	"bytes"
	"encoding/hex"
	"strings"
	"time"
)
import "crypto/md5"

// Passport описывает пасспорт
type Passport struct {
	Serial   string `json:"serial,omitempty"`
	Number   string `json:"number,omitempty"`
	IssuedAt string `json:"issuedAt,omitempty"`
	IssuedBy string `json:"issuedBy,omitempty"`
}

// Patient описывает пациента
type Patient struct {
	Md5      string `json:"md5"`
	Fname    string `json:"fname,omitempty"`
	Lname    string `json:"lname,omitempty"`
	Mname    string `json:"mname,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Card     string `json:"card,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
	*Passport
}

// ComputeMd5 returns hash
func (p *Patient) ComputeMd5() {
	b := new(bytes.Buffer)
	b.WriteString(p.Fname)
	b.WriteString(p.Lname)
	b.WriteString(p.Mname)
	b.WriteString(p.Phone)
	b.WriteString(p.Card)
	b.WriteString(p.Passport.Serial)
	b.WriteString(p.Passport.Number)
	hash := md5.Sum(b.Bytes())
	p.Md5 = hex.EncodeToString(hash[:])
}

const (
	idxLname     = iota
	idxFname     = iota
	idxMname     = iota
	idxSex       = iota
	idxBirthday  = iota
	idxCard      = iota
	_            = iota
	_            = iota
	_            = iota
	_            = iota
	idxDocType   = iota
	idxDocSerial = iota
	idxDocNumber = iota
	idxIssuedAt  = iota
	idxIssuedBy  = iota
	_            = iota
	_            = iota
	idxPhone     = iota
	_            = iota
	_            = iota
	idxAddrStart = iota
)

// NewPatient ...
func NewPatient(line []string) *Patient {
	const timeFormat = "01/_2/2006"

	pat := Patient{
		Lname: line[idxLname],
		Fname: line[idxFname],
		Mname: line[idxMname],
		Sex:   line[idxSex],
		Card:  line[idxCard],
		Phone: line[idxPhone],
		Passport: &Passport{
			Serial: line[idxDocSerial],
			Number: line[idxDocNumber],
		},
	}

	// Адрес
	address := bytes.NewBufferString("")
	for _, x := range line[idxAddrStart:] {
		if len(x) != 0 {
			address.WriteString(x)
			address.WriteString(", ")
		}
	}
	pat.Address = address.String()

	// Пол
	if strings.Contains(strings.ToLower(pat.Sex), "муж") {
		pat.Sex = "male"
	} else {
		pat.Sex = "female"
	}

	// День рождения
	if val := line[idxBirthday]; len(val) != 0 {
		norm := make([]string, 3)
		for i, val := range strings.Split(val, "/") {
			if len(val) == 1 {
				norm[i] = "0" + val
			} else {
				norm[i] = val
			}
		}

		if t, err := time.Parse(timeFormat, strings.Join(norm, "/")); err == nil {
			pat.Birthday = t.Format("2006-01-02")
		}
	}

	return &pat
}
