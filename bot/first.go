package bot

import (
	"dima_bot/primitives"
	"dima_bot/store"
	"dima_bot/vkapi"
	"fmt"
	"log"
	"strconv"
)

const ADMIN_VKID = 141017808;

var datesMap = map[int]string{
	1: "В течение месяца",
	2: "Через несколько месяцев",
	3: "В течение года",
}

var collectionsMap = map[int]string{
	1: "Базовая комплектация",
	2: "Теплый контур",
	3: "Под ключ",
}

var typesMap = map[int]string{
	1: "Клееный брус по ГОСТу — Сосна / Ель",
	2: "Клееный брус с утеплителем Белтермо",
	3: "Безусадочный CLT клееный брус с утеплителем Белтермо",
}

var storeysMap = map[int]string{
	1: "Один этаж",
	2: "Два этажа",
	3: "Свой проект",
}

var projectNames = map[int]string{
	1: "АРГУС ДОМ-БАНЯ [71м²]",
	2: "БЕГОНИЯ [88м²]",
	3: "ДЖЕК [124м²]",
	4: "КОЛИБРИ [97м²]",
	5: "САПСАН [134м²]",
	6: "СОЙКА [56м²]",
	7: "АГАТА [219м²]",
	8: "АЗАЛИЯ ЛЮКС [207м²]",
	9: "АКАЦИЯ [119м²]",
	10: "БАРБАРИС [194м²]",
	11: "БЕКАС [125м²]",
	12: "ГАРМОНИЯ МИНИ [140м²]",
	13: "ИВОЛГА [140м²]",
	14: "ЛОРИ [220м²]",
}

var projectPhotos = map[int]string{
	1: "photo-216949929_457239112,photo-216949929_457239113,photo-216949929_457239114,photo-216949929_457239115,photo-216949929_457239116,photo-216949929_457239117",
	2: "photo-216949929_457239109,photo-216949929_457239110,photo-216949929_457239111",
	3: "photo-216949929_457239118,photo-216949929_457239119,photo-216949929_457239120,photo-216949929_457239121,photo-216949929_457239122",
	4: "photo-216949929_457239123,photo-216949929_457239124,photo-216949929_457239125",
	5: "photo-216949929_457239126,photo-216949929_457239127,photo-216949929_457239128,photo-216949929_457239129",
	6: "photo-216949929_457239130,photo-216949929_457239131,photo-216949929_457239132",
	7: "АГАТА [219м²]",
	8: "АЗАЛИЯ ЛЮКС [207м²]",
	9: "АКАЦИЯ [119м²]",
	10: "БАРБАРИС [194м²]",
	11: "БЕКАС [125м²]",
	12: "ГАРМОНИЯ МИНИ [140м²]",
	13: "ИВОЛГА [140м²]",
	14: "ЛОРИ [220м²]",
}

