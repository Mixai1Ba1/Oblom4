package main

import (
	"fmt"
	"sort"
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

type Match struct {
	Word  string
	Topic string
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Справочная система")
	myWindow.Resize(fyne.NewSize(1000, 600))

	modeLabel := widget.NewLabel("Режим: Посимвольный поиск")
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Введите запрос...")

	contentList := container.NewVBox()
	contentScroll := container.NewVScroll(contentList)

	sortedKeys := getSortedKeys()
	leftPanel := container.NewVBox(widget.NewLabelWithStyle("Темы", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

	for _, topic := range sortedKeys {
		t := topic
		btn := widget.NewButton(t, func() {
			contentList.Objects = []fyne.CanvasObject{
				widget.NewRichTextFromMarkdown(fmt.Sprintf("📘 **%s**", t)),
				widget.NewRichTextFromMarkdown(topics[t]),
			}
			contentList.Refresh()
		})
		leftPanel.Add(btn)
	}
	leftPanel.Resize(fyne.NewSize(300, 600))

	searchEntry.OnChanged = func(text string) {
		if modeLabel.Text != "Режим: Посимвольный поиск" {
			return
		}

		text = strings.TrimSpace(strings.ToLower(text))
		if text == "" {
			contentList.Objects = []fyne.CanvasObject{widget.NewRichTextFromMarkdown("🔍 Введите запрос для поиска.")}
			contentList.Refresh()
			return
		}

		var matches []Match
		for topic, desc := range topics {
			words := strings.Fields(desc)
			for _, word := range words {
				clean := strings.Trim(strings.ToLower(word), ".,!?(){}[]\"'")
				if strings.Contains(clean, text) {
					matches = append(matches, Match{Word: word, Topic: topic})
				}
			}
		}

		if len(matches) == 0 {
			contentList.Objects = []fyne.CanvasObject{widget.NewRichTextFromMarkdown("❌ Нет совпадений.")}
			contentList.Refresh()
			return
		}

		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Word < matches[j].Word
		})

		var cells []fyne.CanvasObject
		for _, match := range matches {
			highlightedWord := widget.NewRichTextFromMarkdown(highlightText(match.Word, text))
			context := widget.NewRichTextFromMarkdown(fmt.Sprintf("🔹 **%s**\n%s", match.Topic, highlightText(topics[match.Topic], text)))

			item := container.NewVBox(
				highlightedWord,
				context,
				widget.NewSeparator(),
			)
			cells = append(cells, item)
		}

		contentList.Objects = cells
		contentList.Refresh()
	}

	searchEntry.OnSubmitted = func(text string) {
		if modeLabel.Text != "Режим: Поиск по шаблону" {
			return
		}

		text = strings.TrimSpace(text)
		found := false
		for topic, desc := range topics {
			if strings.Contains(strings.ToLower(desc), strings.ToLower(text)) {
				contentList.Objects = []fyne.CanvasObject{
					widget.NewRichTextFromMarkdown(fmt.Sprintf("🔍 **Найдено в \"%s\"**", topic)),
					widget.NewRichTextFromMarkdown(highlightText(desc, text)),
				}
				contentList.Refresh()
				found = true
				break
			}
		}
		if !found {
			contentList.Objects = []fyne.CanvasObject{widget.NewRichTextFromMarkdown("❌ Ничего не найдено.")}
			contentList.Refresh()
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
		contentScroll,
	)
	rightPanel.Offset = 0.25

	myWindow.SetContent(container.NewHSplit(leftPanel, rightPanel))
	myWindow.ShowAndRun()
}

func getSortedKeys() []string {
	keys := make([]string, 0, len(topics))
	for k := range topics {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func highlightText(text, query string) string {
	lowered := strings.ToLower(text)
	queryLower := strings.ToLower(query)
	var result strings.Builder

	start := 0
	for {
		index := strings.Index(lowered[start:], queryLower)
		if index == -1 {
			result.WriteString(text[start:])
			break
		}
		index += start
		result.WriteString(text[start:index])
		result.WriteString("**")
		result.WriteString(text[index : index+len(query)])
		result.WriteString("**")
		start = index + len(query)
	}
	return result.String()
}