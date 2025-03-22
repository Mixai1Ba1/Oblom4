package main

import (
	"fmt"
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

	content := widget.NewRichTextFromMarkdown("Выберите тему слева или используйте поиск.")
	content.Wrapping = fyne.TextWrapWord

	sortedKeys := getSortedKeys()
	leftPanel := container.NewVBox(widget.NewLabelWithStyle("Темы", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

	for _, topic := range sortedKeys {
		t := topic
		btn := widget.NewButton(t, func() {
			content.ParseMarkdown(fmt.Sprintf("📘 **%s**\n\n%s", t, topics[t]))
		})
		leftPanel.Add(btn)
	}
	leftPanel.Resize(fyne.NewSize(300, 600))

	// === Поиск на лету (посимв) + по Enter (шаблон) ===
	searchEntry.OnChanged = func(text string) {
		if modeLabel.Text == "Режим: Посимвольный поиск" {
			text = strings.TrimSpace(strings.ToLower(text))
			if text == "" {
				content.ParseMarkdown("🔍 Введите запрос для поиска.")
				return
			}

			var result string
			for topic, desc := range topics {
				if strings.Contains(strings.ToLower(topic), text) || strings.Contains(strings.ToLower(desc), text) {
					preview := highlightText(desc, text)
					result += fmt.Sprintf("🔹 **%s**\n%s\n\n", topic, preview)
				}
			}

			if result == "" {
				content.ParseMarkdown("❌ Нет совпадений.")
			} else {
				content.ParseMarkdown(result)
			}
		}
	}

	searchEntry.OnSubmitted = func(text string) {
		if modeLabel.Text == "Режим: Поиск по шаблону" {
			text = strings.TrimSpace(text)
			found := false
			for topic, desc := range topics {
				if strings.Contains(strings.ToLower(desc), strings.ToLower(text)) {
					highlighted := highlightText(desc, text)
					content.ParseMarkdown(fmt.Sprintf("🔍 **Найдено в \"%s\"**\n\n%s", topic, highlighted))
					found = true
					break
				}
			}
			if !found {
				content.ParseMarkdown("❌ Ничего не найдено.")
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

func getSortedKeys() []string {
	keys := make([]string, 0, len(topics))
	for k := range topics {
		keys = append(keys, k)
	}
	return keys
}

// Подсветка совпавшего текста **жирным**
func highlightText(text, query string) string {
	lowered := strings.ToLower(text)
	loweredQuery := strings.ToLower(query)

	index := strings.Index(lowered, loweredQuery)
	if index == -1 {
		return text
	}

	original := text[index : index+len(query)]
	return text[:index] + "**" + original + "**" + text[index+len(query):]
}
