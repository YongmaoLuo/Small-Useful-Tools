package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"./v2rayAPI"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var api_address string
var api_port int

type Respon struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value"`
}

func (c *Respon) ProtoMessage()  { *c = Respon{} }
func (c *Respon) Reset()         {}
func (c *Respon) String() string { return proto.CompactTextString(c) }

func Search(conn *grpc.ClientConn) {
	father := context.Background()
	var enter []string = make([]string, 2)
	var exit bool

	for true {
		fmt.Println("Please enter command. If you need some help, please type command <help>.")
		fmt.Print("#")
		fmt.Scanln(&enter[0], &enter[1])
		exit = false
		var lowerEnter0 string = strings.ToLower(enter[0])
		switch lowerEnter0 {
		case "help":
			fmt.Println("<exit> exit the program.")
			fmt.Println("<statsService 'method'> get one user or users states.")
			fmt.Println("\tmethod: ")
			fmt.Println("\t\tGetStats: get one user's state with specific 'name' and 'reset':")
			fmt.Println("\t\tQueryStats: get all users' states with ",
				"specific 'pattern' and 'reset':")
			fmt.Println("'request':")
			fmt.Println("\t\tname(User):user>>>[email]>>>traffic>>>uplink/downlink")
			fmt.Println("\t\tname(Global):inbound>>>[tag]>>>traffic>>>uplink/downlink")
			fmt.Println("reset:true/false")
			break
		case "exit":
			exit = true
			break
		case "statsservice":
			isGet := true
			wrongCommand := false
			lowerEnter1 := strings.ToLower(enter[1])
			switch lowerEnter1 {
			case "getstats":
				isGet = true
			case "querystats":
				isGet = false
			default:
				fmt.Println("Wrong Method!")
				wrongCommand = true
			}
			if wrongCommand == true {
				break
			}
			inputReader := bufio.NewReader(os.Stdin)
			fmt.Print("Request content:")
			input, err := inputReader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			ctx, cancelFunc := context.WithCancel(father)
			response, err := v2rayAPI.CallStatsService(ctx, conn, enter[1], input)
			if err != nil {
				cancelFunc() //delete the child context
				panic(err)
			}
			cancelFunc()
			switch isGet {
			case true:
				OutputWrapping(response)
			default:
				fmt.Println("response is:", response)
			}
		default:
			fmt.Println("Wrong command!")
		}
		if exit == true {
			break
		}
	}
}

func OutputWrapping(response string) {
	response_bytes := []byte(response)
	var i int
	for i = 0; i < len(response_bytes); i++ {
		if response_bytes[i] == '\n' {
			break
		}
	}
	var j int
	for j = len(response_bytes) - 2; j >= 0; j-- {
		if response_bytes[j] == '\n' {
			break
		}
	}
	var input_bytes []byte = make([]byte, j-i-1)
	point := i + 1
	k := 0
	for point < j {
		input_bytes[k] = response_bytes[point]
		point++
		k++
	}
	input := string(input_bytes)
	output := &Respon{}
	err := proto.UnmarshalText(input, output)
	if err != nil {
		fmt.Println("protobuf UnmarshalText error: ", err)
	}

	split_name := strings.Split(output.Name, ">>>")
	if len(split_name) != 4 {
		fmt.Println("Wrong Name from Server!")
	}
	fmt.Println("Name:", split_name[1], "Traffic:", split_name[3])
	fmt.Println("The value is:", output.Value)
}

func main() {
	inputReader := bufio.NewReader(os.Stdin)
	for v2rayAPI.CheckIPv4(api_address) == false || api_port <= 0 {
		fmt.Print("Please enter your vps address(127.0.0.1:8080): ")
		input_bytes, err := inputReader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Wrong address! Please type again.")
		}
		input_bytes = bytes.TrimRight(input_bytes, "\r\n")
		var delim []byte = []byte(":")
		address := bytes.Split(input_bytes, delim)
		if len(address) != 2 {
			continue
		}
		api_address = string(address[0])
		api_port, _ = strconv.Atoi(string(address[1]))

	}
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", api_address, api_port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	Search(cmdConn)
}
