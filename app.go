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
	BallyPlus      = "Начислить баллы 🌟"
	BallyMinus     = "Списать баллы ⛔️"
	BallyState     = "Проверить баллы 🚩"
	Help           = "Справка 🆘"
	BackToMainMenu = "Вернуться в главное меню ⬅️"
	Back           = "Вернуться ⬅️"
)

var (
	token = os.Getenv("TOKEN")
	kbd   string
)

func main() {
	// Инициализация бота с токеном
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("Ошибка инициализации бота: %s", err)
	}

	fmt.Printf("Авторизован успешно от имени %s\n", bot.Self.UserName)

	bot.Debug = true

	// Настройка обновлений (получаем сообщения от пользователей)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		fmt.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, update.Message.Chat.ID, "Привет! Я ваш новый семйный Telegram-бот для удобвтсва работы с заданиями и начилсения баллов. 😉")
				sendStartMenu(bot, update.Message.Chat.ID)
			default:
				continue
			}
		} else {
			switch update.Message.Text {
			case Help:
				sendMessage(bot, update.Message.Chat.ID, "Хорошие дела и баллы!\n\n1. Помощь по дому 🏠\nУборка в своей комнате — 5 баллов\nМытьё посуды — 5 баллов\nВынос мусора — 3 балла\nПомощь в готовке еды — 7 баллов\nПомощь в уборке по дому - 7 баллов. \n\n2. Учёба и саморазвитие 👨‍🎓\nВыполненные домашние задания без напоминаний — 10 баллов\nЧтение книги (каждые 30 минут) — 7 баллов\nВыученное новое стихотворение или песня — 10 баллов\nУлучшение оценок по предмету — 15 баллов\n\n3. Доброта и хорошее поведение 👍\nПомощь брату или сестре — 7 баллов\nВежливость и отсутствие конфликтов весь день — 5 баллов\nДобрые поступки (например, помочь соседу или другу) — 10 баллов\n\n4. Физическая активность 🔋\nУтренняя зарядка или упражнения — 5 баллов\nПрогулка на свежем воздухе (более 1 часа) — 5 баллов\nУчастие в спортивных занятиях или соревнованиях — 10 баллов\n\n---\n\nНаграды за баллы! 🏆\n\n1. Маленькие награды (до 30 баллов):\nДополнительное время на планшет (30 минут) — 20 баллов\nНаклейка или маленькая игрушка — 25 баллов\nСладость (например, шоколадка) — 30 баллов\n\n2. Средние награды (30–70 баллов):\nПоход в кино или кафе — 50 баллов\nВыбор мультфильма или фильма для вечера — 50 баллов\nЛего или настольная игра — 70 баллов\n\n3. Крупные награды (70+ баллов):\nБольшая игрушка — 100 баллов\nПоездка в парк аттракционов — 150 баллов\nДень без домашних обязанностей — 70 баллов.")
			case BallyPlus:
				switch update.Message.From.UserName {
				case "arsyacrk", "Alyona_270688":
					kbd = "plus"
					sendPlusBally(bot, update.Message.Chat.ID)
				default:
					sendMessage(bot, update.Message.Chat.ID, "Вам зарпещено начислять баллы!")
					continue
				}
			case BallyMinus:
				switch update.Message.From.UserName {
				case "arsyacrk", "Alyona_270688":
					kbd = "minus"
					sendMinusBally(bot, update.Message.Chat.ID)
				default:
					sendMessage(bot, update.Message.Chat.ID, "Вам зарпещено списывать баллы!")
					continue
				}
			case BackToMainMenu:
				sendStartMenu(bot, update.Message.Chat.ID)
			case BallyState:
				loadScores(bot, update.Message.Chat.ID)
			case Bogdan, Veronika:
				kbdSectionsPlus(bot, update.Message.Chat.ID)
			case Back:
				switch kbd {
				case "plus":
					sendPlusBally(bot, update.Message.Chat.ID)
				case "minus":
					sendMinusBally(bot, update.Message.Chat.ID)
				default:
					sendStartMenu(bot, update.Message.Chat.ID)
				}
			default:
				continue
			}
		}
	}

}

// Функция отправки сообщения
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция вызова стартового меню
func sendStartMenu(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")

	// Создаем кастомную клавиатуру
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BallyPlus),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BallyMinus),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BallyState),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Help),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

//Клавиатура начиления баллов

func kbdSectionsPlus(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите, сколько начислить баллов и за что: ")

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
			tgbotapi.NewKeyboardButton(Back),
			tgbotapi.NewKeyboardButton(BackToMainMenu),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

// Функция вызова меню начилсения баллов
func sendPlusBally(bot *tgbotapi.BotAPI, chatID int64) {
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
func sendMinusBally(bot *tgbotapi.BotAPI, chatID int64) {
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
