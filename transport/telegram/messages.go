package telegram

const (
	GET_TOKEN_VK_TEXT         = "Перейдите по ссылке: https://oauth.vk.com/authorize?client_id=8157407&redirect_uri=https://oauth.vk.com/blank.html&scope=friends,photos,offline&response_type=token&v=5.131 \nПройдите авторизацию вк, вас перенаправит на страницу, где в адресной строке указан ваш токен.\nСкопируйте эту ссылку и отправьте сообщением."
	GET_SUCCESS_VK_TOKEN      = "Ваш токен был сохранен успешно. Если вы сделали все верно, то теперь вы можете продолжать работу со слежкой, иначе ничего работать не будет."
	GET_UNSUCCESS_VK_TOKEN    = "Вы неправильно ввели строку (в строке должны находиться access_token= и expires_in= и user_id=. Попробуйте еще раз скопировать адрессную строку и вставить сюда."
	ADD_BY_NAME_TEXT          = "Этот человек должен быть у вас в друзьях. Введите имя фамилию его с вк. Далее выберите из списка того, за кем хотите следить."
	ADD_BY_LINK_VK_TEXT       = "Скопируйте ссылку профиля вк и отправьте сообщением."
	FRIENDS_BY_NAME_EMPTY     = "Друга с таким именем не было найдено."
	FRIENDS_BY_NAME           = "Успешно"
	USER_BY_LINK_VK_NOT_FOUND = "Пользователь по данной ссылке не найден"
	ADD_BY_ID_TEXT            = "Введите id пользователя VK."
	MAIN_IS_TOKEN_TEXT        = ""
	MAIN_NO_TOKEN_TEXT        = "Для слежки вам необходимо получить токен"
	ID_NOT_INT_ERROR          = "Вы неправильно ввели Id, оно должно быть целым числом"
)
