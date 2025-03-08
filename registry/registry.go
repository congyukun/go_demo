package registry

import "go_demo/controllers"

type RegistryController struct {
	ArticleController *controllers.ArticleController
	LoginController   *controllers.LoginController
}

func NewControllerRegistry() *RegistryController {
	return &RegistryController{
		ArticleController: &controllers.ArticleController{},
		LoginController:   &controllers.LoginController{},
	}
}
