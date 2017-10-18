package main

import (
	"fmt"
	"os/exec"
	"math/rand"
	"time"
	"github.com/nsf/termbox-go"
)
/**
 行，列
 */
const Rows,Cols = 5,5

/**
 格子
 */
var Board [Rows][Cols]int

/**
 名字
 */
var Name = [...]string{"","2","4","8","16","32","64","128","256","512","1024","2048"}
var Color = [...]string{"","37","36","35","34","33","32","31","30","29","28","27"} //"\x1b[%dm%s\x1b[0m"

/**
 是否结束
 */
var GameOver bool
var IsMove bool

/**
 画棋盘
 */
func drawBoard(){
	//下面三行时清屏功能
	cmd := exec.Command("clear")
	result,_ := cmd.Output()
	fmt.Print(string(result))
	for r:=0;r<Rows;r++{
		fmt.Print("|")
		for c:=0;c<Cols;c++{
			//fmt.Printf("%6s|",Name[Board[r][c]]);
			fmt.Printf("\x1b[%sm%6s\x1b[0m|",Color[Board[r][c]],Name[Board[r][c]]);
		}
		fmt.Println("\n-----------------------------------------------");
	}
	if ErrorMsg !=""{
		fmt.Println(ErrorMsg)
	}
}

/**
 获取空白点
 */
func getBlankBox()([2]int,bool){
	var blanks [][2]int
	for r:=0;r<Rows;r++{
		for c:=0;c<Cols;c++{
			if Board[r][c] == 0{
				blanks = append(blanks,[2]int{r,c})
			}
		}
	}
	count := len(blanks)
	//fmt.Println(blanks)
	if count > 0{
		rands := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := rands.Intn(count)
		//fmt.Println("index",index,count)
		return blanks[index],false
	} else {
		return [2]int{0,0},true
	}
}

/**
 添加一个子
 */
func addPoint(){
	_point,GameOver := getBlankBox()
	if GameOver{
		ErrorMsg = "Over"
		return;
	}
	point := [2]int(_point)
	rands := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rands.Intn(2) == 1{
		Board[point[0]][point[1]] = 2
	} else {
		Board[point[0]][point[1]] = 1
	}
	//fmt.Println(8)
}

const (
	KeyUp = "up"
	KeyDown = "down"
	KeyLeft = "left"
	KeyRight = "right"
)
var ErrorMsg string
var Key string
var KeyContent chan string

func KeyInput(){
	var k string
	fmt.Scanln(&k)
	KeyCode := []byte(k)
	Key = ""
	if len(KeyCode) == 3 && KeyCode[0] == 27 && KeyCode[1] == 91{
		switch KeyCode[2] {
		case 65:
			Key = KeyUp
		case 66:
			Key = KeyDown
		case 67:
			Key = KeyLeft
		case 68:
			Key = KeyRight
		}
		ErrorMsg = ""
		ActionBoard()
	} else{
		ErrorMsg = "按键错误，请按方向键后回车"
	}
	drawBoard()
	KeyInput()
}

func ActionBoard(){
	switch Key {
	case KeyUp:
		RateBoardLeft()
		RateBoardLeft()
		AlignBottom()
		RateBoardLeft()
		RateBoardLeft()
	case KeyDown:
		AlignBottom()
	case KeyLeft:
		RateBoardLeft()
		AlignBottom()
		RateBoardLeft()
		RateBoardLeft()
		RateBoardLeft()
	case KeyRight:
		RateBoardLeft()
		RateBoardLeft()
		RateBoardLeft()
		AlignBottom()
		RateBoardLeft()
	}
	Key = ""
	if IsMove {
		addPoint()
		IsMove = false
	}
}
/**
 旋转
 */
func RateBoardLeft(){
	var CopyBoard [Cols][Rows]int
	for r :=0; r< Rows;r++{
		for c:=0;c<Cols;c++{
			CopyBoard[c][Rows - r -1] = Board[r][c]
		}
	}
	Board = CopyBoard
}
/**
 向下对齐
 */
func AlignBottom(){
	_alignBottom()
	//对齐完后进行合并
	MergeBoard()
	_alignBottom()
}
func _alignBottom(){
	for i := 1;i < Rows;i ++{
		for r := (Rows-1);r>0;r-- {
			for c := (Cols-1); c >= 0; c-- {
				if Board[r][c] == 0 && Board[r-1][c] != 0{
					//可以下移动
					Board[r][c] = Board[r-1][c]
					Board[r-1][c] = 0
					IsMove = true
				}
			}
		}
	}
	fmt.Println(Board)
}
func MergeBoard(){
	for i := 1;i < Rows;i ++{
		for r := (Rows-i);r>0;r-- {
			for c := (Cols-1); c >= 0; c-- {
				if  Board[r][c] != 0 && Board[r][c] ==  Board[r-1][c]{
					//可以下移动
					Board[r][c] ++
					Board[r-1][c] = 0
				}
			}
		}
	}
}

/**
 开始
 */
func Start(){
	//fmt.Println(Rows,Cols);
	addPoint()
	addPoint()
	drawBoard()
	//KeyInput()
	//RateBooard();
}

func main(){
	err := termbox.Init()
	if err != nil{
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse | termbox.InputCurrent)
	termbox.Clear(termbox.ColorDefault,termbox.ColorDefault)

	go Start()



//Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			//case termbox.KeyEsc, termbox.KeyF1:
			//	ErrorMsg = "Haha!"
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowLeft:
				Key = KeyRight
				ErrorMsg = ""
			case termbox.KeyArrowRight:
				Key = KeyLeft
				ErrorMsg = ""
			case termbox.KeyArrowUp:
				Key = KeyUp
				ErrorMsg = ""
			case termbox.KeyArrowDown:
				Key = KeyDown
				ErrorMsg = ""
			case termbox.KeyCtrlQ:
				ErrorMsg = "退出游戏!"
				goto Quit
			default:
				//开启这个会退出，因为退出到Loop 这一层，也就是for 循环
				//break Loop
			}
		default:
			//ErrorMsg = fmt.Sprintf("%v",ev.Key)
			//break Loop
		}
		ActionBoard()
		drawBoard()
	}
Quit:
}
