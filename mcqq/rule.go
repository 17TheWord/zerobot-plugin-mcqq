package mcqq

import zero "github.com/wdvxdr1123/ZeroBot"

var GroupList []int64

func GroupRule(ctx *zero.Ctx) bool {
	return zero.OnlyGroup(ctx) && contains(GroupList, ctx.Event.GroupID)
}
