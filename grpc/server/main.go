package main

import (
	"context"
	"log"
	"net"
	"strings"

	pb "github.com/ratnadeepb/grpc/customer"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// Server is used to implement customer.CustomServer
type Server struct {
	savedCustomers []*pb.CustomerRequest
}

// CreateCustomer creates a new customer
func (s *Server) CreateCustomer(ctx context.Context, in *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	s.savedCustomers = append(s.savedCustomers, in)
	return &pb.CustomerResponse{Id: in.Id, Success: true}, nil
}

// GetCustomers returns all customers by given filter
func (s *Server) GetCustomers(filter *pb.CustomerFilter, stream pb.Customer_GetCustomersServer) error {
	for _, customer := range s.savedCustomers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
			if err := stream.Send(customer); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, &Server{})
	s.Serve(lis)
}
