package store

import (
	"sync"
	"testing"
)

func TestClient_ping(t *testing.T) {
	type fields struct {
		Mutex sync.Mutex
		mgo   *mongo.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Mutex: tt.fields.Mutex,
				mgo:   tt.fields.mgo,
			}
			if err := c.ping(); (err != nil) != tt.wantErr {
				t.Errorf("ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}