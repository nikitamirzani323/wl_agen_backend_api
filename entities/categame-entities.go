package entities

type Model_categame struct {
	Categame_id         string      `json:"categame_id"`
	Categame_name       string      `json:"categame_name"`
	Categame_list       interface{} `json:"categame_list"`
	Categame_status     string      `json:"categame_status"`
	Categame_status_css string      `json:"categame_status_css"`
	Categame_create     string      `json:"categame_create"`
	Categame_update     string      `json:"categame_update"`
}
type Model_game struct {
	Game_id            string `json:"game_id"`
	Game_idcategame    string `json:"game_idcategame"`
	Game_idprovider    int    `json:"game_idprovider"`
	Game_nmprovider    string `json:"game_nmprovider"`
	Game_name          string `json:"game_name"`
	Game_img           string `json:"game_img"`
	Game_multiplier    int    `json:"game_multiplier"`
	Game_urlstaging    string `json:"game_urlstaging"`
	Game_urlproduction string `json:"game_urlproduction"`
	Game_status        string `json:"game_status"`
	Game_status_css    string `json:"game_status_css"`
	Game_create        string `json:"game_create"`
	Game_update        string `json:"game_update"`
}
type Controller_categamesave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Categame_id     string `json:"categame_id"`
	Categame_name   string `json:"categame_name" validate:"required"`
	Categame_status string `json:"categame_status" validate:"required"`
}
type Controller_gamesave struct {
	Page               string `json:"page" validate:"required"`
	Sdata              string `json:"sdata" validate:"required"`
	Game_id            string `json:"game_id"`
	Game_idcategame    string `json:"game_idcategame" validate:"required"`
	Game_idprovider    int    `json:"game_idprovider" validate:"required"`
	Game_name          string `json:"game_name" validate:"required"`
	Game_img           string `json:"game_img" validate:"required"`
	Game_multiplier    int    `json:"game_multiplier" validate:"required"`
	Game_urlstaging    string `json:"game_urlstaging" validate:"required"`
	Game_urlproduction string `json:"game_urlproduction" validate:"required"`
	Game_status        string `json:"game_status" validate:"required"`
}
