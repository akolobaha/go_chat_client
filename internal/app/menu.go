package app

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"go_chat_client/internal/connection"
)

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

		for {
			renderMessageFrame(chatId)
		}

	}

	return nil
}
