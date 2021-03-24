package utils

import (
	"errors"
	"github.com/jordy2254/indoormaprestapi/pkg/model"
)

type PairPoint2f struct {
	first  model.Point2f
	second model.Point2f
}

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

func CalculatePolygonPoints(room model.Room) []model.Point2f {
	edgeData := CalculatePolygonEdgePairs(room)

	var firstPair *PairPoint2f
	var lastPair *PairPoint2f
	var panic int

	var pointers []*model.Point2f

	for panic < len(edgeData)*2 {
		panic++
		if firstPair == nil {
			firstPair = edgeData[0]
			lastPair = firstPair
			pointers = append(pointers, &firstPair.first)
			continue
		}

		if model.PointsEqual(lastPair.second, firstPair.first) {
			break
		}

		var current *PairPoint2f

		for _, tmp := range edgeData {
			if model.PointsEqual(tmp.first, lastPair.second) {
				current = tmp
				break
			}
		}

		if current == nil {
			break
		}

		pointers = append(pointers, &current.first)
		lastPair = current

	}
	var returnVal []model.Point2f
	for _, value := range pointers {
		returnVal = append(returnVal, *value)
	}
	return returnVal
}

func CalculatePolygonEdgePairs(room model.Room) []*PairPoint2f {
	var roomEdges []*PairPoint2f
	var indentEdges []*PairPoint2f

	roomEdges = append(roomEdges, calculateRectangleEdgePairs(0, 0, *room.Dimensions.X, *room.Dimensions.Y)...)

	for _, indent := range room.Indents {
		_, startPoints := calculateStartPointsOfIndent(room, indent)

		indentEdges = append(indentEdges, calculateRectangleEdgePairs(*startPoints.X, *startPoints.Y, *indent.Dimensions.X, *indent.Dimensions.Y)...)
	}

	var toRemove []*PairPoint2f

	for _, indentWall := range indentEdges {
		var found *PairPoint2f
		for _, roomWall := range roomEdges {
			if linesIntersect(roomWall.first, roomWall.second, indentWall.first, indentWall.second) {
				found = roomWall
				break
			}
		}

		if found != nil {
			tmp := (*found).second
			(*found).second = (*indentWall).first

			if !model.PointsEqual((*indentWall).second, tmp) {
				roomEdges = append(roomEdges, &PairPoint2f{indentWall.second, tmp})
			}

			if model.PointsEqual(found.first, found.second) {
				roomEdges = removePoint(roomEdges, found)
			}

			toRemove = append(toRemove, indentWall)
		}
	}

	for _, remove := range toRemove {
		indentEdges = removePoint(indentEdges, remove)
		roomEdges = removePoint(roomEdges, remove)
	}

	for _, tmpWall := range indentEdges {
		allWalls := append(roomEdges, indentEdges...)
		var firstEdges []model.Point2f
		var secondEdges []model.Point2f

		for _, wall := range allWalls {
			firstEdges = append(firstEdges, wall.first)
			secondEdges = append(secondEdges, wall.second)
		}

		var (
			firstCount1  int
			secondCount1 int

			firstCount2  int
			secondCount2 int
		)

		for _, edge := range firstEdges {
			if model.PointsEqual(edge, tmpWall.first) {
				firstCount1++
			}
			if model.PointsEqual(edge, tmpWall.second) {
				firstCount2++
			}
		}

		for _, edge := range secondEdges {
			if model.PointsEqual(edge, tmpWall.first) {
				secondCount1++
			}
			if model.PointsEqual(edge, tmpWall.second) {
				secondCount2++
			}
		}

		if firstCount1 != secondCount1 || firstCount2 != secondCount2 {
			tmp := tmpWall.first
			tmpWall.first = tmpWall.second
			tmpWall.second = tmp
		}
	}

	return append(roomEdges, indentEdges...)
}

func calculateStartPointsOfIndent(room model.Room, indent model.Indent) (error, *model.Point2f) {
	if indent.WallKeyA != "" && indent.WallKeyB != "" {
		var xStart float64 = 0
		var yStart float64 = 0

		if indent.WallKeyA == "BOTTOM" {
			yStart = *room.Dimensions.Y - *indent.Dimensions.Y
		}

		if indent.WallKeyB == "RIGHT" {
			xStart = *room.Dimensions.X - *indent.Dimensions.X
		}

		return nil, &model.Point2f{X: &xStart, Y: &yStart}
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
			xStart = indent.Location
			yStart = *room.Dimensions.X - *indent.Dimensions.X
			break
		default:
			return errors.New("No indent location found for" + indent.WallKeyA), nil
		}

		return nil, &model.Point2f{X: &xStart, Y: &yStart}
	}

	return errors.New("No indent location found"), nil
}

func calculateRectangleEdgePairs(x, y, width, height float64) []*PairPoint2f {
	var points []*PairPoint2f = make([]*PairPoint2f, 4)

	xPWidth := x + width
	yPHeight := y + height

	tl := model.Point2f{X: &x, Y: &y}
	tr := model.Point2f{X: &xPWidth, Y: &y}
	bl := model.Point2f{X: &x, Y: &yPHeight}
	br := model.Point2f{X: &xPWidth, Y: &yPHeight}

	points[0] = &PairPoint2f{first: tl, second: tr}
	points[1] = &PairPoint2f{first: tr, second: br}
	points[2] = &PairPoint2f{first: br, second: bl}
	points[3] = &PairPoint2f{first: bl, second: tl}

	return points
}

func linesIntersect(p1, p2, p3, p4 model.Point2f) bool {
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
