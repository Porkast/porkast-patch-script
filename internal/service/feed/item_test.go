package feed

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestSetZHItemTotalCountToCache(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set Chinese feed item count to cache",
			args: args{
				ctx: gctx.New(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetZHItemTotalCountToCache(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetZHItemTotalCountToCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
