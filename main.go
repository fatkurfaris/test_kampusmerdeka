package main

import (
	"log"
	"ruang_belajar/auth"
	"ruang_belajar/handler"
	myclasses "ruang_belajar/models/MyClasses"
	"ruang_belajar/models/articles"
	"ruang_belajar/models/classes"
	"ruang_belajar/models/learners"
	"ruang_belajar/models/tutors"
	"ruang_belajar/repository/database"
	"ruang_belajar/repository/drivers/mysql"
	"ruang_belajar/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func DbMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&learners.Learner{}, &tutors.Tutor{}, &articles.Article{}, &classes.Class{}, &myclasses.MyClass{})
	if err != nil {
		panic(err)
	}
}

func main() {
	mysqlConfig := mysql.ConfigDb{
		DB_Username: viper.GetString(`databases.mysql.user`),
		DB_Password: viper.GetString(`databases.mysql.password`),
		DB_Host:     viper.GetString(`databases.mysql.host`),
		DB_Port:     viper.GetString(`databases.mysql.port`),
		DB_Database: viper.GetString(`databases.mysql.dbname`),
	}

	db := mysqlConfig.InitialDb()
	DbMigrate(db)

	configJWT := viper.GetString(`jwt.SECRET_KEY`)

	tutorRepository := database.NewTutorRepository(db)
	learnerRepository := database.NewLearnerRepository(db)
	classRepository := database.NewClassRepository(db)
	myclassRepository := database.NewMyClassRepository(db)
	articleRepository := database.NewArticleRepository(db)

	authService := auth.NewService(configJWT)
	tutorService := service.NewTutorService(tutorRepository)
	learnerService := service.NewLeranerService(learnerRepository)
	authMiddleware := auth.AuthMiddleware(authService, tutorService, learnerService)
	tutor := auth.Permission(&auth.Role{Roles: "tutor"})
	learner := auth.Permission(&auth.Role{Roles: "learner"})
	classService := service.NewClassService(classRepository, *tutorService)
	myclassService := service.NewMyClassService(myclassRepository)
	articleService := service.NewArticleService(articleRepository, *tutorService)

	userHandler := handler.NewUserHandler(tutorService, learnerService, authService)
	tutorHandler := handler.NewTutorHandler(tutorService)
	learnerHandler := handler.NewLearnerHandler(learnerService)
	classHandler := handler.NewClassHandler(classService)
	myclassHandler := handler.NewMyClassHandler(myclassService, classService)
	articleHandler := handler.NewArticleHandler(articleService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)

	api.PUT("/tutors/:id", authMiddleware, tutor, tutorHandler.UpdateTutor)
	api.GET("/tutors", authMiddleware, tutor, tutorHandler.FetchTutor)

	api.PUT("/learners/:id", authMiddleware, learner, learnerHandler.UpdateLearner)
	api.GET("/learners", authMiddleware, learner, learnerHandler.FetchLearner)

	api.POST("/classes", authMiddleware, tutor, classHandler.CreateClass)
	api.GET("/classes", authMiddleware, classHandler.GetAll)

	api.POST("/myclasses", authMiddleware, learner, myclassHandler.CreateMyClass)
	api.GET("/myclasses", authMiddleware, learner, myclassHandler.GetAllMyClass)

	api.POST("/articles", authMiddleware, tutor, articleHandler.CreateArticle)
	api.GET("/articles", authMiddleware, articleHandler.GetAll)
	api.DELETE("/articles/:id", authMiddleware, tutor, articleHandler.Delete)

	router.Run()
}
