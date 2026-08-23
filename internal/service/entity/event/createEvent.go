package event

import (
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"mime/multipart"
	"time"
)

type (
	CreateEventRequest struct {
		Name        string
		Description string
		ImageFile   *multipart.FileHeader
		Price       float64
		Quota       int64
		StartDate   string
		EndDate     string
		// FormFields opsional: definisi formulir pendaftaran kustom (§5.3).
		FormFields []FormFieldInput
	}

	// FormFieldInput adalah satu field formulir yang disusun author.
	FormFieldInput struct {
		Label     string
		FieldType string
		Required  bool
		Options   []string
		Position  int
	}
)

// ValidateFormFields memastikan definisi formulir masuk akal sebelum disimpan:
// label wajib, tipe valid, dan field pilihan wajib punya minimal satu opsi.
func (r CreateEventRequest) ValidateFormFields() error {
	for i, f := range r.FormFields {
		if f.Label == "" {
			return fmt.Errorf("form field #%d: label wajib diisi", i+1)
		}
		switch f.FieldType {
		case entity.EventFieldTypeText:
			// text tidak butuh options.
		case entity.EventFieldTypeSelect, entity.EventFieldTypeCheckbox:
			if len(f.Options) == 0 {
				return fmt.Errorf("form field %q: tipe %s wajib punya minimal satu opsi", f.Label, f.FieldType)
			}
		default:
			return fmt.Errorf("form field %q: tipe %q tidak dikenal", f.Label, f.FieldType)
		}
	}
	return nil
}

// ToObFormFields memetakan input formulir ke entity outbound untuk sebuah event.
func (r CreateEventRequest) ToObFormFields(eventID int64) []entity.EventFormField {
	if len(r.FormFields) == 0 {
		return nil
	}
	fields := make([]entity.EventFormField, 0, len(r.FormFields))
	for _, f := range r.FormFields {
		var options entity.StringList
		if len(f.Options) > 0 {
			options = entity.StringList(f.Options)
		}
		fields = append(fields, entity.EventFormField{
			EventID:   eventID,
			Label:     f.Label,
			FieldType: f.FieldType,
			Required:  f.Required,
			Options:   options,
			Position:  f.Position,
		})
	}
	return fields
}

func (r CreateEventRequest) ToObEntity(parsedStart, parsedEnd time.Time) entity.Event {
	var imageFile string
	if r.ImageFile != nil {
		imageFile = r.ImageFile.Filename
	}

	return entity.Event{
		Name:        r.Name,
		Description: r.Description,
		ImageFile:   imageFile,
		Price:       r.Price,
		Quota:       r.Quota,
		QuotaRemaining: r.Quota,
		StartDate:      parsedStart,
		EndDate:        parsedEnd,
	}
}
