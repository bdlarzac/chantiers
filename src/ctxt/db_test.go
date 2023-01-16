package ctxt

import "testing"

func TestAjusteDbURL(t *testing.T) {
	for _, u := range []struct {
		a string
		b string
	}{
		{
			"postgresql://user:pass@host:port/name",
			"postgresql://user:pass@host:port/name?&search_path=schema"},
		{
			"postgresql://user:pass@host:port/name?search_path=sp",
			"postgresql://user:pass@host:port/name?search_path=sp",
		},
		{
			"postgresql://user:pass@host:port/name?sslmode=prefer",
			"postgresql://user:pass@host:port/name?&search_path=schema",
		},
		{
			"postgresql://user:pass@host:port/name?sslmode=prefer&param2=p2",
			"postgresql://user:pass@host:port/name?&param2=p2&search_path=schema",
		},
	} {
		res := ajusteDbURL(u.a, "schema")
		if res != u.b {
			t.Fatalf("url %s\nok:  %s\nnok: %s", u.a, u.b, res)
		}
	}
}
