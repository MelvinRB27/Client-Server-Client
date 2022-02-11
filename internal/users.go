package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MelvinRB27/Client-Server/models"
	"github.com/MelvinRB27/Client-Server/utils"
	"github.com/gorilla/websocket"
	"tawesoft.co.uk/go/dialog"
)

type UserChat struct {
	channels   *Channel
	Username   string
	Connection *websocket.Conn
}

func newUserChat(channels *Channel, Username string, Conn *websocket.Conn) *UserChat {
	return &UserChat{
		channels:   channels,
		Username:   Username,
		Connection: Conn,
	}
}

//constantes escucha de msj entrantes por la connection
func (u *UserChat) OnlineList() {
	for {
		if _, message, err := u.Connection.ReadMessage(); err != nil {
			log.Println("Error al leer el mensaje:", err.Error())
			break // si hay un error romper la conexion
		} else {
			msj := &models.Message{}
			fmt.Println("Data: ", string(message))
			if err := json.Unmarshal(message, msj); err != nil {
				log.Printf("No se pudo leer el mensaje: user %s\n, err: %s", u.Username, err.Error())
			} else {
				log.Println(msj)
				u.channels.messageChannel <- msj
			}
		}
	}
	u.channels.leaveChannel <- u
}

func (u *UserChat) SendMessageToClient(message *models.Message) error {
	message.ID = utils.GenRandomID()
	if data, err := json.Marshal(message); err != nil {
		return err
	} else {
		err := u.Connection.WriteMessage(websocket.TextMessage, data)
		log.Printf("Sent message: from %s to %s", message.SenderUsername, message.TargetUsername)
		dialog.Alert("Enviado correctamente")
		return err
	}
}