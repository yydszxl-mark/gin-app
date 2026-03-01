package repository

import "gin-app-start/internal/model"

var RepoGroupApp = new(RepoGroup)

type RepoGroup struct {
	UserRepository
	DeviceRepository
}

func init() {
	RepoGroupApp.UserRepository = &userRepository{BaseRepository: &BaseRepository[model.User]{}}
	RepoGroupApp.DeviceRepository = &deviceRepository{BaseRepository: &BaseRepository[model.Device]{}}
}
