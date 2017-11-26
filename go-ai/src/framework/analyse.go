/**
* 高性能分析子系统
* 分析草丛情况、各部分元素占比
*
* 推荐只在开局调用一次
*
* author: linxingchen
*/
package framework

import (
    "container/list"
    "math"
    "fmt"
)

type MapAnalysis struct {
    Ocnt, Fcnt, Wcnt  int   // 空地、森林、墙数据
    Forests []Forest         // 森林相关

    TmpId       int         // 草丛标记
}

type Forest struct {
    Id          int         // 草丛唯一标示
    ForestMap   map[Position]bool // 草丛地图
    Area        int         // 面积
    Center      Position    // 中心点
    Nearest     Position    // 离出生点最近的入口
    Entrance    int         // 出口数量
    XLength     int         // 森林X长度
    YLength     int         // 森林Y长度
    HasFlag     bool        // 是否含旗
}

func (m *MapAnalysis) Analysis (state *GameState) {
    tmpExecedForest := make(map[Position]bool)

    for y := 1; y < state.Terain.Height - 1; y++ {
        for x := 1; x < state.Terain.Width - 1; x++ {
            switch state.Terain.Get(x, y) {
            case TerainEmpty:
                m.Ocnt++
            case TerainObstacle:
                m.Wcnt++
            case TerainForest:
                m.Fcnt++

                // 假如遇到草丛并且不在处理过的map里面，则进行BFT遍历处理
                if tmpExecedForest[Position{X:x, Y:y,}] == false {
                    // 进行BFS遍历草丛
                    m.Forests = append(m.Forests, m.bftForest(state, Position{X:x, Y:y}, &tmpExecedForest))
                }
            }
        }
    }
}

