package draw

import (
	"fmt"
	"testing"
	"time"
)

func TestDraw(*testing.T) {
	draw := GetLineDrawer()
	draw.SetPath("/tmp/line.html")

	var klineVec = [...]KlineData{
		{Time: "2018/1/24 09:09:09", Data: [4]float32{2320.26, 2320.26, 2287.3, 2362.94}},
		{Time: "2018/1/25 09:09:09", Data: [4]float32{2300, 2291.3, 2288.26, 2308.38}},
		{Time: "2018/1/28 09:09:09", Data: [4]float32{2295.35, 2346.5, 2295.35, 2346.92}},
		{Time: "2018/1/29 09:09:09", Data: [4]float32{2347.22, 2358.98, 2337.35, 2363.8}},
		{Time: "2018/1/30 09:09:09", Data: [4]float32{2360.75, 2382.48, 2347.89, 2383.76}},
		{Time: "2018/1/31 09:09:09", Data: [4]float32{2383.43, 2385.42, 2371.23, 2391.82}},
		{Time: "2018/2/1 09:09:09", Data: [4]float32{2377.41, 2419.02, 2369.57, 2421.15}},
		{Time: "2018/2/4 09:09:09", Data: [4]float32{2425.92, 2428.15, 2417.58, 2440.38}},
		{Time: "2018/2/5 09:09:09", Data: [4]float32{2411, 2433.13, 2403.3, 2437.42}},
		{Time: "2018/2/6 09:09:09", Data: [4]float32{2432.68, 2434.48, 2427.7, 2441.73}},
		{Time: "2018/2/7 09:09:09", Data: [4]float32{2430.69, 2418.53, 2394.22, 2433.89}},
		{Time: "2018/2/8 09:09:09", Data: [4]float32{2416.62, 2432.4, 2414.4, 2443.03}},
	}
	for _, kline := range klineVec {
		draw.PlotKLine(kline)
	}
	/*
		for _, line := range lineVec {
			draw.PlotLine("mm5", line)
		}
		for _, line := range lineVec {
			draw.PlotLine("mm3", line)
		}
	*/
	if err := draw.Draw(); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	time.Sleep(30 * time.Minute)
}
