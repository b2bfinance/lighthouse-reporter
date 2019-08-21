package main

type (
	endpointResult struct {
		Path  string    `json:"path"`
		Score *scoreSet `json:"score"`
	}

	endpointResults []endpointResult
)

func (er endpointResults) Mean() *scoreSet {
	inLen := len(er)
	p := make([]int64, inLen)
	a := make([]int64, inLen)
	b := make([]int64, inLen)
	s := make([]int64, inLen)

	for i, v := range er {
		p[i] = v.Score.Performance
		a[i] = v.Score.Accessibility
		b[i] = v.Score.BestPractises
		s[i] = v.Score.SEO
	}

	return &scoreSet{
		Performance:   intSliceMean(p),
		Accessibility: intSliceMean(a),
		BestPractises: intSliceMean(b),
		SEO:           intSliceMean(s),
	}
}

func intSliceMean(in []int64) int64 {
	sum := int64(0)
	for _, i := range in {
		sum += i
	}
	return int64(sum / int64(len(in)))
}
