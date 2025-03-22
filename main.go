package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var topics = map[string]string{
	"Закон Хика": `Закон Хика описывает время, необходимое для принятия решения, как логарифм от количества альтернатив.
Подробнее: https://ru.wikipedia.org/wiki/Закон_Хика`,
	"Меню в интерфейсах": `Меню представляют собой иерархию пунктов, позволяющих пользователю выбирать команды.`,
	"Эвристики Юзабилити": `Набор принципов для оценки удобства интерфейса, например, рекомендации Нильсена.
Ссылка: https://www.nngroup.com/articles/ten-usability-heuristics/`,
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

	modeLabel := widget.NewLabel("Режим: Поиск по шаблону")
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Введите запрос...")

	content := widget.NewRichTextFromMarkdown("Выберите тему слева или используйте поиск.")
	content.Wrapping = fyne.TextWrapWord

	sortedKeys := getSortedKeys()
	topicList := widget.NewList(
		func() int { return len(sortedKeys) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(sortedKeys[i])
		},
	)
	topicList.OnSelected = func(id int) {
		key := sortedKeys[id]
		content.ParseMarkdown("📘 **" + key + "**\n\n" + topics[key])
	}

	modeToggle := widget.NewButton("🔁 Переключить режим", func() {
		if modeLabel.Text == "Режим: Поиск по шаблону" {
			modeLabel.SetText("Режим: Посимвольный поиск")
		} else {
			modeLabel.SetText("Режим: Поиск по шаблону")
		}
	})

	// searchEntry.OnChanged = func(text string) {
	// 	if modeLabel.Text == "Режим: Посимвольный поиск" {
	// 		filtered := ""
	// 		for topic := range topics {
	// 			if strings.Contains(strings.ToLower(topic), strings.ToLower(text)) {
	// 				filtered += "• " + topic + "\n"
	// 			}
	// 		}
	// 		if filtered == "" {
	// 			content.ParseMarkdown("❌ Нет совпадений.")
	// 		} else {
	// 			content.ParseMarkdown("🔡 Найдено:\n" + filtered)
	// 		}
	// 	}
	// }
	searchEntry.OnChanged = func(text string) {
		if modeLabel.Text == "Режим: Посимвольный поиск" {
			filtered := ""
			for topic := range topics {
				if strings.Contains(strings.ToLower(topic), strings.ToLower(text)) {
					filtered += "• " + topic + "\n"
				}
			}
			if filtered == "" {
				content.ParseMarkdown("❌ Нет совпадений.")
			} else {
				content.ParseMarkdown("🔡 Найдены темы:\n" + filtered)
			}
		}
	}

	// searchEntry.OnSubmitted = func(text string) {
	// 	if modeLabel.Text == "Режим: Поиск по шаблону" {
	// 		found := false
	// 		for topic, desc := range topics {
	// 			if strings.Contains(strings.ToLower(desc), strings.ToLower(text)) {
	// 				highlighted := strings.ReplaceAll(desc, text, "**"+text+"**")
	// 				content.ParseMarkdown("🔍 **Найдено в \"" + topic + "\"**\n\n" + highlighted)
	// 				found = true
	// 				break
	// 			}
	// 		}
	// 		if !found {
	// 			content.ParseMarkdown("❌ Ничего не найдено.")
	// 		}
	// 	}
	// }
	searchEntry.OnSubmitted = func(text string) {
		if modeLabel.Text == "Режим: Поиск по шаблону" {
			found := false
			for topic, desc := range topics {
				index := strings.Index(strings.ToLower(desc), strings.ToLower(text))
				if index != -1 {
					original := desc[index : index+len(text)]
					highlighted := desc[:index] + "**" + original + "**" + desc[index+len(text):]
					content.ParseMarkdown("🔍 **Найдено в \"" + topic + "\"**\n\n" + highlighted)
					found = true
					break
				}
			}
			if !found {
				content.ParseMarkdown("❌ Ничего не найдено.")
			}
		}
	}

	// leftScroll := container.NewVScroll(container.NewVBox(
	// 	widget.NewLabelWithStyle("Темы", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	// 	topicList,
	// ))
	// leftScroll.SetMinSize(fyne.NewSize(300, 550))
	// Статичный список тем слева
	topicButtons := make([]fyne.CanvasObject, 0)

	for _, topic := range sortedKeys {
		t := topic // копия, чтобы не залипло
		btn := widget.NewButton(topic, func() {
			content.ParseMarkdown("📘 **" + t + "**\n\n" + topics[t])
		})
		topicButtons = append(topicButtons, btn)
	}

	leftPanel := container.NewVBox(
		widget.NewLabelWithStyle("Темы", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	leftPanel.Add(container.NewVBox(topicButtons...))
	leftPanel.Resize(fyne.NewSize(300, 600)) // фиксированная высота и ширина

	// Справа: поиск и текст
	rightPanel := container.NewVSplit(
		container.NewVBox(modeLabel, searchEntry, modeToggle),
		content,
	)
	rightPanel.Offset = 0.25

	// myWindow.SetContent(container.NewHSplit(leftScroll, rightPanel))
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
