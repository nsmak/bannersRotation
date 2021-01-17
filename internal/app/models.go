package app

// Slot - место на сайте, на котором мы показываем баннер.
type Slot struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

// Banner - рекламный/информационный элемент, который показывается в слоте.
type Banner struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

// SociodemographicGroup - это группа пользователей сайта со схожими интересами,
// например "девушки 20-25" или "дедушки 80+".
type SociodemographicGroup struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}
