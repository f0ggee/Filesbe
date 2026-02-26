package Handlers

import "Kaban/iternal/DomainLevel"

type HandlerPack struct {
	Tokens    DomainLevel.ManageTokens
	Database  DomainLevel.UserServer
	TokenImpl DomainLevel.ManageTokensImpl
}

type HandlerPackCollect struct {
	S HandlerPack
}

func CollectorPack(s HandlerPack) *HandlerPackCollect {
	return &HandlerPackCollect{S: s}
}
