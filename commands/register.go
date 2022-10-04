package commands

import "github.com/jviguy/speedycmds/command"

func RegisterAll(m *command.Map) {
	m.RegisterCommands([]command.Command{
		Ban{},
		Unban{},
		Mute{},
		Unmute{},
		Kick{},
	}, true)
}
