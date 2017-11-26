/**
 *  进草条件: 长宽大于两个出口
 *
 *
 *
**/
package tactics

import (
	f "framework"
	"fmt"
)

type Snake struct {
    mapanalysis    *f.MapAnalysis
    inforest       bool           // 是否进草战斗
}

func NewSnake() *Snake{
    return &Snake {
        mapanalysis: &f.MapAnalysis{},
        inforest: false,
    }
}

func (s *Snake) Init(state *f.GameState) {
    s.mapanalysis.Analysis(state)
    fmt.Printf("mapanalysis: %+v\n", s.mapanalysis)

    // 寻找最大草丛
    maxarea := -1
    var maxf f.Forest
    for _, f := range s.mapanalysis.Forests {
        if maxarea < 0 || maxarea < f.Area {
            maxf = f
        }
    }
    // 最大草丛是否适合战斗
    if maxf != (f.Forest{}) {
        distx := maxf.Center.X - state.Params.FlagX
        disty := maxf.Center.Y - state.Params.FlagY
        if distx < 0 {
            distx = -distx
        }
        if disty < 0 {
            disty = -disty
        }
        // 进草的最小宽度
        minlen  := state.Params.BulletSpeed * 2    // 草丛两边都大于此值
        maxlen  := (state.Params.BulletSpeed + 2) * 2    // 草丛两边至少有一边大于或等于此值
        minarea := (state.Params.BulletSpeed + 2) * (state.Params.BulletSpeed + 2)
        // 如果最大的森林在地图中央 & 且面积较大 & 长宽适合战斗
        if distx <= state.Params.BulletSpeed && disty <= state.Params.BulletSpeed && maxf.XLength > minlen  && maxf.YLength > minlen && (maxf.XLength >= maxlen || maxf.YLength >= maxlen) && maxf.Area > minarea {
            s.inforest = true
        } else {
            s.inforest = false
        }
    }
}

func (s *Snake) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {


}

func (s *Snake) End(state *f.GameState) {

}
