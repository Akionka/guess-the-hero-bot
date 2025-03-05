package data

type Position string

const (
	PositionCarry       Position = "Carry"
	PositionMid         Position = "Mid"
	PositionOfflane     Position = "Offlane"
	PositionSoftSupport Position = "Soft Support"
	PositionHardSupport Position = "Hard Support"
	PositionUnknown     Position = "Unknown"
)

func (p Position) ToEmoji() string {
	switch p {
	case PositionCarry:
		return "🗡"
	case PositionMid:
		return "🏹"
	case PositionOfflane:
		return "🛡"
	case PositionSoftSupport:
		return "🔮"
	case PositionHardSupport:
		return "✨"
	default:
		return "❌"
	}
}
func (p Position) String() string {
	switch p {
	case PositionCarry:
		return "Carry"
	case PositionMid:
		return "Mid"
	case PositionOfflane:
		return "Offlane"
	case PositionSoftSupport:
		return "Soft Support"
	case PositionHardSupport:
		return "Hard Support"
	default:
		return "Unknown"
	}
}