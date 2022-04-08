package akash

import (
	"bufio"
	"fmt"
	"net/rpc"
	"path/filepath"

	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	cinpuit "github.com/cosmos/cosmos-sdk/client/input"
)

const (
	flagAkashGoRPCServerHost = "akash-go-rpc-server-host"
	flagAkashGoRPCServerPort = "akash-go-rpc-server-port"
)

// NewScaffold returns a command that scafolds the config for a network, default is AKASH
func NewLaunch() *cobra.Command {
	c := &cobra.Command{
		Use:   "launch testnet",
		Short: "launch testnet",
		Args:  cobra.NoArgs,
		RunE:  akashLaunchHandler,
	}

	c.AddCommand()
	c.Flags().String(flagAkashGoRPCServerHost, "localhost", " akash go rpc server host")
	c.Flags().Int(flagAkashGoRPCServerPort, 8080, " akash go rpc server port")
	return c
}

type Args struct{}

func akashLaunchHandler(cmd *cobra.Command, args []string) error {
	hostname, err := cmd.Flags().GetString(flagAkashGoRPCServerHost)
	if err != nil {
		return err
	}

	port, err := cmd.Flags().GetInt(flagAkashGoRPCServerPort)
	if err != nil {
		return err
	}

	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Fatal("Couldn't connect to the AkashGo RPCServer: Is the Server running? \n error: ", err)
	}

	rpcArgs := Args{}
	reply := AccountResponse{}
	err = client.Call("AkashGoRPCService.CreateAccount", rpcArgs, &reply)
	if err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("Address: ", reply.Address)
	log.Println("Mnemonic: ", reply.Mnemonic)

	accountFunded, err := cinpuit.GetConfirmation("funds transferred", bufio.NewReader(os.Stdin), os.Stderr)
	if err != nil {
		panic(err)
	}

	if !accountFunded {
		log.Fatal("Account not funded")
	}

	log.Println("Account funded")
	log.Println("Starting Certificate Creation")

	re := ""
	err = client.Call("AkashGoRPCService.CreateCertificate", rpcArgs, &re)
	if err != nil {
		log.Fatal("Certificate creation failed with error: ", err)
	}

	log.Println("Certificate Creation Done")

	sdlFilePathForWeb := "./akash/SDL/deploy-web.yml"
	absoluteSDLPathForWeb, err := filepath.Abs(sdlFilePathForWeb)
	if err == nil {
		fmt.Println("Absolute:", absoluteSDLPathForWeb)
	}

	log.Printf("Staring Deployment using SDL file(with full path): %s", absoluteSDLPathForWeb)
	createDeploymentRequest := CreateDeploymentRequest{
		SDLFilePath: absoluteSDLPathForWeb,
	}

	de := DeploymentResponse{}
	err = client.Call("AkashGoRPCService.CreateDeployment", createDeploymentRequest, &de)
	if err != nil {
		log.Fatal("Deployment creation failed with error: ", err)
	}

	log.Println("Deployment creation done")
	log.Println("Save your DSEQ ID: ", de.DSeq)

	dseq := de.DSeq

	for {
		bidrequest := BidRequest{
			DSeq:  dseq,
			GSeq:  1,
			OSeq:  1,
			State: "open",
		}

		bidResponse := BidResponse{}
		err = client.Call("AkashGoRPCService.GetBids", bidrequest, &bidResponse)
		if err != nil {
			log.Fatal("fetching bids failed", err)
		}

		log.Println("all bids:")
		log.Println(bidResponse)

		if len(bidResponse.Bids) == 0 {
			log.Println("no bids, waiting for 2 seconds")
			time.Sleep(time.Second * 2)
			continue
		}

		for i, bid := range bidResponse.Bids {
			log.Printf("(%d)- %s: %s \n", i, bid.Provider, bid.Amount)
		}

		log.Println("Select the bid provider(type 0,1,2...)")

		providerIndexStr, err := cinpuit.GetString("bid provider", bufio.NewReader(os.Stdin))
		if err != nil {
			log.Fatal("error getting the input: ", err)
		}

		providerIndex, err := strconv.Atoi(providerIndexStr)
		if err != nil {
			log.Fatal("error converting the input: ", err)
		}

		if providerIndex < 0 || providerIndex >= len(bidResponse.Bids) {
			log.Fatal("invalid index")
		}

		provider := bidResponse.Bids[providerIndex].Provider
		log.Printf("Selected provider: %s", provider)

		log.Println("Starting lease creation")
		leaseCreateRequest := LeaseCreateRequest{
			Provider: provider,
			DSeq:     dseq,
			OSeq:     1,
			GSeq:     1,
		}

		leaseCreateResponse := LeaseCreateResponse{}

		err = client.Call("AkashGoRPCService.CreateLease", leaseCreateRequest, &leaseCreateResponse)
		if err != nil {
			log.Fatal("failed to create the lease", err)
		}

		log.Println("lease created")

		log.Println("Waiting for 2 seconds")
		time.Sleep(time.Second * 2)

		log.Println("Starting Manifest Send")

		sendManifestRequest := ManifestSendRequest{
			DSeq:        dseq,
			SDLFilePath: absoluteSDLPathForWeb,
		}

		sendManifestResponse := ManifestSendResponse{}

		err = client.Call("AkashGoRPCService.SendManifest", sendManifestRequest, &sendManifestResponse)
		if err != nil {
			log.Fatal("failed to send manifest", err)
		}

		log.Println("Manifest sent")

		log.Println("Waiting for 2 seconds")
		time.Sleep(time.Second * 2)

		log.Println("Checking Lease Status")

		lsr := LeaseStatusRequest{
			DSeq:     dseq,
			Provider: provider,
			GSeq:     1,
			OSeq:     1,
		}

		leaseResponse := LeaseStatusResponse{}

		err = client.Call("AkashGoRPCService.GetLeaseStatus", lsr, &leaseResponse)
		if err != nil {
			log.Fatal("failed to get lease status manifest", err)
		}

		log.Println(leaseResponse)
		log.Printf("Open http://%s:%d/ to view the deployment", leaseResponse.Host, leaseResponse.ExternalPort)
		return nil
	}
}

// Todo: generate these from the proto
type ManifestSendRequest struct {
	SDLFilePath string
	DSeq        uint64
}

type ManifestSendResponse struct {
}

type AccountResponse struct {
	Address  string
	Mnemonic string
}

type DeploymentResponse struct {
	DSeq uint64
}

type CreateDeploymentRequest struct {
	SDLFilePath string
}

type BidRequest struct {
	DSeq  uint64
	GSeq  uint32
	OSeq  uint32
	State string
}

type Bid struct {
	Amount   string
	Provider string
}

type BidResponse struct {
	Bids []Bid
}

type LeaseCreateRequest struct {
	Provider string
	DSeq     uint64
	OSeq     uint32
	GSeq     uint32
}

type LeaseCreateResponse struct {
}

type LeaseStatusRequest struct {
	Provider string
	DSeq     uint64
	GSeq     uint32
	OSeq     uint32
}

type LeaseStatusResponse struct {
	Host         string
	Port         uint16
	ExternalPort uint16
	Proto        string
	Available    int32
	Name         string
}
