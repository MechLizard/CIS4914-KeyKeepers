package main

import (
	"fmt"
	"github.com/lmittmann/tint"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"log/slog"
)

type ResponseActor struct{}

func (state *ResponseActor) sendRequest(msg *HashRequest, ctx actor.Context) {
	remotePID := actor.NewPID(msg.GetRemoteAddress(), msg.GetRemotename()) // Remote address + actor name
	ctx.Request(remotePID, msg.GetRequest())
}

func (state *ResponseActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *HashRequest:
		state.sendRequest(msg, ctx)
	case *HashResult:
		if msg.Sha1Hash == "" {
			fmt.Println("File not found in malware database.")
			os.Exit(1)
		}

		fmt.Println("Hash found:")
		fmt.Println("File Name: ", msg.GetFileName())
		if msg.AltName1 != "" {
			fmt.Println("AltName1: ", msg.GetAltName1())
		}
		if msg.AltName2 != "" {
			fmt.Println("AltName2: ", msg.GetAltName2())
		}
		if msg.AltName3 != "" {
			fmt.Println("AltName3: ", msg.GetAltName3())
		}
		if msg.AltName4 != "" {
			fmt.Println("AltName4: ", msg.GetAltName4())
		}
		if msg.AltName5 != "" {
			fmt.Println("AltName5: ", msg.GetAltName5())
		}
		if msg.AltName6 != "" {
			fmt.Println("AltName6: ", msg.GetAltName6())
		}
		if msg.AltName7 != "" {
			fmt.Println("AltName7: ", msg.GetAltName7())
		}
		if msg.AltName8 != "" {
			fmt.Println("AltName8: ", msg.GetAltName8())
		}
		if msg.AltName9 != "" {
			fmt.Println("AltName9: ", msg.GetAltName9())
		}
		if msg.AltName10 != "" {
			fmt.Println("AltName10: ", msg.GetAltName10())
		}
		if msg.AltName11 != "" {
			fmt.Println("AltName11: ", msg.GetAltName11())
		}
		if msg.AltName12 != "" {
			fmt.Println("AltName12: ", msg.GetAltName12())
		}
		if msg.AltName13 != "" {
			fmt.Println("AltName13: ", msg.GetAltName13())
		}
		fmt.Println()
		fmt.Println("First Seen UTC: ", msg.GetFirstSeenUtc())
		fmt.Println("File Type Guess: ", msg.GetFileTypeGuess())
		fmt.Println("Mime Type: ", msg.GetMimeType())
		fmt.Printf("VirusTotal Percent: %v%%\n", msg.GetVtpercent())
		fmt.Println("Signature: ", msg.GetSignature())
		fmt.Println("Reporter: ", msg.GetReporter())
		fmt.Println()
		fmt.Println("Sha256 Hash: ", msg.GetSha256Hash())
		fmt.Println("Md5 Hash: ", msg.GetMd5Hash())
		fmt.Println("Sha1 Hash: ", msg.GetSha1Hash())
		fmt.Println()
		fmt.Println("ssdeep: ", msg.GetSsdeep())
		fmt.Println("tlsh: ", msg.GetTlsh())

		os.Exit(0)
	default:
		//fmt.Printf("Received unknown message: %#v\n", msg)
	}

}

func handleArgs(arg1, arg2, arg3, arg4 string) (string, int, string, string) {
	//fmt.Printf("arg1: %v\n", arg1)
	//fmt.Printf("arg2: %v\n", arg2)
	//fmt.Printf("arg3: %v\n", arg3)
	//fmt.Printf("arg4: %v\n", arg4)

	// IP address
	var remoteIP string
	if net.ParseIP(arg1) == nil {
		ips, err := net.LookupIP(arg1)
		if err != nil {
			fmt.Println("[System] Error resolving remote host:", err)
		}
		remoteIP = ips[0].String()
		fmt.Printf("The IP address of %v resolved to %v\n", arg1, remoteIP)
	} else {
		remoteIP = arg1
	}

	// Port
	port, err := strconv.Atoi(arg2) // Atoi = ASCII to Integer
	if err != nil {
		fmt.Println("[System] Conversion error:", err)
		os.Exit(1)
	}
	if port < 0 || port > 65535 {
		print("[System] Port must be between 0 and 65535")
		printHelp()
		os.Exit(1)
	}

	// Node Name
	nodeName := arg3

	// hash
	if len(arg4) != 40 {
		print("[System] sha1 hash must be 40 characters long.")
		printHelp()
		os.Exit(1)
	}
	hash := arg4
	return remoteIP, port, nodeName, hash
}

func printHelp() {
	fmt.Println(`Usage: keykeeper.exe <Node IP> <Node Port> <Node Name> <hash>

<Node IP>      The IP address of the node to query.
<Node Port>    The port number of the node to query.
<Node Name>    The name of the node to query.
<hash>         The hash of the file in sha1.


Example:
    keykeeper.exe 127.0.0.1 8080 node1 29812108115024db02ae79ac853743d31c606455`)
	os.Exit(1)
}

func findAvailablePort(startPort, maxPort int) (int, error) {
	for port := startPort; port <= maxPort; port++ {
		addr := fmt.Sprintf("0.0.0.0:%d", port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			_ = ln.Close() // Close it so Proto.Actor can use it
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found between %d and %d", startPort, maxPort)
}

func main() {
	if len(os.Args) != 5 {
		printHelp()
		os.Exit(1)
	}

	remoteIP, remotePort, remoteName, sha1Hash := handleArgs(os.Args[1], os.Args[2], os.Args[3], os.Args[4])

	// Create new logger that only logs errors, so client isn't spammed.
	logger := func(system *actor.ActorSystem) *slog.Logger {
		w := os.Stderr
		return slog.New(tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelError,
			TimeFormat: time.Kitchen,
		})).With("lib", "Proto.Actor").
			With("system", system.ID)
	}

	loggingConfig := actor.WithLoggerFactory(logger)
	// Setup actor system and remoting
	system := actor.NewActorSystem(loggingConfig)

	// NOTE: This doesn't open a TCP connection with the node.
	// If this client is in a NAT network they won't get a response back if port forwarding isn't set up.
	port, err := findAvailablePort(8000, 9000)
	if err != nil {
		println("[System] Error finding available port to listen for response:", err)
		os.Exit(1)
	}
	config := remote.Configure("0.0.0.0", port)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	// Start the local actor to receive replies
	props := actor.PropsFromProducer(func() actor.Actor { return &ResponseActor{} })
	clientPID := system.Root.Spawn(props)

	remoteAddress := remoteIP + ":" + strconv.Itoa(remotePort)
	request := &HashRequest{RemoteAddress: remoteAddress, Remotename: remoteName, Request: &CheckHash{Hash: sha1Hash, Type: HashType_sha1}}
	system.Root.Send(clientPID, request)

	// Keep the client alive long enough to receive the response
	time.Sleep(30 * time.Second)
}
