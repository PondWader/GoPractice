package protocol

type VersionData struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PlayersData struct {
	Max    int   `json:"max"`
	Online int   `json:"online"`
	Sample []any `json:"sample"`
}

type ServerDescriptionData struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Version     VersionData           `json:"version"`
	Players     PlayersData           `json:"players"`
	Description ServerDescriptionData `json:"description"`
	Favicon     string                `json:"favicon"`
}

type StatusResponsePacket struct {
	Resp StatusResponse `type:"JSON"`
}

func (client *ProtocolClient) status() {
	client.state = "status"

	packetId, data, err := client.readPacket()
	if err != nil {
		return
	}

	if packetId == 0 {
		data = Serialize(&StatusResponsePacket{
			Resp: StatusResponse{
				Version: VersionData{
					Name:     "GoPractice 1.8.9",
					Protocol: 47,
				},
				Players: PlayersData{
					Max:    1,
					Online: client.server.GetPlayerCount(),
				},
				Description: ServerDescriptionData{
					Text: client.server.GetConfig().Motd,
				},
				Favicon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAMAAACdt4HsAAAB0VBMVEVHcExRprFKl6EtXGIYISEAAAAAAAAfMzUHCgoJEhJev8xTqbQAAAAlTFElT1RdvclJlZ8MGhxQo64cOT0AAABGkJkgQ0hNnqkAAAAmT1Q/gIhauMRbTjxYs79Wsb0vYGY6eIA2b3ZANiqmjmwsXGJVrblGg4lFj5lXsb1cu8c8fISRfGA3cHcyaG9CiZK3nHjZuY4+MiVDipNIk51qWEPLrYV5Z08yVFakjGs0a3KfiGlKPzBwaVdoWUQ8e4OulHJ+blZZfXpZqbKOeV3F////2Odq1+X///9q2OZr2edkzNlp1+VhxdJm0N5q2ehn0uBjyNVs2+llzdtgws9nvMdjxtJXrrln1OJn1uTKsYy8y810srhjwMzv7+/xzZ761aT+/f60xsmCtbsNDg6rwMMBAQDV2tp1tb1RpbD5+Pja39/o6elmtb59oZeLuL2gnpvFqIFcusf00KDsyZuhvcH19PSKsbb8+/rh4+PiwZWZu7+1tbWsq6sVLTFMm6WRppe8rJBqqaquqIzgvpKStruBvMNpuMLFxcVCQkLNzMwbOTwvLy9ISEjM0tI6d4BlaGHFuax9fX2JiYkePUFun6VXk5uzv8E0a3JTSjsAlIVHcExNb34bAAAAm3RSTlMA4MxfJQUMLRcb/OQCSFb5wx/bNQm+P9cPUKr2VvDucpN4Pcho6tv96/mbxX15xN/8MbTYau58j7yG7Em4dKX6u+X2pHn/////////////////////////////////////////////////////////////////////////////////////////////////////////////////ANfG23gAAAPUSURBVFjD5ZfnW9tIEIfdsCxjHMfYYAwEQu+d9HZ3z7PsSpYl94p7g9ANxHTSSC+X3CX3554skzuT4F0Tf7ubT9p9NO/u/HZ2NJLJ/ttWp9bQtMZi0v2Ms5pS3rSOTRuNzQqz9aZeVXchbx111TDgcvoARAhB4HNq5fcaTNX5msSlqKFeJ2u3QyAZ9DM2l5NrvUepekaUGrx7Y6d5qqPH6BROslk3IxGgJ7QZDG6GnNx0n8C5FI04/wY5Z+cPvVw2XpifX9x0s6J/eHd+TrSH2wzrzbmgz4zZg0YB1nO5dV+iUHSZiwXdEDLPHkmDuYchP/v+wOvQ6isDWrTM62S/3Z2JLUg+j1Y4uFTIvCiNUh4kfO53MKOVAUqbcPCa8ScKH96+knyCHn/i1ZvkqjTIuBHKrTu0SoyEvY7coZ3fXU0nf5d83oV9K/vpN59i3wAuLz+G0UA32uq1Qd/KXvLlvgSICzD77u3LF6WBBwKu30phc1ffxwO4lNlblZTfSvgdnvh8IVYS0SeeyTQhD+rM4kvM9lZJw1T48Ov7P4NFDWOLiWJaoF4KD6CNqJg625mtha3F3XB/Mp0+8CZSwfizE66YVkjbggdQV1Ax+Tj38+fZEyH0IZlO/rVsYzwCV8priDvEojUNotKLrJ+1nwQia/fX8vkjJzq9FgDwQ3iAXou+vQrY0E50I5/fiAbcZZNT+HuttP2zlijF0fHH1Y/Hx0fOskmFpWoAQOtfYnsLcS9fPtesqR4AhVRsfyEY/jcAESC/jNfAVQYA/lDqj9Q2Uz4F61VYQMOl8uUAb+NdTPkEQARAY/sZAIAAwrMAQgiX688CfjAkp7EATTOsDaA2gNoAJjMRgNdANsUSANhTMHXN9BEA9l/vzFQMQvPb7NNlP2EHS4Gd7kpF5W4gGiEB7EuBaOBuBcDt+xPdpBDQ4J3ZW12VCmJbG008Ba2emsT0C2oFIZGgS0nIREIqQ09PbXcBCh34qjxAAjBXL3SdzwGM4Mt6a407+K4iXVwDMsDWUxuA9HEkaoBam6r4uP50PRATSU4AAIO6tprIdRI63U4OH8ElQoMha8EfA1QQeqRis4rNglFit67UYrbAKmjyj4bVWbGmwN4WGdnoMeb8PUDWe72an5/J4eUwPGcTiHv8ZLyN7N81Hsk/eWwD6ExbgSB35ZdAdGN8kijBjYlIdOLBdUO7wANUMuhjBps7myZvRJ6uXSMBLMNrG7MzapmlsWPIYKwfaB+oNyqsIw20GL2qe+cWMQZd13D37dNHtUZFUZSKtpxKp1Nda6vi789ikf3P7G8f4AvstOFwrgAAAABJRU5ErkJggg==",
			},
		})
		if err := client.WritePacket(0x00, data); err != nil {
			return
		}

		packetId, data, err = client.readPacket()
		if err != nil {
			return
		}
		if packetId != 1 {
			client.Disconnect("Bad packet ID")
			return
		}

		client.handleStatusPing(data)
	} else if packetId == 1 {
		client.handleStatusPing(data)
	} else {
		client.Disconnect("Bad packet ID")
	}
}

func (client *ProtocolClient) handleStatusPing(data []byte) {
	client.WritePacket(1, data)
	client.Disconnect("Status request completed")
}
