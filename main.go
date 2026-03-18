package main

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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

	// 결과 영역: 성공 시 복사 가능한 Entry, 실패 시 빨간 글씨 Text
	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetMinRowsVisible(2)
	resultEntry.Disable()

	errorText := canvas.NewText("", color.RGBA{R: 220, G: 30, B: 30, A: 255})
	errorText.TextSize = 14
	errorText.Hide()

	resultContainer := container.NewStack(resultEntry, errorText)

	showResult := func(text string) {
		errorText.Hide()
		resultEntry.Enable()
		resultEntry.SetText(text)
		resultEntry.Disable()
		resultEntry.Show()
	}

	showError := func(text string) {
		resultEntry.SetText("")
		resultEntry.Hide()
		errorText.Text = text
		errorText.Color = color.RGBA{R: 220, G: 30, B: 30, A: 255}
		errorText.Show()
		errorText.Refresh()
	}

	clearResult := func() {
		resultEntry.SetText("")
		resultEntry.Show()
		errorText.Hide()
	}

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
