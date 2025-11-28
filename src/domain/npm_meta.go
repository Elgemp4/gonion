package domain

import (
	"fmt"
)

type NpmMeta struct{
	Name 		string
}
/*
func (meta* NPMMeta) isCached() bool{
	_, err := os.Stat(meta.CachePath); 
	return err == nil
}*/

func NewNpmMeta(name string) *NpmMeta {
	return &NpmMeta{
		Name: name,
	}
}

func (m *NpmMeta) RelativeCachedPath() string {
	return fmt.Sprintf("/meta/%s.json", m.Name)
}

func (m *NpmMeta) RessourceName() string {
	return m.Name
}

/*
func (meta* NPMMeta) Load() (string, error){
	if meta.isCached() {
		fmt.Println("Returned by cache")
		return meta.CachePath, nil
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

			writeErr := os.WriteFile(file_path, ([]byte(meta_content)), 0644)
			if(writeErr != nil){
				return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error while attempting to update upstream in %s", pkg_path)})
			}
		}

		return c.SendFile(file_path)
	}
}*/