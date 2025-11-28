package controllers

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type NpmController struct {
	Client 			*resty.Client
	StoragePath 	string
	NpmRegistry 	string
	LocalRegistry 	string
}

func NewNpmController(client *resty.Client, file string, npm_registry string, local_registry string) *NpmController {
	return &NpmController{
		Client: client,
		StoragePath: file,
		NpmRegistry: npm_registry,
		LocalRegistry: local_registry,
	}
}

func (n* NpmController) HandleRequest(c* fiber.Ctx) error {
	raw_path := c.Params("*");
	pkg_path, err := url.QueryUnescape(raw_path)

	if(err != nil) {
		return c.Status(403).JSON(fiber.Map{"message": "Badly formed request url"})
	}


	var file_path string;

	if(strings.HasSuffix(pkg_path, ".tgz")){
		println("Package asked")
		file_path = fmt.Sprintf("%s/packages/%s", n.StoragePath, pkg_path)
	}else{
		println("Meta asked")
		file_path = fmt.Sprintf("%s/meta/%s.json", n.StoragePath, pkg_path)
	}


	if _, err := os.Stat(file_path); err == nil {
		fmt.Println("Returned by cache")
		return c.SendFile(file_path)
	}else{
		fmt.Println("Returned by fetch")
		os.MkdirAll(filepath.Dir(file_path), 0755);
		resp, err := n.Client.R().
					SetOutput(file_path).
					Get(pkg_path)

		if(err != nil || resp.StatusCode() != 200){
			os.Remove(file_path)
			return c.Status(500).JSON(fiber.Map{"message": "There was an error while querying the upstream package"})
		}
		if(strings.Contains(resp.Header().Get("Content-Type"), "application/json")){
			meta_bytes, readErr := os.ReadFile(file_path)

			if(readErr != nil){
				os.Remove(file_path)
				return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error while attempting to update upstream in %s", pkg_path)})
			}

			meta_content := string(meta_bytes)
			meta_content = strings.ReplaceAll(meta_content, n.NpmRegistry, n.LocalRegistry)

			writeErr := os.WriteFile(file_path, ([]byte(meta_content)), 0755)
			if(writeErr != nil){
				return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error while attempting to update upstream in %s", pkg_path)})
			}
		}

		return c.SendFile(file_path)
	}
}
/*
func (n* NpmController) HandleMeta(c* fiber.Ctx) error {
	if _,err := os.Stat(file_path); err == nil{
		return c.SendFile(file_path, false)
	} else if errors.Is(err, os.ErrNotExist){
		resp, _ := n.Client.R().
				SetHeader("Accept", "application/json").
				Get(pkg_name)

		pkg_meta := strings.ReplaceAll(resp.String(), n.NpmRegistry, n.LocalRegistry) 

		file, _ := os.Create(file_path)
		file.Write([]byte(pkg_meta))
		fmt.Println("Returned file by fetch")

		c.Set("Content-Type", "application/json")
		return c.SendString(pkg_meta) 
	}else{
		return c.Status(500).SendString("Internal server error")
	}
}

func (n* NpmController) HandleTarBall(c* fiber.Ctx) error {
	pkg_path := fmt.Sprintf("%s/%s", c.Params("package"), c.Params("*"))
	full_path := fmt.Sprintf("%s/%s", n.StoragePath, pkg_path);

	os.MkdirAll(filepath.Dir(full_path),0755)

	resp, err := n.Client.R().
			SetOutput(full_path).
			Get(pkg_path)

	if err != nil || resp.StatusCode() != 200 {
		os.Remove(full_path)
	}

	return c.SendFile(full_path)
}*/