/*
 * 全图雷达
 *
 * 全图威胁只考虑子弹，不考虑坦克。2回合子弹经过是1,6回合接近0
 * author: linxingchen
 */
package framework

func (self *Radar) fullMapThreat(state *GameState) map[Position]float64 {
    // 不要玄学计算方法，威胁度直接确定在这里
    stepThreat := make(map[int]float64)
    stepThreat[0] = 1
    stepThreat[1] = 1
    stepThreat[2] = 0.8
    stepThreat[3] = 0.6
    stepThreat[4] = 0.4
    stepThreat[5] = 0.2

    fullmap := make(map[Position]float64)

    // 初始化
    for y := range state.Terain.Data {
        for x := range state.Terain.Data[y] {
            fullmap[Position{X: x, Y: y}] = 0
        }
    }

    // 子弹计算
    var bullets []Bullet
    for _, i := range state.EnemyBullet {
        bullets = append(bullets, i)
    }
    for _, t := range state.MyBullet {
        bullets = append(bullets, t)
    }
    loop:
    for _, b := range bullets {
        for step := 0; step < 6; step++ {
            if b.Pos.Direction == DirectionUp {
                // 向上 Y 骤减
                for y := b.Pos.Y - step * state.Params.BulletSpeed; y > b.Pos.Y - (step + 1) * state.Params.BulletSpeed; y-- {
                    // 检查墙
                    if y < 0 || state.Terain.Get(b.Pos.X, y) == TerainObstacle {
                        // 撞墙则停止
                        continue loop
                    }

                    // 根据步数进行赋值
                    fullmap[Position{X:b.Pos.X, Y:y, }] = stepThreat[step]
                }
            }

            if b.Pos.Direction == DirectionDown {
                // 向下 Y 骤增
                for y := b.Pos.Y + step * state.Params.BulletSpeed; y < b.Pos.Y + (step + 1) * state.Params.BulletSpeed; y++ {
                    // 检查墙
                    if y >= state.Terain.Height || state.Terain.Get(b.Pos.X, y) == TerainObstacle {
                        continue loop
                    }
                    fullmap[Position{X:b.Pos.X, Y:y, }] = stepThreat[step]
                }
            }

            if b.Pos.Direction == DirectionLeft {
                // 向左 X 骤减
                for x := b.Pos.X - step * state.Params.BulletSpeed; x > b.Pos.X - (step + 1) * state.Params.BulletSpeed; x-- {
                    if x < 0 || state.Terain.Get(x, b.Pos.Y) == TerainObstacle {
                        continue loop
                    }
                    fullmap[Position{X: x, Y: b.Pos.Y,}] = stepThreat[step]
                }
            }

            if b.Pos.Direction == DirectionRight {
                // 向右 X 骤增
                for x := b.Pos.X + step * state.Params.BulletSpeed; x < b.Pos.X + (step + 1) * state.Params.BulletSpeed; x++ {
                    if x > state.Terain.Width || state.Terain.Get(x, b.Pos.Y) == TerainObstacle {
                        continue loop
                    }
                    fullmap[Position{X:x, Y:b.Pos.Y,}] = stepThreat[step]
                }
            }
        }
    }

    return fullmap
}

