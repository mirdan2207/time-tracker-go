package migrations

import (
	"time"
	"time-tracker-go/models"

	"gorm.io/gorm"
)

func clean(db *gorm.DB) {
	// Удаление данных из таблиц
	db.Exec("DELETE FROM tasks;")
	db.Exec("DELETE FROM users;")
	db.Exec("DELETE FROM peoples;")
}

// Seed populates the database with initial user and task data.
func Seed(db *gorm.DB) {
	clean(db)

	users := []models.User{
		{PassportNumber: "1234 567890", Surname: "Ivanov", Name: "Ivan", Patronymic: "Ivanovich", Address: "г. Москва, ул. Ленина, д. 5, кв. 1"},
		{PassportNumber: "2345 678901", Surname: "Petrov", Name: "Petr", Patronymic: "Petrovich", Address: "г. Санкт-Петербург, Невский проспект, д. 10, кв. 2"},
		{PassportNumber: "3456 789012", Surname: "Sidorov", Name: "Sidr", Patronymic: "Sidorovich", Address: "г. Казань, ул. Баумана, д. 15, кв. 3"},
		{PassportNumber: "4567 890123", Surname: "Smirnov", Name: "Sergey", Patronymic: "Sergeevich", Address: "г. Новосибирск, ул. Красный проспект, д. 20, кв. 4"},
		{PassportNumber: "5678 901234", Surname: "Kuznetsov", Name: "Nikolay", Patronymic: "Nikolaevich", Address: "г. Екатеринбург, ул. Ленина, д. 25, кв. 5"},
		{PassportNumber: "6789 012345", Surname: "Popov", Name: "Aleksey", Patronymic: "Alexeevich", Address: "г. Нижний Новгород, ул. Горького, д. 30, кв. 6"},
		{PassportNumber: "7890 123456", Surname: "Vasiliev", Name: "Dmitry", Patronymic: "Dmitrievich", Address: "г. Самара, ул. Ленина, д. 35, кв. 7"},
		{PassportNumber: "8901 234567", Surname: "Mikhailov", Name: "Mikhail", Patronymic: "Mikhailovich", Address: "г. Омск, ул. Ленина, д. 40, кв. 8"},
		{PassportNumber: "9012 345678", Surname: "Fedorov", Name: "Fedor", Patronymic: "Fedorovich", Address: "г. Казань, ул. Ленина, д. 45, кв. 9"},
		{PassportNumber: "0123 456789", Surname: "Kovalev", Name: "Vladimir", Patronymic: "Vladimirovich", Address: "г. Челябинск, ул. Ленина, д. 50, кв. 10"},
		{PassportNumber: "1111 111111", Surname: "Nikolaev", Name: "Alexandr", Patronymic: "Alexandrovich", Address: "г. Уфа, ул. Ленина, д. 55, кв. 11"},
		{PassportNumber: "2222 222222", Surname: "Semenov", Name: "Semen", Patronymic: "Semenovich", Address: "г. Волгоград, ул. Ленина, д. 60, кв. 12"},
		{PassportNumber: "3333 333333", Surname: "Stepanov", Name: "Stepan", Patronymic: "Stepanovich", Address: "г. Пермь, ул. Ленина, д. 65, кв. 13"},
		{PassportNumber: "4444 444444", Surname: "Pavlov", Name: "Pavel", Patronymic: "Pavlovich", Address: "г. Красноярск, ул. Ленина, д. 70, кв. 14"},
		{PassportNumber: "5555 555555", Surname: "Bogdanov", Name: "Bogdan", Patronymic: "Bogdanovich", Address: "г. Саратов, ул. Ленина, д. 75, кв. 15"},
		{PassportNumber: "6666 666666", Surname: "Novikov", Name: "Nikita", Patronymic: "Nikitich", Address: "г. Воронеж, ул. Ленина, д. 80, кв. 16"},
		{PassportNumber: "7777 777777", Surname: "Gerasimov", Name: "Gerasim", Patronymic: "Gerasimovich", Address: "г. Тольятти, ул. Ленина, д. 85, кв. 17"},
		{PassportNumber: "8888 888888", Surname: "Malyshev", Name: "Malysh", Patronymic: "Malyshovich", Address: "г. Пенза, ул. Ленина, д. 90, кв. 18"},
		{PassportNumber: "9999 999999", Surname: "Gusev", Name: "Gusev", Patronymic: "Gusevich", Address: "г. Киров, ул. Ленина, д. 95, кв. 19"},
		{PassportNumber: "1010 101010", Surname: "Kiselev", Name: "Kisel", Patronymic: "Kiselich", Address: "г. Новокузнецк, ул. Ленина, д. 100, кв. 20"},
	}

	// Save users in the database
	for _, user := range users {
		db.Create(&user)

		tasks := []models.Task{
			{UserID: user.ID, Description: "Task 1", StartTime: time.Now().Add(-10 * time.Hour), EndTime: time.Now().Add(-9 * time.Hour), Duration: 60},
			{UserID: user.ID, Description: "Task 2", StartTime: time.Now().Add(-9 * time.Hour), EndTime: time.Now().Add(-7 * time.Hour), Duration: 120},
			{UserID: user.ID, Description: "Task 3", StartTime: time.Now().Add(-2 * time.Hour), EndTime: time.Now().Add(-5 * time.Hour), Duration: 180},
		}

		for _, task := range tasks {
			db.Create(&task)
		}
	}

	people := []models.People{
		{PassportSeries: 1001, PassportNumber: 100001, Surname: "Smith", Name: "John", Patronymic: "Johnson", Address: "г. Москва, ул. Пушкина, д. 1, кв. 1"},
		{PassportSeries: 1002, PassportNumber: 100002, Surname: "Johnson", Name: "Jane", Patronymic: "Janet", Address: "г. Санкт-Петербург, ул. Пушкина, д. 2, кв. 2"},
		{PassportSeries: 1003, PassportNumber: 100003, Surname: "Brown", Name: "Charlie", Patronymic: "Charles", Address: "г. Казань, ул. Пушкина, д. 3, кв. 3"},
		{PassportSeries: 1004, PassportNumber: 100004, Surname: "Davis", Name: "Alice", Patronymic: "Alicia", Address: "г. Новосибирск, ул. Пушкина, д. 4, кв. 4"},
		{PassportSeries: 1005, PassportNumber: 100005, Surname: "Miller", Name: "Robert", Patronymic: "Roberts", Address: "г. Екатеринбург, ул. Пушкина, д. 5, кв. 5"},
		{PassportSeries: 1006, PassportNumber: 100006, Surname: "Wilson", Name: "James", Patronymic: "Jim", Address: "г. Нижний Новгород, ул. Пушкина, д. 6, кв. 6"},
		{PassportSeries: 1007, PassportNumber: 100007, Surname: "Moore", Name: "Lisa", Patronymic: "Lilian", Address: "г. Самара, ул. Пушкина, д. 7, кв. 7"},
		{PassportSeries: 1008, PassportNumber: 100008, Surname: "Taylor", Name: "David", Patronymic: "Dave", Address: "г. Омск, ул. Пушкина, д. 8, кв. 8"},
		{PassportSeries: 1009, PassportNumber: 100009, Surname: "Anderson", Name: "Eve", Patronymic: "Evans", Address: "г. Челябинск, ул. Пушкина, д. 9, кв. 9"},
		{PassportSeries: 1010, PassportNumber: 100010, Surname: "Thomas", Name: "Frank", Patronymic: "Franklin", Address: "г. Уфа, ул. Пушкина, д. 10, кв. 10"},
		{PassportSeries: 1011, PassportNumber: 100011, Surname: "Jackson", Name: "George", Patronymic: "Georgievich", Address: "г. Волгоград, ул. Пушкина, д. 11, кв. 11"},
		{PassportSeries: 1012, PassportNumber: 100012, Surname: "White", Name: "Mary", Patronymic: "Maryam", Address: "г. Пермь, ул. Пушкина, д. 12, кв. 12"},
		{PassportSeries: 1013, PassportNumber: 100013, Surname: "Harris", Name: "William", Patronymic: "Williams", Address: "г. Красноярск, ул. Пушкина, д. 13, кв. 13"},
		{PassportSeries: 1014, PassportNumber: 100014, Surname: "Martin", Name: "Patricia", Patronymic: "Patricius", Address: "г. Саратов, ул. Пушкина, д. 14, кв. 14"},
		{PassportSeries: 1015, PassportNumber: 100015, Surname: "Thompson", Name: "Richard", Patronymic: "Richards", Address: "г. Воронеж, ул. Пушкина, д. 15, кв. 15"},
	}

	// Save people in the database
	for _, person := range people {
		db.Create(&person)
	}
}
