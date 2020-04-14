package web

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	spc "github.com/doodles526/syncplayBot/client"
	"github.com/gorilla/mux"
)

type Args struct {
	ListenAddr         *net.TCPAddr
	SyncplayClientArgs *spc.Args
}

func ServeForBot(args *Args) error {
	c, err := spc.NewClient(args.SyncplayClientArgs)
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/chat", chatMessageHandlerFactory(c)).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    args.ListenAddr.String(),
	}
	return srv.ListenAndServe()
}

func chatMessageHandlerFactory(c *spc.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var chatReq ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Could not decode chat message")
		}

		if err := c.SendMessage(spc.ClientChatMsg(chatReq.Message)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Error sending message to syncplay: "+err.Error())
		}
	}
}
