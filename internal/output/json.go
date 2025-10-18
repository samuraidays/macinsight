package output

import (
	"encoding/json"
	"io"

	"github.com/samuraidays/macinsight/pkg/types"
)

// レポートをJSONで出力（整形）
func WriteJSON(w io.Writer, r types.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
