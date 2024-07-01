// (c) Jisin0

package plugins

import (
	"fmt"
	"regexp"

	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var cbEditPattern *regexp.Regexp = regexp.MustCompile(`edit\((.+)\)`)

func Start(bot *gotgbot.Bot, update *ext.Context) error {
	go DB.AddUser(update.EffectiveMessage.From.Id)

	_, err := bot.SendMessage(
		update.Message.Chat.Id,
		fmt.Sprintf(utils.TEXT["START"], update.Message.From.FirstName, bot.FirstName),
		&gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: utils.BUTTONS["START"],
			},
			ReplyParameters: &gotgbot.ReplyParameters{
				AllowSendingWithoutReply: true,
			},
		})
	if err != nil {
		fmt.Printf("start: %v\n", err)
	}

	return nil
}

func Stats(bot *gotgbot.Bot, update *ext.Context) error {
	_, err := update.EffectiveMessage.Reply(bot, DB.Stats(), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if err != nil {
		fmt.Printf("stats: %v\n", err)
	}

	return nil
}

func GetId(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		text   string
		update = ctx.Message
	)

	if update.ReplyToMessage != nil {
		text += fmt.Sprintf("\nReplied to user : <code>%v</code>", update.ReplyToMessage.From.Id)

		if f := update.ReplyToMessage.ForwardOrigin; f.GetDate() != 0 {
			if f.MergeMessageOrigin().Chat != nil {
				text += fmt.Sprintf("\nForwarded from : <code>%v</code>", f.MergeMessageOrigin().Chat.Id)
			} else if f.MergeMessageOrigin().SenderChat != nil {
				text += fmt.Sprintf("\nForwarded from : <code>%v</code>", f.MergeMessageOrigin().SenderChat.Id)
			} else if f.MergeMessageOrigin().SenderUser != nil {
				text += fmt.Sprintf("\nForwarded from : <code>%v</code>", f.MergeMessageOrigin().SenderUser.Id)
			}

		}
	}

	text += fmt.Sprintf("\nUser id : <code>%v</code>", update.From.Id)

	if update.Chat.Type != gotgbot.ChatTypePrivate {
		text += fmt.Sprintf("\nChat id : <code>%v</code>", update.Chat.Id)
	}

	_, err := update.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyParameters: &gotgbot.ReplyParameters{MessageId: update.MessageId}})
	if err != nil {
		fmt.Printf("getid: %v\n", err)
	}

	return nil
}

func CbStats(bot *gotgbot.Bot, update *ext.Context) error {
	_, _, err := update.CallbackQuery.Message.EditText(bot, DB.Stats(), &gotgbot.EditMessageTextOpts{
		ChatId:      update.CallbackQuery.Message.GetChat().Id,
		MessageId:   update.CallbackQuery.Message.GetMessageId(),
		ParseMode:   gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: utils.BUTTONS["STATS"]},
	})

	if err != nil {
		fmt.Printf("cbstats: %v\n", err)
	}

	return nil
}

func FilterHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	// What R U Lookking At Its Just a Pro Function ;)
	go MFilter(b, ctx)
	go GFilter(b, ctx)

	return nil
}

// Function to handle edit() callbacks from the Start, About and Help menus
func CbEdit(bot *gotgbot.Bot, update *ext.Context) error {
	key := cbEditPattern.FindStringSubmatch(update.CallbackQuery.Data)[1]

	markup, ok := utils.BUTTONS[key]
	if !ok {
		markup = [][]gotgbot.InlineKeyboardButton{{{Text: "⤝ Bᴀᴄᴋ", CallbackData: "edit(HELP)"}}}
	}

	options := gotgbot.EditMessageTextOpts{
		ChatId:    update.CallbackQuery.Message.GetChat().Id,
		MessageId: update.CallbackQuery.Message.GetMessageId(),
		ParseMode: gotgbot.ParseModeHTML,
		LinkPreviewOptions: &gotgbot.LinkPreviewOptions{
			IsDisabled: true,
		},
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: markup},
	}

	var text string

	if key == "START" {
		text = fmt.Sprintf(utils.TEXT["START"], update.CallbackQuery.From.FirstName, bot.FirstName)
	} else {
		text = utils.TEXT[key]
	}

	_, _, err := update.CallbackQuery.Message.EditText(bot,
		text,
		&options,
	)
	if err != nil {
		fmt.Printf("cbedit: %v\n", err)
	}

	return nil
}

func About(b *gotgbot.Bot, update *ext.Context) error {
	_, err := update.Message.Reply(b, utils.TEXT["ABOUT"], &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: utils.BUTTONS["ABOUT"]}})
	if err != nil {
		fmt.Printf("about: %v\n", err)
	}

	return nil
}

func Help(b *gotgbot.Bot, update *ext.Context) error {
	_, err := update.Message.Reply(b, utils.TEXT["HELP"], &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: utils.BUTTONS["HELP"]}})
	if err != nil {
		fmt.Printf("help: %v\n", err)
	}

	return nil
}
