package summarizer

import (
	"fmt"
	"net/http/httptest"
	"time"

	"logbook/internal/web/registration/reception/captured"
)

func ExampleSummarizer_post() {
	s := New(false)
	w := captured.New(httptest.NewRecorder())
	w.Write([]byte("lorem ipsum dolor sit amet"))
	fmt.Println(s.Post(w, time.Now().Add(-20*time.Millisecond)))
	// Output: 200* 20000µs 26B
}
