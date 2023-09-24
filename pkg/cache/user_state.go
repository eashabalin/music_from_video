package cache

type UserState int

const (
	WaitingForAudio UserState = iota
	WaitingForVideo
)

func (us UserState) String() string {
	switch us {
	case WaitingForAudio:
		return "Waiting for audio"
	case WaitingForVideo:
		return "Waiting for video"
	}
	return "Unknown"
}

func NewUserStateFromString(s string) *UserState {
	var state UserState
	if s == "Waiting for audio" {
		state = WaitingForAudio
		return &state
	} else if s == "Waiting for video" {
		state = WaitingForVideo
		return &state
	}
	return nil
}
