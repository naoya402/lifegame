package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

const width = 90//77
const height = 40//26
const MAXage = 500

type cells [width][height]int
var currentCells cells
var nextCells cells
var age = 0

func getGameState() string {
    var state string
    state += fmt.Sprintf("世代: %d<br>", age)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            if currentCells[x][y] == 1 {
                state += "■"
            } else {
                state += "□"
            }
        }
        state += "<br>" // HTMLの改行タグを使用
    }
    return state
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>ライフゲーム</title>
        <style>
            #game { 
                font-family: monospace; 
                white-space: pre;
                line-height: 1.2; /* 行間を調整 */
            }
        </style>
    </head>
    <body>
        <div id="game"></div>
        <script>
            const eventSource = new EventSource('/game');
            eventSource.onmessage = function(event) {
                const gameElement = document.getElementById('game');
                gameElement.innerHTML = event.data;
            };
        </script>
    </body>
    </html>
    `
    fmt.Fprint(w, html)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }

    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            gameState := getGameState()
            fmt.Fprintf(w, "data: %s\n\n", gameState)
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}

func initCells() {
    // 配列をクリア
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            currentCells[x][y] = 0
            nextCells[x][y] = 0
        }
    }
    
    // 初期パターンを設定
    pattern := [][2]int{
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
    for _, p := range pattern {
        currentCells[p[0]][p[1]] = 1
    }
}

func nextGeneration() {
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            count := 0 
            //セルの境界に達した時反対側に移動
            for dx := -1; dx <= 1; dx++ {
                for dy := -1; dy <= 1; dy++ {
                    if dx != 0 || dy != 0 {
                        nx := (x + dx + width) % width
                        ny := (y + dy + height) % height
                        count += currentCells[nx][ny]
                    }
                }
            }
            // // セルの境界に達した時に消去
            // for dx := -1; dx <= 1; dx++ {
            //     for dy := -1; dy <= 1; dy++ {
            //         if dx != 0 || dy != 0 {
            //             nx := x + dx
            //             ny := y + dy
            //             if nx >= 0 && nx < width && ny >= 0 && ny < height {
            //                 count += currentCells[nx][ny]
            //             }
            //         }
            //     }
            // }
            if currentCells[x][y] == 0 {
                //死んでいるセル
                if count == 3 {
                    //誕生
                    nextCells[x][y] = 1
                } else {
                    nextCells[x][y] = 0
                }
            } else {
                //生きているセル
                if count >= 4 {
                    //過密
                    nextCells[x][y] = 0
                } else if count == 2 || count == 3 {
                    //生存
					nextCells[x][y] = 1
				} else {
                    //過疎
					nextCells[x][y] = 0
				}
            }
        }
    }
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            currentCells[x][y] = nextCells[x][y]
        }
    }
}

func runGame() {
    initCells()
    for {
        nextGeneration()
        age++
        if age == MAXage {
            initCells()
            age = 0
        }
        time.Sleep(50 * time.Millisecond) /* スリープ間隔を調整 */
    }
}

func main() {
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/game", handleGame)
    
    go runGame()

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}