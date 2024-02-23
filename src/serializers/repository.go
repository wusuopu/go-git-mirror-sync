package serializers

import "app/models"

type RepositorySerializer struct {}

func (s *RepositorySerializer) Show(data models.Repository) map[string]interface{} {
	var ret = make(map[string]interface{})
	ret["ID"] = data.ID
	ret["Name"] = data.Name
	ret["Alias"] = data.Alias
	ret["Url"] = data.Url
	ret["AuthType"] = data.AuthType
	ret["Username"] = data.Username
	ret["Password"] = data.Password
	ret["SSHKey"] = data.SSHKey
	ret["PulledAt"] = data.PulledAt
	ret["InitedAt"] = data.InitedAt
	ret["LastError"] = data.LastError
	ret["CreatedAt"] = data.CreatedAt
	ret["UpdatedAt"] = data.UpdatedAt
	return ret

}
func (s *RepositorySerializer) List(data []models.Repository) []map[string]interface{} {
	var ret []map[string]interface{}
	for _, item := range data {
		el := s.Show(item)
		delete(el, "CreatedAt")
		delete(el, "UpdatedAt")
		ret = append(ret, el)
	}
	return ret
}