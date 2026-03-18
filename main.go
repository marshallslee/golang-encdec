package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("암복호화 도구")
	w.Resize(fyne.NewSize(500, 400))

	modeEncrypt := true

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("키를 입력하세요")

	inputLabel := widget.NewLabel("평문")
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("암호화할 문장을 입력하세요")
	inputEntry.SetMinRowsVisible(3)

	// 결과: 복사 가능한 Entry (비활성화하지 않아 검정색 유지)
	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetMinRowsVisible(2)

	// 에러: RichText로 빨간 글씨 (한글 깨짐 없음)
	errorLabel := widget.NewRichTextWithText("")
	errorLabel.Hide()

	resultContainer := container.NewStack(resultEntry, errorLabel)

	showResult := func(text string) {
		errorLabel.Hide()
		resultEntry.SetText(text)
		resultEntry.Show()
	}

	showError := func(text string) {
		resultEntry.SetText("")
		resultEntry.Hide()
		errorLabel.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					ColorName: "",
					TextStyle: fyne.TextStyle{Bold: true},
				},
			},
		}
		// 빨간색 직접 지정
		errorLabel.Segments[0].(*widget.TextSegment).Style.ColorName = "error"
		errorLabel.Show()
		errorLabel.Refresh()
	}

	clearResult := func() {
		resultEntry.SetText("")
		resultEntry.Show()
		errorLabel.Hide()
	}

	_ = showError // suppress unused warning during build

	radio := widget.NewRadioGroup([]string{"암호화", "복호화"}, func(selected string) {
		modeEncrypt = selected == "암호화"
		clearResult()
		if modeEncrypt {
			inputLabel.SetText("평문")
			inputEntry.SetPlaceHolder("암호화할 문장을 입력하세요")
		} else {
			inputLabel.SetText("암호문")
			inputEntry.SetPlaceHolder("복호화할 문장을 입력하세요")
		}
	})
	radio.Horizontal = true
	radio.SetSelected("암호화")

	runBtn := widget.NewButton("Run", func() {
		key := keyEntry.Text
		input := inputEntry.Text

		if key == "" {
			showError("실패: 키가 입력되지 않았습니다")
			return
		}
		if input == "" {
			if modeEncrypt {
				showError("실패: 평문이 입력되지 않았습니다")
			} else {
				showError("실패: 암호문이 입력되지 않았습니다")
			}
			return
		}

		if modeEncrypt {
			result, err := Encrypt(key, input)
			if err != nil {
				showError("암호화 실패: " + err.Error())
			} else {
				showResult(result)
			}
		} else {
			inner := input
			if strings.HasPrefix(inner, "ENC(") && strings.HasSuffix(inner, ")") {
				inner = inner[4 : len(inner)-1]
			}
			result, err := Decrypt(key, inner)
			if err != nil {
				showError("복호화 실패: " + err.Error())
			} else {
				showResult(result)
			}
		}
	})

	content := container.NewVBox(
		radio,
		widget.NewSeparator(),
		widget.NewLabel("키"),
		keyEntry,
		inputLabel,
		inputEntry,
		layout.NewSpacer(),
		runBtn,
		widget.NewSeparator(),
		widget.NewLabel("결과"),
		resultContainer,
	)

	w.SetContent(container.NewPadded(content))
	w.ShowAndRun()
}
