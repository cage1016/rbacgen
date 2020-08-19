package gen

import (
	"fmt"
	"os"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Action struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdateAt    time.Time
}

func NewAction(name string) Action {
	id, _ := ID()
	return Action{
		ID:          id,
		Name:        name,
		Description: name,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
}

func (a Action) String() string {
	return fmt.Sprintf(`INSERT INTO public.opa_actions (id, name, description, created_at, updated_at) VALUES ('%s', '%s','%s', '%s', '%s');`, a.ID, a.Name, a.Description, a.CreatedAt.Format(time.RFC3339), a.UpdateAt.Format(time.RFC3339))
}

type Resource struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (a Resource) String() string {
	return fmt.Sprintf(`INSERT INTO public.opa_resources (id, name, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s');`, a.ID, a.Name, a.CreatedAt.Format(time.RFC3339), a.UpdateAt.Format(time.RFC3339))
}

func NewResource(name string) Resource {
	id, _ := ID()
	return Resource{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}

type Permission struct {
	ID         string
	ActionID   string
	ResourceID string
	CreatedAt  time.Time
	UpdateAt   time.Time
	Key        string
}

func (a Permission) String() string {
	return fmt.Sprintf(`INSERT INTO public.opa_permissions (id, action_id, resource_id, created_at, updated_at, key) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');`, a.ID, a.ActionID, a.ResourceID, a.CreatedAt.Format(time.RFC3339), a.UpdateAt.Format(time.RFC3339), a.Key)
}

func NewPermission(action Action, resource Resource) Permission {
	id, _ := ID()
	return Permission{
		ID:         id,
		ActionID:   action.ID,
		ResourceID: resource.ID,
		CreatedAt:  time.Now(),
		UpdateAt:   time.Now(),
		Key:        fmt.Sprintf("%s:%s", action.Name, resource.Name),
	}
}

type Role struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (a Role) String() string {
	return fmt.Sprintf(`INSERT INTO public.opa_roles (id, name, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s');`, a.ID, a.Name, a.CreatedAt.Format(time.RFC3339), a.UpdateAt.Format(time.RFC3339))
}

func NewRole(name string) Role {
	id, _ := ID()
	return Role{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}

type RolePermission struct {
	RoleID       string
	PermissionID string
}

func (a RolePermission) String() string {
	return fmt.Sprintf(`INSERT INTO public.opa_role_permissions (role_id, permission_id) VALUES ('%s', '%s');`, a.RoleID, a.PermissionID)
}

func NewRolePermission(roleID, permissionID string) RolePermission {
	return RolePermission{roleID, permissionID}
}

func ID() (string, error) {
	return gonanoid.Nanoid(21)
}

func Gen(f *os.File, data map[string]interface{}) error {
	amap := map[string]Action{}
	rmap := map[string]Resource{}
	rolesMap := map[string][]string{}

	if data != nil {
		// actions
		for k, v := range data["actions"].(map[string]interface{}) {
			amap[strings.ToUpper(k)] = NewAction(v.(string))
		}

		// resources
		for _, v := range data["resources"].([]interface{}) {
			rmap[v.(string)] = NewResource(v.(string))
		}

		for k, v := range data["roles"].(map[string]interface{}) {
			for _, item := range v.([]interface{}) {
				rolesMap[k] = append(rolesMap[k], item.(string))
			}
		}
	} else {
		amap["R"] = NewAction("read")
		amap["W"] = NewAction("write")
		amap["E"] = NewAction("edit")
		amap["D"] = NewAction("delete")

		rmap["rbac"] = NewResource("rbac")
		rmap["organization"] = NewResource("organization")
		rmap["invitation"] = NewResource("invitation")
		rmap["report"] = NewResource("report")
		rmap["store"] = NewResource("store")

		rolesMap["owner"] = []string{
			"rbac:RWED",
			"organization:RWED",
			"invitation:R",
			"report:R",
			"store:RWED",
		}
		rolesMap["editor"] = []string{
			"rbac:RWE",
			"organization:RWE",
			"invitation:R",
			"report:R",
			"store:RWE",
		}
		rolesMap["viewer"] = []string{
			"rbac:R",
			"organization:R",
			"invitation:R",
			"report:R",
			"store:R",
		}
		rolesMap["reporter"] = []string{
			"report:R",
		}
	}

	permissions := []Permission{}
	role_permissions := []RolePermission{}
	roles := []Role{}

	pkey := map[string]Permission{}
	for rk, rules := range rolesMap {
		r := NewRole(rk)
		roles = append(roles, r)
		for _, rule := range rules {
			ss := strings.Split(rule, ":")
			pk, actions := ss[0], ss[1]
			for _, char := range actions {
				action := amap[fmt.Sprintf("%c", char)]
				resource := rmap[pk]
				buf := fmt.Sprintf("%s:%s", action.Name, resource.Name)

				var p Permission
				if value, ok := pkey[buf]; ok {
					p = value
				} else {
					p = NewPermission(action, resource)
					permissions = append(permissions, p)
					pkey[buf] = p
				}

				role_permissions = append(role_permissions, NewRolePermission(r.ID, p.ID))
			}
		}
	}

	for _, m := range amap {
		_, err := f.Write([]byte(fmt.Sprintf("%v\n", m)))
		if err != nil {
			return err
		}
	}
	f.WriteString("\n")

	for _, m := range rmap {
		_, err := f.Write([]byte(fmt.Sprintf("%v\n", m)))
		if err != nil {
			return err
		}
	}
	f.WriteString("\n")

	for _, m := range permissions {
		_, err := f.Write([]byte(fmt.Sprintf("%v\n", m)))
		if err != nil {
			return err
		}
	}
	f.WriteString("\n")

	for _, m := range roles {
		_, err := f.Write([]byte(fmt.Sprintf("%v\n", m)))
		if err != nil {
			return err
		}
	}
	f.WriteString("\n")

	for _, m := range role_permissions {
		_, err := f.Write([]byte(fmt.Sprintf("%v\n", m)))
		if err != nil {
			return err
		}
	}

	return nil
}
