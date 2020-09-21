package main

import (
    "github.com/thoj/go-ircevent"
    "fmt"
    "strconv"
    "strings"
    "time"
    "log"
    "net"
    "bufio"
)
type Server struct {
    Addr string
}

var serverName = "IRC.VHFDX.RU:6667"
var userName = "UA3MQJ_KO98KA"
var roomName = "#vhfdx"
var telnetAddr Server

var callsigns map[string]int

// опрос тех, кто в комнате
func f(roomName string, con *irc.Connection) {
    for {
        time.Sleep(1 * time.Second)
        // fmt.Println("Check " + roomName)
        con.Who(roomName)
    }
}

// телнет серверная часть
func (srv Server) ListenAndServe() error {
    addr := srv.Addr
    if addr == "" {
        addr = ":8080"
    }
    log.Printf("starting server on %v\n", addr)
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    defer listener.Close()
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("error accepting connection %v", err)
            continue
        }
        log.Printf("accepted connection from %v", conn.RemoteAddr())
        handle(conn) //TODO: Implement me
    }
}

func handle(conn net.Conn) error {
    defer func() {
        log.Printf("closing connection from %v", conn.RemoteAddr())
        conn.Close()
    }()
    r := bufio.NewReader(conn)
    w := bufio.NewWriter(conn)
    scanr := bufio.NewScanner(r)
    for {
        scanned := scanr.Scan()
        if !scanned {
            if err := scanr.Err(); err != nil {
                log.Printf("%v(%v)", err, conn.RemoteAddr())
                return err
            }
            break
        }
        w.WriteString(strings.ToUpper(scanr.Text()) + "\n")
        w.Flush()
    }
    return nil
}

func connect_to_irc() {
    callsigns := make(map[string]int)

    con := irc.IRC(userName, userName)
    err := con.Connect(serverName)
    if err != nil {
        fmt.Println("Failed connecting")
        return
    }

    fmt.Println("Connected to irc server: " + serverName)

    con.AddCallback("001", func (e *irc.Event) {
        con.Join(roomName)
    })

    con.AddCallback("002", func (e *irc.Event) {
        fmt.Println("002: " + strconv.Quote(e.Message()))
    })

    con.AddCallback("003", func (e *irc.Event) {
        fmt.Println("003: " + strconv.Quote(e.Message()))
    })

    con.AddCallback("004", func (e *irc.Event) {
        fmt.Println("004: " + strconv.Quote(e.Message()))
    })

    con.AddCallback("005", func (e *irc.Event) {
        fmt.Println("005: " + strconv.Quote(e.Message()))
    })


    con.AddCallback("352", func (e *irc.Event) {
        // fmt.Println("352: " + e.Raw)

        var callsign string
        
        callsign = strings.Split(e.Raw, " ")[7]

        if _, exist := callsigns[callsign]; exist {
            // fmt.Println("Key found value is: ", value)
        } else {
            fmt.Println("New call " + callsign)
            callsigns[callsign] = 1
        }

    })

    con.AddCallback("315", func (e *irc.Event) {
        // fmt.Println("315: " + strconv.Quote(e.Message())) // 315: "End of /WHO list."
    })

    con.AddCallback("JOIN", func (e *irc.Event) {
        fmt.Println("Joined to room: " + roomName)
        con.Who(roomName)
        go f(roomName, con)
    })

    con.AddCallback("PRIVMSG", func (e *irc.Event) {
        fmt.Println("PRIVMSG: " + strconv.Quote(e.Message()))
    })


    con.Loop()
}

func main() {

    // connect_to_irc()
    
    Server.ListenAndServe(Server{Addr: ":8080"})

}
