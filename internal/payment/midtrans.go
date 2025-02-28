package payment

import (
	"fmt"

	"github.com/Auxesia23/toko-online/internal/env"
	"github.com/Auxesia23/toko-online/internal/models"
	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func CreateMidtransPayment(order *models.Order) (string, error) {

	serverKey := env.GetString("MIDTRANS_SERVER_KEY", "")
	midclient := snap.Client{}
	midclient.New(serverKey, midtrans.Sandbox)

	var items []midtrans.ItemDetails
	for _, item := range order.OrderItems {
		items = append(items, midtrans.ItemDetails{
			ID:    item.ProductID.String(),
			Name:  item.Product.Name,
			Price: int64(item.Price),
			Qty:   int32(item.Quantity),
		})
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.ID.String(),
			GrossAmt: int64(order.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.User.Name,
			Email: order.User.Email,
		},
		Items: &items,
	}

	snapResp, err := midclient.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	fmt.Println("Midtrans Snap URL:", snapResp.RedirectURL)
	return snapResp.Token, nil
}
