package mcqq

import zero "github.com/wdvxdr1123/ZeroBot"

var groupIdSet = make(map[int64]struct{})

func GroupRule(ctx *zero.Ctx) bool {
	if !zero.OnlyGroup(ctx) {
		return false
	}
	groupId := ctx.Event.GroupID
	_, ok := groupIdSet[groupId]
	return ok
}
