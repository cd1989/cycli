package dag

type Edge struct {
	From      string
	To        string
	Decorator func(format string, a ...interface{}) string
}

func (e *Edge) decorate(element string) string {
	if e.Decorator == nil {
		return element
	}

	return e.Decorator(element)
}

func (e *Edge) v() string {
	return e.decorate("│")
}

func (e *Edge) h() string {
	return e.decorate("─")
}

func (e *Edge) lb() string {
	return e.decorate("╰")
}

func (e *Edge) lt() string {
	return e.decorate("╭")
}

func (e *Edge) rt() string {
	return e.decorate("╮")
}

func (e *Edge) rb() string {
	return e.decorate("╯")
}

func (e *Edge) lArrow() string {
	return e.decorate("<")
}

func (e *Edge) rArrow() string {
	return e.decorate(">")
}

func (e *Edge) tArrow() string {
	return e.decorate("^")
}

func (e *Edge) bArrow() string {
	return e.decorate("v")
}
