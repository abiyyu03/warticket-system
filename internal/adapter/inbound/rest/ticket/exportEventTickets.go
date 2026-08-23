package ticket

import (
	"context"
	"fmt"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// ExportEventTickets menghasilkan file .xlsx berformat rapi berisi daftar tiket
// sebuah event untuk diunduh author.
func (h *Handler) ExportEventTickets(fctx *fiber.Ctx) error {
	ctx := context.Background()

	eventID, err := strconv.ParseInt(fctx.Params("id"), 10, 64)
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse("event id tidak valid", fiber.StatusBadRequest, nil, nil),
		)
	}

	data, err := h.Service.Ticket.ListEventTickets(ctx, eventID)
	if err != nil {
		return err
	}

	buf, err := buildTicketsXLSX(data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("tickets-event-%d.xlsx", eventID)
	fctx.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	fctx.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filename))
	return fctx.Status(fiber.StatusOK).Send(buf)
}

// buildTicketsXLSX menyusun workbook: judul, header berwarna, border, freeze
// header, auto-filter, dan lebar kolom yang nyaman dibaca.
func buildTicketsXLSX(data ucEntity.ListEventTicketsResponse) ([]byte, error) {
	const sheet = "Tiket"
	f := excelize.NewFile()
	defer f.Close()
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No", "Kode Tiket", "Email Peserta", "User ID", "Status", "Berlaku Sampai", "Diterbitkan"}
	widths := map[string]float64{"A": 6, "B": 40, "C": 28, "D": 10, "E": 14, "F": 20, "G": 20}
	for col, w := range widths {
		_ = f.SetColWidth(sheet, col, col, w)
	}

	// --- judul & subjudul (baris 1-2, di-merge) ---
	_ = f.MergeCell(sheet, "A1", "G1")
	_ = f.MergeCell(sheet, "A2", "G2")
	_ = f.SetCellValue(sheet, "A1", fmt.Sprintf("Daftar Tiket — Event #%d", data.EventID))
	_ = f.SetCellValue(sheet, "A2", fmt.Sprintf("Total: %d tiket · Diekspor %s",
		data.Total, time.Now().Format("2006-01-02 15:04")))

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 15, Color: "1F4E78"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	subStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "808080"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	_ = f.SetCellStyle(sheet, "A1", "A1", titleStyle)
	_ = f.SetCellStyle(sheet, "A2", "A2", subStyle)
	_ = f.SetRowHeight(sheet, 1, 24)

	// --- header (baris 3) ---
	const headerRow = 3
	border := []excelize.Border{
		{Type: "top", Color: "BFBFBF", Style: 1},
		{Type: "bottom", Color: "BFBFBF", Style: 1},
		{Type: "left", Color: "BFBFBF", Style: 1},
		{Type: "right", Color: "BFBFBF", Style: 1},
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1F4E78"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    border,
	})
	for i, hName := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		_ = f.SetCellValue(sheet, cell, hName)
	}
	_ = f.SetCellStyle(sheet, "A3", "G3", headerStyle)
	_ = f.SetRowHeight(sheet, headerRow, 20)

	// --- baris data ---
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    border,
	})
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    border,
	})
	// warna status: ACTIVE hijau, REDEEMED biru, lainnya abu.
	statusStyle := func(status string) int {
		color := "808080"
		switch status {
		case "ACTIVE":
			color = "2E7D32"
		case "REDEEMED":
			color = "1565C0"
		}
		id, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: color},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Border:    border,
		})
		return id
	}

	row := headerRow + 1
	for i, tkt := range data.Tickets {
		email := tkt.Email
		if email == "" {
			email = "-"
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), tkt.Code)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), email)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), tkt.UserID)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), tkt.Status)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", row), tkt.ValidUntil.Format("2006-01-02 15:04"))
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", row), tkt.CreatedAt.Format("2006-01-02 15:04"))

		_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), centerStyle)
		_ = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), cellStyle)
		_ = f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), statusStyle(tkt.Status))
		_ = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("G%d", row), centerStyle)
		row++
	}

	// freeze header + auto-filter supaya nyaman di-scroll & difilter.
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      headerRow,
		TopLeftCell: "A4",
		ActivePane:  "bottomLeft",
	})
	_ = f.AutoFilter(sheet, "A3:G3", nil)

	out, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
