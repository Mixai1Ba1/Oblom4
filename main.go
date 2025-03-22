package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var topics = map[string]string{
	"Закон Хика":                `Закон Хика описывает время, необходимое для принятия решения, как логарифм от количества альтернатив. Подробнее: https://ru.wikipedia.org/wiki/Закон_Хика`,
	"Меню в интерфейсах":        `Меню представляют собой иерархию пунктов, позволяющих пользователю выбирать команды.`,
	"Эвристики Юзабилити":       `Набор принципов для оценки удобства интерфейса, например, рекомендации Нильсена. Ссылка: https://www.nngroup.com/articles/ten-usability-heuristics/`,
	"Горячие клавиши":           `Сочетания клавиш, ускоряющие выполнение команд без использования мыши.`,
	"Графические интерфейсы":    `GUI позволяет взаимодействовать с программами через визуальные элементы.`,
	"Логика поиска по шаблону":  `Поиск осуществляется точным совпадением слова в тексте.`,
	"Посимвольный поиск":        `По мере ввода текста список подходящих результатов обновляется.`,
	"Пользовательские сценарии": `Use-case описывает действия пользователя для достижения цели.`,
	"Метрики оценки интерфейса": `Показатели, такие как точность, скорость и удовлетворенность пользователя.`,
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Справочная система")
	myWindow.Resize(fyne.NewSize(1000, 600))

	modeLabel := widget.NewLabel("Режим: Посимвольный поиск")
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Введите запрос...")

	content := widget.NewRichText()
	content.Wrapping = fyne.TextWrapWord
	content.Segments = []widget.RichTextSegment{
		&widget.TextSegment{Text: "Выберите тему слева или используйте поиск."},
	}

	sortedKeys := getSortedKeys()
	leftPanel := container.NewVBox(widget.NewLabelWithStyle("Темы", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

	for _, topic := range sortedKeys {
		t := topic
		btn := widget.NewButton(t, func() {
			content.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text:      "📘 " + t,
					StyleName: widget.RichTextStyleNameStrong,
				},
				&widget.TextSegment{Text: "\n\n" + topics[t]},
			}
			content.Refresh()
		})
		leftPanel.Add(btn)
	}
	leftPanel.Resize(fyne.NewSize(300, 600))

	// === Поиск ===
	searchEntry.OnChanged = func(text string) {
		if modeLabel.Text == "Режим: Посимвольный поиск" {
			text = strings.TrimSpace(strings.ToLower(text))
			if text == "" {
				content.Segments = []widget.RichTextSegment{
					&widget.TextSegment{Text: "🔍 Введите запрос для поиска."},
				}
				content.Refresh()
				return
			}

			var segments []widget.RichTextSegment
			for topic, desc := range topics {
				if strings.Contains(strings.ToLower(topic), text) || strings.Contains(strings.ToLower(desc), text) {
					segments = append(segments,
						&widget.TextSegment{
							Text:      "🔹 " + topic,
							StyleName: widget.RichTextStyleNameStrong,
						},
						&widget.TextSegment{Text: "\n"},
					)
					segments = append(segments, highlightTextSegments(desc, text)...)
					segments = append(segments, &widget.TextSegment{Text: "\n\n"})
				}
			}

			if len(segments) == 0 {
				content.Segments = []widget.RichTextSegment{
					&widget.TextSegment{Text: "❌ Нет совпадений."},
				}
			} else {
				content.Segments = segments
			}
			content.Refresh()
		}
	}

	searchEntry.OnSubmitted = func(text string) {
		if modeLabel.Text == "Режим: Поиск по шаблону" {
			text = strings.TrimSpace(text)
			found := false
			for topic, desc := range topics {
				if strings.Contains(strings.ToLower(desc), strings.ToLower(text)) {
					segments := []widget.RichTextSegment{
						&widget.TextSegment{
							Text:      "🔍 Найдено в \"" + topic + "\"",
							StyleName: widget.RichTextStyleNameStrong,
						},
						&widget.TextSegment{Text: "\n"},
					}
					segments = append(segments, highlightTextSegments(desc, text)...)
					content.Segments = segments
					content.Refresh()
					found = true
					break
				}
			}
			if !found {
				content.Segments = []widget.RichTextSegment{
					&widget.TextSegment{Text: "❌ Ничего не найдено."},
				}
				content.Refresh()
			}
		}
	}

	modeToggle := widget.NewButtonWithIcon("Переключить режим", theme.ViewRefreshIcon(), func() {
		if modeLabel.Text == "Режим: Посимвольный поиск" {
			modeLabel.SetText("Режим: Поиск по шаблону")
		} else {
			modeLabel.SetText("Режим: Посимвольный поиск")
		}
	})

	rightPanel := container.NewVSplit(
		container.NewVBox(modeLabel, searchEntry, modeToggle),
		content,
	)
	rightPanel.Offset = 0.25

	myWindow.SetContent(container.NewHSplit(leftPanel, rightPanel))
	myWindow.ShowAndRun()
}

// Возвращает список ключей в порядке добавления
func getSortedKeys() []string {
	keys := make([]string, 0, len(topics))
	for k := range topics {
		keys = append(keys, k)
	}
	return keys
}

// Разбивает текст на сегменты с подсветкой совпадений
func highlightTextSegments(text, query string) []widget.RichTextSegment {
	var segments []widget.RichTextSegment

	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	start := 0
	for {
		index := strings.Index(lowerText[start:], lowerQuery)
		if index == -1 {
			break
		}
		index += start

		if index > start {
			segments = append(segments, &widget.TextSegment{Text: text[start:index]})
		}

		segments = append(segments, &widget.TextSegment{
			Text:      text[index : index+len(query)],
			StyleName: widget.RichTextStyleNameEmphasis, // синяя/курсивная подсветка
		})

		start = index + len(query)
	}
	if start < len(text) {
		segments = append(segments, &widget.TextSegment{Text: text[start:]})
	}

	return segments
}
