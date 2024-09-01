package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"go_chat_client/internal/connection"
	"os"
	"time"
)

type RawRequest struct {
	Type string     `json:"type"`
	Data ReqMessage `json:"data"`
}

type ReqMessage struct {
	MsgID  string    `json:"msg_id"`
	Body   string    `json:"body"`
	TDate  time.Time `json:"t_date"`
	FromID string    `json:"from_id"`
}

const questionNewChat = "Создать новый чат с другим пользователем"
const questionEnterChat = "Войти в чат с пользователем"

func NewMenu() *cobra.Command {
	var configPath string

	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start chat",
		RunE: func(cmd *cobra.Command, args []string) error {

			err := navigateMenu()
			if err != nil {
				return err
			}

			return nil
		},
	}
	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}

func navigateMenu() error {
	connection.Establish()
	defer connection.Close()

	answers := struct {
		Choice string
	}{}

	var err error

	questions := []*survey.Question{
		{
			Name: "choice",
			Prompt: &survey.Select{
				Message: "Выберите опцию:",
				Options: []string{
					questionNewChat,
					questionEnterChat,
				},
			},
		},
	}

	// Conduct the survey
	err = survey.Ask(questions, &answers)
	if err != nil {
		return fmt.Errorf("ошибка во время опроса: %v", err)
	}

	// Output the selected answer
	fmt.Printf("Вы выбрали: %s\n", answers.Choice)

	// Follow-up questions based on the selected option
	if answers.Choice == questionNewChat {
		var newChatUserId string
		err = survey.AskOne(&survey.Input{
			Message: "Введите айди пользователя, с которым вы бы хотели начать чат или введите return для выхода в предыдущее меню:",
		}, &newChatUserId)
		if err != nil {
			return fmt.Errorf("ошибка во время опроса: %v", err)
		}

		if newChatUserId == "" {
			return navigateMenu()
		}

		createChat(newChatUserId)
		return navigateMenu()

	} else if answers.Choice == questionEnterChat {
		var chatId string
		err = survey.AskOne(&survey.Input{
			Message: "Введите айди чата для начала общения или введите return для выхода в предыдущее меню: ",
		}, &chatId)
		if err != nil {
			return fmt.Errorf("ошибка во время опроса: %v", err)
		}

		if chatId == "" {
			return navigateMenu()
		}

		go ReadMessages()

		fmt.Printf("Вводите сообщения для отправки:\n")

		renderOutputMessageFrame(chatId)
	}

	return nil
}

func renderOutputMessageFrame(chatId string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		inputBufferString, err := reader.ReadString('\n')
		fmt.Print("\033[F")
		fmt.Print("\033[K")

		fmt.Println(">")
		fmt.Println(GetFormattedTime())
		fmt.Print(inputBufferString)

		if err != nil {
			fmt.Println("Error survey.AskOne:", err)
		}
		fmt.Println(">")
		newMessage(inputBufferString, chatId)
	}
}

func ReadMessages() {
	for {
		typ, message, err := connection.WebSocket.ReadMessage()

		if err != nil {
			fmt.Println("ws error", err)
		}
		switch typ {

		case websocket.TextMessage, websocket.BinaryMessage:
			var req RawRequest
			if err = json.Unmarshal(message, &req); err != nil {
				fmt.Println("ws error", err)
				continue
			}

			RenderInputMessageFrame(req.Data)
		default:
			fmt.Println("smth went wrong")
		}

		time.Sleep(1 * time.Second)
	}
}

func RenderInputMessageFrame(message ReqMessage) {
	fmt.Print("\033[F")
	fmt.Println(">")
	fmt.Println(message.FromID, GetFormattedTime())
	fmt.Println(message.Body)
	fmt.Println(">")
}