func (m *MapAnalysis) bftForest(state *GameState, firstForest Position, execedForest *map[Position]bool) Forest{
    // 对草丛进行分析
    // 从每个草堆左上角两个点进行广度优先遍历
    // 初始化
    queue := list.New()
    queue.PushBack(firstForest)
    pointer := queue.Front()
    // 每一个点需要找上、下、左、右
    pushUnique := make(map[Position]bool)
    pushUnique[firstForest] = true
    for pointer.Value != nil {
        // 注意排除已检查过的部分
        if pos, ok := (pointer.Value).(Position); ok && (*execedForest)[Position{X:pos.X, Y:pos.Y}] == false {
            // 上
            if state.Terain.Get(pos.X, pos.Y - 1) == TerainForest && false == pushUnique[Position{X:pos.X, Y:pos.Y - 1}] {
                queue.PushBack(Position{X:pos.X, Y:pos.Y - 1})
                pushUnique[Position{X:pos.X, Y:pos.Y - 1}] = true
            }

            // 左
            if state.Terain.Get(pos.X - 1, pos.Y) == TerainForest && false == pushUnique[Position{X:pos.X - 1, Y:pos.Y}] {
                queue.PushBack(Position{X:pos.X - 1, Y:pos.Y})
                pushUnique[Position{X:pos.X - 1, Y:pos.Y}] = true
            }

            // 右
            if state.Terain.Get(pos.X + 1, pos.Y) == TerainForest && false == pushUnique[Position{X:pos.X + 1, Y:pos.Y}]{
                queue.PushBack(Position{X:pos.X + 1, Y:pos.Y})
                pushUnique[Position{X:pos.X + 1, Y:pos.Y}] = true
            }

            // 下
            if state.Terain.Get(pos.X, pos.Y + 1) == TerainForest && false == pushUnique[Position{X:pos.X, Y:pos.Y + 1}]{
                queue.PushBack(Position{X:pos.X, Y:pos.Y + 1})
                pushUnique[Position{X:pos.X, Y:pos.Y + 1}] = true
            }

            // 添加后注意标记已检查
            (*execedForest)[Position{X:pos.X, Y:pos.Y}] = true
        }

        if pointer.Next() != nil {
            pointer = pointer.Next()
        } else {
            break
        }
    }


    // 开始对森林进行分析
    forest := Forest{}
    forest.Area = queue.Len()
    // 所有坐标值的平均值就是中心
    first := queue.Front()
    // 尽可能一个循环解决，同时求X、Y最大和最小
    xsum := 0
    ysum := 0
    xmin := math.MaxInt32
    ymin := math.MaxInt32
    xmax := -1
    ymax := -1

    // 记录边界情况
    border := make(map[Position]bool)
    border[firstForest] = true

    forest.ForestMap = make(map[Position]bool)

    // 如果只有一个草丛的情况
    for first.Value != nil {
        if pos, ok := (first.Value).(Position); ok {
            forest.ForestMap[pos] = true
            xsum += pos.X
            ysum += pos.Y

            if xmin > pos.X {
                xmin = pos.X
            }
            if ymin > pos.Y {
                ymin = pos.Y
            }
            if xmax < pos.X {
                xmax = pos.X
            }
            if ymax < pos.Y {
                ymax = pos.Y
            }
            // 找边界草丛情况，如果上下左右少于等于3个草，则为边界
            tmpGrass := 0
            if state.Terain.Get(pos.X + 1, pos.Y) == TerainForest {
                tmpGrass++
            }
            if state.Terain.Get(pos.X - 1, pos.Y) == TerainForest {
                tmpGrass++
            }
            if state.Terain.Get(pos.X, pos.Y - 1) == TerainForest {
                tmpGrass++
            }
            if state.Terain.Get(pos.X, pos.Y + 1) == TerainForest {
                tmpGrass++
            }
            if tmpGrass <= 3 {
                border[Position{X:pos.X, Y:pos.Y}] = true
            }
        }

        if first.Next() != nil {
            first = first.Next()
        } else {
            break
        }
    }
    forest.Id       = m.TmpId + 1
    m.TmpId++
    forest.XLength  = xmax - xmin + 1
    forest.YLength  = ymax - ymin + 1
    forest.Center.X = xsum / queue.Len()
    forest.Center.Y = ysum / queue.Len()

    // 判断是否有旗子
    if forest.Center.X == state.Params.FlagX && forest.Center.Y == state.Params.FlagY {
        forest.HasFlag = true
    } else {
        forest.HasFlag = false
    }

    // 森林出口数量使用附近的墙进行计算
    tmpForestExit := make(map[Position]bool)
    tmpForestBorder := make(map[Position]bool)
    for pos := range border {
        // 如果旁边是空地，则是出口
        if state.Terain.Get(pos.X + 1, pos.Y) == TerainEmpty {
            tmpForestExit[Position{X:pos.X + 1, Y:pos.Y}] = true
            tmpForestBorder[Position{X:pos.X, Y:pos.Y}] = true
        }
        if state.Terain.Get(pos.X - 1, pos.Y) == TerainEmpty {
            tmpForestExit[Position{X:pos.X - 1, Y:pos.Y}] = true
            tmpForestBorder[Position{X:pos.X, Y:pos.Y}] = true
        }
        if state.Terain.Get(pos.X, pos.Y + 1) == TerainEmpty {
            tmpForestExit[Position{X:pos.X, Y:pos.Y + 1}] = true
            tmpForestBorder[Position{X:pos.X, Y:pos.Y}] = true
        }
        if state.Terain.Get(pos.X, pos.Y - 1) == TerainEmpty {
            tmpForestExit[Position{X:pos.X, Y:pos.Y - 1}] = true
            tmpForestBorder[Position{X:pos.X, Y:pos.Y}] = true
        }
    }
    forest.Entrance = len(tmpForestExit)
    //遍历离出生点（0,0）最近的那片
    min := math.MaxInt32
    for pos := range tmpForestBorder {
        if min > pos.SDist(Position{X: 0, Y: 0}) {
            min = pos.SDist(Position{X: 0, Y: 0})
            forest.Nearest = pos
        }
    }
    return forest
}

/**
 * 判断点是否在那片草丛，如果命中则返回那片草丛
 */
func (m *MapAnalysis) GetForestByPos(pos Position) Forest{
    pos = Position { X: pos.X, Y: pos.Y }
    for _, f := range m.Forests {
        if f.ForestMap[pos] == true {
            return f
        }
    }
    return Forest{}
}
