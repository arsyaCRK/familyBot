package main

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	Bogdan         = "Богдан"
	Veronika       = "Вероника"
	BackToMainMenu = "Вернуться в главное меню ⬅️"
	Back           = "Вернуться ⬅️"
)

var (
	token = os.Getenv("TOKEN")
	name  string
)

type UserState struct {
	CurrentKeyboard string
	KeyboardHistory []string
}

var userStates = make(map[int64]*UserState)

func main() {
	// Инициализация бота с токеном
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("Ошибка инициализации бота: %s", err)
	}

	fmt.Printf("Авторизован успешно от имени %s\n", bot.Self.UserName)

	bot.Debug = false

	// Настройка обновлений (получаем сообщения от пользователей)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//fmt.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		fmt.Printf("Колбэк: %s", update.CallbackQuery)

		if update.CallbackQuery != nil {

			handleCallbackQuery(bot, update.CallbackQuery)
			continue
		} else if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, update.Message.Chat.ID, "Привет! Я ваш новый семйный Telegram-бот для удобвтсва работы с заданиями и начилсения баллов. 😉")
				sendStartMenu(bot, update.Message.Chat.ID)
			default:
				continue
			}
		} else {
			switch update.Message.Text {
			case "Справка 🆘":
				sendMessage(bot, update.Message.Chat.ID, "Хорошие дела и баллы!\n\n1. Помощь по дому 🏠\nУборка в своей комнате — 5 баллов\nМытьё посуды — 5 баллов\nВынос мусора — 3 балла\nПомощь в готовке еды — 7 баллов\nПомощь в уборке по дому - 7 баллов. \n\n2. Учёба и саморазвитие 👨‍🎓\nВыполненные домашние задания без напоминаний — 10 баллов\nЧтение книги (каждые 30 минут) — 7 баллов\nВыученное новое стихотворение или песня — 10 баллов\nУлучшение оценок по предмету — 15 баллов\n\n3. Доброта и хорошее поведение 👍\nПомощь брату или сестре — 7 баллов\nВежливость и отсутствие конфликтов весь день — 5 баллов\nДобрые поступки (например, помочь соседу или другу) — 10 баллов\n\n4. Физическая активность 🔋\nУтренняя зарядка или упражнения — 5 баллов\nПрогулка на свежем воздухе (более 1 часа) — 5 баллов\nУчастие в спортивных занятиях или соревнованиях — 10 баллов\n\n---\n\nНаграды за баллы! 🏆\n\n1. Маленькие награды (до 30 баллов):\nДополнительное время на планшет (30 минут) — 20 баллов\nНаклейка или маленькая игрушка — 25 баллов\nСладость (например, шоколадка) — 30 баллов\n\n2. Средние награды (30–70 баллов):\nПоход в кино или кафе — 50 баллов\nВыбор мультфильма или фильма для вечера — 50 баллов\nЛего или настольная игра — 70 баллов\n\n3. Крупные награды (70+ баллов):\nБольшая игрушка — 100 баллов\nПоездка в парк аттракционов — 150 баллов\nДень без домашних обязанностей — 70 баллов.")
			case "Начислить баллы 🌟":
				switch update.Message.From.UserName {
				case "arsyacrk", "Alyona_270688":
					pushKeyboardHistory(update.Message.Chat.ID, "menuPlus1")
					sendPlusScores(bot, update.Message.Chat.ID)
				default:
					sendMessage(bot, update.Message.Chat.ID, "Вам зарпещено начислять баллы!")
					continue
				}
			case "Списать баллы ⛔️":
				switch update.Message.From.UserName {
				case "arsyacrk", "Alyona_270688":
					pushKeyboardHistory(update.Message.Chat.ID, "menuMinus1")
					sendMinusScores(bot, update.Message.Chat.ID)
				default:
					sendMessage(bot, update.Message.Chat.ID, "Вам зарпещено списывать баллы!")
					continue
				}
			case BackToMainMenu:
				pushKeyboardHistory(update.Message.Chat.ID, "")
				sendStartMenu(bot, update.Message.Chat.ID)
			case "Проверить баллы 🚩":
				loadScores(bot, update.Message.Chat.ID)
			case Bogdan, Veronika:
				name = update.Message.Text
				pushKeyboardHistory(update.Message.Chat.ID, "menuPlusMinus2")
				kbdSectionsPlus(bot, update.Message.Chat.ID)
			case "Помощь по дому 🏠":
				pushKeyboardHistory(update.Message.Chat.ID, "menuPlus3")
				kbdScoresHomeHelp(bot, update.Message.Chat.ID)
			case "Учёба и саморазвитие 👨‍🎓":
				pushKeyboardHistory(update.Message.Chat.ID, "menuPlus3")
				kbdScoresLearning(bot, update.Message.Chat.ID)
			case "Доброта и хорошее поведение 👍":
				pushKeyboardHistory(update.Message.Chat.ID, "menuPlus3")
				kbdScoresGoodness(bot, update.Message.Chat.ID)
			case "Физическая активность 🔋":
				pushKeyboardHistory(update.Message.Chat.ID, "menuPlus3")
				kbdScoresActivity(bot, update.Message.Chat.ID)
			case "Начислить баллы вручную":
				kbdReplyManual(bot, update.Message.Chat.ID, name)
			case Back:
				if lastKeyboard, ok := popKeyboardHistory(update.Message.Chat.ID); ok {
					switch lastKeyboard {
					case "menuPlus1":
						sendPlusScores(bot, update.Message.Chat.ID)
					case "menuPlusMinus2":
						kbdSectionsPlus(bot, update.Message.Chat.ID)
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нет предыдущего меню")
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				}
			default:
				continue
			}
		}
	}

}

