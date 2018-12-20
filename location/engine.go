package location

const (
	LocationReqFailed = iota
)

type Engine interface {
	GetLocation() (Point2D, error)
}
