package role

type Action string

const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
	ActionManage Action = "manage"
)

func (a Action) IsValid() bool {
	switch a {
	case ActionRead, ActionWrite, ActionDelete, ActionManage:
		return true
	}

	return false
}
