package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var (
	shadowTotal    atomic.Int64
	shadowDisagree atomic.Int64
)

func recordShadowCompare(agree bool) {
	shadowTotal.Add(1)
	if !agree {
		shadowDisagree.Add(1)
	}
}

func shadowStatsSnapshot() map[string]any {
	total := shadowTotal.Load()
	disagree := shadowDisagree.Load()
	agree := total - disagree
	rate := 0.0
	if total > 0 {
		rate = float64(disagree) / float64(total)
	}
	return map[string]any{
		"total":               total,
		"agree":               agree,
		"disagree":            disagree,
		"shadow_disagree_rate": rate,
	}
}

func writeShadowPrometheus(w http.ResponseWriter) {
	total := shadowTotal.Load()
	disagree := shadowDisagree.Load()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP shadow_total Shadow mode comparisons\n# TYPE shadow_total counter\nshadow_total %d\n", total)
	fmt.Fprintf(w, "# HELP shadow_disagree_total Shadow disagreements\n# TYPE shadow_disagree_total counter\nshadow_disagree_total %d\n", disagree)
	if total > 0 {
		rate := float64(disagree) / float64(total)
		fmt.Fprintf(w, "# HELP shadow_disagree_rate Shadow disagree rate\n# TYPE shadow_disagree_rate gauge\nshadow_disagree_rate %f\n", rate)
	} else {
		fmt.Fprint(w, "# HELP shadow_disagree_rate Shadow disagree rate\n# TYPE shadow_disagree_rate gauge\nshadow_disagree_rate 0\n")
	}
}
