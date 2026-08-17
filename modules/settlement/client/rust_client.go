package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	settlementpb "go-backend/modules/settlement/proto"
)

type RustVerificationClient struct {
	conn   *grpc.ClientConn
	client settlementpb.SettlementServiceClient
}

// NewRustVerificationClient creates a gRPC client connection
// to the Rust verification engine.
func NewRustVerificationClient(
	target string,
) (*RustVerificationClient, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create gRPC client for Rust engine at %s: %w",
			target,
			err,
		)
	}

	client := settlementpb.NewSettlementServiceClient(conn)

	return &RustVerificationClient{
		conn:   conn,
		client: client,
	}, nil
}

// VerifyAndSettleTransaction sends the transaction to the
// Rust verification engine for on-chain verification.
//
// The Rust engine is responsible for verifying:
//   - transaction hash
//   - sender address
//   - receiver address
//   - amount
//   - currency
//   - network
//   - transaction status
//   - block/slot
func (c *RustVerificationClient) VerifyAndSettleLinkInvoiceTransaction(
	ctx context.Context,
	invoiceID string,
	txHash string,
	network string,
	amountPaid float64,
	currency string,
	senderAddress string,
	receiverAddress string,
	blockNumber int64,
) (
	bool,
	string,
	string,
	error,
) {
	// Validate client
	if c == nil || c.client == nil {
		return false, "", "",
			fmt.Errorf("rust verification client is not initialized")
	}

	// Build protobuf request
	req := &settlementpb.SettlementRequest{
		InvoiceId:       invoiceID,
		TxHash:          txHash,
		Network:         network,
		AmountPaid:      amountPaid,
		Currency:        currency,
		SenderAddress:   senderAddress,
		ReceiverAddress: receiverAddress,
		BlockNumber:     blockNumber,
	}

	// Call Rust verification engine
	resp, err := c.client.VerifyAndSettleTransaction(
		ctx,
		req,
	)

	if err != nil {
		return false, "", "",
			fmt.Errorf(
				"gRPC verification execution failed: %w",
				err,
			)
	}

	if resp == nil {
		return false, "", "",
			fmt.Errorf(
				"rust verification engine returned empty response",
			)
	}

	// Return:
	// success
	// merchant_id
	// message
	return resp.GetSuccess(),
		resp.GetMerchantId(),
		resp.GetMessage(),
		nil
}
func (c *RustVerificationClient) VerifyAndSettleAPIInvoiceTransaction(
	ctx context.Context,
	invoiceID string,
	txHash string,
	network string,
	amountPaid float64,
	currency string,
	senderAddress string,
	receiverAddress string,
	blockNumber int64,
) (
	bool,
	string,
	string,
	error,
) {
	// Validate client
	if c == nil || c.client == nil {
		return false, "", "",
			fmt.Errorf("rust verification client is not initialized")
	}

	// Build protobuf request
	req := &settlementpb.SettlementRequest{
		InvoiceId:       invoiceID,
		TxHash:          txHash,
		Network:         network,
		AmountPaid:      amountPaid,
		Currency:        currency,
		SenderAddress:   senderAddress,
		ReceiverAddress: receiverAddress,
		BlockNumber:     blockNumber,
	}

	// Call Rust verification engine
	resp, err := c.client.VerifyAndSettleTransaction(
		ctx,
		req,
	)

	if err != nil {
		return false, "", "",
			fmt.Errorf(
				"gRPC verification execution failed: %w",
				err,
			)
	}

	if resp == nil {
		return false, "", "",
			fmt.Errorf(
				"rust verification engine returned empty response",
			)
	}

	// Return:
	// success
	// merchant_id
	// message
	return resp.GetSuccess(),
		resp.GetMerchantId(),
		resp.GetMessage(),
		nil
}

// Close gracefully closes the gRPC connection.
func (c *RustVerificationClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf(
			"failed to close Rust gRPC connection: %w",
			err,
		)
	}

	return nil
}