// Функции для работы с историей клавиатур
func pushKeyboardHistory(chatID int64, keyboardName string) {
	if _, ok := userStates[chatID]; !ok {
		userStates[chatID] = &UserState{
			CurrentKeyboard: keyboardName,
			KeyboardHistory: []string{},
		}
		return
	}

	// Сохраняем текущую клавиатуру в историю
	if userStates[chatID].CurrentKeyboard != "" {
		userStates[chatID].KeyboardHistory = append(userStates[chatID].KeyboardHistory, userStates[chatID].CurrentKeyboard)
	}

	// Устанавливаем новую клавиатуру как текущую
	userStates[chatID].CurrentKeyboard = keyboardName
}

func popKeyboardHistory(chatID int64) (string, bool) {
	if state, ok := userStates[chatID]; ok && len(state.KeyboardHistory) > 0 {
		// Получаем последнюю клавиатуру из истории
		lastKeyboard := state.KeyboardHistory[len(state.KeyboardHistory)-1]

		// Удаляем её из истории
		state.KeyboardHistory = state.KeyboardHistory[:len(state.KeyboardHistory)-1]

		// Устанавливаем её как текущую
		state.CurrentKeyboard = lastKeyboard

		return lastKeyboard, true
	}
	return "", false
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	fmt.Printf("Данные: %s", data)

	// Отправляем ответ на callback (убирает "часики" у кнопки)
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Println(err)
	}

	switch data {
	case "10":
		sendMessage(bot, chatID, "Вы выбрали действие 10")
	case "20":
		sendMessage(bot, chatID, "Вы выбрали Inline действие 20")
	}
}

// Функция отправки сообщения
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция удаления сообщения
func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	// Создаем запрос на удаление сообщения
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)

	// Отправляем запрос
	_, err := bot.Request(deleteConfig)
	if err != nil {
		log.Printf("Ошибка при удалении сообщения: %v", err)
	} else {
		log.Println("Сообщение успешно удалено")
	}
}

