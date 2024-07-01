package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "time"
	"net/http"
	"log"
)

const width =40
const height = 20
const MAXage = 500
const Interv = 0.1
const clear  = "\033[2J"
const head   = "\033[1;1H"

type cells [width][height]int
var currentCells cells
var nextCells cells
var age = 0

func main() {
	// var scanner = bufio.NewScanner(os.Stdin)

	http.HandleFunc("/", handleGame)
    go func() {
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatal(err)
        }
    }()

	initCells()
	drawCells()
	for {
		fmt.Println(" Age =", age )
		// scanner.Scan()
		nextGeneration()
		drawCells()
		age++
		if age == MAXage {
			fmt.Println("Press Enter to continue")
			// scanner.Scan()
			initCells()
			drawCells()
			age = 0
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func initCells() {
	rand.Seed(time.Now().UnixNano())
	// for x := 1; x < width; x++ {
	// 	for y := 1; y < height; y++ {
	// 		currentCells[x][y] = rand.Intn(2)
	// 	}
	// }
	 // グライダー銃のパターン
	 pattern := [40][2]int{
        {1, 5}, {1, 6}, {2, 5}, {2, 6},
        {11, 5}, {11, 6}, {11, 7},
        {12, 4}, {12, 8},
        {13, 3}, {13, 9},
        {14, 3}, {14, 9},
        {15, 6},
        {16, 4}, {16, 8},
        {17, 5}, {17, 6}, {17, 7},
        {18, 6},
        {21, 3}, {21, 4}, {21, 5},
        {22, 3}, {22, 4}, {22, 5},
        {23, 2}, {23, 6},
        {25, 1}, {25, 2}, {25, 6}, {25, 7},
        {35, 3}, {35, 4},
        {36, 3}, {36, 4},
    }
	// パターンをセルに設定
    for _, p := range pattern {
        currentCells[p[0]][p[1]] = 1
    }
}

func drawCells() {
	fmt.Print(clear)
    fmt.Print(head)
	for y := 1; y < height; y++ {
		for x := 1; x < width; x++ {
			if currentCells[x][y] == 1 {
				fmt.Print("■")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func handleGame(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s\n", clear)
    fmt.Fprintf(w, "%s\n", head)
    for y := 1; y < height; y++ {
        for x := 1; x < width; x++ {
            if currentCells[x][y] == 1 {
                fmt.Fprint(w, "■")
            } else {
                fmt.Fprint(w, " ")
            }
        }
        fmt.Fprint(w, "\n")
    }
}


func nextGeneration() {
	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			//まわりのセルの生存数を数える
			count := 0 
			// if x >= width  {
			// 	// fmt.Println("x=", x)
			// 	x =1
			// }else if x < 0 {
			// 	x = width
			// }else if y >= height  {
			// 	// fmt.Println("x=", x)
			// 	y = 1
			// }else if y < 0 {
			// 	y = width
			// }
			for dx := -1; dx <= 1; dx++ {
                for dy := -1; dy <= 1; dy++ {
                    if dx != 0 || dy != 0 {
                        nx := (x + dx + width) % width
                        ny := (y + dy + height) % height
                        count += currentCells[nx][ny]
                    }
                }
            }
			//count := currentCells[x-1][y-1] + currentCells[x][y-1] + currentCells[x+1][y-1] + currentCells[x-1][y] + currentCells[x+1][y] + currentCells[x-1][y+1] + currentCells[x][y+1] + currentCells[x+1][y+1]
			if currentCells[x][y] == 0 {//死んでいるセル
				if count == 3 {//誕生
					nextCells[x][y] = 1
				} else {
					nextCells[x][y] = 0
				}
			} else {//生きているセル
				if count >= 4 { //過密
					nextCells[x][y] = 0
				}else if count == 2 || count == 3 {//生存
					nextCells[x][y] = 1
				} else {//過疎
					nextCells[x][y] = 0
				}
			}
		}
	}
	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			currentCells[x][y] = nextCells[x][y]
		}
	}
}

func end() {
    bufio.NewScanner(os.Stdin).Scan()
    fmt.Print(clear)
    fmt.Print(head)
}