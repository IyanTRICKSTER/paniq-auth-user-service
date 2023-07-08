package database

import (
	"gorm.io/gorm"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
	"paniq-auth-user-service/pkg/entities"
	bcryptUtils "paniq-auth-user-service/pkg/utils/bcrypt"
)

func PermissionSeeder(conn *gorm.DB) []error {

	permissions := []entities.PermissionEntity{
		{
			Name:     "View User",
			Code:     permissionCodes.VIEW,
			Resource: apiResources.USER,
			Roles:    nil,
		},
		{
			Name:     "List User",
			Code:     permissionCodes.LIST,
			Resource: apiResources.USER,
			Roles:    nil,
		},
		{
			Name:     "Create User",
			Code:     permissionCodes.CREATE,
			Resource: apiResources.USER,
			Roles:    nil,
		},
		{
			Name:     "Update User",
			Code:     permissionCodes.UPDATE,
			Resource: apiResources.USER,
			Roles:    nil,
		},
		{
			Name:     "Delete User",
			Code:     permissionCodes.DELETE,
			Resource: apiResources.USER,
			Roles:    nil,
		},
		//Post permission
		{
			Name:     "View Post",
			Code:     permissionCodes.VIEW,
			Resource: apiResources.POST,
			Roles:    nil,
		},
		{
			Name:     "List Post",
			Code:     permissionCodes.LIST,
			Resource: apiResources.POST,
			Roles:    nil,
		},
		{
			Name:     "Create Post",
			Code:     permissionCodes.CREATE,
			Resource: apiResources.POST,
			Roles:    nil,
		},
		{
			Name:     "Update Post",
			Code:     permissionCodes.UPDATE,
			Resource: apiResources.POST,
			Roles:    nil,
		},
		{
			Name:     "Delete Post",
			Code:     permissionCodes.DELETE,
			Resource: apiResources.POST,
			Roles:    nil,
		},
		{
			Name:     "Validate Post",
			Code:     permissionCodes.VALIDATE_POST,
			Resource: apiResources.POST,
			Roles:    nil,
		},
	}

	var errs []error

	for i, entity := range permissions {
		entity.ID = uint(i + 1)
		err := conn.Create(&entity).Error
		errs = append(errs, err)
	}
	return errs
}
func RoleSeeder(conn *gorm.DB) []error {

	roles := []entities.RoleEntity{
		{
			Name:  "admin",
			Users: nil,
			Permissions: []entities.PermissionEntity{
				{
					ID:   1,
					Code: permissionCodes.VIEW,
				},
				{
					ID:   2,
					Code: permissionCodes.LIST,
				},
				{
					ID:   3,
					Code: permissionCodes.CREATE,
				},
				{
					ID:   4,
					Code: permissionCodes.UPDATE,
				},
				{
					ID:   5,
					Code: permissionCodes.DELETE,
				},
				{
					ID:   6,
					Code: permissionCodes.VIEW,
				},
				{
					ID:   7,
					Code: permissionCodes.LIST,
				},
				{
					ID:   8,
					Code: permissionCodes.CREATE,
				},
				{
					ID:   9,
					Code: permissionCodes.UPDATE,
				},
				{
					ID:   10,
					Code: permissionCodes.DELETE,
				},
				{
					ID:   11,
					Code: permissionCodes.VALIDATE_POST,
				},
			},
		},
		{
			Name:  "member",
			Users: nil,
			Permissions: []entities.PermissionEntity{
				{
					ID:   1,
					Code: permissionCodes.VIEW,
				},
				{
					ID:   2,
					Code: permissionCodes.LIST,
				},
				{
					ID:   3,
					Code: permissionCodes.CREATE,
				},
				{
					ID:   4,
					Code: permissionCodes.UPDATE,
				},
				{
					ID:   5,
					Code: permissionCodes.DELETE,
				},
				{
					ID:   6,
					Code: permissionCodes.VIEW,
				},
				{
					ID:   7,
					Code: permissionCodes.LIST,
				},
				{
					ID:   8,
					Code: permissionCodes.CREATE,
				},
				{
					ID:   9,
					Code: permissionCodes.UPDATE,
				},
				{
					ID:   10,
					Code: permissionCodes.DELETE,
				},
				{
					ID:   11,
					Code: permissionCodes.VALIDATE_POST,
				},
			},
		},
		{
			Name:  "user",
			Users: nil,
			Permissions: []entities.PermissionEntity{
				{
					ID:   1,
					Code: permissionCodes.VIEW,
				},
				{
					ID:   2,
					Code: permissionCodes.LIST,
				},
				{
					ID:   3,
					Code: permissionCodes.CREATE,
				},
				{
					ID:   4,
					Code: permissionCodes.UPDATE,
				},
				{
					ID:   5,
					Code: permissionCodes.DELETE,
				},
				{
					ID:   6,
					Code: permissionCodes.VIEW,
				},
				{
					ID:   7,
					Code: permissionCodes.LIST,
				},
				{
					ID:   8,
					Code: permissionCodes.CREATE,
				},
				{
					ID:   9,
					Code: permissionCodes.UPDATE,
				},
				{
					ID:   10,
					Code: permissionCodes.DELETE,
				},
				{
					ID:   11,
					Code: permissionCodes.VALIDATE_POST,
				},
			},
		},
	}

	var errs []error

	for i, entity := range roles {
		entity.ID = uint(i + 1)
		err := conn.Create(&entity).Error
		errs = append(errs, err)
	}
	return errs
}
func UserSeeder(conn *gorm.DB) []error {

	nims := []string{
		"11210910000004",
		"11210910000003",
		"11210910000002",
	}

	nips := []string{
		"11310030000004",
		"11310030000003",
	}

	users := []entities.UserEntity{
		{
			Role:       entities.RoleEntity{},
			RoleID:     uint(1),
			Username:   "iyan pratama",
			Password:   bcryptUtils.NewHashFunction().Hash("iyan12345"),
			Email:      "iyanpratama2002@gmail.com",
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        &nims[0],
			NIP:        nil,
			Major:      "CS Degree",
			ResetToken: "",
		},
		{
			Role:       entities.RoleEntity{},
			RoleID:     uint(2),
			Username:   "septian putra pratama",
			Password:   bcryptUtils.NewHashFunction().Hash("iyan12345"),
			Email:      "septianputra.pratama02@gmail.com",
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        &nims[1],
			NIP:        &nips[0],
			Major:      "CS Degree",
			ResetToken: "",
		},
		{
			Role:       entities.RoleEntity{},
			RoleID:     uint(3),
			Username:   "akiyan2002",
			Password:   bcryptUtils.NewHashFunction().Hash("iyan12345"),
			Email:      "akiyan2002@gmail.com",
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        &nims[2],
			NIP:        &nips[1],
			Major:      "CS Degree",
			ResetToken: "",
		},
	}

	var errs []error

	for i, entity := range users {
		entity.ID = uint(i + 1)
		err := conn.Create(&entity).Error
		errs = append(errs, err)
	}

	return errs

}
