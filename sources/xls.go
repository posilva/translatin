package sources

import (
	"fmt"
	"strings"

	"github.com/posilva/translatin/translator"
	"github.com/xuri/excelize/v2"
)

const defaultSourceLang = "en-GB"

type CellAddress struct {
	Row  int
	Col  int
	Name string
}

type Entry struct {
	Key             string
	BaseText        string
	BaseLanguage    string
	TargetLanguages []string
}

type Opener interface {
	Open() error
}

type Closer interface {
	Close() error
}

type Source interface {
	Opener
	Closer
}

type XLSSourceSettings struct {
	FilePath  string
	SheetName string
}

type XLSSource struct {
	settings        XLSSourceSettings
	file            *excelize.File
	sourceLanguage  string
	langs2Address   map[string]CellAddress
	targetLanguages []string
}

// NewXLSSourceWithDefaultLanguage creates a new XLSSource with the default language.
func NewXLSSourceWithDefaultLanguage(settings XLSSourceSettings) *XLSSource {
	return NewXLSSource(defaultSourceLang, settings)
}

// NewXLSSource creates a new XLSSource.
func NewXLSSource(srcLang string, settings XLSSourceSettings) *XLSSource {
	return &XLSSource{
		settings:        settings,
		sourceLanguage:  srcLang,
		langs2Address:   make(map[string]CellAddress),
		targetLanguages: []string{},
	}
}

// Open opens the source file.
func (s *XLSSource) Close() error {
	if !s.IsOpen() {
		return nil
	}
	return s.file.Close()
}

// Open opens the source file.
func (s *XLSSource) Open() error {
	if s.IsOpen() {
		return nil
	}
	var err error
	s.file, err = excelize.OpenFile(s.settings.FilePath)
	return err
}

// IsOpen returns true if the source file is open.
func (s *XLSSource) IsOpen() bool {
	return s.file != nil
}

func (s *XLSSource) BuildQueries() ([]translator.Query, error) {
	if !s.IsOpen() {
		return nil, fmt.Errorf("source not open")
	}

	rows, err := s.file.Rows(s.settings.SheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get columns: %w", err)
		}
		for colIdx, col := range cols {
			if strings.Compare(col, "comment") == 0 {
				continue
			}
			address, err := excelize.CoordinatesToCellName(colIdx+1, 1)
			if err != nil {
				return nil, fmt.Errorf("failed to convert coordinates to cell name: %w", err)
			}
			_, r, err := excelize.SplitCellName(address)
			if err != nil {
				return nil, fmt.Errorf("failed to split cell name: %w", err)
			}
			s.langs2Address[col] = CellAddress{Row: r, Col: colIdx + 1, Name: address}
			if col != s.sourceLanguage && col != "Keyname" && col != "en-US" {
				s.targetLanguages = append(s.targetLanguages, col)
			}
		}

	}

	baseLanguageColIdx := s.langs2Address[s.sourceLanguage].Col - 1
	keyColIdx := s.langs2Address["Keyname"].Col - 1
	var queries []translator.Query
	rc := 1
	for rows.Next() {
		rc++

		cols, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get columns: %w", err)
		}

		for _, t := range s.targetLanguages {
			id, err := excelize.CoordinatesToCellName(s.langs2Address[t].Col, rc)
			if err != nil {
				return nil, fmt.Errorf("failed to convert coordinates to cell name [%s] on row [%d]: %w", t, rc, err)
			}
			qry := translator.Query{
				ID:             id,
				Key:            cols[keyColIdx],
				SourceLanguage: s.sourceLanguage,
				TargetLanguage: t,
				Text:           cols[baseLanguageColIdx],
			}
			queries = append(queries, qry)
		}
	}
	return queries, nil
}
