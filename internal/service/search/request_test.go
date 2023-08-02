package search

import (
	"context"
	"os"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestClient_CreateIndex(t *testing.T) {
	type args struct {
		ctx  context.Context
		body string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create feed_item index for zincsearch",
			args: args{
				ctx:  gctx.New(),
				body: CREATE_FEED_ITEM_INDEX_REQUEST,
			},
			wantErr: false,
		},
		{
			name: "Create feed_channel index for zincsearch",
			args: args{
				ctx:  gctx.New(),
				body: CREATE_FEED_CHANNEL_INDEX_REQUEST,
			},
			wantErr: false,
		},
	}

	os.Setenv("ZINC_BASE_URL", "http://localhost:4080")
	os.Setenv("ZINC_FIRST_ADMIN_USER", "admin")
	os.Setenv("ZINC_FIRST_ADMIN_PASSWORD", "qazxsw")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitClient(tt.args.ctx)
			c := GetClient(tt.args.ctx)
			if err := c.CreateIndex(tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