func init() {
	*primitives.GetNode(Node1) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			kb := &vkapi.Keyboard{
				ButtonsGrid: [][]vkapi.Button{{vkapi.MakeButton("Жми", "", vkapi.KbBlue)}},
			}
			text := fmt.Sprintf("Привет, %s!\nГотовы сделать первый шаг к жизни в новом доме?\n\nЗаполните заявку и получите готовый ориентировочный расчёт.", msg.User.UserName)
			r := msg.Client.SendMessage(msg.PeerID, text, kb.String(), "")
			log.Println("send 1st msg", string(r))
		},
		Default: func(msg *primitives.Message) bool {
			msg.Trans(Node2)
			return true
		},
	}

	*primitives.GetNode(Node2) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			kb := &vkapi.Keyboard{
				ButtonsGrid: [][]vkapi.Button{
					{
						vkapi.MakeButton("В течение месяца", "1", vkapi.KbBlue),
					},
					{
						vkapi.MakeButton("Через несколько месяцев", "2", vkapi.KbBlue),
					},
					{
						vkapi.MakeButton("В течение года", "3", vkapi.KbBlue),
					},
				},
			}
			msg.Client.SendMessage(msg.PeerID, "Когда планируете начать строительство?", kb.String(), "")
		},
		Default: func(msg *primitives.Message) bool {
			var itemID int
			switch msg.Payload {
			case "1":
				itemID = 1
			case "2":
				itemID = 2
			case "3":
				itemID = 3
			default:
				return false
			}

			msg.User.SetSelectDate(itemID)
			msg.Trans(Node3)
			return true
		},
	}

	*primitives.GetNode(Node3) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			text := `
			В какой комлектации комлектации желаете начать строительство дома?
			
			1. Базовая комплектация (фундамент, домокомплект, кровля)
			2. Теплый контур (базовая комплектация + утепление, установка окон, дверей и т.п.)
			3. Под ключ (готовый дом для проживания со всеми инж. коммуникациями)
			`
			kb := &vkapi.Keyboard{
				ButtonsGrid: [][]vkapi.Button{
					{
						vkapi.MakeButton("Базовая комплектация", "1", vkapi.KbBlue),
					},
					{
						vkapi.MakeButton("Теплый контур", "2", vkapi.KbBlue),
					},
					{
						vkapi.MakeButton("Под ключ", "3", vkapi.KbBlue),
					},
				},
			}
			msg.Client.SendMessage(msg.PeerID, text, kb.String(), "")
		},
		Default: func(msg *primitives.Message) bool {
			var itemID int
			switch msg.Payload {
			case "1":
				itemID = 1
			case "2":
				itemID = 2
			case "3":
				itemID = 3
			default:
				return false
			}

			msg.User.SetSelectCollection(itemID)
			msg.Trans(Node4)
			return true
		},
	}

	*primitives.GetNode(Node4) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			text := `
			Какой тип клееного бруса Вы желаете?
			
			1. Клееный брус по ГОСТу — Сосна / Ель
			2. Клееный брус с утеплителем Белтермо
			3. Безусадочный CLT клееный брус с утеплителем Белтермо
			`
			kb := &vkapi.Keyboard{
				ButtonsGrid: [][]vkapi.Button{
					{
						vkapi.MakeButton("1", "1", vkapi.KbBlue),
						vkapi.MakeButton("2", "2", vkapi.KbBlue),
						vkapi.MakeButton("3", "3", vkapi.KbBlue),
					},
				},
			}
			msg.Client.SendMessage(msg.PeerID, text, kb.String(), "")
		},
		Default: func(msg *primitives.Message) bool {
			var itemID int
			switch msg.Payload {
			case "1":
				itemID = 1
			case "2":
				itemID = 2
			case "3":
				itemID = 3
			default:
				return false
			}

			msg.User.SetSelectType(itemID)
			msg.Trans(Node5)
			return true
		},
	}

	*primitives.GetNode(Node5) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			text := "Предполагаемая этажность будущего дома? (Выбор из готовых проектов)"
			kb := &vkapi.Keyboard{
				ButtonsGrid: [][]vkapi.Button{
					{
						vkapi.MakeButton("Один этаж", "1", vkapi.KbBlue),
					},
					{
						vkapi.MakeButton("Два этажа", "2", vkapi.KbBlue),
					},
					{
						vkapi.MakeButton("Свой проект", "3", vkapi.KbBlue),
					},
				},
			}
			msg.Client.SendMessage(msg.PeerID, text, kb.String(), "")
		},
		Default: func(msg *primitives.Message) bool {
			var itemID int
			switch msg.Payload {
			case "1":
				itemID = 1
			case "2":
				itemID = 2
			case "3":
				itemID = 3
			default:
				return false
			}

			msg.User.SetSelectStoreys(itemID)
			if itemID == 3 {
				msg.Trans(Node7)
			}

			msg.Trans(Node6)
			return true
		},
	}

	*primitives.GetNode(Node6) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			// если 1, отправляем одноэтажные кнопки + назад
			// если 2, отправляет двухэтажные кнопки + назад
			u := msg.User
			msg.Client.SendMessage(u.VKID, GetProjectsText(msg.User), GetProjectsKB(u).String(), GetProjectsPhotos(u))
		},
		Default: func(msg *primitives.Message) bool {
			// если нажата кнопка назад, то на ноду 5
			if msg.Payload == "back" {
				msg.Trans(Node5)
				return true
			}

			// если проект, то скидываем фотки и устанавливаем проект + (выбрать данный проект?)
			// todo: msg.Client.SetProject(itemID)
			itemID, err := strconv.Atoi(msg.Payload)
			if err != nil {
				return false
			}
			if _, ok := projectPhotos[itemID]; !ok {
				return false
			}

			msg.User.SetSelectProject(itemID)
			msg.Client.SendMessage(msg.User.VKID, "Используйте кнопку чтобы подтвердить выбор или вернуться назад", GetProjectsKB(msg.User).String(), projectPhotos[msg.User.SelectProject])
			return true
		},
	}

	*primitives.GetNode(Node7) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			text := `
			Укажите локацию (регион и населенный пункт) стройки. Отправьте информацию текстом.
			
			Доставка: грузовой транспорт с КМУ.
			От 50 рублей за км.`
			msg.Client.SendMessage(msg.PeerID, text, "", "")
		},
		Default: func(msg *primitives.Message) bool {
			msg.User.SetData(msg.Text)
			msg.Trans(Node8)
			return true
		},
	}

	*primitives.GetNode(Node8) = primitives.AnswerNode{
		Performance: func(msg *primitives.Message) {
			u := msg.User
			text := `Спасибо за заявку!
			Укажите свой номер телефона, если хотите, чтобы мы вам перезвонили.`
			msg.Client.SendMessage(msg.PeerID, text, "", "")

			tmp := `@id%d заявка:
			%s
			%s
			%s
			%s
			%s
			`
			project := ""
			if u.SelectProject != 0 {
				project = projectNames[u.SelectProject]
			}
			text = fmt.Sprintf(tmp, msg.User.VKID, datesMap[u.SelectDate], collectionsMap[u.SelectCollection], typesMap[u.SelectType], storeysMap[u.SelectStoreys], project)
			msg.Client.SendMessage(ADMIN_VKID, text, "", "")
		},
		Default: func(msg *primitives.Message) bool {
			msg.Client.SendMessage(ADMIN_VKID, fmt.Sprintf("@id%d: %s", msg.User.VKID, msg.Text), "", "")
			return true
		},
	}
}