// Функция вызова стартового меню
func sendStartMenu(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")

	// Создаем кастомную клавиатуру
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начислить баллы 🌟"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Списать баллы ⛔️"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Проверить баллы 🚩"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Справка 🆘"),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func kbdReplyManual(bot *tgbotapi.BotAPI, chatID int64, messageText string) {
	messageText = fmt.Sprintf("Укажите сколько баллов начислить %s: ", messageText)
	msg := tgbotapi.NewMessage(chatID, messageText)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5", "10"),
			tgbotapi.NewInlineKeyboardButtonData("10", "20"),
			tgbotapi.NewInlineKeyboardButtonData("15", "20"),
			tgbotapi.NewInlineKeyboardButtonData("20", "20"),
			tgbotapi.NewInlineKeyboardButtonData("30", "20"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("40", "10"),
			tgbotapi.NewInlineKeyboardButtonData("50", "20"),
			tgbotapi.NewInlineKeyboardButtonData("60", "20"),
			tgbotapi.NewInlineKeyboardButtonData("70", "20"),
			tgbotapi.NewInlineKeyboardButtonData("80", "20"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("90", "10"),
			tgbotapi.NewInlineKeyboardButtonData("100", "20"),
			tgbotapi.NewInlineKeyboardButtonData("150", "20"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Клавиатура начиления баллов
func kbdSectionsPlus(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Помощь по дому 🏠"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Учёба и саморазвитие 👨‍🎓"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Доброта и хорошее поведение 👍"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Физическая активность 🔋"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начислить баллы вручную"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Back),
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Клавиатура для начиления баллов за помощь по дому
func kbdScoresHomeHelp(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, сколько начислить баллов и за что: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Уборка в своей комнате — 5 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Мытьё посуды — 5 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вынос мусора — 3 балла"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Помощь в готовке еды — 7 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Помощь в уборке по дому - 7 баллов."),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Back),
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция клавиатуры за доброту и хорошее поведение
func kbdScoresGoodness(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, сколько начислить баллов и за что: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Помощь брату или сестре — 7 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вежливость и отсутствие конфликтов весь день — 5 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добрые поступки (например, помочь соседу или другу) — 10 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Back),
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Клавиатура для начиления баллов за Учёба и саморазвитие
func kbdScoresLearning(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, сколько начислить баллов и за что: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Домашние задания без напоминаний — 10 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Чтение книги (каждые 30 минут) — 7 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Выученное новое стихотворение или песня — 10 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Улучшение оценок по предмету — 15 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Back),
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Клавиатура для начиления баллов за физическую активность
func kbdScoresActivity(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, сколько начислить баллов и за что: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Утренняя зарядка или упражнения — 5 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Прогулка на свежем воздухе (более 1 часа) — 5 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Участие в спортивных занятиях или соревнованиях — 10 баллов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Back),
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция вызова меню начилсения баллов
func sendPlusScores(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, кому начилсить былллы: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Bogdan),
			tgbotapi.NewKeyboardButton(Veronika),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция вызова меню списания баллов
func sendMinusScores(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, у кого списать баллы: ")

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Bogdan),
			tgbotapi.NewKeyboardButton(Veronika),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция загружает баллы из файла
func loadScores(bot *tgbotapi.BotAPI, chatID int64) {

	file, err := os.Open("scores.txt")
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			name := parts[0]
			score := parts[1]

			state := fmt.Sprintf("%s: %s", name, score)
			voiceFileName := fmt.Sprintf("voice%s.mp3", name)
			textToSpeech(bot, chatID, state, voiceFileName)

			msg := tgbotapi.NewMessage(chatID, state)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}

func textToSpeech(bot *tgbotapi.BotAPI, chatID int64, text string, filename string) {
	t := text
	outputFile := filename
	cmd := exec.Command("espeak-ng", "-v", "ru", "-w", outputFile, t)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды:", err)
		os.Exit(1)
		return
	}

	fmt.Printf("Аудиофайл успешно создан: %s\n", outputFile)

	msg := tgbotapi.NewVoice(chatID, tgbotapi.FilePath(outputFile))

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
