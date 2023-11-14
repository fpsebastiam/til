// Copyright 2023 Alexey Kutepov <reximkut@gmail.com>

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"log"
	"net"
	"time"
	"fmt"
	"unicode/utf8"
)

const (
	Port = "8080"
	SafeMode = true
	MessageRate = 1.0
	BanLimit = 10 * 60.0
	StrikeLimit = 10
)

func sensitive(message string) string {
	if SafeMode {
		return "[REDACTED]"
	} else {
		return message
	}
}

type MessageType int
const (
	ClientConnected MessageType = iota + 1
	ClientDisconnected
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

type Client struct {
	Conn net.Conn
	LastMessage time.Time
	StrikeCount int
}

func server(messages chan Message) {
	clients := map[string]*Client{}
	bannedClients := map[string]time.Time{}
	for {
		msg := <- messages
		authorAddr := msg.Conn.RemoteAddr().(*net.TCPAddr)
		switch msg.Type {
		case ClientConnected:
			bannedAt, banned := bannedClients[authorAddr.IP.String()]
			now := time.Now()
			if banned {
				if now.Sub(bannedAt).Seconds() >= BanLimit {
					delete(bannedClients, authorAddr.IP.String())
					banned = false
				}
			}
			if !banned {
				log.Printf("Client %s connected", sensitive(authorAddr.String()))
				clients[msg.Conn.RemoteAddr().String()] = &Client{
					Conn: msg.Conn,
					LastMessage: time.Now(),
				}
			} else {
				msg.Conn.Write([]byte(fmt.Sprintf("You are banned: %f secs left\n", BanLimit - now.Sub(bannedAt).Seconds())))
				msg.Conn.Close()
			}
		case ClientDisconnected:
			log.Printf("Client %s disconnected", sensitive(authorAddr.String()))
		case NewMessage:
			author := clients[authorAddr.String()]
			now := time.Now()
			if author != nil {
				if now.Sub(author.LastMessage).Seconds() >= MessageRate {
					if utf8.ValidString(msg.Text) {
						author.LastMessage = now
						author.StrikeCount = 0
						log.Printf("client %s sent message %s", sensitive(authorAddr.String()), msg.Text);
						for _, client := range clients {
							if client.Conn.RemoteAddr().String() != authorAddr.String() {
								client.Conn.Write([]byte(msg.Text))
							}
						}
					} else {
						author.StrikeCount += 1
						if author.StrikeCount >= StrikeLimit{
							bannedClients[authorAddr.IP.String()] = now
							author.Conn.Write([]byte("You are banned\n"))
							author.Conn.Close()
						}
					}
				} else {
					author.StrikeCount += 1
					if author.StrikeCount >= StrikeLimit {
						bannedClients[authorAddr.IP.String()] = now
						author.Conn.Write([]byte("You are banned\n"))
						author.Conn.Close()
					}
				}
			} else {
				msg.Conn.Close()
			}
		}
	}
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 64)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: ClientDisconnected,
				Conn: conn,
			}
			return
		}
		text := string(buffer[0:n])
		messages <- Message{
			Type: NewMessage,
			Text: text,
			Conn: conn,
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s: %s\n", Port, sensitive(err.Error()))
	}
	log.Printf("Listening to TCP connections on port %s ...", Port)

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept a connection: %s\n", sensitive(err.Error()))
			continue
		}
		messages <- Message{Type: ClientConnected, Conn: conn}
		go client(conn, messages)
	}
}
