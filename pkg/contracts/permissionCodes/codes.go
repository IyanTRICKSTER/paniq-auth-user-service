package permissionCodes

type PermissionCode string

const (
	CREATE        PermissionCode = "c"
	LIST          PermissionCode = "l"
	UPDATE        PermissionCode = "u"
	DELETE        PermissionCode = "d"
	VIEW          PermissionCode = "v"
	VALIDATE_POST PermissionCode = "vp"
)
