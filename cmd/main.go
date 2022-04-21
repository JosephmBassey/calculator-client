package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/josephmbassey/calculator-service/rpc/client"
	"github.com/josephmbassey/calculator-service/rpc/proto/calculatorpb"
)

func main() {
	// GRPC Client connections
	ctx := context.Background()
	timeout := time.Duration(10) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	logger := log.NewLogfmtLogger(os.Stdout)
	defer cancel()

	calculatorGRPCServer := os.Getenv("CALCULATOR_GRPC_SERVER")
	if calculatorGRPCServer == "" {
		calculatorGRPCServer = "0.0.0.0:9081"
	}
	calculatorAPIClient, err := connectToCalculatorAPI(ctx, logger, calculatorGRPCServer)
	if err != nil {
		logger.Log("msg", "error connecting to calculator service", "err", err)
		os.Exit(1)
	}

	op := flag.String("method", "", "The Operator to use")
	numberOne := flag.Float64("a", 0, "The first operand")
	numberTwo := flag.Float64("b", 0, "The second operand ")
	flag.Parse()

	operator := *op
	firstOperand := *numberOne
	secondOperand := *numberTwo

	computeRequest := &calculatorpb.CalculateRequest{
		Operands: &calculatorpb.OPERANDS{
			Number_1: firstOperand,
			Number_2: secondOperand,
		},
	}

	switch operator {
	case "add":
		computeRequest.Operator = calculatorpb.OPERATOR_OPERATOR_ADD
	case "sub":
		computeRequest.Operator = calculatorpb.OPERATOR_OPERATOR_SUBTRACT
	case "mul":
		computeRequest.Operator = calculatorpb.OPERATOR_OPERATOR_MULTIPLY
	case "div":
		computeRequest.Operator = calculatorpb.OPERATOR_OPERATOR_DIVIDE
	default:
		logger.Log("msg", `invalid operator, expected any of "add, sub, mul, div"`)
		os.Exit(1)
	}

	result, err := calculatorAPIClient.Calculator(ctx, computeRequest)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println(result.GetResult())
}

func connectToCalculatorAPI(ctx context.Context, logger log.Logger, calculatorServiceGRPCEndpoint string) (*client.CalculatorClient, error) {
	dopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Minute,
			Timeout:             30 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	calculatorClientConn, err := grpc.DialContext(ctx, calculatorServiceGRPCEndpoint, dopts...)
	if err != nil {
		return nil, err
	}

	calculatorSrvClient := client.NewCalculatorServiceClient(logger, calculatorClientConn)

	return calculatorSrvClient, nil
}
