package databox

type Dialect interface {
	Qt(exp string) string
	QtAll(exp ...string) []string
	QtPtr(exp *string) *string
}

type DialectConfig struct {
	quote string
}

func (d *DialectConfig) Qt(exp string) string {
	return d.quote + exp + d.quote
}

func (d *DialectConfig) QtPtr(exp *string) *string {
	if exp == nil {
		return nil
	}
	*exp = d.quote + *exp + d.quote
	return exp
}

func (d *DialectConfig) QtAll(exp ...string) []string {
	qts := []string{}
	for _, e := range exp {
		qts = append(qts, d.Qt(e))
	}
	return qts
}

func NewDialect() Dialect {
	return &DialectConfig{
		quote: `"`,
	}
}
