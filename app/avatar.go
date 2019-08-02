package main

import "errors"

var ErrNoAvatar = errors.New("Unable to get avatar")

type Avatar interface {
	GetAvatarURL(Desenvolvedor) (string, error)
}

type FileAvatar struct{}

type ProviderAvatar struct{}

type UseAvatar []Avatar
