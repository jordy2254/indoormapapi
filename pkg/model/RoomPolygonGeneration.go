package model

import (
	"errors"
	"fmt"
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

	var firstPair *PairPoint2f
	var lastPair *PairPoint2f
	var panic int

	var pointers []*Point2f

	for panic < len(edgeData)*2 {
		panic++
		if firstPair == nil {
			firstPair = edgeData[0]
			lastPair = firstPair
			pointers = append(pointers, &firstPair.First)
			continue
		}

		if PointsEqual(lastPair.Second, firstPair.First) {
			break
		}

		var current *PairPoint2f

		for _, tmp := range edgeData {
			if PointsEqual(tmp.First, lastPair.Second) {
				current = tmp
				break
			}
		}

		if current == nil {
			break
		}

		pointers = append(pointers, &current.First)
		lastPair = current

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
		for _, entrance := range room.Entrances {
			var(
				x1 float64
				y1 float64
				x2 float64
				y2 float64
			)

			if entrance.WallKey == "LEFT" || entrance.WallKey == "RIGHT" {
				y1 = entrance.Location
				y2 = entrance.Location + entrance.Length
				if entrance.WallKey == "LEFT"{
					x1 = 0
					x2 = 0
				}

				if entrance.WallKey == "RIGHT"{
					x1 = *room.Dimensions.X
					x2 = *room.Dimensions.X
				}
			}

			if entrance.WallKey == "TOP" || entrance.WallKey == "BOTTOM" {
				x1 = entrance.Location
				x2 = entrance.Location + entrance.Length
				if entrance.WallKey == "TOP"{
					y1 = 0
					y2 = 0
				}
				if entrance.WallKey == "BOTTOM"{
					y1 = *room.Dimensions.Y
					y2 = *room.Dimensions.Y
				}
			}


			entranceEdges = append(entranceEdges, &PairPoint2f{
				First:  Point2f{
					X: &x1,
					Y: &y1,
				},
				Second: Point2f{
					X: &x2,
					Y: &y2,
				},
			})
		}
	}

	var toRemove []*PairPoint2f

	for _, indentWall := range indentEdges {
		var found *PairPoint2f
		for _, roomWall := range roomEdges {
			if linesIntersect(roomWall.First, roomWall.Second, indentWall.First, indentWall.Second) {
				found = roomWall
				break
			}
		}

		if found != nil {
			tmp := (*found).Second
			(*found).Second = (*indentWall).First

			if !PointsEqual((*indentWall).Second, tmp) {
				fmt.Printf("Adding wall: %d:%d, %d:%d\n", *indentWall.Second.X, *indentWall.Second.Y, *tmp.X, *tmp.Y)
				roomEdges = append(roomEdges, &PairPoint2f{indentWall.Second, tmp})
			}

			if PointsEqual(found.First, found.Second) {
				roomEdges = removePoint(roomEdges, found)
			}

			toRemove = append(toRemove, indentWall)
		}
	}

	for _, remove := range toRemove {
		indentEdges = removePoint(indentEdges, remove)
		roomEdges = removePoint(roomEdges, remove)
	}

	if excludeEntrances {
		sortWallsForPoly(indentEdges, roomEdges)
	}


	walls := append(roomEdges, indentEdges...)

	if excludeEntrances{
		return walls
	}

	sortWallsInPointOrder(walls)
	sortWallsInPointOrder(entranceEdges)

	removeWalls := []*PairPoint2f{}

	for _, entrance := range entranceEdges {

		isVerticalEntrance := *entrance.First.X == *entrance.Second.X

		for _, wall := range walls {
			isVerticalWall := *wall.First.X == *wall.Second.X

			//skip if wall and entrance are not on the same plane
			if isVerticalEntrance != isVerticalWall {
				continue
			}

			if !linesIntersect(wall.First, wall.Second, entrance.First, entrance.Second) {
				continue
			}

			if !PointsEqual(wall.Second, entrance.Second){
				walls = append(walls, &PairPoint2f{
					First:  entrance.Second,
					Second: wall.Second,
				})
			}

			wall.Second = entrance.First

			if PointsEqual(wall.First, wall.Second){
				removeWalls = append(removeWalls, wall)
			}

			_ = wall
			_ = entrance
		}
	}

	for _, wall := range removeWalls {
		walls = removePoint(walls, wall)
	}
	return walls
}

func sortWallsInPointOrder(walls []*PairPoint2f) {
	for _, wall := range walls {
		isVerticalWall := *wall.First.X == *wall.Second.X

		swap := false
		if isVerticalWall {
			swap = doSwap(wall.First.Y, wall.Second.Y)
		}

		if !isVerticalWall {
			swap = doSwap(wall.First.X, wall.Second.X)
		}

		if swap {
			tmp := wall.Second
			wall.Second = wall.First
			wall.First = tmp
		}
	}
}

func doSwap(fst *float64, snd *float64) bool{
	if *fst > *snd {
		return true
	}
	return false
}

func sortWallsForPoly(indentEdges []*PairPoint2f, roomEdges []*PairPoint2f) {
	for _, tmpWall := range indentEdges {
		allWalls := append(roomEdges, indentEdges...)
		var firstEdges []Point2f
		var secondEdges []Point2f

		for _, wall := range allWalls {
			firstEdges = append(firstEdges, wall.First)
			secondEdges = append(secondEdges, wall.Second)
		}

		var (
			firstCount1  int
			secondCount1 int

			firstCount2  int
			secondCount2 int
		)

		for _, edge := range firstEdges {
			if PointsEqual(edge, tmpWall.First) {
				firstCount1++
			}
			if PointsEqual(edge, tmpWall.Second) {
				firstCount2++
			}
		}

		for _, edge := range secondEdges {
			if PointsEqual(edge, tmpWall.First) {
				secondCount1++
			}
			if PointsEqual(edge, tmpWall.Second) {
				secondCount2++
			}
		}

		if firstCount1 != secondCount1 || firstCount2 != secondCount2 {
			tmp := tmpWall.First
			tmpWall.First = tmpWall.Second
			tmpWall.Second = tmp
		}
	}
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
			xStart = indent.Location
			yStart = *room.Dimensions.X - *indent.Dimensions.X
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
