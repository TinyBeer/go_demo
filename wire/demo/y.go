package demo

type Y struct {
	Val int
}

func NewY(x X) Y {
	return Y{
		Val: x.Val + 1,
	}
}
