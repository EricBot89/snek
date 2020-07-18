# Snek
Multiplayer snake over tcp written in go

## usage

Building Snek:
    run
	go install https://github.com/EricBot89/snek



To start a server: 
    run snek_game -s 
    By default this serves on all ips on port 8080. To configure the port to serve on use the -p flag

To start a client session: 
    run snek_game -c. 
    By default this connects to localhost (127.0.0.1 on port 8080) with the name PLAYER. To configure player name use the -n flag, to configure the host ip use the -ip flag, to configure the host port use the -p flag.


## Credits
Snek is built on top of the excellent
[termbox-go](https://github.com/nsf/termbox-go) library.

[Applied Go](https://appliedgo.net/networking/) provided an excelent article on tcp networking which was of great help.