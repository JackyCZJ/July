package store

import (
	"testing"
)

func TestCartAdd(t *testing.T) {
	type args struct {
		id      int32
		product Product
		count   int
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	tests = append(tests, struct {
		name    string
		args    args
		wantErr bool
	}{
		name: "test1", args: args{
			id: 31209,
			product: Product{
				ProductId: "5e832e5b041c8d2f73392b20",
			},
			count: 3,
		}, wantErr: false})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CartAdd(int32(tt.args.id), tt.args.product.ProductId, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("CartAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartDel(t *testing.T) {
	err := CartDel(29677, "5e8ebf344777eab5936c351d")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCartClear(t *testing.T) {
	err := CartClear(23086)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCartList(t *testing.T) {
	ca, err := CartList(29677)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ca)
}
