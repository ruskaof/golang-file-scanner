package preprocessing

import (
	"biocadTask/internal/storage"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"strconv"
)

// GenerateProcessedFileReports generates a pdf report for a message in a file and saves it
// with guid as a filename
func GenerateProcessedFileReports(
	entities []storage.DeviceEntity,
	guid uuid.UUID, outputDirectoryPath string,
) error {
	pdf := setupStandardPdf()

	var stringEntities []string
	for _, entity := range entities {
		stringEntities = append(stringEntities, entity.String())
	}

	pdf.MultiCell(0, 10, "Успешно обработано", "", "C", false)
	for _, pdfMessage := range stringEntities {
		pdf.MultiCell(0, 10, pdfMessage, "", "", false)
		pdf.MultiCell(0, 10, "", "", "", false)
	}

	err := pdf.OutputFileAndClose(outputDirectoryPath + "/" + guid.String() + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

// GenerateProcessedFileReports generates a pdf report for an error and saves it in format "error<id>"
func GenerateErrorReport(errMessage error, id int64, outputDirectoryPath string) error {
	pdf := setupStandardPdf()

	pdf.MultiCell(0, 10, "Ошибка обработки", "", "C", false)
	pdf.MultiCell(0, 10, errMessage.Error(), "", "", false)

	err := pdf.OutputFileAndClose(outputDirectoryPath + "/" + "error" + strconv.FormatInt(id, 10) + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func setupStandardPdf() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("roboto", "", "fonts/Roboto-Regular.ttf")
	pdf.SetFont("roboto", "", 16)
	return pdf
}
