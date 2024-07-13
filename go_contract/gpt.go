package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type KYCRegistry struct {
	*KYCRegistrySession
}

type KYCRegistrySession struct {
	contract     *bind.BoundContract
	client       *ethclient.Client
	from         common.Address
	contractAddr common.Address
}

// NewKYCRegistry creates a new instance of KYCRegistrySession
func NewKYCRegistry(client *ethclient.Client, contractAddr common.Address, from common.Address) (*KYCRegistry, error) {
	// Define the ABI for the contract
	contractABI := `[{"inputs":[{"internalType":"address","name":"_address","type":"address"}],"name":"isKYCed","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"kycedAddresses","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_address","type":"address"},{"internalType":"bool","name":"_status","type":"bool"}],"name":"setKYCStatus","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	// Parse the ABI
	abiParsed, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	// Bind the contract
	contract := bind.NewBoundContract(contractAddr, abiParsed, nil, nil, nil)

	// Create session instance
	session := &KYCRegistrySession{
		contract:     contract,
		client:       client,
		from:         from,
		contractAddr: contractAddr,
	}

	return &KYCRegistry{session}, nil
}

// IsKYCed checks if an address is KYCed
func (s *KYCRegistrySession) IsKYCed(ctx context.Context, address common.Address) (bool, error) {
	callOpts := bind.CallOpts{
		Pending: false,
		Context: ctx,
		From:    s.from,
	}

	// Prepare the method name and arguments
	method := "isKYCed"
	args := []interface{}{address}

	// Prepare an empty interface for receiving the result
	var result bool

	err := s.contract.Call(&callOpts, &result, method, args...)
	if err != nil {
		return false, err
	}

	return result, nil
}

func main() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	// Replace with your contract address and from address
	contractAddr := common.HexToAddress("0xA957D512DCc88ABBF23eF46b555cF7fbC79745f9")
	fromAddr := common.HexToAddress("0x5D9BD6F09052372b1ECA8D72e917AEb19566a33b")

	// Create a new KYCRegistry instance
	kycRegistry, err := NewKYCRegistry(client, contractAddr, fromAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with the address to check
	address := common.HexToAddress("0x5D9BD6F09052372b1ECA8D72e917AEb19566a33b")

	// Check if the address is KYCed
	kycStatus, err := kycRegistry.IsKYCed(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Address %s KYC status: %v\n", address.Hex(), kycStatus)
}

