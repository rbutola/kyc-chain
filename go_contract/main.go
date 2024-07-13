package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/crypto"
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
	if client == nil || contractAddr == (common.Address{}) {
		return nil, errors.New("client or contract address is nil")
	}

	// Define the ABI for the contract
	contractABI := `[{"inputs":[{"internalType":"address","name":"_address","type":"address"}],"name":"isKYCed","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"kycedAddresses","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_address","type":"address"},{"internalType":"bool","name":"_status","type":"bool"}],"name":"setKYCStatus","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	// Parse the ABI
	abiParsed, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	// Bind the contract
	contract := bind.NewBoundContract(contractAddr, abiParsed, client, client, nil)

	fmt.Println("This is my contract structure after initialization", contract)
	// Create session instance
	session := &KYCRegistrySession{
		contract:     contract,
		client:       client,
		from:         from,
		contractAddr: contractAddr,
	}

	fmt.Println("this is my session", session)
	return &KYCRegistry{session}, nil
}

// IsKYCed checks if an address is KYCed
func (s *KYCRegistrySession) IsKYCed(ctx context.Context, address common.Address) (bool, error) {
	if s.contract == nil {
		return false, errors.New("contract is nil")
	}

	callOpts := bind.CallOpts{
		Pending: false,
		Context: ctx,
		From: s.from,
	}

	// Prepare the method name and arguments
	method := "isKYCed"
	args := []interface{}{address}
	var result []interface{}

	fmt.Println("this the address", address)
	fmt.Println("These are the values that I am passing - ", callOpts, result, method, args)
	fmt.Println("DISCO These are the values that I am passing - ", callOpts, result, method, args)
	err := s.contract.Call(&callOpts, &result, method, args...)
	if err != nil {
		return false, err
	}

	return result[0].(bool), nil

}
/*
// SetKYCStatus sets the KYC status for an address
func (s *KYCRegistrySession) SetKYCStatus(ctx context.Context, address common.Address, status bool) (common.Hash, error) {
	transactOpts := bind.NewKeyedTransactor(s.from)

	tx, err := s.contract.Transact(ctx, "setKYCStatus", transactOpts, address, status)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}
*/
func (s *KYCRegistrySession) SetKYCStatus(ctx context.Context, address common.Address) (common.Hash, error) {
    // Create transaction options with the sender's private key
    // TODO: replace from with private key
    privateKeyS := "22281c4386f53038002c539f4d835ac409998a2b8a736a7d64ff80fe9868d6e1"
    privateKey, err := crypto.HexToECDSA(privateKeyS) // Replace with your private key
    if err != nil {
	    return common.Hash{}, fmt.Errorf("failed to parse private key: %v", err)
    }
    var nonce int = 42069
    // Convert int to *big.Int
    nonceBig := big.NewInt(int64(nonce))
    transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, nonceBig)
    if err != nil {
        return common.Hash{}, fmt.Errorf("failed to create transactor: %v", err)
    }
    // Set gas limit and gas price (optional, adjust as needed)
    transactOpts.GasLimit = 300000 // Example gas limit
    transactOpts.GasPrice = big.NewInt(int64(21000))   // Use suggested gas price
    statusC := true
    //param := interface{}{address, statusC}
    //param2 := interface{}{statusC}
    // Perform the transaction
    fmt.Println("this is my contract", s.contract)
    fmt.Println("this is my contract", s.contract)
    fmt.Println("this is my contract", s.contract)
    tx, err := s.contract.Transact(transactOpts, "setKYCStatus", address, statusC)
    if err != nil {
	    return common.Hash{}, fmt.Errorf("failed to send transaction: %v", err)
    }

    // Return the transaction hash
    return tx.Hash(), nil
}

func main() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("This is my client Connection", client)
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

	// Set address is KYCed
	newAddr := common.HexToAddress("0x5D9BD6F09052372b1ECA8D72e917AEb19566a33b")
	hash, err := kycRegistry.SetKYCStatus(context.Background(), newAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Address set kyc %s KYC status: %v\n", newAddr.Hex(), hash)
}
