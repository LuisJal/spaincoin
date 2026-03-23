// Package market provides a deterministic price simulator for the SPC/EUR pair.
// The price is derived from block height, so all clients see the same value.
package market

import (
	"math"
)

// PricePoint represents a single price data point.
type PricePoint struct {
	Price     float64 `json:"price"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    float64 `json:"volume"`
	Height    uint64  `json:"height"`
	Timestamp int64   `json:"timestamp"`
}

// Simulator generates deterministic prices based on block height.
type Simulator struct {
	basePrice float64
}

// NewSimulator creates a new price simulator with a base price.
func NewSimulator(basePrice float64) *Simulator {
	return &Simulator{basePrice: basePrice}
}

// PriceAtHeight returns the SPC/EUR price at a given block height.
// Uses a combination of sinusoidal waves to create realistic-looking movement.
func (s *Simulator) PriceAtHeight(height uint64) float64 {
	h := float64(height)

	// Multiple overlapping waves at different frequencies
	wave1 := math.Sin(h/120.0) * 0.012 // slow trend
	wave2 := math.Sin(h/37.0) * 0.008  // medium oscillation
	wave3 := math.Sin(h/11.0) * 0.003  // fast noise
	wave4 := math.Sin(h/200.0) * 0.015 // macro trend

	// Slight upward drift over time (growth)
	drift := math.Log1p(h/50000.0) * 0.01

	price := s.basePrice + wave1 + wave2 + wave3 + wave4 + drift

	// Clamp to never go below 0.01
	if price < 0.01 {
		price = 0.01
	}

	// Round to 6 decimals
	return math.Round(price*1_000_000) / 1_000_000
}

// VolumeAtHeight returns simulated 24h volume at a given block height.
func (s *Simulator) VolumeAtHeight(height uint64) float64 {
	h := float64(height)
	base := 5000.0
	wave := math.Sin(h/80.0)*2000.0 + math.Sin(h/23.0)*800.0
	vol := base + wave
	if vol < 500 {
		vol = 500
	}
	return math.Round(vol*100) / 100
}

// CurrentPrice returns the price point at the given height.
func (s *Simulator) CurrentPrice(height uint64) PricePoint {
	price := s.PriceAtHeight(height)
	return PricePoint{
		Price:  price,
		High:   s.high24h(height),
		Low:    s.low24h(height),
		Volume: s.VolumeAtHeight(height),
		Height: height,
	}
}

// PriceHistory returns historical price points up to the given height.
// step controls how many blocks between each point.
func (s *Simulator) PriceHistory(currentHeight uint64, points int, step int) []PricePoint {
	if step < 1 {
		step = 1
	}
	result := make([]PricePoint, 0, points)

	startHeight := int64(currentHeight) - int64(points)*int64(step)
	if startHeight < 0 {
		startHeight = 0
	}

	for i := 0; i < points; i++ {
		h := uint64(startHeight + int64(i)*int64(step))
		if h > currentHeight {
			break
		}
		p := s.PriceAtHeight(h)

		// Compute candle high/low over the step range
		high, low := p, p
		for j := 0; j < step && h+uint64(j) <= currentHeight; j++ {
			pp := s.PriceAtHeight(h + uint64(j))
			if pp > high {
				high = pp
			}
			if pp < low {
				low = pp
			}
		}

		result = append(result, PricePoint{
			Price:  p,
			High:   high,
			Low:    low,
			Volume: s.VolumeAtHeight(h),
			Height: h,
		})
	}

	return result
}

// high24h returns the highest price in the last ~17280 blocks (24h at 5s/block).
func (s *Simulator) high24h(currentHeight uint64) float64 {
	blocks24h := uint64(17280)
	start := uint64(0)
	if currentHeight > blocks24h {
		start = currentHeight - blocks24h
	}
	high := 0.0
	// Sample every 100 blocks for performance
	for h := start; h <= currentHeight; h += 100 {
		p := s.PriceAtHeight(h)
		if p > high {
			high = p
		}
	}
	// Also check current
	p := s.PriceAtHeight(currentHeight)
	if p > high {
		high = p
	}
	return high
}

// low24h returns the lowest price in the last ~17280 blocks (24h at 5s/block).
func (s *Simulator) low24h(currentHeight uint64) float64 {
	blocks24h := uint64(17280)
	start := uint64(0)
	if currentHeight > blocks24h {
		start = currentHeight - blocks24h
	}
	low := math.MaxFloat64
	for h := start; h <= currentHeight; h += 100 {
		p := s.PriceAtHeight(h)
		if p < low {
			low = p
		}
	}
	p := s.PriceAtHeight(currentHeight)
	if p < low {
		low = p
	}
	return low
}

// Change24h returns the percentage change over the last 24h.
func (s *Simulator) Change24h(currentHeight uint64) float64 {
	blocks24h := uint64(17280)
	prevHeight := uint64(0)
	if currentHeight > blocks24h {
		prevHeight = currentHeight - blocks24h
	}
	pricePrev := s.PriceAtHeight(prevHeight)
	priceNow := s.PriceAtHeight(currentHeight)
	if pricePrev == 0 {
		return 0
	}
	change := ((priceNow - pricePrev) / pricePrev) * 100
	return math.Round(change*100) / 100
}
