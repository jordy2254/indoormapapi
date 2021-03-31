package model

import (
	"errors"
)

func removePoint(array []*PairPoint2f, point *PairPoint2f) []*PairPoint2f {
	index := -1

	for i, value := range array {
		if value == point {
			index = i
			break
		}
	}
	if index != -1 {
		var returnVal []*PairPoint2f

		if index > 0 {
			returnVal = append(returnVal, array[0:index]...)
		}

		if index <= len(array) {
			returnVal = append(returnVal, array[index+1:len(array)]...)
		}

		return returnVal
	}

	return array
}

func CalculatePolygonPoints(room Room) []Point2f {
	edgeData := CalculatePolygonEdgePairs(room, true)

	var firstPoint *Point2f
	var lastPoint *Point2f

	var panic int

	var pointers []*Point2f
	var indexes = make([]bool, len(edgeData))

	for i := range indexes {
		indexes[i] = false
	}
	for panic < len(edgeData)*2 {
		panic++
		if firstPoint == nil {
			firstPoint = &edgeData[0].First
			lastPoint =  &edgeData[0].Second
			indexes[0] = true
			pointers = append(pointers, firstPoint)
			pointers = append(pointers, lastPoint)
			continue
		}

		if PointsEqual(*firstPoint, *lastPoint) {
			break
		}

		var current *Point2f

		for i, tmp := range edgeData {
			if indexes[i]{
				continue
			}
			if PointsEqual(*lastPoint, tmp.First) {
				current = &tmp.Second
				indexes[i] = true
				break
			}

			if PointsEqual(*lastPoint, tmp.Second) {
				current = &tmp.First
				indexes[i] = true
				break
			}
		}

		if current == nil {
			break
		}

		pointers = append(pointers, current)
		lastPoint = current
	}

	var returnVal []Point2f
	for _, value := range pointers {
		returnVal = append(returnVal, *value)
	}
	return returnVal
}

func CalculatePolygonEdgePairs(room Room, excludeEntrances bool) []*PairPoint2f {
	var roomEdges []*PairPoint2f
	var indentEdges []*PairPoint2f
	var entranceEdges []*PairPoint2f

	roomEdges = append(roomEdges, calculateRectangleEdgePairs(0, 0, *room.Dimensions.X, *room.Dimensions.Y)...)

	for _, indent := range room.Indents {
		_, startPoints := calculateStartPointsOfIndent(room, indent)

		indentEdges = append(indentEdges, calculateRectangleEdgePairs(*startPoints.X, *startPoints.Y, *indent.Dimensions.X, *indent.Dimensions.Y)...)
	}

	if !excludeEntrances {
		entranceEdges = createEntranceEdges(room, entranceEdges)
	}

	allWalls := cutoutAndMergeEdges(roomEdges, indentEdges)

	if excludeEntrances {
		return allWalls
	}

	return cutoutAndMergeEdges(allWalls, entranceEdges)
}

func createEntranceEdges(room Room, entranceEdges []*PairPoint2f) []*PairPoint2f {
	for _, entrance := range room.Entrances {
		var (
			x1 float64
			y1 float64
			x2 float64
			y2 float64
		)

		if entrance.WallKey == "LEFT" || entrance.WallKey == "RIGHT" {
			y1 = entrance.Location
			y2 = entrance.Location + entrance.Length
			if entrance.WallKey == "LEFT" {
				x1 = 0
				x2 = 0
			}

			if entrance.WallKey == "RIGHT" {
				x1 = *room.Dimensions.X
				x2 = *room.Dimensions.X
			}
		}

		if entrance.WallKey == "TOP" || entrance.WallKey == "BOTTOM" {
			x1 = entrance.Location
			x2 = entrance.Location + entrance.Length
			if entrance.WallKey == "TOP" {
				y1 = 0
				y2 = 0
			}
			if entrance.WallKey == "BOTTOM" {
				y1 = *room.Dimensions.Y
				y2 = *room.Dimensions.Y
			}
		}

		entranceEdges = append(entranceEdges, &PairPoint2f{
			First: Point2f{
				X: &x1,
				Y: &y1,
			},
			Second: Point2f{
				X: &x2,
				Y: &y2,
			},
		})
	}
	return entranceEdges
}

