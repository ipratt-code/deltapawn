package main

func engine() (frEng, toEng chan string){
    tell("Hello from deltapawn's engine!")
    frEng = make(chan string)
    toEng = make(chan string)
    go func(){
        for cmd := range toEng{
            switch cmd{
                case "stop":
                case "quit":
            }
        }
    }()
    return
}