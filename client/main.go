package main

import (
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/midoblgsm/mySimpleGRPC/customer"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {
	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {
	// calling the streaming API
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}
func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    11,
		Name:  "Mohamed Mohamed",
		Email: "email@xyz.com",
		Phone: "111-222-3333",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "650 Harry Road",
				City:              "San Jose",
				State:             "CA",
				Zip:               "95123",
				IsShippingAddress: false,
			},
			&pb.CustomerRequest_Address{
				Street:            "My second nice address",
				City:              "My nice city",
				State:             "ST",
				Zip:               "12345",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)

	customer = &pb.CustomerRequest{
		Id:    12,
		Name:  "Good Friend",
		Email: "anotheremail@xyz.com",
		Phone: "444-555-6666",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Another Address",
				City:              "Same City",
				State:             "ST",
				Zip:               "12345",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)
	// Filter with an empty Keyword
	filter := &pb.CustomerFilter{Keyword: ""}
	getCustomers(client, filter)
}
