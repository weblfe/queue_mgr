package utils

import "math"

const (
	radius = 6378.137
	rad    = math.Pi / 180.0
)

type Point struct {
	Lat float64
	Lng float64
}

func MakePoint(lat, lng float64) Point {
	return Point{
		Lat: lat, Lng: lng,
	}
}

func NewPoint(lat, lng float64) *Point {
	return &Point{
		Lat: lat, Lng: lng,
	}
}

func (point Point) RadLat() float64 {
	return point.Lat * rad
}

func (point Point) RadLng() float64 {
	return point.Lng * rad
}

// GetDistance  计算2金纬度之间距离
func (point Point) GetDistance(p Point) float64 {
	var (
		lat1, lat2, theta, dist float64
	)
	lat2 = p.RadLat()
	lat1 = point.RadLat()
	theta = p.RadLat() - point.RadLng()
	dist = math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}
