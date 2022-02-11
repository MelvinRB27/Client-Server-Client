package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MelvinRB27/Client-Server/models"
	"github.com/MelvinRB27/Client-Server/utils"
	"github.com/gorilla/websocket"
)

var upGraderWebSocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// para revisar si la concexion se puede realizar o no
func checkOrigin(r *http.Request) bool {
	log.Printf(" %s %s %s %v ", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

type MessageChannel chan *models.Message
type UserChannel chan *UserChat

//crear la comunicacion del usario y la app central del channel
type Channel struct {
	messageChannel MessageChannel
	leaveChannel   UserChannel
}

//lógica del chat
type websocketChat struct {
	users       map[string]*UserChat
	joinChannel UserChannel
	channel     *Channel
}

func NewWebSocketChat() *websocketChat {
	return &websocketChat{
		users:       make(map[string]*UserChat),
		joinChannel: make(UserChannel),
		channel: &Channel{
			messageChannel: make(MessageChannel),
			leaveChannel:   make(UserChannel),
		},
	}
}

//handler para la conexiones de los usuarios
func (w *websocketChat) HandlerUserConnection(rw http.ResponseWriter, r *http.Request) {
	connection, err := upGraderWebSocket.Upgrade(rw, r, nil)
	if err != nil {
		log.Panicln("No se pudo conectar", r.Host)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	//para saber la llave de la cabecera de la URL con
	keys := r.URL.Query()
	username := strings.TrimSpace(keys.Get("username"))
	if username == "" {
		username = fmt.Sprintf("user-%d", utils.GenRandomID())
	}

	//creando nuevo usuario
	u := newUserChat(w.channel, username, connection)
	w.joinChannel <- u // un nuevo usuario se acaba de conectar
	u.OnlineList() // poner al usuario en escucha si tiene nuevo msj
}

//proceso de leer los channels
func (w *websocketChat) UserChatManager() {
	for {
		select {
		case userChat := <-w.joinChannel:
			w.addUser(userChat)
		case message := <-w.channel.messageChannel:
			w.SendMessage(message)
		case userchat := <-w.channel.leaveChannel:
			w.LeaveChat(userchat.Username)
		}
	}
}

//creando usuario
func (w *websocketChat) addUser(userchat *UserChat) {
	if user, ok := w.users[userchat.Username]; ok {
		user.Connection = userchat.Connection //actualizando su nueva conexion
		log.Printf("Reconexion usuario: %s \n", userchat.Username)
	} else {
		w.users[userchat.Username] = userchat
		log.Printf("Nuevo usuario: %s \n", userchat.Username)
	}
}

func (w *websocketChat) SendMessage(message *models.Message) {
	if user, ok := w.users[message.TargetUsername]; ok {
		if err := user.SendMessageToClient(message); err != nil {
			log.Printf("No se pudo conectar con el cliente %s, Error: %v", message.TargetUsername, err)
		}
	}
}

func (w *websocketChat) LeaveChat(username string) {
	if user, ok := w.users[username]; ok {
		defer user.Connection.Close()
		delete(w.users, username)
		log.Printf("User: %s, dejó el chat\n", username)
	}
}