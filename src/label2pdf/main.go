package main

import (
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"os"
)

type Paper struct {
	Name        string  `json:"name,omitempty"`
	Unit        string  `json:"unit,omitempty"`
	Width       float64 `json:"width,omitempty"`
	Height      float64 `json:"height,omitempty"`
	Top         float64 `json:"top,omitempty"`
	Bottom      float64 `json:"bottom,omitempty"`
	Left        float64 `json:"left,omitempty"`
	Right       float64 `json:"right,omitempty"`
	HSpace      float64 `json:hspace,omitempty"`
	VSpace      float64 `json:vspace,omitempty"`
	Columns     int     `json:"columns,omitempty"`
	Rows        int     `json:"rows,omitempty"`
	LabelWidth  float64 `json:"label_width,omitempty"`
	LabelHeight float64 `json:"label_height,omitempty"`
}

type Page struct {
	Model       string `json:"model,omitempty"`
	Description string `json:"description,omitempty"`
	Paper       Paper  `json:"paper,omitempty"`
}

type PdfInfo struct {
	Model       string     `json:"model,omitempty"`
	ImageType   string     `json:"image_type,omitempty"`
	ImageWidth  float64    `json:"image_width,omitempty"`
	ImageHeight float64    `json:"image_height,omitempty"`
	Images      [][]string `json:"images,omitempty"`
	Repeat      bool       `json:"repeat,omitempty"`
	Orientation string     `json:"orientation,omitempty"`
}

func readPageSettings(file string) ([]Page, error) {

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	var pages []Page
	err = json.Unmarshal(bytes, &pages)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return pages, nil
}

func readPdfInfo(file string) (*PdfInfo, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	var imageinfo PdfInfo
	err = json.Unmarshal(bytes, &imageinfo)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return &imageinfo, nil
}

func findPage(pages []Page, model string) *Page {
	for _, p := range pages {
		if p.Model == model {
			return &p
		}
	}

	return nil
}

/*
func getSpacing(page *Page) (int, int) {

	vspace := (page.Paper.Height - float64(page.Paper.Rows) * page.Paper.LabelHeight - page.Paper.Top - page.Paper.Bottom)/ float64(page.Paper.Rows-1)
	hspace := (page.Paper.Width - float64(page.Paper.Columns) * page.Paper.LabelWidth - page.Paper.Left - page.Paper.Right)/ float64(page.Paper.Columns-1)

	return hspace, vspace
}
*/

func checkSize(page *Page) bool {
	height := float64(page.Paper.Rows)*page.Paper.LabelHeight + page.Paper.Top + page.Paper.Bottom + page.Paper.VSpace*float64(page.Paper.Rows-1)
	width := float64(page.Paper.Columns)*page.Paper.LabelWidth + page.Paper.Left + page.Paper.Right + page.Paper.HSpace*float64(page.Paper.Columns-1)

	if page.Paper.Width != width || page.Paper.Height != height {
		fmt.Printf("Paper size not match (%d, %d) != (%d, %d)\n", page.Paper.Width, page.Paper.Height, width, height)
		return false
	}

	return true
}

func adjustCenter(info *PdfInfo, lwidth float64, lheight float64) (deltax float64, deltay float64) {
	deltax = (lwidth - info.ImageWidth) / 2
	deltay = (lheight - info.ImageHeight) / 2
	return deltax, deltay
}

func createPdf(info *PdfInfo, currpage *Page, outfile string) error {

	page := *currpage
	pdf := gofpdf.New(info.Orientation, page.Paper.Unit, page.Paper.Name, "")
	pdf.AddPage()

	if info.Orientation == "L" {
		page.Paper.Left = currpage.Paper.Top
		page.Paper.Right = currpage.Paper.Bottom
		page.Paper.Top = currpage.Paper.Right
		page.Paper.Bottom = currpage.Paper.Left
		page.Paper.VSpace = currpage.Paper.HSpace
		page.Paper.HSpace = currpage.Paper.VSpace
		page.Paper.Columns = currpage.Paper.Rows
		page.Paper.Rows = currpage.Paper.Columns
		page.Paper.Width = currpage.Paper.Height
		page.Paper.Height = currpage.Paper.Width
		page.Paper.LabelWidth = currpage.Paper.LabelHeight
		page.Paper.LabelHeight = currpage.Paper.LabelWidth
	}

	x := page.Paper.Left
	y := page.Paper.Top

	if info.Repeat {
		// get the first image and repeat
		img := info.Images[0][0]
		for i := 0; i < page.Paper.Rows; i++ {

			for j := 0; j < page.Paper.Columns; j++ {

				deltax, deltay := adjustCenter(info, page.Paper.LabelWidth, page.Paper.LabelHeight)
				newx := x + deltax
				newy := y + deltay
				pdf.ImageOptions(img, newx, newy, info.ImageWidth, info.ImageHeight, false,
					gofpdf.ImageOptions{ImageType: info.ImageType, ReadDpi: true}, 0, "")

				x += page.Paper.HSpace + page.Paper.LabelWidth
			}
			x = page.Paper.Left
			y += page.Paper.LabelHeight + page.Paper.VSpace
		}

	} else {
		for i, a := range info.Images {
			if i >= page.Paper.Rows {
				break
			}

			for j, img := range a {
				if j >= page.Paper.Columns {
					break
				}

				if img != "" {
					// fmt.Printf("I file=%s, x=%d, y=%d\n", img, x,y)
					deltax, deltay := adjustCenter(info, page.Paper.LabelWidth, page.Paper.LabelHeight)
					newx := x + deltax
					newy := y + deltay
					pdf.ImageOptions(img, newx, newy, info.ImageWidth, info.ImageHeight, false,
						gofpdf.ImageOptions{ImageType: info.ImageType, ReadDpi: true}, 0, "")
				}

				x += page.Paper.HSpace + page.Paper.LabelWidth

			}
			x = page.Paper.Left
			y += page.Paper.LabelHeight + page.Paper.VSpace
		}

	}

	err := pdf.OutputFileAndClose(outfile)
	if err != nil {
		panic(err)
	}

	return nil
}

func main() {

	args := os.Args

	if len(args) != 4 {
		fmt.Println("usage: label2pdf page.json label.json outfile.pdf")
		return
	}

	pagefile := args[1]
	labelfile := args[2]
	outfile := args[3]

	pages, err := readPageSettings(pagefile)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	// fmt.Println(pages)

	pdfinfo, err := readPdfInfo(labelfile)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	// fmt.Println(pdfinfo)

	p := findPage(pages, pdfinfo.Model)
	if p == nil {
		fmt.Println("error:", "Invalid Model name. Please refer to page settings JSON")
	}

	// fmt.Println(p)

	if !checkSize(p) {
		return
	}

	err = createPdf(pdfinfo, p, outfile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s pdf file is created\n", outfile)
}