func cutoutAndMergeEdges(mainEdges []*PairPoint2f, toCut []*PairPoint2f) []*PairPoint2f {
	var toRemove []*PairPoint2f

	for _, cutWall := range toCut {
		var found *PairPoint2f
		for _, mainWall := range mainEdges {
			if linesIntersect(mainWall.First, mainWall.Second, cutWall.First, cutWall.Second) {
				found = mainWall
				break
			}
		}

		if found != nil {
			tmp := (*found).Second
			(*found).Second = (*cutWall).First

			if !PointsEqual((*cutWall).Second, tmp) {
				mainEdges = append(mainEdges, &PairPoint2f{cutWall.Second, tmp})
			}

			if PointsEqual(found.First, found.Second) {
				mainEdges = removePoint(mainEdges, found)
			}

			toRemove = append(toRemove, cutWall)
		}
	}

	for _, remove := range toRemove {
		toCut = removePoint(toCut, remove)
		mainEdges = removePoint(mainEdges, remove)
	}

	return append(mainEdges, toCut...)
}

func calculateStartPointsOfIndent(room Room, indent Indent) (error, *Point2f) {
	if indent.WallKeyA != "" && indent.WallKeyB != "" {
		var xStart float64 = 0
		var yStart float64 = 0

		if indent.WallKeyA == "BOTTOM" {
			yStart = *room.Dimensions.Y - *indent.Dimensions.Y
		}

		if indent.WallKeyB == "RIGHT" {
			xStart = *room.Dimensions.X - *indent.Dimensions.X
		}

		return nil, &Point2f{X: &xStart, Y: &yStart}
	} else if indent.WallKeyA != "" {
		var xStart float64 = 0
		var yStart float64 = 0

		switch indent.WallKeyA {
		case "TOP":
			xStart = indent.Location
			yStart = 0
			break
		case "BOTTOM":
			xStart = indent.Location
			yStart = *room.Dimensions.Y - *indent.Dimensions.Y
			break
		case "LEFT":
			xStart = 0
			yStart = indent.Location
			break
		case "RIGHT":
			xStart = *room.Dimensions.X - *indent.Dimensions.X
			yStart = indent.Location
			break
		default:
			return errors.New("No indent location found for" + indent.WallKeyA), nil
		}

		return nil, &Point2f{X: &xStart, Y: &yStart}
	}

	return errors.New("No indent location found"), nil
}

func calculateRectangleEdgePairs(x, y, width, height float64) []*PairPoint2f {
	var points []*PairPoint2f = make([]*PairPoint2f, 4)

	xPWidth := x + width
	yPHeight := y + height

	tl := Point2f{X: &x, Y: &y}
	tr := Point2f{X: &xPWidth, Y: &y}
	bl := Point2f{X: &x, Y: &yPHeight}
	br := Point2f{X: &xPWidth, Y: &yPHeight}

	points[0] = &PairPoint2f{First: tl, Second: tr}
	points[1] = &PairPoint2f{First: tr, Second: br}
	points[2] = &PairPoint2f{First: br, Second: bl}
	points[3] = &PairPoint2f{First: bl, Second: tl}

	return points
}

func linesIntersect(p1, p2, p3, p4 Point2f) bool {
	//vertical line
	if *p1.X == *p3.X && *p2.X == *p3.X && *p4.X == *p1.X {
		//posative increase
		if *p1.Y-*p2.Y < 0 {
			return *p3.Y >= *p1.Y && *p3.Y <= *p2.Y
		}

		//negative increase
		if *p1.Y-*p2.Y > 0 {
			return *p3.Y <= *p1.Y && *p3.Y >= *p2.Y
		}
	}

	//Horazontal line
	if *p1.Y == *p3.Y && *p2.Y == *p3.Y && *p4.Y == *p1.Y {
		//posative increase
		if *p1.X-*p2.X < 0 {
			return *p3.X >= *p1.X && *p3.X <= *p2.X
		}

		//negative increase
		if *p1.X-*p2.X > 0 {
			return *p3.X <= *p1.X && *p3.X >= *p2.X
		}
	}
	return false
}