func GetProjectsText(u *store.User) string {
	if u.SelectStoreys == 1 {
		return "Готовые одноэтажные проекты:"
	}

	return "Готовые двухэтажные проекты:"
}

func GetProjectsPhotos(u *store.User) string {
	if u.SelectStoreys == 1 {
		return "photo-216949929_457239112,photo-216949929_457239109,photo-216949929_457239118,photo-216949929_457239123,photo-216949929_457239126,photo-216949929_457239130"
	}

	return "photo-216949929_457239133,photo-216949929_457239136,photo-216949929_457239139,photo-216949929_457239145,photo-216949929_457239148,photo-216949929_457239152,photo-216949929_457239154,photo-216949929_457239157"
}

func GetProjectsKB(u *store.User) *vkapi.Keyboard {
	var kb *vkapi.Keyboard
	if u.SelectStoreys == 1 {
		kb = &vkapi.Keyboard{
			ButtonsGrid: [][]vkapi.Button{
				{
					vkapi.MakeButton("Аргус", "1", vkapi.KbBlue),
					vkapi.MakeButton("Бегония", "2", vkapi.KbBlue),
					vkapi.MakeButton("Джек", "3", vkapi.KbBlue),
				},
				{
					vkapi.MakeButton("Колибри", "4", vkapi.KbBlue),
					vkapi.MakeButton("Сапсан", "5", vkapi.KbBlue),
					vkapi.MakeButton("Сойка", "6", vkapi.KbBlue),
				},
				{
					vkapi.MakeButton("Назад", "back", vkapi.KbRed),
				},
			},
		}
	} else {
		kb = &vkapi.Keyboard{
			ButtonsGrid: [][]vkapi.Button{
				{
					vkapi.MakeButton("АГАТА", "7", vkapi.KbBlue),
					vkapi.MakeButton("АЗАЛИЯ ЛЮКС", "8", vkapi.KbBlue),
				},
				{
					vkapi.MakeButton("АКАЦИЯ", "9", vkapi.KbBlue),
					vkapi.MakeButton("БАРБАРИС", "10", vkapi.KbBlue),
				},
				{
					vkapi.MakeButton("БЕКАС", "11", vkapi.KbBlue),
					vkapi.MakeButton("ГАРМОНИЯ МИНИ", "12", vkapi.KbBlue),

				},
				{
					vkapi.MakeButton("ИВОЛГА", "13", vkapi.KbBlue),
					vkapi.MakeButton("ЛОРИ", "14", vkapi.KbBlue),
				},
				{
					vkapi.MakeButton("Назад", "back", vkapi.KbRed),
				},
			},
		}
	}

	if u.SelectProject != 0 {
		kb.ButtonsGrid[len(kb.ButtonsGrid) - 1] = append(kb.ButtonsGrid[len(kb.ButtonsGrid)-1], vkapi.MakeButton("Выбрать данный проект", "select", vkapi.KbGreen))
	}

	return kb
}