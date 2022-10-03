package query

import (
	"github.com/sandertv/gophertunnel/query"
	"strconv"
	"strings"
)

type Request struct {
	Ip   string
	Port uint16
}

type Response struct {
	Serversoftware string
	Plugins        string
	Version        string
	Whitelist      string
	Players        []string
	PlayerCount    string
	MaxPlayers     string
	GameName       string
	GameMode       string
	MapName        string
	HostName       string
	HostIp         string
	HostPort       string
}

func Query(request Request) (Response, error) {
	data, err := query.Do(request.Ip + ":" + strconv.Itoa(int(request.Port)))
	return Response{
		Players:        strings.Split(data["players"], ", "),
		Serversoftware: data["server_engine"],
		Plugins:        data["plugins"],
		Whitelist:      data["whitelist"],
		Version:        data["version"],
		PlayerCount:    data["numplayers"],
		MaxPlayers:     data["maxplayers"],
		MapName:        data["map"],
		HostPort:       data["hostport"],
		HostName:       data["hostname"],
		HostIp:         data["hostip"],
		GameMode:       data["gametype"],
		GameName:       data["game_id"],
	}, err
}
