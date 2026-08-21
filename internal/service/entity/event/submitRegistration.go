package event

import (
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"strconv"
	"strings"
)

type (
	SubmitRegistrationRequest struct {
		UserID  int64
		EventID int64
		Answers []AnswerInput
	}

	AnswerInput struct {
		FieldID int64
		Value   []string
	}
)

// nonEmpty menyaring nilai kosong/whitespace.
func nonEmpty(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}

func inOptions(value string, options []string) bool {
	for _, o := range options {
		if o == value {
			return true
		}
	}
	return false
}

// Validate mencocokkan jawaban terhadap definisi field:
// - field wajib harus terjawab,
// - jawaban select/checkbox harus termasuk opsi yang tersedia,
// - jawaban untuk field yang tidak dikenal ditolak.
func (r SubmitRegistrationRequest) Validate(fields []entity.EventFormField) error {
	byID := make(map[int64]entity.EventFormField, len(fields))
	for _, f := range fields {
		byID[f.ID] = f
	}

	answered := make(map[int64][]string, len(r.Answers))
	for _, a := range r.Answers {
		if _, ok := byID[a.FieldID]; !ok {
			return fmt.Errorf("jawaban untuk field_id %d tidak dikenal pada event ini", a.FieldID)
		}
		answered[a.FieldID] = nonEmpty(a.Value)
	}

	for _, f := range fields {
		vals := answered[f.ID]

		if f.Required && len(vals) == 0 {
			return fmt.Errorf("field %q wajib diisi", f.Label)
		}
		if len(vals) == 0 {
			continue // opsional & tidak diisi
		}

		switch f.FieldType {
		case entity.EventFieldTypeText:
			if len(vals) > 1 {
				return fmt.Errorf("field %q hanya menerima satu jawaban", f.Label)
			}
		case entity.EventFieldTypeSelect:
			if len(vals) != 1 {
				return fmt.Errorf("field %q hanya boleh memilih satu opsi", f.Label)
			}
			if !inOptions(vals[0], f.Options) {
				return fmt.Errorf("field %q: opsi %q tidak valid", f.Label, vals[0])
			}
		case entity.EventFieldTypeCheckbox:
			for _, v := range vals {
				if !inOptions(v, f.Options) {
					return fmt.Errorf("field %q: opsi %q tidak valid", f.Label, v)
				}
			}
		}
	}

	return nil
}

// ToObEntity menyusun baris registrasi; answers disimpan keyed by field_id.
func (r SubmitRegistrationRequest) ToObEntity() entity.UserRegistration {
	answers := make(entity.AnswerMap, len(r.Answers))
	for _, a := range r.Answers {
		vals := nonEmpty(a.Value)
		if len(vals) == 0 {
			continue
		}
		answers[strconv.FormatInt(a.FieldID, 10)] = vals
	}

	return entity.UserRegistration{
		UserID:  r.UserID,
		EventID: r.EventID,
		Answers: answers,
	}
}
