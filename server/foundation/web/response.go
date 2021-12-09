package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteResponse(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	//trying to get value from context
	_, ok := ctx.Value(ContextKeyValue).(*ContextValue)
	if !ok {
		return fmt.Errorf("failed to retrieve value from context")
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("ContentType", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}
